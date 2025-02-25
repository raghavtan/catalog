package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/motain/fact-collector/internal/modules/metric/repository"
	"github.com/motain/fact-collector/internal/modules/metric/resources"
	"github.com/motain/fact-collector/internal/services/compassservice"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repository := repository.NewRepository(mockCompass)

	type args struct {
		ctx    context.Context
		metric resources.Metric
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func()
		want      string
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				metric: resources.Metric{
					Name:        "test-metric",
					Description: "test-description",
					Format: resources.MetricFormat{
						Unit: "unit",
					},
				},
			},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id").Times(1)
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, arg1, arg2, arg3 interface{}) {
					response := arg3.(*struct {
						Compass struct {
							CreateMetricDefinition struct {
								Success                bool `json:"success"`
								CreateMetricDefinition struct {
									ID string `json:"id"`
								} `json:"createdMetricDefinition"`
							} `json:"createMetricDefinition"`
						} `json:"compass"`
					})
					response.Compass.CreateMetricDefinition.Success = true
					response.Compass.CreateMetricDefinition.CreateMetricDefinition.ID = "metric-id"
				}).Return(nil).Times(1)
			},
			want:    "metric-id",
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				ctx: context.TODO(),
				metric: resources.Metric{
					Name:        "test-metric",
					Description: "test-description",
					Format: resources.MetricFormat{
						Unit: "unit",
					},
				},
			},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id").Times(1)
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to create metric")).Times(1)
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repository.Create(tt.args.ctx, tt.args.metric)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repository := repository.NewRepository(mockCompass)

	type args struct {
		ctx    context.Context
		metric resources.Metric
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				metric: resources.Metric{
					ID:          stringPtr("metric-id"),
					Name:        "test-metric",
					Description: "test-description",
					Format: resources.MetricFormat{
						Unit: "unit",
					},
				},
			},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id").Times(1)
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, query string, variables map[string]interface{}, response interface{}) {
					resp := response.(*struct {
						Compass struct {
							UpdateMetricDefinition struct {
								Success bool `json:"success"`
							} `json:"updateMetricDefinition"`
						} `json:"compass"`
					})
					resp.Compass.UpdateMetricDefinition.Success = true
				}).Return(nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				ctx: context.TODO(),
				metric: resources.Metric{
					ID:          stringPtr("metric-id"),
					Name:        "test-metric",
					Description: "test-description",
					Format: resources.MetricFormat{
						Unit: "unit",
					},
				},
			},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id").Times(1)
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to update metric")).Times(1)
			},
			wantErr: true,
		},
		{
			name: "update not successful",
			args: args{
				ctx: context.TODO(),
				metric: resources.Metric{
					ID:          stringPtr("metric-id"),
					Name:        "test-metric",
					Description: "test-description",
					Format: resources.MetricFormat{
						Unit: "unit",
					},
				},
			},
			mockSetup: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id").Times(1)
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, query string, variables map[string]interface{}, response interface{}) {
					resp := response.(*struct {
						Compass struct {
							UpdateMetricDefinition struct {
								Success bool `json:"success"`
							} `json:"updateMetricDefinition"`
						} `json:"compass"`
					})
					resp.Compass.UpdateMetricDefinition.Success = false
				}).Return(nil).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repository.Update(tt.args.ctx, tt.args.metric)
			if tt.wantErr {
				assert.Error(t, err)
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
	repository := repository.NewRepository(mockCompass)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				id:  "metric-id",
			},
			mockSetup: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, query string, variables map[string]interface{}, response interface{}) {
					resp := response.(*struct {
						Compass struct {
							DeleteMetricDefinition struct {
								Success bool `json:"success"`
							} `json:"deleteMetricDefinition"`
						} `json:"compass"`
					})
					resp.Compass.DeleteMetricDefinition.Success = true
				}).Return(nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				ctx: context.TODO(),
				id:  "metric-id",
			},
			mockSetup: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to delete metric")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repository.Delete(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRepository_CreateMetricSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repository := repository.NewRepository(mockCompass)

	type args struct {
		ctx         context.Context
		metricID    string
		componentID string
		intentifier string
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func()
		want      string
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx:         context.TODO(),
				metricID:    "metric-id",
				componentID: "component-id",
				intentifier: "external-id",
			},
			mockSetup: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, query string, variables map[string]interface{}, response interface{}) {
					resp := response.(*struct {
						Compass struct {
							CreateMetricSource struct {
								Success            bool `json:"success"`
								CreateMetricSource struct {
									ID string `json:"id"`
								} `json:"createdMetricSource"`
							} `json:"createMetricSource"`
						} `json:"compass"`
					})
					resp.Compass.CreateMetricSource.Success = true
					resp.Compass.CreateMetricSource.CreateMetricSource.ID = "metric-source-id"
				}).Return(nil).Times(1)
			},
			want:    "metric-source-id",
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				ctx:         context.TODO(),
				metricID:    "metric-id",
				componentID: "component-id",
				intentifier: "external-id",
			},
			mockSetup: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to create metric source")).Times(1)
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := repository.CreateMetricSource(tt.args.ctx, tt.args.metricID, tt.args.componentID, tt.args.intentifier)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRepository_DeleteMetricSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repository := repository.NewRepository(mockCompass)

	type args struct {
		ctx            context.Context
		metricSourceID string
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx:            context.TODO(),
				metricSourceID: "metric-source-id",
			},
			mockSetup: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(ctx context.Context, query string, variables map[string]interface{}, response interface{}) {
					resp := response.(*struct {
						Compass struct {
							DeleteMetricSource struct {
								Success bool `json:"success"`
							} `json:"deleteMetricSource"`
						} `json:"compass"`
					})
					resp.Compass.DeleteMetricSource.Success = true
				}).Return(nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				ctx:            context.TODO(),
				metricSourceID: "metric-source-id",
			},
			mockSetup: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("failed to delete metric source")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repository.DeleteMetricSource(tt.args.ctx, tt.args.metricSourceID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestRepository_Push(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repository := repository.NewRepository(mockCompass)

	type args struct {
		ctx            context.Context
		metricSourceID string
		value          float64
		recordedAt     time.Time
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func()
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				ctx:            context.TODO(),
				metricSourceID: "metric-source-id",
				value:          123.45,
				recordedAt:     time.Now(),
			},
			mockSetup: func() {
				mockCompass.EXPECT().SendMetric(gomock.Any()).Do(func(requestBody map[string]string) {
					assert.Equal(t, "metric-source-id", requestBody["metricSourceId"])
					assert.Equal(t, "123.450000", requestBody["value"])
					assert.NotEmpty(t, requestBody["timestamp"])
				}).Times(1).Return("", nil)
			},
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				ctx:            context.TODO(),
				metricSourceID: "metric-source-id",
				value:          123.45,
				recordedAt:     time.Now(),
			},
			mockSetup: func() {
				mockCompass.EXPECT().SendMetric(gomock.Any()).Return("", errors.New("failed to send metric")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := repository.Push(tt.args.ctx, tt.args.metricSourceID, tt.args.value, tt.args.recordedAt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
