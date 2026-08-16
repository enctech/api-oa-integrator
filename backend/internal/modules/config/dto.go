package config

import (
	"api-oa-integrator/database"
	"api-oa-integrator/internal/plaza"
)

type SnbConfig struct {
	Id         string   `json:"id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Username   string   `json:"username,omitempty"`
	Password   string   `json:"password,omitempty"`
	Endpoint   string   `json:"endpoint"`
	Facilities []string `json:"facilities"`
	Devices    []string `json:"devices"`
}

type IntegratorConfig struct {
	Id                string `json:"id,omitempty"`
	ServiceProviderId string `json:"serviceProviderId,omitempty"`
	Name              string `json:"name,omitempty"`
	DisplayName       string `json:"displayName,omitempty"`
	IntegratorName    string `json:"integratorName,omitempty"`
	Url               string `json:"url,omitempty"`
	// Groups replace the config-level clientId and providerId. One config
	// covers several sites - `name` is uniquely indexed, so sites sharing an
	// integrator cannot be separate rows - and each site owns its own pair.
	Groups             []plaza.Group          `json:"groups"`
	InsecureSkipVerify bool                   `json:"insecureSkipVerify,omitempty"`
	Extra              map[string]string      `json:"extra,omitempty"`
	TaxRate            float64                `json:"taxRate"`
	Surcharge          float64                `json:"surcharge"`
	SurchargeType      database.SurchargeType `json:"surchargeType,omitempty"`
}
