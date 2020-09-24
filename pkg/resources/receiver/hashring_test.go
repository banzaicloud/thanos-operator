package receiver

import (
	"encoding/json"
	"testing"
)

func TestHashRingMarshal(t *testing.T) {
	input := []HashRingGroup{
		{
			HashRing:  "test",
			Endpoints: []string{"endpoint1", "endpoint2"},
			Tenants:   []string{"test"},
		},
	}
	expected := `[{"hashring":"test","endpoints":["endpoint1","endpoint2"],"tenants":["test"]}]`
	result, err := json.Marshal(input)
	if err != nil {
		t.Error(err)
	}
	if expected != string(result) {
		t.Errorf("expected: %s != %s", expected, result)
	}
}

func TestGenerateHashring(t *testing.T) {
}
