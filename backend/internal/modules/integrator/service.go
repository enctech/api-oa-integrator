package integrator

import (
	"api-oa-integrator/internal/database"
	"api-oa-integrator/tng"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"
)

var Integrators = []string{"tng"}

type Process interface {
	VerifyVehicle(plateNumber, entryLane string) error
	PerformTransaction(locationId, plateNumber, entryLane, exitLane string, entryAt time.Time, amount float64) error
}

func getConfigFromIntegratorBasedOnIntegrator(client, locationId string) (Process, error) {
	cfg, err := database.New(database.D()).GetIntegratorConfigByName(context.Background(), sql.NullString{String: client, Valid: true})
	if err != nil {
		return nil, err
	}
	var plazaIdMap map[string]any
	err = json.Unmarshal(cfg.PlazaIDMap.RawMessage, &plazaIdMap)
	if err != nil {
		return nil, err
	}
	switch cfg.IntegratorName.String {
	case "tng":
		return tng.Config{IntegratorConfig: cfg, PlazaId: fmt.Sprintf("%v", plazaIdMap[locationId])}, nil
	default:
		return nil, errors.New("invalid integrator name")
	}
}

func VerifyVehicle(client, locationId, plateNumber, lane string) error {
	if plateNumber == "" {
		return errors.New("empty plate number")
	}
	integratorConfig, err := getConfigFromIntegratorBasedOnIntegrator(client, locationId)
	if err != nil {
		return err
	}
	err = integratorConfig.VerifyVehicle(plateNumber, lane)
	if err != nil {
		return err
	}
	return nil
}

func PerformTransaction(client, locationId, plateNumber, entryLane, exitLane string, entryAt time.Time, amount float64) error {
	zap.L().Sugar().With("plateNumber", plateNumber).Info("PerformTransaction")
	if plateNumber == "" {
		return errors.New("empty plate number")
	}
	integratorConfig, err := getConfigFromIntegratorBasedOnIntegrator(client, locationId)
	if err != nil {
		return err
	}
	err = integratorConfig.PerformTransaction(locationId, plateNumber, entryLane, exitLane, entryAt, amount)
	if err != nil {
		return err
	}
	return nil
}
