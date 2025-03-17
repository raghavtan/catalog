package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	ID    string
	Value string
}

func fromStateToConfig(state *testStruct, conf *testStruct) {
	state.ID = conf.ID
}

func isEqual(a, b *testStruct) bool {
	return a.Value == b.Value
}

func TestDetect(t *testing.T) {
	tests := []struct {
		name              string
		stateMap          map[string]*testStruct
		configMap         map[string]*testStruct
		expectedCreate    map[string]*testStruct
		expectedUpdate    map[string]*testStruct
		expectedDelete    map[string]*testStruct
		expectedUnchanged map[string]*testStruct
	}{
		{
			name: "all unchanged",
			stateMap: map[string]*testStruct{
				"1": {ID: "1", Value: "a"},
				"2": {ID: "2", Value: "b"},
			},
			configMap: map[string]*testStruct{
				"1": {ID: "1", Value: "a"},
				"2": {ID: "2", Value: "b"},
			},
			expectedCreate: map[string]*testStruct{},
			expectedUpdate: map[string]*testStruct{},
			expectedDelete: map[string]*testStruct{},
			expectedUnchanged: map[string]*testStruct{
				"1": {ID: "1", Value: "a"},
				"2": {ID: "2", Value: "b"},
			},
		},
		{
			name: "one created, one updated, one deleted",
			stateMap: map[string]*testStruct{
				"1": {ID: "1", Value: "a"},
				"2": {ID: "2", Value: "b"},
			},
			configMap: map[string]*testStruct{
				"1": {ID: "1", Value: "a"},
				"2": {ID: "2", Value: "c"},
				"3": {ID: "3", Value: "d"},
			},
			expectedCreate: map[string]*testStruct{
				"3": {ID: "3", Value: "d"},
			},
			expectedUpdate: map[string]*testStruct{
				"2": {ID: "2", Value: "c"},
			},
			expectedDelete: map[string]*testStruct{},
			expectedUnchanged: map[string]*testStruct{
				"1": {ID: "1", Value: "a"},
			},
		},
		{
			name:     "all created",
			stateMap: map[string]*testStruct{},
			configMap: map[string]*testStruct{
				"1": {ID: "1", Value: "a"},
				"2": {ID: "2", Value: "b"},
			},
			expectedCreate: map[string]*testStruct{
				"1": {ID: "1", Value: "a"},
				"2": {ID: "2", Value: "b"},
			},
			expectedUpdate:    map[string]*testStruct{},
			expectedDelete:    map[string]*testStruct{},
			expectedUnchanged: map[string]*testStruct{},
		},
		{
			name: "all deleted",
			stateMap: map[string]*testStruct{
				"1": {ID: "1", Value: "a"},
				"2": {ID: "2", Value: "b"},
			},
			configMap:      map[string]*testStruct{},
			expectedCreate: map[string]*testStruct{},
			expectedUpdate: map[string]*testStruct{},
			expectedDelete: map[string]*testStruct{
				"1": {ID: "1", Value: "a"},
				"2": {ID: "2", Value: "b"},
			},
			expectedUnchanged: map[string]*testStruct{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, updated, deleted, unchanged := Detect(tt.stateMap, tt.configMap, fromStateToConfig, isEqual)
			matchElements(t, "created", tt.expectedCreate, created)
			matchElements(t, "updated", tt.expectedUpdate, updated)
			matchElements(t, "deleted", tt.expectedDelete, deleted)
			matchElements(t, "unchanged", tt.expectedUnchanged, unchanged)
		})
	}
}

func matchElements(t *testing.T, mapName string, expected, actual map[string]*testStruct) {
	for key, expectedValue := range expected {
		actualValue, exists := actual[key]
		assert.True(t, exists, "expected key %s to be in map %s", key, mapName)
		assert.Equal(t, expectedValue, actualValue, "expected value for key %s to match in %s", key, mapName)
	}
}
