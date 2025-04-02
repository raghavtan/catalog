package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/motain/of-catalog/internal/modules/metric/repository"
	"github.com/motain/of-catalog/internal/modules/metric/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/metric/resources"
	compassservice "github.com/motain/of-catalog/internal/services/compassservice/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)
	ctx := context.Background()

	tests := []struct {
		name           string
		metric         resources.Metric
		setupMocks     func()
		expectedID     string
		expectedError  bool
		errorSubstring string
	}{
		{
			name: "successful creation",
			metric: resources.Metric{
				Name: "test-metric",
			},
			setupMocks: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-123")
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(_ context.Context, input *dtos.CreateMetricInput, output *dtos.CreateMetricOutput) error {
					output.Compass.CreateMetric.Definition.ID = "metric-123"
					output.Compass.CreateMetric.Success = true
					return nil
				})
			},
			expectedID:    "metric-123",
			expectedError: false,
		},
		{
			name: "error during creation",
			metric: resources.Metric{
				Name: "test-metric",
			},
			setupMocks: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-123")
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("compass error"))
			},
			expectedID:     "",
			expectedError:  true,
			errorSubstring: "Create error for",
		},
		// We're not testing the pre-validation function here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			id, err := repo.Create(ctx, tt.metric)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorSubstring)
				assert.Equal(t, tt.expectedID, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}
func TestRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)
	ctx := context.Background()

	tests := []struct {
		name           string
		metric         resources.Metric
		setupMocks     func()
		expectedError  bool
		errorSubstring string
	}{
		{
			name: "successful update",
			metric: resources.Metric{
				ID:   "metric-123",
				Name: "updated-metric",
			},
			setupMocks: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-123")
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(_ context.Context, input *dtos.UpdateMetricInput, output *dtos.UpdateMetricOutput) error {
					output.Compass.UpdateMetric.Success = true
					return nil
				})
			},
			expectedError: false,
		},
		{
			name: "error during update",
			metric: resources.Metric{
				ID:   "metric-123",
				Name: "updated-metric",
			},
			setupMocks: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-123")
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("compass error"))
			},
			expectedError:  true,
			errorSubstring: "Update error for",
		},
		{
			name: "metric with empty ID",
			metric: resources.Metric{
				Name: "metric-without-id",
			},
			setupMocks: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-123")
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(_ context.Context, input *dtos.UpdateMetricInput, output *dtos.UpdateMetricOutput) error {
					output.Compass.UpdateMetric.Success = true
					return nil
				})
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := repo.Update(ctx, tt.metric)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorSubstring)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)
	ctx := context.Background()

	tests := []struct {
		name           string
		metricID       string
		setupMocks     func()
		expectedError  bool
		errorSubstring string
	}{
		{
			name:     "successful deletion",
			metricID: "metric-123",
			setupMocks: func() {
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(_ context.Context, input *dtos.DeleteMetricInput, output *dtos.DeleteMetricOutput) error {
					output.Compass.DeleteMetric.Success = true
					return nil
				})
			},
			expectedError: false,
		},
		{
			name:     "error during deletion",
			metricID: "metric-123",
			setupMocks: func() {
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("compass error"))
			},
			expectedError:  true,
			errorSubstring: "Delete error for",
		},
		{
			name:     "empty metric ID",
			metricID: "",
			setupMocks: func() {
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("empty metric ID"))
			},
			expectedError:  true,
			errorSubstring: "Delete error for",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := repo.Delete(ctx, tt.metricID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorSubstring)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestRepository_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)
	ctx := context.Background()

	tests := []struct {
		name           string
		metric         resources.Metric
		setupMocks     func()
		expectedMetric *resources.Metric
		expectedError  bool
		errorSubstring string
	}{
		{
			name: "successful search",
			metric: resources.Metric{
				Name: "test-metric",
			},
			setupMocks: func() {
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(_ context.Context, input *dtos.SearchMetricsInput, output *dtos.SearchMetricsOutput) error {
					output.Compass.Definitions.Nodes = []dtos.Metric{
						{
							ID:   "metric-123",
							Name: "test-metric",
						},
					}
					return nil
				})
			},
			expectedMetric: &resources.Metric{
				ID: "metric-123",
			},
			expectedError: false,
		},
		{
			name: "metric not found",
			metric: resources.Metric{
				Name: "non-existent-metric",
			},
			setupMocks: func() {
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(_ context.Context, input *dtos.SearchMetricsInput, output *dtos.SearchMetricsOutput) error {
					output.Compass.Definitions.Nodes = []dtos.Metric{
						{
							ID:   "other-metric-123",
							Name: "other-metric",
						},
					}
					return nil
				})
			},
			expectedMetric: nil,
			expectedError:  true,
			errorSubstring: "metric not found",
		},
		{
			name: "service error",
			metric: resources.Metric{
				Name: "test-metric",
			},
			setupMocks: func() {
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("compass error"))
			},
			expectedMetric: nil,
			expectedError:  true,
			errorSubstring: "Search error for",
		},
		{
			name: "empty nodes response",
			metric: resources.Metric{
				Name: "test-metric",
			},
			setupMocks: func() {
				mockCompass.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).DoAndReturn(func(_ context.Context, input *dtos.SearchMetricsInput, output *dtos.SearchMetricsOutput) error {
					output.Compass.Definitions.Nodes = []dtos.Metric{}
					return nil
				})
			},
			expectedMetric: nil,
			expectedError:  true,
			errorSubstring: "metric not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			metric, err := repo.Search(ctx, tt.metric)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorSubstring)
				assert.Nil(t, metric)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedMetric.ID, metric.ID)
			}
		})
	}
}
