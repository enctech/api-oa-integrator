package oa

import (
	"api-oa-integrator/database"
	"api-oa-integrator/internal/modules/integrator"
	"api-oa-integrator/logger"
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
	"maps"
	"net/http"
	"strconv"
	"sync"
	"time"
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
		logger.LogData("error", fmt.Sprintf("Error Marshal %v", err), nil)
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
		logger.LogData("error", fmt.Sprintf("Error create oa transaction %v", err), nil)
		go sendEmptyFinalMessage(metadata)
		return
	}

	configs, err := database.New(database.D()).GetIntegratorConfigs(context.Background())
	if err != nil {
		logger.LogData("error", fmt.Sprintf("Error create oa transaction %v", err), nil)
		go sendEmptyFinalMessage(metadata)
		return
	}

	var wg sync.WaitGroup
	vendorChannel := make(chan string, 1)

	for i := range configs {
		wg.Add(1)
		go func(vendorName string) {
			defer wg.Done()
			err = integrator.VerifyVehicle(vendorName, metadata.facility, lpn, lane)

			if err != nil {
				logger.LogData("error", fmt.Sprintf("Error integrator.VerifyVehicle %v", err), nil)
				jsonStr, err := json.Marshal(map[string]any{
					"steps": "identification_entry_error",
					"error": err.Error(),
				})
				if err != nil {
					logger.LogData("error", fmt.Sprintf("Error marshal %v", err), nil)
					go sendEmptyFinalMessage(metadata)
					return
				}
				_, err = database.New(database.D()).UpdateOATransaction(context.Background(), database.UpdateOATransactionParams{
					Businesstransactionid: data.Businesstransactionid,
					Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
				})
			} else {
				select {
				case vendorChannel <- vendorName:
				default: // Ignore if already filled
				}
			}
		}(configs[i].Name.String)
	}

	wg.Wait()
	close(vendorChannel)

	var successfulVendor string
	select {
	case successfulVendor = <-vendorChannel:
		fmt.Printf("Success from: %s\n", successfulVendor)
	default:
		sendEmptyFinalMessage(metadata)
		return
	}

	go func() {
		for _, cfg := range configs {
			vendor := cfg.Name.String
			if vendor != successfulVendor {
				go integrator.CancelEntry(vendor, metadata.facility)
			}
		}
	}()

	go func() {
		if successfulVendor == "" {
			go sendEmptyFinalMessage(metadata)
			return
		}

		jsonStr, err := json.Marshal(map[string]any{
			"steps": "identification_entry_done",
		})
		if err != nil {
			logger.LogData("error", fmt.Sprintf("Error Marshal %v", err), nil)
			go sendEmptyFinalMessage(metadata)
			return
		}

		config, err := utils.FirstWhere(configs, func(config database.IntegratorConfig) bool {
			return config.Name.String == successfulVendor
		})

		_, err = database.New(database.D()).UpdateOATransaction(context.Background(), database.UpdateOATransactionParams{
			Businesstransactionid: data.Businesstransactionid,
			Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
			IntegratorID: uuid.NullUUID{
				UUID: config.ID, Valid: true,
			},
		})
		sendFinalMessageCustomer(metadata, FMCReq{
			Identifier:          Identifier{Name: lpn},
			BusinessTransaction: &BusinessTransaction{ID: btid},
			CustomerInformation: &CustomerInformation{
				Customer: Customer{
					CustomerId:    data.Customerid.String,
					CustomerGroup: successfulVendor,
				},
			},
			PaymentInformation: BuildPaymentInformation(nil),
		}, successfulVendor)
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

	btid := job.BusinessTransaction.ID
	oaTxn, err := database.New(database.D()).GetLatestOATransaction(context.Background(), btid)

	_, _ = database.New(database.D()).CreateOATransaction(context.Background(), database.CreateOATransactionParams{
		Businesstransactionid: btid,
		Device:                sql.NullString{String: metadata.device, Valid: true},
		Facility:              sql.NullString{String: metadata.facility, Valid: true},
		Jobid:                 sql.NullString{String: metadata.jobId, Valid: true},
		Lpn:                   sql.NullString{String: oaTxn.Lpn.String, Valid: true},
		Customerid:            sql.NullString{String: oaTxn.Customerid.String, Valid: true},
		Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
		EntryLane:             sql.NullString{String: oaTxn.EntryLane.String, Valid: true},
		IntegratorID:          oaTxn.IntegratorID,
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
	btid := job.BusinessTransaction.ID
	oaTxn, err := database.New(database.D()).GetLatestOATransaction(context.Background(), btid)
	config, err := database.New(database.D()).GetIntegratorConfig(context.Background(), oaTxn.IntegratorID.UUID)

	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	_, _ = database.New(database.D()).CreateOATransaction(context.Background(), database.CreateOATransactionParams{
		Businesstransactionid: btid,
		Device:                sql.NullString{String: metadata.device, Valid: true},
		Facility:              sql.NullString{String: metadata.facility, Valid: true},
		Jobid:                 sql.NullString{String: metadata.jobId, Valid: true},
		Lpn:                   sql.NullString{String: oaTxn.Lpn.String, Valid: true},
		Customerid:            sql.NullString{String: oaTxn.Customerid.String, Valid: true},
		Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
		EntryLane:             sql.NullString{String: oaTxn.EntryLane.String, Valid: true},
		ExitLane:              sql.NullString{String: lane, Valid: true},
		IntegratorID:          oaTxn.IntegratorID,
	})

	go func() {
		cfg, err := database.New(database.D()).GetIntegratorConfigByName(context.Background(), sql.NullString{String: config.Name.String, Valid: true})
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
		}, cfg.Name.String)
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
	btid := job.BusinessTransaction.ID

	jsonStr, err := json.Marshal(map[string]any{
		"steps": "payment_exit_start",
	})

	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	lpn := job.MediaDataList.Identifier.Name

	oaTxn, err := database.New(database.D()).GetLatestOATransaction(context.Background(), btid)

	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	var extra map[string]string
	_ = json.Unmarshal(oaTxn.Extra.RawMessage, &extra)

	_, _ = database.New(database.D()).CreateOATransaction(context.Background(), database.CreateOATransactionParams{
		Businesstransactionid: btid,
		Device:                sql.NullString{String: metadata.device, Valid: true},
		Facility:              sql.NullString{String: metadata.facility, Valid: true},
		Jobid:                 sql.NullString{String: metadata.jobId, Valid: true},
		Lpn:                   sql.NullString{String: oaTxn.Lpn.String, Valid: true},
		Customerid:            sql.NullString{String: oaTxn.Customerid.String, Valid: true},
		Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
		EntryLane:             sql.NullString{String: oaTxn.EntryLane.String, Valid: true},
		ExitLane:              sql.NullString{String: oaTxn.ExitLane.String, Valid: true},
		IntegratorID:          oaTxn.IntegratorID,
	})

	cfg, err := database.New(database.D()).GetIntegratorConfig(context.Background(), oaTxn.IntegratorID.UUID)

	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	customerInformation := &CustomerInformation{Customer: Customer{
		CustomerId:    oaTxn.Customerid.String,
		CustomerGroup: cfg.Name.String,
	}}

	sendZeroAmount := func() {
		sendFinalMessageCustomer(metadata, FMCReq{
			Identifier:          Identifier{Name: oaTxn.Lpn.String},
			BusinessTransaction: &BusinessTransaction{ID: oaTxn.Businesstransactionid},
			CustomerInformation: customerInformation,
			PaymentInformation: BuildPaymentInformation(&PaymentData{
				OriginalAmount: OriginalAmount{
					VatRate: job.PaymentData.OriginalAmount.VatRate,
					Amount:  "0",
				},
			}),
		}, cfg.Name.String)
	}

	amount, err := strconv.ParseFloat(job.PaymentData.OriginalAmount.Amount, 64)
	if err != nil {
		logger.LogData("error", fmt.Sprintf("failed to ParseFloat %v", err), nil)
		go sendZeroAmount()
		return
	}
	amountConv := amount / 100

	if extra["steps"] != "identification_exit_done" && extra["steps"] != "leave_loop_entry_done" {
		logger.LogData("error", "previous step is not either identification_exit_done or leave_loop_entry_done", nil)
		go sendZeroAmount()
		return
	}

	arg := integrator.TransactionArg{
		LPN:                   lpn,
		EntryLane:             oaTxn.EntryLane.String,
		ExitLane:              oaTxn.ExitLane.String,
		Amount:                amountConv,
		EntryAt:               oaTxn.CreatedAt,
		BusinessTransactionId: btid,
		Client:                cfg.Name.String,
		Facility:              metadata.facility,
	}
	err = integrator.PerformTransaction(arg)
	if err != nil {
		jsonStr, err = json.Marshal(map[string]any{
			"steps": "payment_exit_error",
			"error": err.Error(),
		})
		go sendZeroAmount()
		_, _ = database.New(database.D()).CreateOATransaction(context.Background(), database.CreateOATransactionParams{
			Businesstransactionid: btid,
			Device:                sql.NullString{String: metadata.device, Valid: true},
			Facility:              sql.NullString{String: metadata.facility, Valid: true},
			Jobid:                 sql.NullString{String: metadata.jobId, Valid: true},
			Lpn:                   sql.NullString{String: oaTxn.Lpn.String, Valid: true},
			Customerid:            sql.NullString{String: oaTxn.Customerid.String, Valid: true},
			Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
			EntryLane:             sql.NullString{String: oaTxn.EntryLane.String, Valid: true},
			ExitLane:              sql.NullString{String: oaTxn.ExitLane.String, Valid: true},
			IntegratorID:          oaTxn.IntegratorID,
		})
		return
	}
	jsonStr, err = json.Marshal(map[string]any{
		"steps":     "payment_exit_done",
		"paymentAt": time.Now().UTC(),
	})

	if err != nil {
		go sendZeroAmount()
		return
	}

	_, _ = database.New(database.D()).CreateOATransaction(context.Background(), database.CreateOATransactionParams{
		Businesstransactionid: btid,
		Device:                sql.NullString{String: metadata.device, Valid: true},
		Facility:              sql.NullString{String: metadata.facility, Valid: true},
		Jobid:                 sql.NullString{String: metadata.jobId, Valid: true},
		Lpn:                   sql.NullString{String: oaTxn.Lpn.String, Valid: true},
		Customerid:            sql.NullString{String: oaTxn.Customerid.String, Valid: true},
		Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
		EntryLane:             sql.NullString{String: oaTxn.EntryLane.String, Valid: true},
		ExitLane:              sql.NullString{String: oaTxn.ExitLane.String, Valid: true},
		IntegratorID:          oaTxn.IntegratorID,
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
	}, cfg.Name.String)
}

func handleLeaveLoopExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "LEAVE_LOOP" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	lpn := job.MediaDataList.Identifier.Name

	if job.BusinessTransaction == nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	btid := job.BusinessTransaction.ID

	oaTxn, err := database.New(database.D()).GetLatestOATransaction(context.Background(), btid)

	var extra map[string]any

	err = json.Unmarshal(oaTxn.Extra.RawMessage, &extra)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	cfg, err := database.New(database.D()).GetIntegratorConfig(context.Background(), oaTxn.IntegratorID.UUID)

	if extra["steps"] != "payment_exit_done" {
		arg := integrator.TransactionArg{
			LPN:                   lpn,
			EntryLane:             oaTxn.EntryLane.String,
			ExitLane:              oaTxn.ExitLane.String,
			Amount:                0.00,
			EntryAt:               oaTxn.CreatedAt,
			BusinessTransactionId: btid,
			Client:                cfg.Name.String,
			Facility:              metadata.facility,
		}
		err = integrator.PerformTransaction(arg)
	}

	newExtra := map[string]any{
		"steps":   "exit_leave_loop_done",
		"leaveAt": time.Now().UTC(),
	}

	maps.Copy(extra, newExtra)

	jsonStr, err := json.Marshal(extra)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	_, _ = database.New(database.D()).CreateOATransaction(context.Background(), database.CreateOATransactionParams{
		Businesstransactionid: btid,
		Device:                sql.NullString{String: metadata.device, Valid: true},
		Facility:              sql.NullString{String: metadata.facility, Valid: true},
		Jobid:                 sql.NullString{String: metadata.jobId, Valid: true},
		Lpn:                   sql.NullString{String: oaTxn.Lpn.String, Valid: true},
		Customerid:            sql.NullString{String: oaTxn.Customerid.String, Valid: true},
		Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
		EntryLane:             sql.NullString{String: oaTxn.EntryLane.String, Valid: true},
		ExitLane:              sql.NullString{String: oaTxn.ExitLane.String, Valid: true},
		IntegratorID:          oaTxn.IntegratorID,
	})

	go sendEmptyFinalMessage(metadata)
}

type FMCReq struct {
	PaymentInformation  *PaymentInformation
	BusinessTransaction *BusinessTransaction
	CustomerInformation *CustomerInformation
	Identifier          Identifier
}

func sendFinalMessageCustomer(metadata *RequestMetadata, in FMCReq, vendorName string) {
	config, err := database.New(database.D()).GetSnbConfigByFacilityAndDevice(context.Background(), database.GetSnbConfigByFacilityAndDeviceParams{
		Device:   metadata.device,
		Facility: metadata.facility,
	})

	if err != nil {
		fmt.Println("Error get config", err)
		return
	}

	vendor, err := database.New(database.D()).GetIntegratorConfigByName(context.Background(), sql.NullString{String: vendorName, Valid: true})

	var counting *string = nil
	if in.PaymentInformation != nil {
		_counting := "RESERVED"
		counting = &_counting
	}
	xmlData, err := xml.Marshal(&FinalMessageCustomer{
		PaymentInformation: in.PaymentInformation,
		ProviderInformation: &ProviderInformation{
			Provider: Provider{
				ProviderId:   fmt.Sprintf("%v", vendor.ProviderID.Int32),
				ProviderName: vendor.Name.String,
			},
		},
		Counting: counting,
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

	client := &http.Client{
		Timeout: time.Second * 10,
	}
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
