package integrator

import (
	"api-oa-integrator/internal/database"
	"api-oa-integrator/tng"
	"context"
	"database/sql"
	"errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

func VerifyVehicle(txnId, plateNumber, lane string) error {
	if plateNumber == "" {
		return errors.New("empty plate number")
	}
	cfg, err := database.New(database.D()).GetIntegratorConfig(context.Background(), sql.NullString{String: viper.GetString("vendor.id"), Valid: true})
	if err != nil {
		return err
	}
	err = tng.Config{IntegratorConfig: cfg}.VerifyVehicle(plateNumber, lane)
	if err != nil {
		return err
	}
	return nil
}

func PerformTransaction(plateNumber, entryLane, exitLane string, entryAt time.Time, amount float64) error {
	zap.L().Sugar().With("plateNumber", plateNumber).Info("PerformTransaction")
	if plateNumber == "" {
		return errors.New("empty plate number")
	}
	cfg, err := database.New(database.D()).GetIntegratorConfig(context.Background(), sql.NullString{String: viper.GetString("vendor.id"), Valid: true})
	if err != nil {
		return err
	}
	err = tng.Config{IntegratorConfig: cfg}.PerformTransaction(tng.TransactionArg{
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
