package config

import (
	"api-oa-integrator/internal/database"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"go.uber.org/zap"
)

func createSnbConfig(ctx context.Context, in SnbConfig) (SnbConfig, error) {
	txn, _ := database.D().Begin()
	config, err := database.New(database.D()).WithTx(txn).CreateSnbConfig(ctx, database.CreateSnbConfigParams{
		Endpoint: sql.NullString{String: in.Endpoint, Valid: true},
		Facility: in.Facilities,
		Device:   in.Devices,
		Name:     sql.NullString{String: in.Name, Valid: true},
		Username: sql.NullString{String: in.Username, Valid: true},
		Password: sql.NullString{String: in.Password, Valid: true},
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
		Username: sql.NullString{String: in.Username, Valid: in.Username != ""},
		Password: sql.NullString{String: in.Password, Valid: in.Password != ""},
	})
	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return SnbConfig{}, err
	}
	err = txn.Commit()
	return SnbConfig{
		Name:       config.Name.String,
		Endpoint:   config.Endpoint.String,
		Facilities: config.Facility,
		Devices:    config.Device,
		Username:   config.Username.String,
		Password:   config.Password.String,
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
		Name:       config.Name.String,
		Endpoint:   config.Endpoint.String,
		Facilities: config.Facility,
		Devices:    config.Device,
		Username:   config.Username.String,
		Password:   config.Password.String,
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

func createIntegratorConfig(ctx context.Context, in IntegratorConfig) (IntegratorConfig, error) {
	txn, _ := database.D().Begin()
	jsonString, _ := json.Marshal(in.PlazaIdMap)
	config, err := database.New(database.D()).WithTx(txn).CreateIntegratorConfig(ctx, database.CreateIntegratorConfigParams{
		ClientID:           sql.NullString{String: in.ClientId, Valid: true},
		ProviderID:         sql.NullInt32{Int32: in.ProviderId, Valid: true},
		SpID:               sql.NullString{String: in.ServiceProviderId, Valid: true},
		Name:               sql.NullString{String: in.Name, Valid: true},
		InsecureSkipVerify: sql.NullBool{Bool: in.InsecureSkipVerify, Valid: true},
		PlazaIDMap:         pqtype.NullRawMessage{Valid: true, RawMessage: jsonString},
		Url:                sql.NullString{String: in.Url, Valid: true},
	})
	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return IntegratorConfig{}, err
	}
	err = txn.Commit()
	return IntegratorConfig{
		Id:                 config.ID.String(),
		ClientId:           config.ClientID.String,
		ProviderId:         config.ProviderID.Int32,
		ServiceProviderId:  config.SpID.String,
		Name:               config.Name.String,
		InsecureSkipVerify: config.InsecureSkipVerify.Bool,
		PlazaIdMap:         in.PlazaIdMap,
		Url:                config.Url.String,
	}, nil
}

func getIntegratorConfigs(ctx context.Context, in IntegratorConfig) ([]IntegratorConfig, error) {
	txn, _ := database.D().Begin()
	configs, err := database.New(database.D()).WithTx(txn).GetIntegratorConfigs(ctx)

	var out []IntegratorConfig
	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return out, err
	}
	err = txn.Commit()

	for _, config := range configs {
		out = append(out, IntegratorConfig{
			Id:                 config.ID.String(),
			ClientId:           config.ClientID.String,
			ProviderId:         config.ProviderID.Int32,
			ServiceProviderId:  config.SpID.String,
			Name:               config.Name.String,
			InsecureSkipVerify: config.InsecureSkipVerify.Bool,
			PlazaIdMap:         in.PlazaIdMap,
			Url:                config.Url.String,
		})
	}

	return out, nil
}
