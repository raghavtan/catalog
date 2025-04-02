package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/motain/of-catalog/internal/modules/scorecard/repository"
	"github.com/motain/of-catalog/internal/modules/scorecard/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/scorecard/resources"
	compassmocks "github.com/motain/of-catalog/internal/services/compassservice/mocks"
)

func TestRepository_Create(t *testing.T) {
	// Define test cases
	testcases := []struct {
		name            string
		scorecard       resources.Scorecard
		setupMocks      func(compassMock *compassmocks.MockCompassServiceInterface)
		expectedID      string
		expectedCritMap map[string]string
		expectError     bool
		errorMessage    string
	}{
		{
			name: "successful creation",
			scorecard: resources.Scorecard{
				ID:   stringPtr("test-id"),
				Name: "Test Scorecard",
			},
			setupMocks: func(compassMock *compassmocks.MockCompassServiceInterface) {
				compassMock.EXPECT().GetCompassCloudId().Return("cloud-id")
				compassMock.EXPECT().RunWithDTOs(
					gomock.Any(),
					&dtos.CreateScorecardInput{
						CompassCloudID: "cloud-id",
						Scorecard: resources.Scorecard{
							ID:   stringPtr("test-id"),
							Name: "Test Scorecard",
						},
					},
					gomock.Any(),
				).DoAndReturn(func(_ context.Context, _ *dtos.CreateScorecardInput, output *dtos.CreateScorecardOutput) error {
					output.Compass.CreateScorecard.Scorecard = dtos.ScorecardDetails{
						ID: "returned-id",
						Criteria: []dtos.Criterion{
							{ID: "criterion-id-1", Name: "Criterion 1"},
							{ID: "criterion-id-2", Name: "Criterion 2"},
						},
					}
					return nil
				})
			},
			expectedID: "returned-id",
			expectedCritMap: map[string]string{
				"Criterion 1": "criterion-id-1",
				"Criterion 2": "criterion-id-2",
			},
			expectError: false,
		},
		{
			name: "creation error",
			scorecard: resources.Scorecard{
				ID:   stringPtr("test-id"),
				Name: "Test Scorecard",
			},
			setupMocks: func(compassMock *compassmocks.MockCompassServiceInterface) {
				compassMock.EXPECT().GetCompassCloudId().Return("cloud-id")
				compassMock.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("compass error"))
			},
			expectError:  true,
			errorMessage: "Create error for test-id: compass error",
		},
	}

	// Run test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			compassMock := compassmocks.NewMockCompassServiceInterface(ctrl)
			tc.setupMocks(compassMock)

			repo := repository.NewRepository(compassMock)

			// Execute
			id, critMap, err := repo.Create(context.Background(), tc.scorecard)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.errorMessage, err.Error())
				assert.Empty(t, id)
				assert.Empty(t, critMap)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, id)
				assert.Equal(t, tc.expectedCritMap, critMap)
			}
		})
	}
}
func TestRepository_Update(t *testing.T) {
	// Define test cases
	testcases := []struct {
		name           string
		scorecard      resources.Scorecard
		createCriteria []*resources.Criterion
		updateCriteria []*resources.Criterion
		deleteCriteria []string
		setupMocks     func(compassMock *compassmocks.MockCompassServiceInterface)
		expectError    bool
		errorMessage   string
	}{
		{
			name: "successful update",
			scorecard: resources.Scorecard{
				ID:   stringPtr("test-id"),
				Name: "Updated Scorecard",
			},
			createCriteria: []*resources.Criterion{
				{HasMetricValue: resources.MetricValue{Name: "New Criterion"}},
			},
			updateCriteria: []*resources.Criterion{
				{HasMetricValue: resources.MetricValue{ID: "criterion-id-1", Name: "Updated Criterion"}},
			},
			deleteCriteria: []string{"criterion-id-2"},
			setupMocks: func(compassMock *compassmocks.MockCompassServiceInterface) {
				compassMock.EXPECT().RunWithDTOs(
					gomock.Any(),
					&dtos.UpdateScorecardInput{
						Scorecard: resources.Scorecard{
							ID:   stringPtr("test-id"),
							Name: "Updated Scorecard",
						},
						CreateCriteria: []*resources.Criterion{
							{HasMetricValue: resources.MetricValue{Name: "New Criterion"}},
						},
						UpdateCriteria: []*resources.Criterion{
							{HasMetricValue: resources.MetricValue{ID: "criterion-id-1", Name: "Updated Criterion"}},
						},
						DeleteCriteria: []string{"criterion-id-2"},
					},
					gomock.Any(),
				).Return(nil)
			},
			expectError: false,
		},
		{
			name: "update error",
			scorecard: resources.Scorecard{
				ID:   stringPtr("test-id"),
				Name: "Updated Scorecard",
			},
			createCriteria: []*resources.Criterion{},
			updateCriteria: []*resources.Criterion{},
			deleteCriteria: []string{},
			setupMocks: func(compassMock *compassmocks.MockCompassServiceInterface) {
				compassMock.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("compass error"))
			},
			expectError:  true,
			errorMessage: "Update error for test-id: compass error",
		},
	}

	// Run test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			compassMock := compassmocks.NewMockCompassServiceInterface(ctrl)
			tc.setupMocks(compassMock)

			repo := repository.NewRepository(compassMock)

			// Execute
			err := repo.Update(
				context.Background(),
				tc.scorecard,
				tc.createCriteria,
				tc.updateCriteria,
				tc.deleteCriteria,
			)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestRepository_Delete(t *testing.T) {
	// Define test cases
	testcases := []struct {
		name         string
		scorecardID  string
		setupMocks   func(compassMock *compassmocks.MockCompassServiceInterface)
		expectError  bool
		errorMessage string
	}{
		{
			name:        "successful deletion",
			scorecardID: "test-id",
			setupMocks: func(compassMock *compassmocks.MockCompassServiceInterface) {
				compassMock.EXPECT().RunWithDTOs(
					gomock.Any(),
					&dtos.DeleteScorecardInput{
						ScorecardID: "test-id",
					},
					gomock.Any(),
				).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "deletion error",
			scorecardID: "test-id",
			setupMocks: func(compassMock *compassmocks.MockCompassServiceInterface) {
				compassMock.EXPECT().RunWithDTOs(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errors.New("compass error"))
			},
			expectError:  true,
			errorMessage: "Delete error for test-id: compass error",
		},
	}

	// Run test cases
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			compassMock := compassmocks.NewMockCompassServiceInterface(ctrl)
			tc.setupMocks(compassMock)

			repo := repository.NewRepository(compassMock)

			// Execute
			err := repo.Delete(context.Background(), tc.scorecardID)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.errorMessage, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
