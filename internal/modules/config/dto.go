package config

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
	ClientId          string `json:"clientId,omitempty"`
	ProviderId        int    `json:"providerId,omitempty"`
	ServiceProviderId int    `json:"serviceProviderId,omitempty"`
	Name              string `json:"name,omitempty"`
	Url               string `json:"url,omitempty"`
}
