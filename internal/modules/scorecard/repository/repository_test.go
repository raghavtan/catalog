package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/motain/of-catalog/internal/modules/scorecard/repository"
	"github.com/motain/of-catalog/internal/modules/scorecard/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/scorecard/resources"
	compassservice "github.com/motain/of-catalog/internal/services/compassservice/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)

	tests := []struct {
		name           string
		setupMocks     func()
		inputScorecard resources.Scorecard
		expectedID     string
		expectedMap    map[string]string
		expectedErr    error
	}{
		{
			name: "success",
			setupMocks: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(
					func(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) {
						dto := resp.(*dtos.CreateScorecardOutput)
						dto.Compass.CreateScorecard.Success = true
						dto.Compass.CreateScorecard.Scorecard = dtos.ScorecardDetails{
							ID: "scorecard-id",
							Criteria: []dtos.Criterion{
								{Name: "criterion1", ID: "criterion1-id"},
								{Name: "criterion2", ID: "criterion2-id"},
							},
						}
					},
				)
			},
			inputScorecard: resources.Scorecard{},
			expectedID:     "scorecard-id",
			expectedMap: map[string]string{
				"criterion1": "criterion1-id",
				"criterion2": "criterion2-id",
			},
			expectedErr: nil,
		},
		{
			name: "compass run error",
			setupMocks: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			inputScorecard: resources.Scorecard{},
			expectedID:     "",
			expectedMap:    nil,
			expectedErr:    errors.New("compass error"),
		},
		{
			name: "unsuccessful response",
			setupMocks: func() {
				mockCompass.EXPECT().GetCompassCloudId().Return("cloud-id")
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
						dto := resp.(*dtos.CreateScorecardOutput)
						dto.Compass.CreateScorecard.Scorecard = dtos.ScorecardDetails{}
						return nil
					},
				)
			},
			inputScorecard: resources.Scorecard{},
			expectedID:     "",
			expectedMap:    nil,
			expectedErr:    errors.New("failed to create scorecard"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			id, criteriaMap, err := repo.Create(context.Background(), tt.inputScorecard)

			assert.Equal(t, tt.expectedID, id)
			assert.Equal(t, tt.expectedMap, criteriaMap)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
func TestRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)
	scorecardID := "scorecard-id"

	tests := []struct {
		name           string
		setupMocks     func()
		inputScorecard resources.Scorecard
		createCriteria []*resources.Criterion
		updateCriteria []*resources.Criterion
		deleteCriteria []string
		expectedErr    error
	}{
		{
			name: "success",
			setupMocks: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
						dto := resp.(*dtos.UpdateScorecard)
						dto.Compass.UpdateScorecard.Success = true
						return nil
					},
				)
			},
			inputScorecard: resources.Scorecard{ID: &scorecardID},
			createCriteria: []*resources.Criterion{
				{
					HasMetricValue: resources.MetricValue{
						Weight:             1,
						Name:               "new-criterion",
						MetricDefinitionId: "new-criterion-id",
						ComparatorValue:    1,
						Comparator:         "eq",
					},
				},
			},
			updateCriteria: []*resources.Criterion{
				{
					HasMetricValue: resources.MetricValue{
						Weight:             1,
						Name:               "updated-criterion",
						MetricDefinitionId: "updated-criterion-id",
						ComparatorValue:    1,
						Comparator:         "eq",
					},
				},
			},
			deleteCriteria: []string{"delete-criterion-id"},
			expectedErr:    nil,
		},
		{
			name: "compass run error",
			setupMocks: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			inputScorecard: resources.Scorecard{ID: &scorecardID},
			createCriteria: []*resources.Criterion{
				{
					HasMetricValue: resources.MetricValue{
						Weight:             1,
						Name:               "new-criterion",
						MetricDefinitionId: "new-criterion-id",
						ComparatorValue:    1,
						Comparator:         "eq",
					},
				},
			},
			updateCriteria: []*resources.Criterion{
				{
					HasMetricValue: resources.MetricValue{
						Weight:             1,
						Name:               "updated-criterion",
						MetricDefinitionId: "updated-criterion-id",
						ComparatorValue:    1,
						Comparator:         "eq",
					},
				},
			},
			deleteCriteria: []string{"delete-criterion-id"},
			expectedErr:    errors.New("compass error"),
		},
		{
			name: "unsuccessful response",
			setupMocks: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
						dto := resp.(*dtos.UpdateScorecard)
						dto.Compass.UpdateScorecard.Success = false
						return nil
					},
				)
			},
			inputScorecard: resources.Scorecard{ID: &scorecardID},
			createCriteria: []*resources.Criterion{
				{
					HasMetricValue: resources.MetricValue{
						Weight:             1,
						Name:               "new-criterion",
						MetricDefinitionId: "new-criterion-id",
						ComparatorValue:    1,
						Comparator:         "eq",
					},
				},
			},
			updateCriteria: []*resources.Criterion{
				{
					HasMetricValue: resources.MetricValue{
						Weight:             1,
						Name:               "updated-criterion",
						MetricDefinitionId: "updated-criterion-id",
						ComparatorValue:    1,
						Comparator:         "eq",
					},
				},
			},
			deleteCriteria: []string{"delete-criterion-id"},
			expectedErr:    errors.New("failed to update scorecard"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := repo.Update(context.Background(), tt.inputScorecard, tt.createCriteria, tt.updateCriteria, tt.deleteCriteria)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
func TestRepository_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCompass := compassservice.NewMockCompassServiceInterface(ctrl)
	repo := repository.NewRepository(mockCompass)
	scorecardID := "scorecard-id"

	tests := []struct {
		name        string
		setupMocks  func()
		inputID     string
		expectedErr error
	}{
		{
			name: "success",
			setupMocks: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
						dto := resp.(*dtos.DeleteScorecard)
						dto.Compass.DeleteScorecard.Success = true
						return nil
					},
				)
			},
			inputID:     scorecardID,
			expectedErr: nil,
		},
		{
			name: "compass run error",
			setupMocks: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("compass error"))
			},
			inputID:     scorecardID,
			expectedErr: errors.New("compass error"),
		},
		{
			name: "unsuccessful response",
			setupMocks: func() {
				mockCompass.EXPECT().Run(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, query string, variables map[string]interface{}, resp interface{}) error {
						dto := resp.(*dtos.DeleteScorecard)
						dto.Compass.DeleteScorecard.Success = false
						return nil
					},
				)
			},
			inputID:     scorecardID,
			expectedErr: errors.New("failed to delete scorecard"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := repo.Delete(context.Background(), tt.inputID)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
