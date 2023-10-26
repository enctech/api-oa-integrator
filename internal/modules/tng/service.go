package tng

import (
	"errors"
	"go.uber.org/zap"
)

func VerifyVehicle(plateNumber string) error {
	zap.L().Sugar().With("plateNumber", plateNumber).Info("VerifyVehicle")
	if plateNumber == "" {
		return errors.New("empty plate number")
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
