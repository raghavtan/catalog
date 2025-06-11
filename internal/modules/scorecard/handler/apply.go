package handler

import (
	"context"
	"fmt"
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

	// DEBUG: Log scorecard config
	fmt.Printf("DEBUG: Found %d config scorecards\n", len(configScorecards))

	// Read metrics from split metric state files
	stateMetrics, errMetricState := yaml.Parse(yaml.GetMetricStateInput(), metricdtos.GetMetricUniqueKey)
	if errMetricState != nil {
		log.Fatalf("error: %v", errMetricState)
	}

	// DEBUG: Log metrics found in state
	fmt.Printf("DEBUG: Found %d state metrics\n", len(stateMetrics))
	for name, metric := range stateMetrics {
		fmt.Printf("DEBUG: State metric '%s' has ID: '%s'\n", name, metric.Spec.ID)
	}

	// DEBUG: Check metric assignment
	for scorecardName, scorecard := range configScorecards {
		fmt.Printf("DEBUG: Processing scorecard: %s\n", scorecardName)
		for i, criterion := range scorecard.Spec.Criteria {
			metricName := criterion.HasMetricValue.MetricName
			fmt.Printf("DEBUG: Criterion %d needs metric: '%s'\n", i, metricName)

			if stateMetric, exists := stateMetrics[metricName]; exists {
				fmt.Printf("DEBUG: Found metric '%s' with ID: '%s'\n", metricName, stateMetric.Spec.ID)
				criterion.HasMetricValue.MetricDefinitionId = stateMetric.Spec.ID
				fmt.Printf("DEBUG: Assigned MetricDefinitionId: '%s'\n", criterion.HasMetricValue.MetricDefinitionId)
			} else {
				fmt.Printf("DEBUG: ERROR - Metric '%s' not found in state metrics!\n", metricName)
				fmt.Printf("DEBUG: Available metrics: ")
				for availableName := range stateMetrics {
					fmt.Printf("'%s' ", availableName)
				}
				fmt.Println()
			}
		}
	}

	// Read scorecards from split scorecard state files
	stateScorecards, errState := yaml.Parse(yaml.GetScorecardStateInput(), dtos.GetScorecardUniqueKey)
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	// DEBUG: Log scorecard state
	fmt.Printf("DEBUG: Found %d state scorecards\n", len(stateScorecards))

	created, updated, deleted, unchanged := drift.Detect(
		stateScorecards,
		configScorecards,
		dtos.FromStateToConfig,
		dtos.IsScoreCardEqual,
	)

	// DEBUG: Log drift detection
	fmt.Printf("DEBUG: Drift detection - Created: %d, Updated: %d, Deleted: %d, Unchanged: %d\n",
		len(created), len(updated), len(deleted), len(unchanged))

	result := make([]*dtos.ScorecardDTO, 0)
	h.handleDeleted(ctx, deleted)
	result = h.handleUnchanged(ctx, result, unchanged)
	result = h.handleCreated(ctx, result, created)
	result = h.handleUpdated(ctx, result, updated, stateScorecards)

	// DEBUG: Check final result
	fmt.Printf("DEBUG: Final result has %d scorecards\n", len(result))
	for _, scorecard := range result {
		fmt.Printf("DEBUG: Scorecard '%s' has %d criteria\n", scorecard.Spec.Name, len(scorecard.Spec.Criteria))
		for i, criterion := range scorecard.Spec.Criteria {
			fmt.Printf("DEBUG: Criterion %d: ID='%s', MetricDefinitionId='%s'\n",
				i, criterion.HasMetricValue.ID, criterion.HasMetricValue.MetricDefinitionId)
		}
	}

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

func (h *ApplyHandler) handleUnchanged(ctx context.Context, result []*dtos.ScorecardDTO, scorecards map[string]*dtos.ScorecardDTO) []*dtos.ScorecardDTO {
	fmt.Printf("DEBUG: handleUnchanged processing %d scorecards\n", len(scorecards))
	for name, scorecardDTO := range scorecards {
		fmt.Printf("DEBUG: Adding unchanged scorecard: %s\n", name)
		result = append(result, scorecardDTO)
	}
	return result
}

func (h *ApplyHandler) handleCreated(ctx context.Context, result []*dtos.ScorecardDTO, scorecards map[string]*dtos.ScorecardDTO) []*dtos.ScorecardDTO {
	fmt.Printf("DEBUG: handleCreated processing %d scorecards\n", len(scorecards))
	for name, scorecardDTO := range scorecards {
		fmt.Printf("DEBUG: Creating scorecard: %s\n", name)
		scorecard := h.scorecardDTOToResource(scorecardDTO)

		id, criteriaMap, errScorecard := h.repository.Create(ctx, scorecard)
		if errScorecard != nil {
			panic(errScorecard)
		}

		fmt.Printf("DEBUG: Created scorecard with ID: %s\n", id)
		fmt.Printf("DEBUG: Criteria map has %d entries\n", len(criteriaMap))
		for criteriaName, criteriaID := range criteriaMap {
			fmt.Printf("DEBUG: Criteria '%s' -> ID '%s'\n", criteriaName, criteriaID)
		}

		scorecardDTO.Spec.ID = &id
		for i, criterion := range scorecardDTO.Spec.Criteria {
			criteriaID := criteriaMap[criterion.HasMetricValue.Name]
			fmt.Printf("DEBUG: Setting criterion %d ID from '%s' to '%s'\n",
				i, criterion.HasMetricValue.ID, criteriaID)
			criterion.HasMetricValue.ID = criteriaID
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
	fmt.Printf("DEBUG: handleUpdated processing %d scorecards\n", len(scorecards))
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
