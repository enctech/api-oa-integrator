package oa

import (
	"api-oa-integrator/internal/database"
	"api-oa-integrator/internal/modules/tng"
	"api-oa-integrator/internal/utils"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"net/http"
)

func handleIdentificationEntry(job *Job, metadata *RequestMetadata) {
	if job.JobType != "IDENTIFICATION" || job.TimeAndPlace.Device.DeviceType != "ENTRY" {
		return
	}

	err := utils.LogToDb("TnG", "Verify Vehicle Entry", []byte("{}"))
	if err != nil {
		return
	}
	lpn := job.MediaDataList.Identifier.Name
	err = tng.VerifyVehicle(lpn)
	if err != nil {
		data, err := json.Marshal(map[string]any{
			"status":      "fail",
			"error":       err.Error(),
			"plateNumber": lpn,
		})
		err = utils.LogToDb("TnG", "Verify Vehicle Entry Error", data)
		if err != nil {
			return
		}
		go sendEmptyFinalMessage(metadata)
		return
	}

	go sendFinalMessageCustomer(metadata, FMCReq{
		Identifier: Identifier{Name: lpn},
		BusinessTransaction: &BusinessTransaction{
			ID: uuid.New().String(),
		},
		CustomerInformation: &CustomerInformation{
			Customer: Customer{
				CustomerId:    encryptLpn(lpn),
				CustomerGroup: viper.GetString("vendor.name"),
			},
		},
	})

	data, err := json.Marshal(map[string]any{
		"status":      "success",
		"error":       err.Error(),
		"plateNumber": lpn,
	})
	err = utils.LogToDb("TnG", "Verify Vehicle Entry Error", data)
	if err != nil {
		return
	}
}

func handleLeaveLoopEntry(job *Job, metadata *RequestMetadata) {
	if job.JobType != "LEAVE_LOOP" || job.TimeAndPlace.Device.DeviceType != "ENTRY" {
		return
	}

	err := utils.LogToDb("TnG", "Leave Loop Entry", []byte("{}"))
	if err != nil {
		return
	}
	lpn := job.MediaDataList.Identifier.Name
	err = tng.CreateSession(lpn)
	if err != nil {
		data, err := json.Marshal(map[string]any{
			"status":      "fail",
			"error":       err.Error(),
			"plateNumber": lpn,
		})
		err = utils.LogToDb("TnG", "Leave Loop Entry Error", data)
		if err != nil {
			return
		}
		go sendEmptyFinalMessage(metadata)
		return
	}
	data, err := json.Marshal(map[string]any{
		"status":      "success",
		"error":       err.Error(),
		"plateNumber": lpn,
	})
	err = utils.LogToDb("TnG", "Leave Loop Entry Success", data)
	if err != nil {
		return
	}
}

func handleIdentificationExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "IDENTIFICATION" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	err := utils.LogToDb("TnG", "Verify Vehicle Exit", []byte("{}"))
	if err != nil {
		return
	}
	lpn := job.MediaDataList.Identifier.Name
	err = tng.VerifyVehicle(lpn)
	if err != nil {
		data, err := json.Marshal(map[string]any{
			"status":      "fail",
			"error":       err.Error(),
			"plateNumber": lpn,
		})
		err = utils.LogToDb("TnG", "Verify Vehicle Exit Error", data)
		if err != nil {
			return
		}
		go sendEmptyFinalMessage(metadata)
		return
	}
	data, err := json.Marshal(map[string]any{
		"status":      "success",
		"error":       err.Error(),
		"plateNumber": lpn,
	})
	err = utils.LogToDb("TnG", "Verify Vehicle Exit Success", data)
	if err != nil {
		return
	}
}

func handlePaymentExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "PAYMENT" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	err := utils.LogToDb("TnG", "PAYMENT", []byte("{}"))
	if err != nil {
		return
	}
	lpn := job.MediaDataList.Identifier.Name
	err = tng.EndSession(lpn)
	if err != nil {
		data, err := json.Marshal(map[string]any{
			"status":      "fail",
			"error":       err.Error(),
			"plateNumber": lpn,
		})
		err = utils.LogToDb("TnG", "PAYMENT Error", data)
		if err != nil {
			return
		}
		go sendEmptyFinalMessage(metadata)
		return
	}
	data, err := json.Marshal(map[string]any{
		"status":      "success",
		"error":       err.Error(),
		"plateNumber": lpn,
	})
	err = utils.LogToDb("TnG", "PAYMENT Success", data)
	if err != nil {
		return
	}
}

func handleLeaveLoopExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "LEAVE_LOOP" && job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	err := utils.LogToDb("TnG", "Leave Loop Exit", []byte("{}"))
	if err != nil {
		return
	}
	lpn := job.MediaDataList.Identifier.Name
	err = tng.EndSession(lpn)
	if err != nil {
		data, err := json.Marshal(map[string]any{
			"status":      "fail",
			"error":       err.Error(),
			"plateNumber": lpn,
		})
		err = utils.LogToDb("TnG", "Leave Loop Exit Error", data)
		if err != nil {
			return
		}
		go sendEmptyFinalMessage(metadata)
		return
	}
	data, err := json.Marshal(map[string]any{
		"status":      "success",
		"error":       err.Error(),
		"plateNumber": lpn,
	})
	err = utils.LogToDb("TnG", "Leave Loop Exit Success", data)
	if err != nil {
		return
	}
}

type FMCReq struct {
	BusinessTransaction *BusinessTransaction
	CustomerInformation *CustomerInformation
	Identifier          Identifier
}

func sendFinalMessageCustomer(metadata *RequestMetadata, in FMCReq) {
	config, err := database.New(database.D()).GetConfig(context.Background(), database.GetConfigParams{
		Device:   sql.NullString{String: metadata.device, Valid: true},
		Facility: sql.NullString{String: metadata.facility, Valid: true},
	})

	if err != nil {
		fmt.Println("Error get config", err)
		return
	}

	counting := "RESERVED"
	xmlData, err := xml.Marshal(&FinalMessageCustomer{
		PaymentInformation: &PaymentInformation{
			PaymentLocation: "PAY_LOCAL",
		},
		ProviderInformation: &ProviderInformation{
			Provider: Provider{
				ProviderId:   viper.GetString("vendor.id"),
				ProviderName: viper.GetString("vendor.name"),
			},
		},
		Reservation: &Reservation{
			ReservationTariff: ReservationTariff{
				TariffName:   "Tariff OnlineManipulation",
				TariffNumber: 34,
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

	req, err := http.NewRequest("PUT", config.Endpoint.String, bytes.NewBuffer(xmlData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	client.Transport = &utils.LoggingRoundTripper{
		Transport: http.DefaultTransport,
	}

	data, err := json.Marshal(map[string]any{
		"message":  "sending empty final message",
		"device":   metadata.device,
		"jobId":    metadata.jobId,
		"facility": metadata.facility,
	})
	if err != nil {
		return
	}
	err = utils.LogToDb("OA", "sendEmptyFinalMessage", data)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		data, err := json.Marshal(map[string]any{
			"message":  "sending empty final message error",
			"device":   metadata.device,
			"jobId":    metadata.jobId,
			"facility": metadata.facility,
			"error":    err.Error(),
		})
		if err != nil {
			return
		}
		err = utils.LogToDb("OA", "sendEmptyFinalMessage", data)
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}

func sendEmptyFinalMessage(metadata *RequestMetadata) {
	config, err := database.New(database.D()).GetConfig(context.Background(), database.GetConfigParams{
		Device:   sql.NullString{String: metadata.device, Valid: true},
		Facility: sql.NullString{String: metadata.facility, Valid: true},
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

	req, err := http.NewRequest("PUT", config.Endpoint.String, bytes.NewBuffer(xmlData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	client.Transport = &utils.LoggingRoundTripper{
		Transport: http.DefaultTransport,
	}

	data, err := json.Marshal(map[string]any{
		"message":  "sending empty final message",
		"device":   metadata.device,
		"jobId":    metadata.jobId,
		"facility": metadata.facility,
	})
	if err != nil {
		return
	}
	err = utils.LogToDb("OA", "sendEmptyFinalMessage", data)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		data, err := json.Marshal(map[string]any{
			"message":  "sending empty final message error",
			"device":   metadata.device,
			"jobId":    metadata.jobId,
			"facility": metadata.facility,
			"error":    err.Error(),
		})
		if err != nil {
			return
		}
		err = utils.LogToDb("OA", "sendEmptyFinalMessage", data)
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
}
