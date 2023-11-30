package integrator

import (
	"api-oa-integrator/internal/database"
	"api-oa-integrator/tng"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"go.uber.org/zap"
	"time"
)

var Integrators = []string{"tng"}

type Process interface {
	VerifyVehicle(plateNumber, entryLane string) error
	PerformTransaction(locationId, plateNumber, entryLane, exitLane string, entryAt time.Time, amount float64) (map[string]any, map[string]any, error)
}

func getConfigFromIntegratorBasedOnIntegrator(client, locationId string) (Process, database.IntegratorConfig, error) {
	cfg, err := database.New(database.D()).GetIntegratorConfigByName(context.Background(), sql.NullString{String: client, Valid: true})
	if err != nil {
		return nil, database.IntegratorConfig{}, err
	}
	var plazaIdMap map[string]any
	err = json.Unmarshal(cfg.PlazaIDMap.RawMessage, &plazaIdMap)
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

func VerifyVehicle(client, locationId, plateNumber, lane string) error {
	if plateNumber == "" {
		return errors.New("empty plate number")
	}
	integratorConfig, _, err := getConfigFromIntegratorBasedOnIntegrator(client, locationId)
	if err != nil {
		return err
	}
	err = integratorConfig.VerifyVehicle(plateNumber, lane)
	if err != nil {
		return err
	}
	return nil
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

func PerformTransaction(arg TransactionArg) error {
	zap.L().Sugar().With("plateNumber", arg.LPN).Info("PerformTransaction")
	if arg.LPN == "" {
		return errors.New("empty plate number")
	}
	integratorProcess, integratorConfig, err := getConfigFromIntegratorBasedOnIntegrator(arg.Client, arg.Facility)
	if err != nil {
		return err
	}
	data, taxData, err := integratorProcess.PerformTransaction(arg.Facility, arg.LPN, arg.EntryLane, arg.ExitLane, arg.EntryAt, arg.Amount)
	jsonStr, err := json.Marshal(data)
	taxJsonStr, err := json.Marshal(taxData)

	status := "success"
	errorMessage := ""
	if err != nil {
		status = "fail"
		errorMessage = err.Error()
	}

	_, err = database.New(database.D()).CreateIntegratorTransaction(context.Background(), database.CreateIntegratorTransactionParams{
		Lpn:                   sql.NullString{String: arg.LPN, Valid: true},
		BusinessTransactionID: uuid.MustParse(arg.BusinessTransactionId),
		Extra:                 pqtype.NullRawMessage{Valid: err == nil, RawMessage: jsonStr},
		TaxData:               pqtype.NullRawMessage{Valid: err == nil, RawMessage: taxJsonStr},
		ID:                    integratorConfig.ID,
		Error:                 sql.NullString{String: errorMessage, Valid: err != nil},
		Amount:                sql.NullString{String: fmt.Sprintf("%v", arg.Amount), Valid: true},
		Status:                sql.NullString{String: status, Valid: true},
	})

	return err
}
