package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	componentutils "github.com/motain/fact-collector/internal/modules/component/utils"
	"github.com/motain/fact-collector/internal/modules/metric/dtos"
	"github.com/motain/fact-collector/internal/modules/metric/handler/factinterpreter"
	"github.com/motain/fact-collector/internal/modules/metric/repository"
	"github.com/motain/fact-collector/internal/utils/yaml"
)

type TrackHandler struct {
	repository      repository.RepositoryInterface
	factInterpreter factinterpreter.FactInterpreterInterface
}

func NewTrackHandler(
	repository repository.RepositoryInterface,
	factInterpreter factinterpreter.FactInterpreterInterface,
) *TrackHandler {
	return &TrackHandler{repository: repository, factInterpreter: factInterpreter}
}

func (h *TrackHandler) Track(componentType, componentName, metricName string) string {

	stateMetricSource, errMSState := yaml.Parse[dtos.MetricSourceDTO](yaml.State, dtos.GetMetricSourceUniqueKey)
	if errMSState != nil {
		log.Fatalf("error: %v", errMSState)
	}

	componentSlug := componentutils.GetSlug(componentName, componentType)
	metricSourceName := fmt.Sprintf("%s-%s", metricName, componentSlug)

	var metricSource *dtos.MetricSourceDTO
	for _, metricSourceDTO := range stateMetricSource {
		if metricSourceDTO.Spec.Name == metricSourceName {
			metricSource = metricSourceDTO
			break
		}
	}

	metricValue := h.factInterpreter.ProcessFacts(metricSource.Metadata.Facts)
	h.repository.Push(context.Background(), *metricSource.Spec.ID, metricValue, time.Now())

	return ""
}
