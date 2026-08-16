// Package plaza decodes the integrator_config.plaza_id_map column.
//
// One integrator_config row covers several sites, because `name` is the
// config identity and is uniquely indexed - four sites all named TNG cannot
// be four rows. Each site is a group here: it owns the providerId that S&B
// expects (S&B accepts a single providerId per server), the clientId the
// vendor issued, and the vendor location those map to. It lists the OA
// facilities that belong to it.
//
// Only the grouped format is understood. Migration 000019 converts the two
// older layouts, so nothing reaching this code is still flat.
package plaza

import "encoding/json"

type Group struct {
	ProviderId int32  `json:"providerId,omitempty"`
	ClientId   string `json:"clientId,omitempty"`
	// VendorLocationId is per site, not per facility: every facility on a
	// site reports to the same vendor location.
	VendorLocationId string `json:"vendorLocationId,omitempty"`
	// Facilities are the OA facility IDs belonging to this site.
	Facilities []string `json:"facilities"`
}

type Map struct {
	Groups []Group `json:"groups"`
}

// Parse returns an empty Map rather than an error: a config with an
// unreadable plaza map should fail the facility lookup below, which callers
// already handle, instead of taking down every other config in the fan-out.
func Parse(raw []byte) Map {
	var m Map
	if err := json.Unmarshal(raw, &m); err != nil {
		return Map{}
	}
	return m
}

// Find returns the group owning facility. Groups are few (one per site) so a
// scan is cheaper than an index.
func (m Map) Find(facility string) (Group, bool) {
	for _, g := range m.Groups {
		for _, f := range g.Facilities {
			if f == facility {
				return g, true
			}
		}
	}
	return Group{}, false
}
