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
	"net/http"
)

func handleIdentificationEntry(job *Job, metadata *RequestMetadata) {
	if job.JobType != "IDENTIFICATION" || job.TimeAndPlace.Device.DeviceType != "ENTRY" {
		return
	}

	err := tng.VerifyVehicle(job.MediaDataList.Identifier.Name)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
	}
}

func handleLeaveLoopEntry(job *Job, metadata *RequestMetadata) {
	if job.JobType != "LEAVE_LOOP" || job.TimeAndPlace.Device.DeviceType != "ENTRY" {
		return
	}

	err := tng.CreateSession(job.MediaDataList.Identifier.Name)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
	}
}

func handleIdentificationExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "IDENTIFICATION" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	err := tng.VerifyVehicle(job.MediaDataList.Identifier.Name)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
	}
}

func handlePaymentExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "PAYMENT" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	err := tng.EndSession(job.MediaDataList.Identifier.Name)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
	}
}

func handleLeaveLoopExit(job *Job, metadata *RequestMetadata) {
	if job.JobType != "LEAVE_LOOP" && job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	err := tng.AcknowledgeUserExit(job.MediaDataList.Identifier.Name)
	if err != nil {
		go sendEmptyFinalMessage(metadata)
	}
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
