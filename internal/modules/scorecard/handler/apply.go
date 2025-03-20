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

	stateMetrics, errMetricState := yaml.Parse(yaml.GetStateInput(stateRootLocation), metricdtos.GetMetricUniqueKey)
	if errMetricState != nil {
		log.Fatalf("error: %v", errMetricState)
	}

	for _, scorecard := range configScorecards {
		for _, criterion := range scorecard.Spec.Criteria {
			criterion.HasMetricValue.MetricDefinitionId = stateMetrics[criterion.HasMetricValue.MetricName].Spec.ID
		}
	}

	stateScorecards, errState := yaml.Parse(yaml.GetStateInput(stateRootLocation), dtos.GetScorecardUniqueKey)
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	created, updated, deleted, unchanged := drift.Detect(
		stateScorecards,
		configScorecards,
		dtos.FromStateToConfig,
		dtos.IsScoreCardEqual,
	)

	result := make([]*dtos.ScorecardDTO, 0)
	h.handleDeleted(ctx, deleted)
	result = h.handleUnchanged(ctx, result, unchanged)
	result = h.handleCreated(ctx, result, created)
	result = h.handleUpdated(ctx, result, updated, stateScorecards)

	err := yaml.WriteState(result)
	if err != nil {
		log.Fatalf("error writing scorecards to file: %v", err)
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
	for _, scorecardDTO := range scorecards {
		result = append(result, scorecardDTO)
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
