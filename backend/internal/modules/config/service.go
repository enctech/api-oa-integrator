package config

import (
	"api-oa-integrator/database"
	"api-oa-integrator/logger"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
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
	extraData, err := json.Marshal(in.Extra)
	config, err := database.New(database.D()).CreateIntegratorConfig(ctx, database.CreateIntegratorConfigParams{
		SpID:               sql.NullString{String: in.ServiceProviderId, Valid: in.ServiceProviderId != ""},
		Name:               sql.NullString{String: in.Name, Valid: in.Name != ""},
		DisplayName:        sql.NullString{String: in.DisplayName, Valid: in.DisplayName != ""},
		InsecureSkipVerify: sql.NullBool{Bool: in.InsecureSkipVerify, Valid: true},
		Url:                sql.NullString{String: in.Url, Valid: in.Url != ""},
		IntegratorName:     sql.NullString{String: in.IntegratorName, Valid: in.Url != ""},
		Extra:              pqtype.NullRawMessage{RawMessage: extraData, Valid: extraData != nil || len(extraData) > 0},
		Surcharge:          sql.NullString{String: fmt.Sprintf("%f", in.Surcharge), Valid: true},
		TaxRate:            sql.NullString{String: fmt.Sprintf("%f", in.TaxRate), Valid: true},
		SurchangeType:      database.NullSurchargeType{SurchargeType: in.SurchargeType, Valid: true},
	})
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error create integrator config %v", err), nil)
		return IntegratorConfig{}, err
	}

	if err = replaceSites(ctx, config.ID, in.Groups); err != nil {
		logger.LogData("error", fmt.Sprintf("error create sites %v", err), nil)
		return IntegratorConfig{}, err
	}

	groups := sitesFor(ctx, config.ID)
	var extra map[string]string
	_ = json.Unmarshal(config.Extra.RawMessage, &extra)
	surchRes, err := strconv.ParseFloat(strings.TrimSpace(config.Surcharge.String), 64)
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error parse float surcharge %v", err), nil)
	}
	taxRateRes, err := strconv.ParseFloat(strings.TrimSpace(config.TaxRate.String), 64)
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error parse float TaxRate %v", err), nil)
	}

	return IntegratorConfig{
		IntegratorName:     config.IntegratorName.String,
		Id:                 config.ID.String(),
		ServiceProviderId:  config.SpID.String,
		Name:               config.Name.String,
		InsecureSkipVerify: config.InsecureSkipVerify.Bool,
		Groups:             groups,
		Url:                config.Url.String,
		Extra:              extra,
		SurchargeType:      config.SurchangeType.SurchargeType,
		Surcharge:          surchRes,
		TaxRate:            taxRateRes,
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
		groups := sitesFor(ctx, config.ID)

		var extra map[string]string
		_ = json.Unmarshal(config.Extra.RawMessage, &extra)
		out = append(out, IntegratorConfig{
			DisplayName:        config.DisplayName.String,
			IntegratorName:     config.IntegratorName.String,
			Id:                 config.ID.String(),
			ServiceProviderId:  config.SpID.String,
			Name:               config.Name.String,
			InsecureSkipVerify: config.InsecureSkipVerify.Bool,
			Groups:             groups,
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
	groups := sitesFor(ctx, config.ID)

	var extra map[string]string
	_ = json.Unmarshal(config.Extra.RawMessage, &extra)

	surcharge := 0.0
	if s, err := strconv.ParseFloat(strings.TrimSpace(config.Surcharge.String), 64); err == nil {
		surcharge = s
	}

	taxRate := 0.0
	if s, err := strconv.ParseFloat(strings.TrimSpace(config.TaxRate.String), 64); err == nil {
		taxRate = s
	}
	return IntegratorConfig{
		DisplayName:        config.DisplayName.String,
		IntegratorName:     config.IntegratorName.String,
		Id:                 config.ID.String(),
		ServiceProviderId:  config.SpID.String,
		Name:               config.Name.String,
		InsecureSkipVerify: config.InsecureSkipVerify.Bool,
		Groups:             groups,
		Url:                config.Url.String,
		Extra:              extra,
		SurchargeType:      config.SurchangeType.SurchargeType,
		Surcharge:          surcharge,
		TaxRate:            taxRate,
	}, nil
}

func updateIntegratorConfig(ctx context.Context, id uuid.UUID, in IntegratorConfig) (IntegratorConfig, error) {
	extraData, err := json.Marshal(in.Extra)
	config, err := database.New(database.D()).UpdateIntegratorConfig(ctx, database.UpdateIntegratorConfigParams{
		ID:                 id,
		SpID:               sql.NullString{String: in.ServiceProviderId, Valid: in.ServiceProviderId != ""},
		Name:               sql.NullString{String: in.Name, Valid: in.Name != ""},
		DisplayName:        sql.NullString{String: in.DisplayName, Valid: in.DisplayName != ""},
		InsecureSkipVerify: sql.NullBool{Bool: in.InsecureSkipVerify, Valid: true},
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

	if err = replaceSites(ctx, config.ID, in.Groups); err != nil {
		logger.LogData("error", fmt.Sprintf("error update sites %v", err), nil)
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
		ServiceProviderId:  config.SpID.String,
		Name:               config.Name.String,
		InsecureSkipVerify: config.InsecureSkipVerify.Bool,
		Groups:             sitesFor(ctx, config.ID),
		Url:                config.Url.String,
		Extra:              extra,
		SurchargeType:      config.SurchangeType.SurchargeType,
		Surcharge:          surchRes,
		TaxRate:            taxRateRes,
	}, nil
}

// ErrIntegratorConfigNotFound is returned when the config does not exist or
// has already been retired.
var ErrIntegratorConfigNotFound = errors.New("integrator config not found")

// deleteIntegratorConfig retires a config by setting deleted_at rather than
// removing the row. integrator_transactions and oa_transactions reference it
// and hold amounts and tax data, so a hard delete would either fail on those
// foreign keys or destroy financial history.
func deleteIntegratorConfig(ctx context.Context, id uuid.UUID) error {
	res, err := database.New(database.D()).DeleteIntegratorConfig(ctx, id)
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error delete integrator config %v", err), nil)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error reading rows affected %v", err), nil)
		return err
	}
	if rows == 0 {
		return ErrIntegratorConfigNotFound
	}
	return nil
}

// sitesFor reads a config's sites. Errors become an empty list: a config that
// cannot list its sites should still render in the UI rather than fail the
// whole request.
func sitesFor(ctx context.Context, configID uuid.UUID) []PlazaGroup {
	rows, err := database.New(database.D()).GetSitesByConfig(ctx, configID)
	if err != nil {
		logger.LogData("error", fmt.Sprintf("error get sites for config %v: %v", configID, err), nil)
		return []PlazaGroup{}
	}
	out := make([]PlazaGroup, 0, len(rows))
	for _, r := range rows {
		out = append(out, PlazaGroup{
			ProviderId:       r.ProviderID.Int32,
			ClientId:         r.ClientID.String,
			VendorLocationId: r.VendorLocationID.String,
			Facilities:       r.Facilities,
		})
	}
	return out
}

// replaceSites swaps a config's sites for the submitted set. The form always
// posts every site, so this is a replace rather than a diff; it runs in one
// transaction so a failure part-way cannot leave a config with half its sites.
func replaceSites(ctx context.Context, configID uuid.UUID, groups []PlazaGroup) error {
	tx, err := database.D().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := database.New(database.D()).WithTx(tx)
	if err = q.DeleteSitesByConfig(ctx, configID); err != nil {
		return err
	}
	for _, g := range groups {
		site, err := q.CreateSite(ctx, database.CreateSiteParams{
			IntegratorConfigID: configID,
			ProviderID:         sql.NullInt32{Int32: g.ProviderId, Valid: g.ProviderId != 0},
			ClientID:           sql.NullString{String: g.ClientId, Valid: g.ClientId != ""},
			VendorLocationID:   sql.NullString{String: g.VendorLocationId, Valid: g.VendorLocationId != ""},
		})
		if err != nil {
			return err
		}
		for _, facility := range g.Facilities {
			if facility == "" {
				continue
			}
			// The unique constraint rejects a facility already claimed by
			// another site of this config, which the form does not prevent.
			if err = q.AddSiteFacility(ctx, database.AddSiteFacilityParams{
				SiteID:             site.ID,
				IntegratorConfigID: configID,
				Facility:           facility,
			}); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}
