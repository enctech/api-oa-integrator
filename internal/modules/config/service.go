package config

import (
	"api-oa-integrator/internal/database"
	"context"
	"database/sql"
	"go.uber.org/zap"
)

func createSnbConfig(ctx context.Context, in SnbConfig) (SnbConfig, error) {
	txn, _ := database.D().Begin()
	config, err := database.New(database.D()).WithTx(txn).CreateConfig(ctx, database.CreateConfigParams{
		Endpoint: sql.NullString{String: in.Endpoint, Valid: true},
		Facility: in.Facilities,
		Device:   in.Devices,
	})
	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return SnbConfig{}, err
	}
	err = txn.Commit()
	return SnbConfig{
		Endpoint:   config.Endpoint.String,
		Facilities: config.Facility,
		Devices:    config.Device,
	}, nil
}

func getSnbConfig(ctx context.Context, in SnbConfig) (SnbConfig, error) {
	config, err := database.New(database.D()).GetConfig(ctx, database.GetConfigParams{
		Facility: in.Facilities,
		Device:   in.Devices,
	})
	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return SnbConfig{}, err
	}
	return SnbConfig{
		Endpoint:   config.Endpoint.String,
		Facilities: config.Facility,
		Devices:    config.Device,
	}, nil
}
