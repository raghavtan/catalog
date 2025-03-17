package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/motain/of-catalog/internal/modules/metric/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/metric/resources"
	compassserviceError "github.com/motain/of-catalog/internal/services/compassservice"
	compassservice "github.com/motain/of-catalog/internal/services/compassservice/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := NewRepository(mockCompass)

	tests := []struct {
		name          string
		metric        resources.Metric
		mockSetup     func()
		expectedID    string
		expectedError error
	}{
		{
			name:   "successful creation",
			metric: resources.Metric{Name: "test-metric"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateMetricOutput) error {
						output.Compass.CreateMetric.Definition.ID = "metric-id"
						output.Compass.CreateMetric.Success = true
						return nil
					},
				)
			},
			expectedID:    "metric-id",
			expectedError: nil,
		},
		{
			name:   "creation fails with error",
			metric: resources.Metric{Name: "test-metric"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("run error"))
			},
			expectedID:    "",
			expectedError: errors.New("run error"),
		},
		{
			name:   "creation fails with already exists error",
			metric: resources.Metric{Name: "test-metric"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateMetricOutput) error {
						output.Compass.CreateMetric.Errors = []compassserviceError.CompassError{{Message: "already exists"}}
						output.Compass.CreateMetric.Success = false
						return nil
					},
				)
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.SearchMetricsOutput) error {
						output.Compass.Definitions.Nodes = []dtos.Metric{{ID: "existing-id", Name: "test-metric"}}
						return nil
					},
				)
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.UpdateMetricOutput) error {
						output.Compass.UpdateMetric.Success = true
						return nil
					},
				)
			},
			expectedID:    "existing-id",
			expectedError: nil,
		},
		{
			name:   "creation fails with unsuccessful response",
			metric: resources.Metric{Name: "test-metric"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.CreateMetricOutput) error {
						output.Compass.CreateMetric.Definition.ID = ""
						return nil
					},
				)
			},
			expectedID:    "",
			expectedError: errors.New("failed to create metric"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			id, err := repo.Create(context.Background(), tt.metric)
			assert.Equal(t, tt.expectedID, id)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := NewRepository(mockCompass)

	tests := []struct {
		name          string
		metric        resources.Metric
		mockSetup     func()
		expectedError error
	}{
		{
			name:   "successful update",
			metric: resources.Metric{ID: "metric-id", Name: "updated-metric"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.UpdateMetricOutput) error {
						output.Compass.UpdateMetric.Success = true
						return nil
					},
				)
			},
			expectedError: nil,
		},
		{
			name:   "update fails with error",
			metric: resources.Metric{ID: "metric-id", Name: "updated-metric"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("run error"))
			},
			expectedError: errors.New("run error"),
		},
		{
			name:   "update fails with unsuccessful response",
			metric: resources.Metric{ID: "metric-id", Name: "updated-metric"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.UpdateMetricOutput) error {
						output.Compass.UpdateMetric.Success = false
						return nil
					},
				)
			},
			expectedError: errors.New("failed to update metric"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.Update(context.Background(), tt.metric)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := NewRepository(mockCompass)

	tests := []struct {
		name          string
		id            string
		mockSetup     func()
		expectedError error
	}{
		{
			name: "successful deletion",
			id:   "metric-id",
			mockSetup: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.DeleteMetricOutput) error {
						output.Compass.DeleteMetric.Success = true
						return nil
					},
				)
			},
			expectedError: nil,
		},
		{
			name: "deletion fails with error",
			id:   "metric-id",
			mockSetup: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("run error"))
			},
			expectedError: errors.New("run error"),
		},
		{
			name: "deletion fails with unsuccessful response",
			id:   "metric-id",
			mockSetup: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.DeleteMetricOutput) error {
						output.Compass.DeleteMetric.Success = false
						return nil
					},
				)
			},
			expectedError: errors.New("failed to delete metric"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repo.Delete(context.Background(), tt.id)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestRepository_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := NewRepository(mockCompass)

	tests := []struct {
		name           string
		metric         resources.Metric
		mockSetup      func()
		expectedMetric *resources.Metric
		expectedError  error
	}{
		{
			name:   "search returns no results",
			metric: resources.Metric{Name: "non-existent-metric"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.SearchMetricsOutput) error {
						output.Compass.Definitions.Nodes = []dtos.Metric{}
						return nil
					},
				)
			},
			expectedMetric: nil,
			expectedError:  errors.New("metric not found"),
		},
		{
			name:   "search returns one result",
			metric: resources.Metric{Name: "existing-metric"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.SearchMetricsOutput) error {
						output.Compass.Definitions.Nodes = []dtos.Metric{
							{ID: "metric-id", Name: "existing-metric"},
						}
						return nil
					},
				)
			},
			expectedMetric: &resources.Metric{ID: "metric-id"},
			expectedError:  nil,
		},
		{
			name:   "search returns multiple results",
			metric: resources.Metric{Name: "duplicate-metric"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, output *dtos.SearchMetricsOutput) error {
						output.Compass.Definitions.Nodes = []dtos.Metric{
							{ID: "metric-id-1", Name: "duplicate-metric"},
							{ID: "metric-id-2", Name: "duplicate-metric"},
						}
						return nil
					},
				)
			},
			expectedMetric: &resources.Metric{ID: "metric-id-1"},
			expectedError:  nil,
		},
		{
			name:   "search fails with error",
			metric: resources.Metric{Name: "error-metric"},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("run error"))
			},
			expectedMetric: nil,
			expectedError:  errors.New("run error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			result, err := repo.Search(tt.metric)
			assert.Equal(t, tt.expectedMetric, result)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
