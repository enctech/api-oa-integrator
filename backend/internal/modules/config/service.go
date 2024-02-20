package config

import (
	"api-oa-integrator/database"
	"api-oa-integrator/logger"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"strconv"
	"strings"
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
		logger.LogData("error", fmt.Sprintf("Error create snb config %v", err), nil)
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
		logger.LogData("error", fmt.Sprintf("Error update snb config: %v", err), nil)

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
		logger.LogData("error", fmt.Sprintf("error get all snb config %v", err), nil)
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
		logger.LogData("error", fmt.Sprintf("error get snb config %v", err), nil)
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
		logger.LogData("error", fmt.Sprintf("error delete snb config %v", err), nil)
		return err
	}
	return nil
}

func createIntegratorConfig(ctx context.Context, in IntegratorConfig) (IntegratorConfig, error) {
	jsonString, err := json.Marshal(in.PlazaIdMap)
	extraData, err := json.Marshal(in.Extra)
	config, err := database.New(database.D()).CreateIntegratorConfig(ctx, database.CreateIntegratorConfigParams{
		ClientID:           sql.NullString{String: in.ClientId, Valid: in.ClientId != ""},
		ProviderID:         sql.NullInt32{Int32: in.ProviderId, Valid: true},
		SpID:               sql.NullString{String: in.ServiceProviderId, Valid: in.ServiceProviderId != ""},
		Name:               sql.NullString{String: in.Name, Valid: in.Name != ""},
		InsecureSkipVerify: sql.NullBool{Bool: in.InsecureSkipVerify, Valid: true},
		PlazaIDMap:         pqtype.NullRawMessage{RawMessage: jsonString, Valid: jsonString != nil},
		Url:                sql.NullString{String: in.Url, Valid: in.Url != ""},
		IntegratorName:     sql.NullString{String: in.IntegratorName, Valid: in.Url != ""},
		Extra:              pqtype.NullRawMessage{RawMessage: extraData, Valid: extraData != nil || len(extraData) > 0},
	})
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error create integrator config %v", err), nil)
		return IntegratorConfig{}, err
	}

	var plazaId map[string]string
	_ = json.Unmarshal(config.PlazaIDMap.RawMessage, &plazaId)
	var extra map[string]string
	_ = json.Unmarshal(config.Extra.RawMessage, &extra)
	return IntegratorConfig{
		IntegratorName:     config.IntegratorName.String,
		Id:                 config.ID.String(),
		ClientId:           config.ClientID.String,
		ProviderId:         config.ProviderID.Int32,
		ServiceProviderId:  config.SpID.String,
		Name:               config.Name.String,
		InsecureSkipVerify: config.InsecureSkipVerify.Bool,
		PlazaIdMap:         plazaId,
		Url:                config.Url.String,
		Extra:              extra,
	}, nil
}

func getIntegratorConfigs(ctx context.Context) ([]IntegratorConfig, error) {
	configs, err := database.New(database.D()).GetIntegratorConfigs(ctx)

	var out []IntegratorConfig
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error get integrator configs %v", err), nil)
		return out, err
	}
	for _, config := range configs {
		var plazaId map[string]string
		_ = json.Unmarshal(config.PlazaIDMap.RawMessage, &plazaId)

		var extra map[string]string
		_ = json.Unmarshal(config.Extra.RawMessage, &plazaId)
		out = append(out, IntegratorConfig{
			IntegratorName:     config.IntegratorName.String,
			Id:                 config.ID.String(),
			ClientId:           config.ClientID.String,
			ProviderId:         config.ProviderID.Int32,
			ServiceProviderId:  config.SpID.String,
			Name:               config.Name.String,
			InsecureSkipVerify: config.InsecureSkipVerify.Bool,
			PlazaIdMap:         plazaId,
			Url:                config.Url.String,
			Extra:              extra,
		})
	}

	return out, nil
}

func getIntegratorConfig(ctx context.Context, id uuid.UUID) (IntegratorConfig, error) {
	config, err := database.New(database.D()).GetIntegratorConfig(ctx, id)

	if err != nil {
		logger.LogData("error", fmt.Sprintf("error get integrator config %v", err), nil)
		return IntegratorConfig{}, err
	}
	var plazaId map[string]string
	_ = json.Unmarshal(config.PlazaIDMap.RawMessage, &plazaId)

	var extra map[string]string
	_ = json.Unmarshal(config.Extra.RawMessage, &extra)

	surcharge := 0.0
	if s, err := strconv.ParseFloat(config.Surcharge.String, 32); err == nil {
		surcharge = s
	}

	taxRate := 0.0
	if s, err := strconv.ParseFloat(config.TaxRate.String, 32); err == nil {
		taxRate = s
	}
	return IntegratorConfig{
		IntegratorName:     config.IntegratorName.String,
		Id:                 config.ID.String(),
		ClientId:           config.ClientID.String,
		ProviderId:         config.ProviderID.Int32,
		ServiceProviderId:  config.SpID.String,
		Name:               config.Name.String,
		InsecureSkipVerify: config.InsecureSkipVerify.Bool,
		PlazaIdMap:         plazaId,
		Url:                config.Url.String,
		Extra:              extra,
		SurchargeType:      config.SurchangeType.SurchargeType,
		Surcharge:          surcharge,
		TaxRate:            taxRate,
	}, nil
}

func updateIntegratorConfig(ctx context.Context, id uuid.UUID, in IntegratorConfig) (IntegratorConfig, error) {
	jsonString, err := json.Marshal(in.PlazaIdMap)
	extraData, err := json.Marshal(in.Extra)
	config, err := database.New(database.D()).UpdateIntegratorConfig(ctx, database.UpdateIntegratorConfigParams{
		ID:                 id,
		ClientID:           sql.NullString{String: in.ClientId, Valid: in.ClientId != ""},
		ProviderID:         sql.NullInt32{Int32: in.ProviderId, Valid: true},
		SpID:               sql.NullString{String: in.ServiceProviderId, Valid: in.ServiceProviderId != ""},
		Name:               sql.NullString{String: in.Name, Valid: in.Name != ""},
		InsecureSkipVerify: sql.NullBool{Bool: in.InsecureSkipVerify, Valid: true},
		PlazaIDMap:         pqtype.NullRawMessage{RawMessage: jsonString, Valid: err == nil},
		Url:                sql.NullString{String: in.Url, Valid: in.Url != ""},
		IntegratorName:     sql.NullString{String: in.IntegratorName, Valid: in.Url != ""},
		Extra:              pqtype.NullRawMessage{RawMessage: extraData, Valid: extraData != nil || len(extraData) > 0},
		Surcharge:          sql.NullString{String: fmt.Sprintf("%f", in.Surcharge), Valid: true},
		TaxRate:            sql.NullString{String: fmt.Sprintf("%f", in.TaxRate), Valid: true},
		SurchangeType:      database.NullSurchargeType{SurchargeType: in.SurchargeType, Valid: true},
	})
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error update integrator config %v", err), nil)
		return IntegratorConfig{}, err
	}

	surchRes, err := strconv.ParseFloat(strings.TrimSpace(config.Surcharge.String), 64)
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error parse float surcharge %v", err), nil)
	}
	taxRateRes, err := strconv.ParseFloat(strings.TrimSpace(config.TaxRate.String), 64)
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error parse float TaxRate %v", err), nil)
	}

	var extra map[string]string
	_ = json.Unmarshal(config.Extra.RawMessage, &extra)
	return IntegratorConfig{
		IntegratorName:     config.IntegratorName.String,
		Id:                 config.ID.String(),
		ClientId:           config.ClientID.String,
		ProviderId:         config.ProviderID.Int32,
		ServiceProviderId:  config.SpID.String,
		Name:               config.Name.String,
		InsecureSkipVerify: config.InsecureSkipVerify.Bool,
		PlazaIdMap:         in.PlazaIdMap,
		Url:                config.Url.String,
		Extra:              extra,
		SurchargeType:      config.SurchangeType.SurchargeType,
		Surcharge:          surchRes,
		TaxRate:            taxRateRes,
	}, nil
}

func deleteIntegratorConfig(ctx context.Context, id uuid.UUID) error {
	_, err := database.New(database.D()).DeleteIntegratorConfig(ctx, id)
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error delete integrator config %v", err), nil)
		return err
	}
	return nil
}
