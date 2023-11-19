package config

import (
	"api-oa-integrator/internal/database"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func createSnbConfig(ctx context.Context, in SnbConfig) (SnbConfig, error) {
	txn, _ := database.D().Begin()
	config, err := database.New(database.D()).WithTx(txn).CreateSnbConfig(ctx, database.CreateSnbConfigParams{
		Endpoint: sql.NullString{String: in.Endpoint, Valid: true},
		Facility: in.Facilities,
		Device:   in.Devices,
		Name:     sql.NullString{String: in.Name, Valid: true},
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

func updateSnbConfig(ctx context.Context, id uuid.UUID, in SnbConfig) (SnbConfig, error) {
	txn, _ := database.D().Begin()
	config, err := database.New(database.D()).WithTx(txn).UpdateSnbConfig(ctx, database.UpdateSnbConfigParams{
		ID:       id,
		Endpoint: sql.NullString{String: in.Endpoint, Valid: true},
		Facility: in.Facilities,
		Device:   in.Devices,
		Name:     sql.NullString{String: in.Name, Valid: in.Name != ""},
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

func getAllSnbConfig(ctx context.Context) ([]SnbConfig, error) {
	configs, err := database.New(database.D()).GetAllSnbConfig(ctx)
	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return []SnbConfig{}, err
	}

	var out []SnbConfig
	for _, config := range configs {
		out = append(out, SnbConfig{
			Id:         config.ID.String(),
			Name:       config.Name.String,
			Endpoint:   config.Endpoint.String,
			Facilities: config.Facility,
			Devices:    config.Device,
		})
	}
	return out, nil
}

func getSnbConfig(ctx context.Context, in uuid.UUID) (SnbConfig, error) {
	config, err := database.New(database.D()).GetSnbConfig(ctx, in)
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

func deleteSnbConfig(ctx context.Context, in uuid.UUID) error {
	_, err := database.New(database.D()).DeleteSnbConfig(ctx, in)
	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return err
	}
	return nil
}
