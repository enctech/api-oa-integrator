package oa

import (
	"api-oa-integrator/database"
	"api-oa-integrator/internal/modules/integrator"
	"api-oa-integrator/logger"
	"api-oa-integrator/tracing"
	"api-oa-integrator/utils"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"maps"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/sqlc-dev/pqtype"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func handleIdentificationEntry(c echo.Context, job *Job, metadata *RequestMetadata) {
	if job.JobType != "IDENTIFICATION" || job.TimeAndPlace.Device.DeviceType != "ENTRY" {
		return
	}

	ctx := c.Request().Context()
	ctx, span := tracing.Tracer().Start(ctx, "oa.handleIdentificationEntry",
		trace.WithAttributes(tracing.JobAttributes(metadata.jobId, metadata.facility, metadata.device, job.MediaDataList.Identifier.Name)...),
	)
	defer span.End()

	lpn := job.MediaDataList.Identifier.Name
	lane := job.TimeAndPlace.Device.DeviceNumber
	btid := uuid.New().String()
	customerId := encryptLpn(lpn)

	span.SetAttributes(tracing.TransactionIDKey.String(btid))
	span.SetAttributes(tracing.CustomerIDKey.String(customerId))

	jsonStr, err := json.Marshal(map[string]any{
		"steps": "identification_entry_start",
	})
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error Marshal %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}
	data, err := database.New(database.TracedD()).CreateOATransaction(ctx, database.CreateOATransactionParams{
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
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error create oa transaction %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}

	configs, err := database.New(database.TracedD()).GetIntegratorConfigs(ctx)
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error get integrator configs %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}

	var wg sync.WaitGroup
	vendorChannel := make(chan string, 1)

	for i := range configs {
		wg.Add(1)
		go func(vendorName string) {
			defer wg.Done()
			err = integrator.VerifyVehicle(ctx, vendorName, metadata.facility, lpn, lane)

			if err != nil {
				logger.LogWithContext(ctx, "error", fmt.Sprintf("Error integrator.VerifyVehicle %v", err), map[string]interface{}{
					"lpn":        lpn,
					"vendorName": vendorName,
					"facility":   metadata.facility,
				})
			} else {
				select {
				case vendorChannel <- vendorName:
				default:
				}
			}
		}(configs[i].Name.String)
	}

	wg.Wait()
	close(vendorChannel)

	var successfulVendor string
	select {
	case successfulVendor = <-vendorChannel:
		logger.LogWithContext(ctx, "info", fmt.Sprintf("Success from: %s", successfulVendor), nil)
		span.SetAttributes(tracing.VendorKey.String(successfulVendor))
	default:
		sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
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
		bgCtx := tracing.DetachedContext(ctx)
		if successfulVendor == "" {
			go sendEmptyFinalMessage(bgCtx, metadata)
			return
		}

		jsonStr, err := json.Marshal(map[string]any{
			"steps": "identification_entry_done",
		})
		if err != nil {
			logger.LogWithContext(bgCtx, "error", fmt.Sprintf("Error Marshal %v", err), nil)
			go sendEmptyFinalMessage(bgCtx, metadata)
			return
		}

		config, err := utils.FirstWhere(configs, func(config database.IntegratorConfig) bool {
			return config.Name.String == successfulVendor
		})

		_, err = database.New(database.TracedD()).UpdateOATransaction(bgCtx, database.UpdateOATransactionParams{
			Businesstransactionid: data.Businesstransactionid,
			Extra:                 pqtype.NullRawMessage{Valid: true, RawMessage: jsonStr},
			IntegratorID: uuid.NullUUID{
				UUID: config.ID, Valid: true,
			},
		})
		sendFinalMessageCustomer(bgCtx, metadata, FMCReq{
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

func handleLeaveLoopEntry(ctx context.Context, job *Job, metadata *RequestMetadata) {
	if job.JobType != "LEAVE_LOOP" || job.TimeAndPlace.Device.DeviceType != "ENTRY" {
		return
	}

	ctx, span := tracing.Tracer().Start(ctx, "oa.handleLeaveLoopEntry",
		trace.WithAttributes(tracing.JobAttributes(metadata.jobId, metadata.facility, metadata.device, "")...),
	)
	defer span.End()

	if job.BusinessTransaction == nil {
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}

	if job.BusinessTransaction.ID == "" {
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}

	jsonStr, err := json.Marshal(map[string]any{
		"steps": "leave_loop_entry_done",
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}

	btid := job.BusinessTransaction.ID
	span.SetAttributes(tracing.TransactionIDKey.String(btid))

	oaTxn, err := database.New(database.TracedD()).GetLatestOATransaction(ctx, btid)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}

	_, _ = database.New(database.TracedD()).CreateOATransaction(ctx, database.CreateOATransactionParams{
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
	go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
}

func handleIdentificationExit(ctx context.Context, job *Job, metadata *RequestMetadata) {
	if job.JobType != "IDENTIFICATION" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	ctx, span := tracing.Tracer().Start(ctx, "oa.handleIdentificationExit",
		trace.WithAttributes(tracing.JobAttributes(metadata.jobId, metadata.facility, metadata.device, "")...),
	)
	defer span.End()

	if job.BusinessTransaction == nil {
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}

	if job.BusinessTransaction.ID == "" {
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}

	jsonStr, err := json.Marshal(map[string]any{
		"steps": "identification_exit_done",
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}

	lane := job.TimeAndPlace.Device.DeviceNumber
	btid := job.BusinessTransaction.ID
	span.SetAttributes(tracing.TransactionIDKey.String(btid))
	span.SetAttributes(tracing.ExitLaneKey.String(lane))

	oaTxn, err := database.New(database.TracedD()).GetLatestOATransaction(ctx, btid)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}

	config, err := database.New(database.TracedD()).GetIntegratorConfig(ctx, oaTxn.IntegratorID.UUID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(tracing.DetachedContext(ctx), metadata)
		return
	}

	_, _ = database.New(database.TracedD()).CreateOATransaction(ctx, database.CreateOATransactionParams{
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
		bgCtx := tracing.DetachedContext(ctx)
		cfg, err := database.New(database.TracedD()).GetIntegratorConfigByName(bgCtx, sql.NullString{String: config.Name.String, Valid: true})
		if err != nil {
			go sendEmptyFinalMessage(bgCtx, metadata)
			return
		}
		sendFinalMessageCustomer(bgCtx, metadata, FMCReq{
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

func handlePaymentExit(ctx context.Context, job *Job, metadata *RequestMetadata) {
	if job.JobType != "PAYMENT" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	ctx, span := tracing.Tracer().Start(ctx, "oa.handlePaymentExit",
		trace.WithAttributes(tracing.JobAttributes(metadata.jobId, metadata.facility, metadata.device, job.MediaDataList.Identifier.Name)...),
	)
	defer span.End()

	bgCtx := tracing.DetachedContext(ctx)

	if job.BusinessTransaction == nil {
		go sendEmptyFinalMessage(bgCtx, metadata)
		return
	}
	btid := job.BusinessTransaction.ID
	span.SetAttributes(tracing.TransactionIDKey.String(btid))

	jsonStr, err := json.Marshal(map[string]any{
		"steps": "payment_exit_start",
	})

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(bgCtx, metadata)
		return
	}

	lpn := job.MediaDataList.Identifier.Name

	oaTxn, err := database.New(database.TracedD()).GetLatestOATransaction(ctx, btid)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(bgCtx, metadata)
		return
	}

	var extra map[string]string
	_ = json.Unmarshal(oaTxn.Extra.RawMessage, &extra)

	_, _ = database.New(database.TracedD()).CreateOATransaction(ctx, database.CreateOATransactionParams{
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

	cfg, err := database.New(database.TracedD()).GetIntegratorConfig(ctx, oaTxn.IntegratorID.UUID)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(bgCtx, metadata)
		return
	}

	customerInformation := &CustomerInformation{Customer: Customer{
		CustomerId:    oaTxn.Customerid.String,
		CustomerGroup: cfg.Name.String,
	}}

	sendZeroAmount := func() {
		sendFinalMessageCustomer(bgCtx, metadata, FMCReq{
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
		logger.LogWithContext(ctx, "error", fmt.Sprintf("failed to ParseFloat %v", err), nil)
		span.RecordError(err)
		go sendZeroAmount()
		return
	}
	amountConv := amount / 100
	span.SetAttributes(tracing.AmountKey.Float64(amountConv))

	if extra["steps"] != "identification_exit_done" && extra["steps"] != "leave_loop_entry_done" {
		logger.LogWithContext(ctx, "error", "previous step is not either identification_exit_done or leave_loop_entry_done", nil)
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
	err = integrator.PerformTransaction(ctx, arg)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		jsonStr, err = json.Marshal(map[string]any{
			"steps": "payment_exit_error",
			"error": err.Error(),
		})
		go sendZeroAmount()
		_, _ = database.New(database.TracedD()).CreateOATransaction(ctx, database.CreateOATransactionParams{
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
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendZeroAmount()
		return
	}

	_, _ = database.New(database.TracedD()).CreateOATransaction(ctx, database.CreateOATransactionParams{
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

	go sendFinalMessageCustomer(bgCtx, metadata, FMCReq{
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

func handleLeaveLoopExit(ctx context.Context, job *Job, metadata *RequestMetadata) {
	if job.JobType != "LEAVE_LOOP" || job.TimeAndPlace.Device.DeviceType != "EXIT" {
		return
	}

	ctx, span := tracing.Tracer().Start(ctx, "oa.handleLeaveLoopExit",
		trace.WithAttributes(tracing.JobAttributes(metadata.jobId, metadata.facility, metadata.device, job.MediaDataList.Identifier.Name)...),
	)
	defer span.End()

	bgCtx := tracing.DetachedContext(ctx)
	lpn := job.MediaDataList.Identifier.Name

	if job.BusinessTransaction == nil {
		go sendEmptyFinalMessage(bgCtx, metadata)
		return
	}

	btid := job.BusinessTransaction.ID
	span.SetAttributes(tracing.TransactionIDKey.String(btid))

	oaTxn, err := database.New(database.TracedD()).GetLatestOATransaction(ctx, btid)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(bgCtx, metadata)
		return
	}

	var extra map[string]any

	err = json.Unmarshal(oaTxn.Extra.RawMessage, &extra)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(bgCtx, metadata)
		return
	}

	cfg, err := database.New(database.TracedD()).GetIntegratorConfig(ctx, oaTxn.IntegratorID.UUID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(bgCtx, metadata)
		return
	}

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
		err = integrator.PerformTransaction(ctx, arg)
		if err != nil {
			span.RecordError(err)
		}
	}

	newExtra := map[string]any{
		"steps":   "exit_leave_loop_done",
		"leaveAt": time.Now().UTC(),
	}

	maps.Copy(extra, newExtra)

	jsonStr, err := json.Marshal(extra)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		go sendEmptyFinalMessage(bgCtx, metadata)
		return
	}

	_, _ = database.New(database.TracedD()).CreateOATransaction(ctx, database.CreateOATransactionParams{
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

	go sendEmptyFinalMessage(bgCtx, metadata)
}

type FMCReq struct {
	PaymentInformation  *PaymentInformation
	BusinessTransaction *BusinessTransaction
	CustomerInformation *CustomerInformation
	Identifier          Identifier
}

func sendFinalMessageCustomer(ctx context.Context, metadata *RequestMetadata, in FMCReq, vendorName string) {
	ctx, span := tracing.Tracer().Start(ctx, "snb.sendFinalMessageCustomer",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			tracing.VendorKey.String(vendorName),
			tracing.FacilityKey.String(metadata.facility),
			tracing.DeviceKey.String(metadata.device),
			tracing.JobIDKey.String(metadata.jobId),
		),
	)
	defer span.End()

	config, err := database.New(database.TracedD()).GetSnbConfigByFacilityAndDevice(ctx, database.GetSnbConfigByFacilityAndDeviceParams{
		Device:   metadata.device,
		Facility: metadata.facility,
	})

	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error get config: %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	vendor, err := database.New(database.TracedD()).GetIntegratorConfigByName(ctx, sql.NullString{String: vendorName, Valid: true})

	var counting *string = nil
	if in.PaymentInformation != nil {
		_counting := "NON-RESERVED"
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
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error marshaling XML data: %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%v/AuthorizationServiceSB/%v/%v/%v/finalmessage", config.Endpoint.String, metadata.facility, metadata.device, metadata.jobId), bytes.NewBuffer(xmlData))
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error creating request: %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/xml")
	req.SetBasicAuth(config.Username.String, config.Password.String)

	resp, err := utils.GlobalInsecureHttpClient.Do(req)
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error sending request: %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	defer resp.Body.Close()
}

func sendEmptyFinalMessage(ctx context.Context, metadata *RequestMetadata) {
	ctx, span := tracing.Tracer().Start(ctx, "snb.sendEmptyFinalMessage",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			tracing.FacilityKey.String(metadata.facility),
			tracing.DeviceKey.String(metadata.device),
			tracing.JobIDKey.String(metadata.jobId),
		),
	)
	defer span.End()

	config, err := database.New(database.TracedD()).GetSnbConfigByFacilityAndDevice(ctx, database.GetSnbConfigByFacilityAndDeviceParams{
		Device:   metadata.device,
		Facility: metadata.facility,
	})

	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error get config: %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	xmlData, err := xml.Marshal(&FinalMessageCustomer{})
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error marshaling XML data: %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	reqCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, "PUT", fmt.Sprintf("%v/AuthorizationServiceSB/%v/%v/%v/finalmessage", config.Endpoint.String, metadata.facility, metadata.device, metadata.jobId), bytes.NewBuffer(xmlData))
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error creating request: %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/xml")
	req.SetBasicAuth(config.Username.String, config.Password.String)

	resp, err := utils.GlobalInsecureHttpClient.Do(req)
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error sending request: %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	defer resp.Body.Close()
}

func CheckSystemAvailability(facility, device string) error {
	ctx, span := tracing.Tracer().Start(context.Background(), "snb.CheckSystemAvailability",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			tracing.FacilityKey.String(facility),
			tracing.DeviceKey.String(device),
		),
	)
	defer span.End()

	config, err := database.New(database.TracedD()).GetSnbConfigByFacilityAndDevice(ctx, database.GetSnbConfigByFacilityAndDeviceParams{
		Device:   device,
		Facility: facility,
	})

	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error get config: %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	xmlOut := []byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><version>
	<customerVersion>%v</customerVersion>
	<sbAuthorizationServiceVersion>2.5.6</sbAuthorizationServiceVersion>
	<configuration>
	</configuration>
	</version>`, viper.GetString("app.version")))

	reqCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "PUT", fmt.Sprintf("%v/AuthorizationServiceSB/version", config.Endpoint.String), bytes.NewBuffer(xmlOut))
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error creating request: %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	req.Header.Set("Content-Type", "application/xml")
	req.SetBasicAuth(config.Username.String, config.Password.String)

	resp, err := utils.GlobalInsecureHttpClient.Do(req)
	if err != nil {
		logger.LogWithContext(ctx, "error", fmt.Sprintf("Error sending request: %v", err), nil)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	defer resp.Body.Close()
	return nil
}
