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

func VerifyVehicle(client, locationId, plateNumber, lane string) error {
	if plateNumber == "" {
		return errors.New("empty plate number")
	}
	cfg, err := database.New(database.D()).GetIntegratorConfigByName(context.Background(), sql.NullString{String: client, Valid: true})
	if err != nil {
		return err
	}
	var plazaIdMap map[string]any
	err = json.Unmarshal(cfg.PlazaIDMap.RawMessage, &plazaIdMap)
	if err != nil {
		return err
	}
	err = tng.Config{IntegratorConfig: cfg, PlazaId: fmt.Sprintf("%v", plazaIdMap[locationId])}.VerifyVehicle(plateNumber, lane)
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
	cfg, err := database.New(database.D()).GetIntegratorConfigByName(context.Background(), sql.NullString{String: client, Valid: true})
	if err != nil {
		return err
	}

	var plazaIdMap map[string]any
	err = json.Unmarshal(cfg.PlazaIDMap.RawMessage, &plazaIdMap)
	if err != nil {
		return err
	}

	err = tng.Config{IntegratorConfig: cfg, PlazaId: fmt.Sprintf("%v", plazaIdMap[locationId])}.PerformTransaction(tng.TransactionArg{
		Amount:    amount,
		ExitLane:  exitLane,
		LPN:       plateNumber,
		EntryLane: entryLane,
		EntryTime: entryAt,
	})
	if err != nil {
		return err
	}
	return nil
}
