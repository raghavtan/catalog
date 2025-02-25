package factinterpreter_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/motain/fact-collector/internal/modules/metric/dtos"
	"github.com/motain/fact-collector/internal/modules/metric/handler/factcollectors"
	"github.com/motain/fact-collector/internal/modules/metric/handler/factinterpreter"
	"github.com/stretchr/testify/assert"
)

func TestProcessFacts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGithubFC := factcollectors.NewMockGithubFactCollectorInterface(ctrl)
	factInterpreter := factinterpreter.NewFactInterpreter(mockGithubFC)

	tests := []struct {
		name           string
		factOperations dtos.FactOperations
		mockSetup      func()
		expected       float64
	}{
		{
			name: "All facts succeed",
			factOperations: dtos.FactOperations{
				All: []dtos.Fact{{Source: "github"}},
				Any: nil,
			},
			mockSetup: func() {
				mockGithubFC.EXPECT().Check(gomock.Any()).Return(true, nil).Times(1)
			},
			expected: 1,
		},
		{
			name: "Any facts succeed",
			factOperations: dtos.FactOperations{
				All: nil,
				Any: []dtos.Fact{{Source: "github"}},
			},
			mockSetup: func() {
				mockGithubFC.EXPECT().Check(gomock.Any()).Return(true, nil).Times(1)
			},
			expected: 1,
		},
		{
			name: "All facts fail",
			factOperations: dtos.FactOperations{
				All: []dtos.Fact{{Source: "github"}},
				Any: nil,
			},
			mockSetup: func() {
				mockGithubFC.EXPECT().Check(gomock.Any()).Return(false, nil).Times(1)
			},
			expected: 0,
		},
		{
			name: "Any facts fail",
			factOperations: dtos.FactOperations{
				All: nil,
				Any: []dtos.Fact{{Source: "github"}},
			},
			mockSetup: func() {
				mockGithubFC.EXPECT().Check(gomock.Any()).Return(false, nil).Times(1)
			},
			expected: 0,
		},
		{
			name: "Mixed facts with all and any",
			factOperations: dtos.FactOperations{
				All: []dtos.Fact{{Source: "github"}},
				Any: []dtos.Fact{{Source: "github"}},
			},
			mockSetup: func() {
				mockGithubFC.EXPECT().Check(gomock.Any()).Return(true, nil).Times(2)
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result := factInterpreter.ProcessFacts(tt.factOperations)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessOperationsAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGithubFC := factcollectors.NewMockGithubFactCollectorInterface(ctrl)
	factInterpreter := factinterpreter.NewFactInterpreter(mockGithubFC)

	tests := []struct {
		name      string
		facts     []dtos.Fact
		mockSetup func()
		expected  bool
	}{
		{
			name:  "All facts succeed",
			facts: []dtos.Fact{{Source: "github"}},
			mockSetup: func() {
				mockGithubFC.EXPECT().Check(gomock.Any()).Return(true, nil).Times(1)
			},
			expected: true,
		},
		{
			name:  "All facts fail",
			facts: []dtos.Fact{{Source: "github"}},
			mockSetup: func() {
				mockGithubFC.EXPECT().Check(gomock.Any()).Return(false, nil).Times(1)
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result := factInterpreter.ProcessOperationsAll(tt.facts)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessOperationsAny(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGithubFC := factcollectors.NewMockGithubFactCollectorInterface(ctrl)
	factInterpreter := factinterpreter.NewFactInterpreter(mockGithubFC)

	tests := []struct {
		name      string
		facts     []dtos.Fact
		mockSetup func()
		expected  bool
	}{
		{
			name:  "Any facts succeed",
			facts: []dtos.Fact{{Source: "github"}},
			mockSetup: func() {
				mockGithubFC.EXPECT().Check(gomock.Any()).Return(true, nil).Times(1)
			},
			expected: true,
		},
		{
			name:  "Any facts fail",
			facts: []dtos.Fact{{Source: "github"}},
			mockSetup: func() {
				mockGithubFC.EXPECT().Check(gomock.Any()).Return(false, nil).Times(1)
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result := factInterpreter.ProcessOperationsAny(tt.facts)
			assert.Equal(t, tt.expected, result)
		})
	}
}
