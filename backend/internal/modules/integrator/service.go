package integrator

import (
	"api-oa-integrator/database"
	"api-oa-integrator/logger"
	"api-oa-integrator/tng"
	"api-oa-integrator/tracing"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var Integrators = []string{"tng"}

type Process interface {
	VerifyVehicle(ctx context.Context, plateNumber, entryLane string) error
	PerformTransaction(ctx context.Context, locationId, plateNumber, entryLane, exitLane string, entryAt time.Time, amount float64) (map[string]any, map[string]any, *string, error)
	VoidTransaction(ctx context.Context, plateNumber, transactionId string) error
	CancelEntry() error
}

func getConfigFromIntegratorBasedOnIntegrator(client, locationId string) (Process, database.IntegratorConfig, error) {
	cfg, err := database.New(database.D()).GetIntegratorConfigByName(context.Background(), sql.NullString{String: client, Valid: true})
	if err != nil {
		return nil, database.IntegratorConfig{}, err
	}
	var plazaIdMap map[string]any
	err = json.Unmarshal(cfg.PlazaIDMap.RawMessage, &plazaIdMap)
	if plazaIdMap[locationId] == nil || plazaIdMap[locationId] == "" {
		return nil, database.IntegratorConfig{}, errors.New(fmt.Sprintf("plazaId not found for locationId %v", locationId))
	}
	if err != nil {
		return nil, database.IntegratorConfig{}, err
	}
	switch cfg.IntegratorName.String {
	case "tng":
		return tng.Config{IntegratorConfig: cfg, PlazaId: fmt.Sprintf("%v", plazaIdMap[locationId])}, cfg, nil
	default:
		return nil, database.IntegratorConfig{}, errors.New(fmt.Sprintf("integrator %v not found", cfg.IntegratorName.String))
	}
}

func VerifyVehicle(ctx context.Context, client, locationId, plateNumber, lane string) error {
	ctx, span := tracing.Tracer().Start(ctx, "integrator.VerifyVehicle",
		trace.WithAttributes(
			tracing.VendorKey.String(client),
			tracing.PlateNumberKey.String(plateNumber),
			tracing.FacilityKey.String(locationId),
			tracing.EntryLaneKey.String(lane),
		),
	)
	defer span.End()

	if plateNumber == "" {
		err := errors.New("empty plate number")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	integratorConfig, _, err := getConfigFromIntegratorBasedOnIntegrator(client, locationId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	err = integratorConfig.VerifyVehicle(ctx, plateNumber, lane)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func CancelEntry(client, locationId string) {
	integratorConfig, _, _ := getConfigFromIntegratorBasedOnIntegrator(client, locationId)
	if integratorConfig != nil {
		_ = integratorConfig.CancelEntry()
	}
}

type TransactionArg struct {
	BusinessTransactionId string
	Client                string
	Facility              string
	LPN                   string
	EntryLane             string
	ExitLane              string
	EntryAt               time.Time
	Amount                float64
}

func PerformTransaction(ctx context.Context, arg TransactionArg) error {
	ctx, span := tracing.Tracer().Start(ctx, "integrator.PerformTransaction",
		trace.WithAttributes(
			tracing.TransactionIDKey.String(arg.BusinessTransactionId),
			tracing.VendorKey.String(arg.Client),
			tracing.FacilityKey.String(arg.Facility),
			tracing.PlateNumberKey.String(arg.LPN),
			tracing.EntryLaneKey.String(arg.EntryLane),
			tracing.ExitLaneKey.String(arg.ExitLane),
			tracing.AmountKey.Float64(arg.Amount),
		),
	)
	defer span.End()

	logger.LogWithContext(ctx, "info", "PerformTransaction", map[string]interface{}{"plateNumber": arg.LPN})
	if arg.LPN == "" {
		err := errors.New("empty plate number")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	integratorProcess, integratorConfig, err := getConfigFromIntegratorBasedOnIntegrator(arg.Client, arg.Facility)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	data, taxData, customStatus, txnErr := integratorProcess.PerformTransaction(ctx, arg.Facility, arg.LPN, arg.EntryLane, arg.ExitLane, arg.EntryAt, arg.Amount)
	status := "success"
	if customStatus != nil {
		status = *customStatus
	}
	errorMessage := ""
	if txnErr != nil {
		status = "fail"
		errorMessage = txnErr.Error()
		span.RecordError(txnErr)
		span.SetStatus(codes.Error, txnErr.Error())
	}
	jsonStr, err := json.Marshal(data)
	taxJsonStr, err := json.Marshal(taxData)

	_, err = database.New(database.TracedD()).CreateIntegratorTransaction(ctx, database.CreateIntegratorTransactionParams{
		Lpn:                   sql.NullString{String: arg.LPN, Valid: true},
		BusinessTransactionID: uuid.MustParse(arg.BusinessTransactionId),
		Extra:                 pqtype.NullRawMessage{Valid: err == nil, RawMessage: jsonStr},
		TaxData:               pqtype.NullRawMessage{Valid: err == nil, RawMessage: taxJsonStr},
		ID:                    integratorConfig.ID,
		Error:                 sql.NullString{String: errorMessage, Valid: txnErr != nil},
		Amount:                sql.NullString{String: fmt.Sprintf("%v", arg.Amount), Valid: true},
		Status:                sql.NullString{String: status, Valid: true},
	})

	return txnErr
}
