package config

import "api-oa-integrator/database"

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
	Id                 string                 `json:"id,omitempty"`
	ClientId           string                 `json:"clientId,omitempty"`
	ProviderId         int32                  `json:"providerId,omitempty"`
	ServiceProviderId  string                 `json:"serviceProviderId,omitempty"`
	Name               string                 `json:"name,omitempty"`
	IntegratorName     string                 `json:"integratorName,omitempty"`
	Url                string                 `json:"url,omitempty"`
	InsecureSkipVerify bool                   `json:"insecureSkipVerify,omitempty"`
	PlazaIdMap         map[string]string      `json:"plazaIdMap,omitempty"`
	Extra              map[string]string      `json:"extra,omitempty"`
	TaxRate            float64                `json:"taxRate"`
	Surcharge          float64                `json:"surcharge"`
	SurchargeType      database.SurchargeType `json:"surchargeType,omitempty"`
}
