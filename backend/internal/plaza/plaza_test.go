package plaza

import "testing"

const fourSites = `{"groups":[
	{"providerId":2,"clientId":"CZBAJL09","vendorLocationId":"AJL",
	 "facilities":["1077","1078","1888","1165"]},
	{"providerId":3,"clientId":"CZBAIL09","vendorLocationId":"AIL",
	 "facilities":["1060","1058","1270"]},
	{"providerId":4,"clientId":"CZBAGL09","vendorLocationId":"AGL",
	 "facilities":["1114","1116"]},
	{"providerId":1,"clientId":"CZBAHL09","vendorLocationId":"AHL",
	 "facilities":["2200","2300","2301","2323"]}]}`

func TestFind(t *testing.T) {
	m := Parse([]byte(fourSites))
	if len(m.Groups) != 4 {
		t.Fatalf("got %d groups, want 4", len(m.Groups))
	}

	tests := []struct {
		facility        string
		wantClientId    string
		wantProviderId  int32
		wantVendorLocId string
	}{
		{"1077", "CZBAJL09", 2, "AJL"},
		{"1165", "CZBAJL09", 2, "AJL"},
		{"1270", "CZBAIL09", 3, "AIL"},
		{"1116", "CZBAGL09", 4, "AGL"},
		{"2323", "CZBAHL09", 1, "AHL"},
	}

	for _, tt := range tests {
		t.Run(tt.facility, func(t *testing.T) {
			g, ok := m.Find(tt.facility)
			if !ok {
				t.Fatalf("facility %v not found", tt.facility)
			}
			if g.ClientId != tt.wantClientId {
				t.Errorf("clientId = %v, want %v", g.ClientId, tt.wantClientId)
			}
			if g.ProviderId != tt.wantProviderId {
				t.Errorf("providerId = %v, want %v", g.ProviderId, tt.wantProviderId)
			}
			if g.VendorLocationId != tt.wantVendorLocId {
				t.Errorf("vendorLocationId = %v, want %v", g.VendorLocationId, tt.wantVendorLocId)
			}
		})
	}
}

func TestFindMisses(t *testing.T) {
	m := Parse([]byte(fourSites))
	if _, ok := m.Find("9999"); ok {
		t.Error("unknown facility reported as found")
	}
	// Facilities are compared as strings, so no numeric coercion applies.
	if _, ok := m.Find("01077"); ok {
		t.Error("zero-padded facility matched")
	}
}

func TestParseRejectsLegacyAndGarbage(t *testing.T) {
	// Migration 000019 converts every older layout, so none should resolve
	// here. Reaching this code means the row was never migrated, and a silent
	// wrong providerId would be worse than a miss.
	for _, raw := range []string{
		`{"1225":"ACL"}`,
		`{"1225":{"vendorLocationId":"ACL","clientId":"X","providerId":2}}`,
		`{"groups":[{"providerId":2,"facilities":{"1225":"ACL"}}]}`,
		`not json`,
		``,
	} {
		m := Parse([]byte(raw))
		if _, ok := m.Find("1225"); ok {
			t.Errorf("Parse(%q) resolved a facility", raw)
		}
	}
}
