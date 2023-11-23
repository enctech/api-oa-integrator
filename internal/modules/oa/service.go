package oa

import (
	"api-oa-integrator/internal/database"
	"api-oa-integrator/internal/modules/integrator"
	"api-oa-integrator/utils"
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/sqlc-dev/pqtype"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func handleIdentificationEntry(c echo.Context, job *Job, metadata *RequestMetadata) {
	if job.JobType != "IDENTIFICATION" || job.TimeAndPlace.Device.DeviceType != "ENTRY" {
		return
	}
	lpn := job.MediaDataList.Identifier.Name
	lane := job.TimeAndPlace.Device.DeviceNumber
	btid := uuid.New().String()
	customerId := encryptLpn(lpn)

	jsonStr, err := json.Marshal(map[string]any{
		"steps": "identification_entry_start",
	})
	if err != nil {
		zap.L().Sugar().Info("Error Marshal ", err)
		go sendEmptyFinalMessage(metadata)
		return
	}
	data, err := database.New(database.D()).CreateOATransaction(context.Background(), database.CreateOATransactionParams{
		Businesstransactionid: btid,
		Device:                sql.NullString{String: metadata.device, Valid: true},
		Facility:              sql.NullString{String: metadata.facility, Valid: true},
		Jobid:                 sql.NullString{String: metadata.jobId, Valid: true},
		Lpn:                   sql.NullString{String: lpn, Valid: true},
		Customerid:            sql.NullString{String: customerId, Valid: true},
		Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
		EntryLane:             sql.NullString{String: lane, Valid: true},
	})

	if err != nil {
		zap.L().Sugar().Info("Error create oa transaction ", err)
		go sendEmptyFinalMessage(metadata)
		return
	}

	err = integrator.VerifyVehicle(metadata.vendor, metadata.facility, lpn, lane)
	if err != nil {
		zap.L().Sugar().Info("Error integrator.VerifyVehicle ", err)
		go sendEmptyFinalMessage(metadata)
		return
	}

	go func() {
		cfg, err := database.New(database.D()).GetIntegratorConfigByName(context.Background(), sql.NullString{String: metadata.vendor, Valid: true})
		if err != nil {
			go sendEmptyFinalMessage(metadata)
			return
		}
		jsonStr, err := json.Marshal(map[string]any{
			"steps": "identification_entry_done",
		})
		if err != nil {
			zap.L().Sugar().Info("Error Marshal ", err)
			go sendEmptyFinalMessage(metadata)
			return
		}
		_, err = database.New(database.D()).UpdateOATransaction(context.Background(), database.UpdateOATransactionParams{
			Businesstransactionid: data.Businesstransactionid,
			Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
		})
		sendFinalMessageCustomer(metadata, FMCReq{
			Identifier:          Identifier{Name: lpn},
			BusinessTransaction: &BusinessTransaction{ID: btid},
			CustomerInformation: &CustomerInformation{
				Customer: Customer{
					CustomerId:    data.Customerid.String,
					CustomerGroup: cfg.Name.String,
				},
			},
			PaymentInformation: BuildPaymentInformation(nil),
		})
	}()

	if c.Request().Header.Get("istest") != "" {
		c.Response().Header().Set("btid", btid)
		c.Response().Header().Set("customerId", data.Customerid.String)
	}
}

func handleLeaveLoopEntry(job *Job, metadata *RequestMetadata) {
	if job.JobType != "LEAVE_LOOP" || job.TimeAndPlace.Device.DeviceType != "ENTRY" {
		return
	}
	if job.BusinessTransaction == nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	if job.BusinessTransaction.ID == "" {
		go sendEmptyFinalMessage(metadata)
		return
	}

	jsonStr, err := json.Marshal(map[string]any{
		"steps": "leave_loop_entry_done",
	})

	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	_, _ = database.New(database.D()).UpdateOATransaction(context.Background(), database.UpdateOATransactionParams{
		Businesstransactionid: job.BusinessTransaction.ID,
		Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
	})
	go sendEmptyFinalMessage(metadata)
}

func handleIdentificationExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "IDENTIFICATION" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	if job.BusinessTransaction == nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	if job.BusinessTransaction.ID == "" {
		go sendEmptyFinalMessage(metadata)
		return
	}

	jsonStr, err := json.Marshal(map[string]any{
		"steps": "identification_exit_done",
	})

	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	lane := job.TimeAndPlace.Device.DeviceNumber
	oaTxn, err := database.New(database.D()).UpdateOATransaction(context.Background(), database.UpdateOATransactionParams{
		Businesstransactionid: job.BusinessTransaction.ID,
		Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
		ExitLane:              sql.NullString{String: lane, Valid: true},
	})

	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	go func() {
		cfg, err := database.New(database.D()).GetIntegratorConfigByName(context.Background(), sql.NullString{String: metadata.vendor, Valid: true})
		if err != nil {
			go sendEmptyFinalMessage(metadata)
			return
		}
		sendFinalMessageCustomer(metadata, FMCReq{
			Identifier:          Identifier{Name: oaTxn.Lpn.String},
			BusinessTransaction: &BusinessTransaction{ID: oaTxn.Businesstransactionid},
			CustomerInformation: &CustomerInformation{Customer: Customer{
				CustomerId:    oaTxn.Customerid.String,
				CustomerGroup: cfg.Name.String,
			}},
			PaymentInformation: BuildPaymentInformation(nil),
		})
	}()
}

func handlePaymentExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "PAYMENT" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	if job.BusinessTransaction == nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	jsonStr, err := json.Marshal(map[string]any{
		"steps": "payment_exit_start",
	})

	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	lpn := job.MediaDataList.Identifier.Name

	oaTxn, err := database.New(database.D()).UpdateOATransaction(context.Background(), database.UpdateOATransactionParams{
		Businesstransactionid: job.BusinessTransaction.ID,
		Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
	})

	amount, err := strconv.ParseFloat(job.PaymentData.OriginalAmount.Amount, 64)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}
	amountConv := amount / 100

	var customerInformation *CustomerInformation
	if job.CustomerInformation != nil && job.CustomerInformation.Customer != (Customer{}) {
		customerInformation = &CustomerInformation{Customer: job.CustomerInformation.Customer}
	}

	sendZeroAmount := func() {
		sendFinalMessageCustomer(metadata, FMCReq{
			Identifier:          Identifier{Name: oaTxn.Lpn.String},
			BusinessTransaction: &BusinessTransaction{ID: oaTxn.Businesstransactionid},
			CustomerInformation: customerInformation,
			PaymentInformation: BuildPaymentInformation(&PaymentData{
				OriginalAmount: OriginalAmount{
					VatRate: job.PaymentData.OriginalAmount.VatRate,
					Amount:  "0.00",
				},
			}),
		})
	}

	err = integrator.PerformTransaction(metadata.vendor, metadata.facility, lpn, oaTxn.EntryLane.String, oaTxn.ExitLane.String, oaTxn.CreatedAt, amountConv)
	if err != nil {
		jsonStr, err = json.Marshal(map[string]any{
			"steps": "payment_exit_error",
			"error": err.Error(),
		})
		go sendZeroAmount()
		return
	}
	jsonStr, err = json.Marshal(map[string]any{
		"steps": "payment_exit_done",
	})

	if err != nil {
		go sendZeroAmount()
		return
	}

	oaTxn, err = database.New(database.D()).UpdateOATransaction(context.Background(), database.UpdateOATransactionParams{
		Businesstransactionid: job.BusinessTransaction.ID,
		Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
	})

	go sendFinalMessageCustomer(metadata, FMCReq{
		Identifier:          Identifier{Name: oaTxn.Lpn.String},
		BusinessTransaction: &BusinessTransaction{ID: oaTxn.Businesstransactionid},
		CustomerInformation: customerInformation,
		PaymentInformation: BuildPaymentInformation(&PaymentData{
			OriginalAmount: OriginalAmount{
				VatRate: job.PaymentData.OriginalAmount.VatRate,
				Amount:  job.PaymentData.OriginalAmount.Amount,
			},
		}),
	})
}

func handleLeaveLoopExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "LEAVE_LOOP" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	lpn := job.MediaDataList.Identifier.Name

	oaTxn, err := database.New(database.D()).GetOATransaction(context.Background(), job.BusinessTransaction.ID)

	var extra map[string]any

	err = json.Unmarshal(oaTxn.Extra.RawMessage, &extra)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	if extra["steps"] != "payment_exit_done" {
		err = integrator.PerformTransaction(metadata.vendor, metadata.facility, lpn, oaTxn.EntryLane.String, oaTxn.ExitLane.String, oaTxn.CreatedAt, 0.00)
	}

	go sendEmptyFinalMessage(metadata)
}

type FMCReq struct {
	PaymentInformation  *PaymentInformation
	BusinessTransaction *BusinessTransaction
	CustomerInformation *CustomerInformation
	Identifier          Identifier
}

func sendFinalMessageCustomer(metadata *RequestMetadata, in FMCReq) {
	config, err := database.New(database.D()).GetSnbConfigByFacilityAndDevice(context.Background(), database.GetSnbConfigByFacilityAndDeviceParams{
		Device:   metadata.device,
		Facility: metadata.facility,
	})

	if err != nil {
		fmt.Println("Error get config", err)
		return
	}

	vendor, err := database.New(database.D()).GetIntegratorConfigByName(context.Background(), sql.NullString{String: metadata.vendor, Valid: true})

	counting := "RESERVED"
	xmlData, err := xml.Marshal(&FinalMessageCustomer{
		PaymentInformation: in.PaymentInformation,
		ProviderInformation: &ProviderInformation{
			Provider: Provider{
				ProviderId:   fmt.Sprintf("%v", vendor.ProviderID.Int32),
				ProviderName: vendor.Name.String,
			},
		},
		Counting: &counting,
		MediaDataList: &[]MediaDataList{
			{
				MediaType:  "LICENSE_PLATE",
				Identifier: in.Identifier,
			},
		},
		CustomerInformation: in.CustomerInformation,
		BusinessTransaction: in.BusinessTransaction,
	})
	if err != nil {
		fmt.Println("Error marshaling XML data:", err)
		return
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%v/AuthorizationServiceSB/%v/%v/%v/finalmessage", config.Endpoint.String, metadata.facility, metadata.device, metadata.jobId), bytes.NewBuffer(xmlData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/xml")
	req.SetBasicAuth(config.Username.String, config.Password.String)

	client := &http.Client{}
	client.Transport = &utils.LoggingRoundTripper{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	return
}

func sendEmptyFinalMessage(metadata *RequestMetadata) {
	config, err := database.New(database.D()).GetSnbConfigByFacilityAndDevice(context.Background(), database.GetSnbConfigByFacilityAndDeviceParams{
		Device:   metadata.device,
		Facility: metadata.facility,
	})

	if err != nil {
		fmt.Println("Error get config", err)
		return
	}

	xmlData, err := xml.Marshal(&FinalMessageCustomer{})
	if err != nil {
		fmt.Println("Error marshaling XML data:", err)
		return
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%v/AuthorizationServiceSB/%v/%v/%v/finalmessage", config.Endpoint.String, metadata.facility, metadata.device, metadata.jobId), bytes.NewBuffer(xmlData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/xml")
	req.SetBasicAuth(config.Username.String, config.Password.String)

	client := &http.Client{}
	client.Transport = &utils.LoggingRoundTripper{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}

func CheckSystemAvailability(facility, device string) error {
	config, err := database.New(database.D()).GetSnbConfigByFacilityAndDevice(context.Background(), database.GetSnbConfigByFacilityAndDeviceParams{
		Device:   device,
		Facility: facility,
	})

	if err != nil {
		fmt.Println("Error get config", err)
		return err
	}

	// Fake request body. Request body is required for this endpoint.
	// We don't really care about the response. We're good if there is response.
	xmlOut := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><version>
	<customerVersion>%v</customerVersion>
	<sbAuthorizationServiceVersion>2.5.6</sbAuthorizationServiceVersion>
	<configuration>
	</configuration>
	</version>`, viper.GetString("app.version")))

	req, err := http.NewRequest("PUT", fmt.Sprintf("%v/AuthorizationServiceSB/version", config.Endpoint.String), bytes.NewBuffer(xmlOut))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/xml")
	req.SetBasicAuth(config.Username.String, config.Password.String)

	client := &http.Client{}
	client.Transport = &utils.LoggingRoundTripper{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
