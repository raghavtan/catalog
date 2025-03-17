package repository

//go:generate mockgen -destination=./mock/mock_repository.go -package=repository github.com/motain/of-catalog/internal/modules/scorecard/repository RepositoryInterface

import (
	"context"
	"errors"
	"log"

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
	scoreCardDto := dtos.CreateScorecardOutput{}
	query := scoreCardDto.GetQuery()
	variables := scoreCardDto.SetVariables(r.compass.GetCompassCloudId(), scorecard)

	if err := r.compass.Run(ctx, query, variables, &scoreCardDto); err != nil {
		log.Printf("Failed to create scorecard: %v", err)
		return "", nil, err
	}

	if !scoreCardDto.IsSuccessful() {
		return "", nil, errors.New("failed to create scorecard")
	}

	scorecardDetails := scoreCardDto.Compass.CreateScorecard.Scorecard
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
	scoreCardDto := dtos.UpdateScorecard{}
	query := scoreCardDto.GetQuery()
	variables := scoreCardDto.SetVariables(scorecard, createCriteria, updateCriteria, deleteCriteria)

	if err := r.compass.Run(ctx, query, variables, &scoreCardDto); err != nil {
		log.Printf("Failed to update scorecard: %v", err)
		return err
	}

	if !scoreCardDto.IsSuccessful() {
		return errors.New("failed to update scorecard")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	scoreCardDto := dtos.DeleteScorecard{}
	query := scoreCardDto.GetQuery()
	variables := scoreCardDto.SetVariables(id)

	if err := r.compass.Run(ctx, query, variables, &scoreCardDto); err != nil {
		log.Printf("failed to delete scorecard: %v", err)
		return err
	}

	if !scoreCardDto.IsSuccessful() {
		return errors.New("failed to delete scorecard")
	}

	return nil
}
