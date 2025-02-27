package repository

//go:generate mockgen -destination=./mock_repository.go -package=repository github.com/motain/fact-collector/internal/modules/component/repository RepositoryInterface

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/motain/fact-collector/internal/modules/scorecard/resources"
	"github.com/motain/fact-collector/internal/services/compassservice"
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

func NewRepository(
	compass compassservice.CompassServiceInterface,
) *Repository {
	return &Repository{compass: compass}
}

func (r *Repository) Create(ctx context.Context, scorecard resources.Scorecard) (string, map[string]string, error) {
	query := `
		mutation createScorecard ($cloudId: ID!, $scorecardDetails: CreateCompassScorecardInput!) {
			compass {
				createScorecard(cloudId: $cloudId, input: $scorecardDetails) {
					success
					scorecardDetails {
						id
						criterias {
							id
							name
						}
					}
					errors {
						message
					}
				}
			}
		}`

	criteria := make([]map[string]map[string]string, len(scorecard.Criteria))
	for i, criterion := range scorecard.Criteria {
		criteria[i] = make(map[string]map[string]string)
		criteria[i]["hasMetricValue"] = make(map[string]string)
		criteria[i]["hasMetricValue"] = map[string]string{
			"weight":             fmt.Sprintf("%d", criterion.HasMetricValue.Weight),
			"name":               criterion.HasMetricValue.Name,
			"metricDefinitionId": criterion.HasMetricValue.MetricDefinitionId,
			"comparatorValue":    fmt.Sprintf("%d", criterion.HasMetricValue.ComparatorValue),
			"comparator":         criterion.HasMetricValue.Comparator,
		}
	}

	variables := map[string]interface{}{
		"cloudId": r.compass.GetCompassCloudId(),
		"scorecardDetails": map[string]interface{}{
			"name":                scorecard.Name,
			"description":         scorecard.Description,
			"state":               scorecard.State,
			"componentTypeIds":    scorecard.ComponentTypeIDs,
			"importance":          scorecard.Importance,
			"scoringStrategyType": scorecard.ScoringStrategyType,
			"criterias":           criteria,
		},
	}

	if scorecard.OwnerID != "" {
		variables["scorecardDetails"].(map[string]interface{})["ownerId"] = scorecard.OwnerID
	}

	var response struct {
		Compass struct {
			CreateScorecard struct {
				Success          bool `json:"success"`
				ScorecardDetails struct {
					ID       string `json:"id"`
					Criteria []struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"criterias"`
				} `json:"scorecardDetails"`
			} `json:"createScorecard"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to create scorecard: %v", err)
		return "", nil, err
	}

	criteriaMap := make(map[string]string, len(response.Compass.CreateScorecard.ScorecardDetails.Criteria))
	for _, criterion := range response.Compass.CreateScorecard.ScorecardDetails.Criteria {
		criteriaMap[criterion.Name] = criterion.ID
	}

	if !response.Compass.CreateScorecard.Success {
		return "", nil, errors.New("failed to create scorecard")
	}

	return response.Compass.CreateScorecard.ScorecardDetails.ID, criteriaMap, nil
}

func (r *Repository) Update(
	ctx context.Context,
	scorecard resources.Scorecard,
	createCriteria []*resources.Criterion,
	updateCriteria []*resources.Criterion,
	deleteCriteria []string,
) error {
	query := `
		mutation updateScorecard ($scorecardId: ID! $scorecardDetails: UpdateCompassScorecardInput!) {
			compass {
				updateScorecard(scorecardId: $scorecardId, input: $scorecardDetails) {
					success
					errors {
						message
					}
				}
			}
		}`

	criteriaToAdd := make([]map[string]map[string]string, len(createCriteria))
	for i, criterion := range createCriteria {
		criteriaToAdd[i] = make(map[string]map[string]string)
		criteriaToAdd[i]["hasMetricValue"] = make(map[string]string)
		criteriaToAdd[i]["hasMetricValue"] = map[string]string{
			"weight":             fmt.Sprintf("%d", criterion.HasMetricValue.Weight),
			"name":               criterion.HasMetricValue.Name,
			"metricDefinitionId": criterion.HasMetricValue.MetricDefinitionId,
			"comparatorValue":    fmt.Sprintf("%d", criterion.HasMetricValue.ComparatorValue),
			"comparator":         criterion.HasMetricValue.Comparator,
		}
	}

	criteriaToUpdate := make([]map[string]map[string]string, len(updateCriteria))
	for i, criterion := range updateCriteria {
		criteriaToUpdate[i] = make(map[string]map[string]string)
		criteriaToUpdate[i]["hasMetricValue"] = make(map[string]string)
		criteriaToUpdate[i]["hasMetricValue"] = map[string]string{
			"id":                 criterion.HasMetricValue.ID,
			"weight":             fmt.Sprintf("%d", criterion.HasMetricValue.Weight),
			"name":               criterion.HasMetricValue.Name,
			"metricDefinitionId": criterion.HasMetricValue.MetricDefinitionId,
			"comparatorValue":    fmt.Sprintf("%d", criterion.HasMetricValue.ComparatorValue),
			"comparator":         criterion.HasMetricValue.Comparator,
		}
	}

	variables := map[string]interface{}{
		"scorecardId": scorecard.ID,
		"scorecardDetails": map[string]interface{}{
			"name":                scorecard.Name,
			"description":         scorecard.Description,
			"state":               scorecard.State,
			"componentTypeIds":    scorecard.ComponentTypeIDs,
			"importance":          scorecard.Importance,
			"scoringStrategyType": scorecard.ScoringStrategyType,
			"createCriteria":      criteriaToAdd,
			"updateCriteria":      criteriaToUpdate,
			"deleteCriteria":      deleteCriteria,
		},
	}

	if scorecard.OwnerID != "" {
		variables["scorecardDetails"].(map[string]interface{})["ownerId"] = scorecard.OwnerID
	}

	var response struct {
		Compass struct {
			UpdateScorecard struct {
				Success bool `json:"success"`
			} `json:"updateScorecard"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to update scorecard: %v", err)
		return err
	}

	if !response.Compass.UpdateScorecard.Success {
		return errors.New("failed to update scorecard")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `
		mutation deleteScorecard($scorecardId: ID!) {
			compass {
				deleteScorecard(scorecardId: $scorecardId) {
					scorecardId
					errors {
						message
					}
					success
				}
			}
		}`

	variables := map[string]interface{}{
		"scorecardId": id,
	}

	var response struct {
		Compass struct {
			DeleteScorecard struct {
				Success bool `json:"success"`
			} `json:"deleteScorecard"`
		} `json:"compass"`
	}

	if err := r.compass.Run(ctx, query, variables, &response); err != nil {
		log.Printf("Failed to create scorecard: %v", err)
		return err
	}

	if !response.Compass.DeleteScorecard.Success {
		return errors.New("failed to delete scorecard")
	}

	return nil
}
