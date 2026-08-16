package oa

import (
	"api-oa-integrator/database"
	"testing"

	"github.com/sqlc-dev/pqtype"
)

// Two sites on one server, the case from the original report: each carries its
// own providerId, and the facility decides which one goes back to S&B.
const twoSites = `{"groups":[
	{"providerId":2,"clientId":"CZBAJL09","vendorLocationId":"AJL",
	 "facilities":["1077","1078"]},
	{"providerId":3,"clientId":"CZBAIL09","vendorLocationId":"AIL",
	 "facilities":["1060"]}]}`

func cfgWithMap(plazaMap string) database.IntegratorConfig {
	var c database.IntegratorConfig
	if plazaMap != "" {
		c.PlazaIDMap = pqtype.NullRawMessage{RawMessage: []byte(plazaMap), Valid: true}
	}
	return c
}

func TestProviderIdFor(t *testing.T) {
	tests := []struct {
		name     string
		plazaMap string
		facility string
		want     int32
	}{
		{"first site", twoSites, "1077", 2},
		{"same site, other facility", twoSites, "1078", 2},
		{"second site on the same server", twoSites, "1060", 3},
		{"unknown facility", twoSites, "9999", 0},
		{"null plaza map", "", "1077", 0},
		{"empty groups", `{"groups":[]}`, "1077", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := providerIdFor(cfgWithMap(tt.plazaMap), tt.facility); got != tt.want {
				t.Errorf("providerIdFor(%q) = %v, want %v", tt.facility, got, tt.want)
			}
		})
	}
}
