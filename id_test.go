package thebus

import (
	"testing"
)

func TestDefaultIDGenerator(t *testing.T) {
	var ids []string
	for i := 0; i < 100; i++ {
		ids = append(ids, DefaultIDGenerator())
	}
	uniqValues := make(map[string]bool)
	for _, id := range ids {
		uniqValues[id] = true
	}
	if len(uniqValues) != len(ids) {
		t.Errorf("uniqValues does not match len(ids): %v != %v", uniqValues, ids)
	}
}
