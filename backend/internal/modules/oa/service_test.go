package oa

import (
	"api-oa-integrator/database"
	"database/sql"
	"testing"

	"github.com/sqlc-dev/pqtype"
)

func cfgWithMap(providerId int32, plazaMap string) database.IntegratorConfig {
	c := database.IntegratorConfig{ProviderID: sql.NullInt32{Int32: providerId, Valid: true}}
	if plazaMap != "" {
		c.PlazaIDMap = pqtype.NullRawMessage{RawMessage: []byte(plazaMap), Valid: true}
	}
	return c
}

func TestProviderIdFor(t *testing.T) {
	tests := []struct {
		name     string
		cfg      database.IntegratorConfig
		facility string
		want     int32
	}{
		{
			name:     "per-facility override wins",
			cfg:      cfgWithMap(10, `{"P01":{"vendorLocationId":"A","providerId":20}}`),
			facility: "P01",
			want:     20,
		},
		{
			name:     "no override for this facility falls back",
			cfg:      cfgWithMap(10, `{"P01":{"vendorLocationId":"A","providerId":20}}`),
			facility: "P02",
			want:     10,
		},
		{
			name:     "override omitted falls back",
			cfg:      cfgWithMap(10, `{"P01":{"vendorLocationId":"A"}}`),
			facility: "P01",
			want:     10,
		},
		{
			name:     "old string format falls back",
			cfg:      cfgWithMap(10, `{"P01":"A"}`),
			facility: "P01",
			want:     10,
		},
		{
			name:     "null plaza map falls back",
			cfg:      cfgWithMap(10, ""),
			facility: "P01",
			want:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := providerIdFor(tt.cfg, tt.facility); got != tt.want {
				t.Errorf("providerIdFor() = %v, want %v", got, tt.want)
			}
		})
	}
}
