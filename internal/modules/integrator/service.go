package integrator

import (
	"api-oa-integrator/internal/database"
	"api-oa-integrator/tng"
	"context"
	"database/sql"
	"errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func VerifyVehicle(txnId, plateNumber, lane string) error {
	zap.L().Sugar().With("plateNumber", plateNumber).Info("VerifyVehicle")
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

func CreateSession(plateNumber string) error {
	zap.L().Sugar().With("plateNumber", plateNumber).Info("CreateSession")
	if plateNumber == "" {
		return errors.New("empty plate number")
	}
	return nil
}

func EndSession(plateNumber string) error {
	zap.L().Sugar().With("plateNumber", plateNumber).Info("EndSession")
	if plateNumber == "" {
		return errors.New("empty plate number")
	}
	return nil
}

func AcknowledgeUserExit(plateNumber string) error {
	zap.L().Sugar().With("plateNumber", plateNumber).Info("AcknowledgeUserExit")
	if plateNumber == "" {
		return errors.New("empty plate number")
	}
	return nil
}
