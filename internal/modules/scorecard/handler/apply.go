package handler

import (
	"context"
	"log"

	metricdtos "github.com/motain/of-catalog/internal/modules/metric/dtos"
	"github.com/motain/of-catalog/internal/modules/scorecard/dtos"
	"github.com/motain/of-catalog/internal/modules/scorecard/repository"
	"github.com/motain/of-catalog/internal/modules/scorecard/resources"
	"github.com/motain/of-catalog/internal/utils/drift"
	"github.com/motain/of-catalog/internal/utils/yaml"
)

type ApplyHandler struct {
	repository repository.RepositoryInterface
}

func NewApplyHandler(
	repository repository.RepositoryInterface,
) *ApplyHandler {
	return &ApplyHandler{repository: repository}
}

func (h *ApplyHandler) Apply(ctx context.Context, configRootLocation string, stateRootLocation string, recursive bool) {
	parseInput := yaml.ParseInput{
		RootLocation: configRootLocation,
		Recursive:    recursive,
	}
	configScorecards, errConfig := yaml.Parse(parseInput, dtos.GetScorecardUniqueKey)
	if errConfig != nil {
		log.Fatalf("error: %v", errConfig)
	}

	// Read metrics from split metric state files
	stateMetrics, errMetricState := yaml.Parse(yaml.GetMetricStateInput(), metricdtos.GetMetricUniqueKey)
	if errMetricState != nil {
		log.Fatalf("error: %v", errMetricState)
	}

	for _, scorecard := range configScorecards {
		for _, criterion := range scorecard.Spec.Criteria {
			criterion.HasMetricValue.MetricDefinitionId = stateMetrics[criterion.HasMetricValue.MetricName].Spec.ID
		}
	}

	// Read scorecards from split scorecard state files
	stateScorecards, errState := yaml.Parse(yaml.GetScorecardStateInput(), dtos.GetScorecardUniqueKey)
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	created, updated, deleted, unchanged := drift.Detect(
		stateScorecards,
		configScorecards,
		dtos.FromStateToConfig,
		dtos.IsScoreCardEqual,
	)

	var result []*dtos.ScorecardDTO
	h.handleDeleted(ctx, deleted)
	result = h.handleUnchanged(ctx, result, unchanged, stateScorecards, configScorecards)
	result = h.handleCreated(ctx, result, created)
	result = h.handleUpdated(ctx, result, updated, stateScorecards)

	// Write each scorecard to its own state file
	err := yaml.WriteScorecardStates(result, dtos.GetScorecardUniqueKey)
	if err != nil {
		log.Fatalf("error writing scorecards to files: %v", err)
	}
}

func (h *ApplyHandler) handleDeleted(ctx context.Context, scorecards map[string]*dtos.ScorecardDTO) {
	for _, scorecardDTO := range scorecards {
		errScorecard := h.repository.Delete(ctx, *scorecardDTO.Spec.ID)
		if errScorecard != nil {
			panic(errScorecard)
		}
	}
}

// FIXED: handleUnchanged now merges state scorecards (with IDs) with config scorecards (with updated MetricDefinitionIds)
func (h *ApplyHandler) handleUnchanged(
	ctx context.Context,
	result []*dtos.ScorecardDTO,
	unchanged map[string]*dtos.ScorecardDTO,
	stateScorecards map[string]*dtos.ScorecardDTO,
	configScorecards map[string]*dtos.ScorecardDTO,
) []*dtos.ScorecardDTO {
	for name := range unchanged {
		// Start with the state scorecard (which has the criterion IDs)
		stateScorecard := stateScorecards[name]
		configScorecard := configScorecards[name]

		// Create a copy of the state scorecard to avoid modifying the original
		mergedScorecard := &dtos.ScorecardDTO{
			APIVersion: stateScorecard.APIVersion,
			Kind:       stateScorecard.Kind,
			Metadata:   stateScorecard.Metadata,
			Spec: dtos.Spec{
				ID:                  stateScorecard.Spec.ID,
				Name:                stateScorecard.Spec.Name,
				Description:         stateScorecard.Spec.Description,
				OwnerID:             stateScorecard.Spec.OwnerID,
				State:               stateScorecard.Spec.State,
				ComponentTypeIDs:    stateScorecard.Spec.ComponentTypeIDs,
				Importance:          stateScorecard.Spec.Importance,
				ScoringStrategyType: stateScorecard.Spec.ScoringStrategyType,
				Criteria:            make([]*dtos.Criterion, len(stateScorecard.Spec.Criteria)),
			},
		}

		// Copy criteria from state (preserving IDs) but update MetricDefinitionIds from config
		for i, stateCriterion := range stateScorecard.Spec.Criteria {
			// Find matching criterion in config by name
			var configCriterion *dtos.Criterion
			for _, c := range configScorecard.Spec.Criteria {
				if c.HasMetricValue.Name == stateCriterion.HasMetricValue.Name {
					configCriterion = c
					break
				}
			}

			// Create merged criterion
			mergedScorecard.Spec.Criteria[i] = &dtos.Criterion{
				HasMetricValue: dtos.MetricValue{
					ID:                 stateCriterion.HasMetricValue.ID, // Keep ID from state
					Weight:             stateCriterion.HasMetricValue.Weight,
					Name:               stateCriterion.HasMetricValue.Name,
					MetricName:         stateCriterion.HasMetricValue.MetricName,
					MetricDefinitionId: configCriterion.HasMetricValue.MetricDefinitionId, // Updated MetricDefinitionId from config
					ComparatorValue:    stateCriterion.HasMetricValue.ComparatorValue,
					Comparator:         stateCriterion.HasMetricValue.Comparator,
				},
			}
		}

		result = append(result, mergedScorecard)
	}

	return result
}

func (h *ApplyHandler) handleCreated(ctx context.Context, result []*dtos.ScorecardDTO, scorecards map[string]*dtos.ScorecardDTO) []*dtos.ScorecardDTO {
	for _, scorecardDTO := range scorecards {
		scorecard := h.scorecardDTOToResource(scorecardDTO)

		id, criteriaMap, errScorecard := h.repository.Create(ctx, scorecard)
		if errScorecard != nil {
			panic(errScorecard)
		}

		scorecardDTO.Spec.ID = &id
		for _, criterion := range scorecardDTO.Spec.Criteria {
			criterion.HasMetricValue.ID = criteriaMap[criterion.HasMetricValue.Name]
		}
		result = append(result, scorecardDTO)
	}

	return result
}

func (h *ApplyHandler) handleUpdated(
	ctx context.Context,
	result []*dtos.ScorecardDTO,
	scorecards map[string]*dtos.ScorecardDTO,
	stateScorecards map[string]*dtos.ScorecardDTO,
) []*dtos.ScorecardDTO {
	for _, scorecardDTO := range scorecards {

		stateScorecard, ok := stateScorecards[scorecardDTO.Spec.Name]
		if !ok {
			continue
		}

		created, updated, deleted, _ := drift.Detect(
			h.mapCriteria(stateScorecard.Spec.Criteria),
			h.mapCriteria(scorecardDTO.Spec.Criteria),
			dtos.FromStateCriteriaToConfig,
			dtos.IsCriterionEqual,
		)

		deletedIDs := make([]string, len(deleted))
		i := 0
		for _, criterion := range deleted {
			deletedIDs[i] = criterion.HasMetricValue.MetricDefinitionId
			i += 1
		}

		scorecard := h.scorecardDTOToResource(scorecardDTO)
		errScorecard := h.repository.Update(
			ctx,
			scorecard,
			h.criteriaDTOToResource(created),
			h.criteriaDTOToResource(updated),
			deletedIDs,
		)
		if errScorecard != nil {
			panic(errScorecard)
		}

		result = append(result, scorecardDTO)
	}

	return result
}

func (h *ApplyHandler) mapCriteria(criteria []*dtos.Criterion) map[string]*dtos.Criterion {
	criteriaMap := make(map[string]*dtos.Criterion)
	for _, criterion := range criteria {
		criteriaMap[criterion.HasMetricValue.Name] = criterion
	}
	return criteriaMap
}

func (h *ApplyHandler) scorecardDTOToResource(scorecardDTO *dtos.ScorecardDTO) resources.Scorecard {
	return resources.Scorecard{
		ID:                  scorecardDTO.Spec.ID,
		Name:                scorecardDTO.Spec.Name,
		Description:         scorecardDTO.Spec.Description,
		OwnerID:             scorecardDTO.Spec.OwnerID,
		State:               scorecardDTO.Spec.State,
		ComponentTypeIDs:    scorecardDTO.Spec.ComponentTypeIDs,
		Importance:          scorecardDTO.Spec.Importance,
		ScoringStrategyType: scorecardDTO.Spec.ScoringStrategyType,
		Criteria:            h.criteriaDTOToResource(h.mapCriteria(scorecardDTO.Spec.Criteria)),
	}
}

func (h *ApplyHandler) criteriaDTOToResource(criteriaDTO map[string]*dtos.Criterion) []*resources.Criterion {
	criteria := make([]*resources.Criterion, len(criteriaDTO))
	i := 0
	for _, criterion := range criteriaDTO {
		criteria[i] = &resources.Criterion{
			HasMetricValue: resources.MetricValue{
				ID:                 criterion.HasMetricValue.ID,
				Weight:             criterion.HasMetricValue.Weight,
				Name:               criterion.HasMetricValue.Name,
				MetricName:         criterion.HasMetricValue.MetricName,
				MetricDefinitionId: criterion.HasMetricValue.MetricDefinitionId,
				ComparatorValue:    criterion.HasMetricValue.ComparatorValue,
				Comparator:         criterion.HasMetricValue.Comparator,
			},
		}
		i += 1
	}
	return criteria
}
