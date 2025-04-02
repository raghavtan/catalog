package repository

//go:generate mockgen -destination=./mocks/mock_repository.go -package=repository github.com/motain/of-catalog/internal/modules/scorecard/repository RepositoryInterface

import (
	"context"
	"fmt"

	"github.com/motain/of-catalog/internal/modules/scorecard/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/scorecard/resources"
	"github.com/motain/of-catalog/internal/services/compassservice"
)

type RepositoryInterface interface {
	Create(ctx context.Context, scorecard resources.Scorecard) (string, map[string]string, error)
	Update(
		ctx context.Context,
		scorecard resources.Scorecard,
		createCriteria []*resources.Criterion,
		updateCriteria []*resources.Criterion,
		deleteCriteria []string,
	) error
	Delete(ctx context.Context, id string) error
}

type Repository struct {
	compass compassservice.CompassServiceInterface
}

func NewRepository(compass compassservice.CompassServiceInterface) *Repository {
	return &Repository{compass: compass}
}

func (r *Repository) Create(ctx context.Context, scorecard resources.Scorecard) (string, map[string]string, error) {
	input := &dtos.CreateScorecardInput{CompassCloudID: r.compass.GetCompassCloudId(), Scorecard: scorecard}
	output := &dtos.CreateScorecardOutput{}
	if runErr := r.compass.RunWithDTOs(ctx, input, output); runErr != nil {
		return "", nil, fmt.Errorf("Create error for %s: %s", *scorecard.ID, runErr)
	}

	scorecardDetails := output.Compass.CreateScorecard.Scorecard
	criteriaMap := make(map[string]string, len(scorecardDetails.Criteria))
	for _, criterion := range scorecardDetails.Criteria {
		criteriaMap[criterion.Name] = criterion.ID
	}

	return scorecardDetails.ID, criteriaMap, nil
}

func (r *Repository) Update(
	ctx context.Context,
	scorecard resources.Scorecard,
	createCriteria []*resources.Criterion,
	updateCriteria []*resources.Criterion,
	deleteCriteria []string,
) error {
	input := &dtos.UpdateScorecardInput{
		Scorecard:      scorecard,
		CreateCriteria: createCriteria,
		UpdateCriteria: updateCriteria,
		DeleteCriteria: deleteCriteria,
	}
	output := &dtos.UpdateScorecardOutput{}
	if runErr := r.compass.RunWithDTOs(ctx, input, output); runErr != nil {
		return fmt.Errorf("Update error for %s: %s", *scorecard.ID, runErr)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	input := &dtos.DeleteScorecardInput{ScorecardID: id}
	output := &dtos.DeleteScorecardOutput{}
	if runErr := r.compass.RunWithDTOs(ctx, input, output); runErr != nil {
		return fmt.Errorf("Delete error for %s: %s", id, runErr)
	}
	return nil
}
