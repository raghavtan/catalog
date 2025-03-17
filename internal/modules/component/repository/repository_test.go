package repository_test

import (
	context "context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/motain/of-catalog/internal/modules/component/repository"
	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	compassserviceError "github.com/motain/of-catalog/internal/services/compassservice"
	compassservice "github.com/motain/of-catalog/internal/services/compassservice/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

	successComponent := dtos.Component{
		ID: "component-id",
		MetricSources: dtos.MetricSources{
			Nodes: []dtos.MetricSource{
				{
					ID: "metric-source-id",
					MetricDefinition: dtos.MetricDefinition{
						ID:   "metric-id",
						Name: "metric-name",
					},
				},
			},
		},
		Links: []dtos.Link{
			{
				ID:   "link-id",
				Type: "link-type",
				Name: "link-name",
				URL:  "link-url",
			},
		},
	}

	tests := []struct {
		name           string
		component      resources.Component
		mockSetup      func()
		expectedResult resources.Component
		expectedError  error
	}{
		{
			name:      "successfully creates a component",
			component: resources.Component{Slug: "test-slug"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateComponentOutput) error {
						output.Compass = dtos.CompassCreatedComponentOutput{
							CreateComponent: dtos.CompassCreateComponentOutput{
								Success: true,
								Details: successComponent,
							},
						}

						return nil
					},
				)
			},
			expectedResult: resources.Component{
				ID:   "component-id",
				Slug: "test-slug",
				MetricSources: map[string]*resources.MetricSource{
					"metric-name": {
						ID:     "metric-source-id",
						Metric: "metric-id",
					},
				},
				Links: []resources.Link{
					{
						ID:   "link-id",
						Type: "link-type",
						Name: "link-name",
						URL:  "link-url",
					},
				},
			},
			expectedError: nil,
		},
		{
			name:      "fails to create component due to compass error",
			component: resources.Component{},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateComponentOutput) error {
						return errors.New("compass error")
					},
				)
			},
			expectedResult: resources.Component{},
			expectedError:  errors.New("compass error"),
		},
		{
			name:      "fails to create component due to unsuccessful response",
			component: resources.Component{},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateComponentOutput) error {
						output.Compass = dtos.CompassCreatedComponentOutput{
							CreateComponent: dtos.CompassCreateComponentOutput{
								Success: false,
								Errors: []compassserviceError.CompassError{
									{Message: "error message"},
								},
							},
						}

						return nil
					},
				)
			},
			expectedResult: resources.Component{},
			expectedError:  errors.New("failed to create component: [{error message}]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			got, err := repo.Create(context.Background(), tt.component)
			// assert.True(t, reflect.DeepEqual(tt.expectedResult, got), "expected: %+v, got: %+v", tt.expectedResult, got)
			assert.Equal(t, tt.expectedResult, got)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
