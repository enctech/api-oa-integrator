package oa

import (
	"api-oa-integrator/internal/database"
	"api-oa-integrator/internal/modules/integrator"
	"api-oa-integrator/utils"
	"bytes"
	"context"
	"database/sql"
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
	lpn := job.MediaDataList.Identifier.Name
	lane := job.TimeAndPlace.Device.DeviceNumber
	btid := uuid.New().String()
	err := integrator.VerifyVehicle(btid, lpn, lane)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}

	go func() {
		sendFinalMessageCustomer(metadata, FMCReq{
			Identifier: Identifier{Name: lpn},
			BusinessTransaction: &BusinessTransaction{
				ID: btid,
			},
			CustomerInformation: &CustomerInformation{
				Customer: Customer{
					CustomerId:    encryptLpn(lpn),
					CustomerGroup: viper.GetString("vendor.name"),
				},
			},
		})
	}()
}

func handleLeaveLoopEntry(job *Job, metadata *RequestMetadata) {
	if job.JobType != "LEAVE_LOOP" || job.TimeAndPlace.Device.DeviceType != "ENTRY" {
		return
	}

	lpn := job.MediaDataList.Identifier.Name
	_ = integrator.CreateSession(lpn)
	go sendEmptyFinalMessage(metadata)
}

func handleIdentificationExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "IDENTIFICATION" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	lpn := job.MediaDataList.Identifier.Name
	lane := job.TimeAndPlace.Device.DeviceNumber
	err := integrator.VerifyVehicle("", lpn, lane)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}
}

func handlePaymentExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "PAYMENT" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	lpn := job.MediaDataList.Identifier.Name
	err := integrator.EndSession(lpn)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
		return
	}
}

func handleLeaveLoopExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "LEAVE_LOOP" && job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	lpn := job.MediaDataList.Identifier.Name
	err := integrator.EndSession(lpn)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
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

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	return
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
