package config

type SnbConfig struct {
	Id         string   `json:"id,omitempty"`
	Endpoint   string   `json:"endpoint"`
	Facilities []string `json:"facilities"`
	Devices    []string `json:"devices"`
}
