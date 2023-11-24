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
	config, err := database.New(database.D()).CreateSnbConfig(ctx, database.CreateSnbConfigParams{
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
	return SnbConfig{
		Endpoint:   config.Endpoint.String,
		Facilities: config.Facility,
		Devices:    config.Device,
	}, nil
}

func updateSnbConfig(ctx context.Context, id uuid.UUID, in SnbConfig) (SnbConfig, error) {
	config, err := database.New(database.D()).UpdateSnbConfig(ctx, database.UpdateSnbConfigParams{
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
	jsonString, err := json.Marshal(in.PlazaIdMap)
	config, err := database.New(database.D()).CreateIntegratorConfig(ctx, database.CreateIntegratorConfigParams{
		ClientID:           sql.NullString{String: in.ClientId, Valid: in.ClientId != ""},
		ProviderID:         sql.NullInt32{Int32: in.ProviderId, Valid: true},
		SpID:               sql.NullString{String: in.ServiceProviderId, Valid: in.ServiceProviderId != ""},
		Name:               sql.NullString{String: in.Name, Valid: in.Name != ""},
		InsecureSkipVerify: sql.NullBool{Bool: *(in.InsecureSkipVerify), Valid: in.InsecureSkipVerify != nil},
		PlazaIDMap:         pqtype.NullRawMessage{RawMessage: jsonString, Valid: jsonString != nil},
		Url:                sql.NullString{String: in.Url, Valid: in.Url != ""},
	})
	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return IntegratorConfig{}, err
	}

	var plazaId map[string]string
	_ = json.Unmarshal(config.PlazaIDMap.RawMessage, &plazaId)
	return IntegratorConfig{
		Id:                 config.ID.String(),
		ClientId:           config.ClientID.String,
		ProviderId:         config.ProviderID.Int32,
		ServiceProviderId:  config.SpID.String,
		Name:               config.Name.String,
		InsecureSkipVerify: &(config.InsecureSkipVerify.Bool),
		PlazaIdMap:         plazaId,
		Url:                config.Url.String,
	}, nil
}

func getIntegratorConfigs(ctx context.Context) ([]IntegratorConfig, error) {
	configs, err := database.New(database.D()).GetIntegratorConfigs(ctx)

	var out []IntegratorConfig
	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return out, err
	}
	for _, config := range configs {
		var plazaId map[string]string
		_ = json.Unmarshal(config.PlazaIDMap.RawMessage, &plazaId)
		out = append(out, IntegratorConfig{
			Id:                 config.ID.String(),
			ClientId:           config.ClientID.String,
			ProviderId:         config.ProviderID.Int32,
			ServiceProviderId:  config.SpID.String,
			Name:               config.Name.String,
			InsecureSkipVerify: &(config.InsecureSkipVerify.Bool),
			PlazaIdMap:         plazaId,
			Url:                config.Url.String,
		})
	}

	return out, nil
}

func getIntegratorConfig(ctx context.Context, id uuid.UUID) (IntegratorConfig, error) {
	config, err := database.New(database.D()).GetIntegratorConfig(ctx, id)

	if err != nil {
		zap.L().Sugar().Errorf("Error create user %v", err)
		return IntegratorConfig{}, err
	}
	var plazaId map[string]string
	_ = json.Unmarshal(config.PlazaIDMap.RawMessage, &plazaId)
	return IntegratorConfig{
		Id:                 config.ID.String(),
		ClientId:           config.ClientID.String,
		ProviderId:         config.ProviderID.Int32,
		ServiceProviderId:  config.SpID.String,
		Name:               config.Name.String,
		InsecureSkipVerify: &(config.InsecureSkipVerify.Bool),
		PlazaIdMap:         plazaId,
		Url:                config.Url.String,
	}, nil
}

func updateIntegratorConfig(ctx context.Context, id uuid.UUID, in IntegratorConfig) (IntegratorConfig, error) {
	jsonString, err := json.Marshal(in.PlazaIdMap)
	config, err := database.New(database.D()).UpdateIntegratorConfig(ctx, database.UpdateIntegratorConfigParams{
		ID:                 id,
		ClientID:           sql.NullString{String: in.ClientId, Valid: in.ClientId != ""},
		ProviderID:         sql.NullInt32{Int32: in.ProviderId, Valid: true},
		SpID:               sql.NullString{String: in.ServiceProviderId, Valid: in.ServiceProviderId != ""},
		Name:               sql.NullString{String: in.Name, Valid: in.Name != ""},
		InsecureSkipVerify: sql.NullBool{Bool: *(in.InsecureSkipVerify), Valid: in.InsecureSkipVerify != nil},
		PlazaIDMap:         pqtype.NullRawMessage{RawMessage: jsonString, Valid: err == nil},
		Url:                sql.NullString{String: in.Url, Valid: in.Url != ""},
	})
	if err != nil {
		zap.L().Sugar().Errorf("Error update integrator %v", err)
		return IntegratorConfig{}, err
	}
	return IntegratorConfig{
		Id:                 config.ID.String(),
		ClientId:           config.ClientID.String,
		ProviderId:         config.ProviderID.Int32,
		ServiceProviderId:  config.SpID.String,
		Name:               config.Name.String,
		InsecureSkipVerify: &(config.InsecureSkipVerify.Bool),
		PlazaIdMap:         in.PlazaIdMap,
		Url:                config.Url.String,
	}, nil
}

func deleteIntegratorConfig(ctx context.Context, id uuid.UUID) error {
	_, err := database.New(database.D()).DeleteIntegratorConfig(ctx, id)
	if err != nil {
		zap.L().Sugar().Errorf("Error update integrator %v", err)
		return err
	}
	return nil
}
