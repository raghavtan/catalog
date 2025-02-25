package drift_test

import (
	"testing"

	"github.com/motain/fact-collector/internal/utils/drift"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	ID    string
	Value string
}

func TestDetect(t *testing.T) {
	tests := []struct {
		name              string
		stateList         []*testStruct
		configList        []*testStruct
		expectedCreate    []*testStruct
		expectedUpdate    []*testStruct
		expectedDelete    []*testStruct
		expectedUnchanged []*testStruct
	}{
		{
			name: "all unchanged",
			stateList: []*testStruct{
				{ID: "1", Value: "a"},
				{ID: "2", Value: "b"},
			},
			configList: []*testStruct{
				{ID: "1", Value: "a"},
				{ID: "2", Value: "b"},
			},
			expectedCreate: []*testStruct{},
			expectedUpdate: []*testStruct{},
			expectedDelete: []*testStruct{},
			expectedUnchanged: []*testStruct{
				{ID: "1", Value: "a"},
				{ID: "2", Value: "b"},
			},
		},
		{
			name: "one created, one updated, one deleted",
			stateList: []*testStruct{
				{ID: "1", Value: "a"},
				{ID: "2", Value: "b"},
			},
			configList: []*testStruct{
				{ID: "1", Value: "a1"},
				{ID: "3", Value: "c"},
			},
			expectedCreate: []*testStruct{
				{ID: "3", Value: "c"},
			},
			expectedUpdate: []*testStruct{
				{ID: "1", Value: "a1"},
			},
			expectedDelete: []*testStruct{
				{ID: "2", Value: "b"},
			},
			expectedUnchanged: []*testStruct{},
		},
		{
			name:      "all created",
			stateList: []*testStruct{},
			configList: []*testStruct{
				{ID: "1", Value: "a"},
				{ID: "2", Value: "b"},
			},
			expectedCreate: []*testStruct{
				{ID: "1", Value: "a"},
				{ID: "2", Value: "b"},
			},
			expectedUpdate:    []*testStruct{},
			expectedDelete:    []*testStruct{},
			expectedUnchanged: []*testStruct{},
		},
		{
			name: "all_deleted",
			stateList: []*testStruct{
				{ID: "1", Value: "a"},
				{ID: "2", Value: "b"},
			},
			configList:     []*testStruct{},
			expectedCreate: []*testStruct{},
			expectedUpdate: []*testStruct{},
			expectedDelete: []*testStruct{
				{ID: "1", Value: "a"},
				{ID: "2", Value: "b"},
			},
			expectedUnchanged: []*testStruct{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, updated, deleted, unchanged := drift.Detect(
				tt.stateList,
				tt.configList,
				func(t *testStruct) string { return t.ID },
				func(t *testStruct) string { return t.ID },
				func(t *testStruct, id string) { t.ID = id },
				func(t1, t2 *testStruct) bool { return t1.Value == t2.Value },
			)

			assert.ElementsMatch(t, tt.expectedCreate, created)
			assert.ElementsMatch(t, tt.expectedUpdate, updated)
			assert.ElementsMatch(t, tt.expectedDelete, deleted)
			assert.ElementsMatch(t, tt.expectedUnchanged, unchanged)
		})
	}
}
