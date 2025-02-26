package handler

import (
	"context"
	"log"
	"time"

	"github.com/motain/fact-collector/internal/modules/metric/dtos"
	"github.com/motain/fact-collector/internal/modules/metric/handler/factinterpreter"
	"github.com/motain/fact-collector/internal/modules/metric/repository"
	"github.com/motain/fact-collector/internal/modules/metric/utils"
	"github.com/motain/fact-collector/internal/utils/yaml"
)

type ComputeHandler struct {
	repository      repository.RepositoryInterface
	factInterpreter factinterpreter.FactInterpreterInterface
}

func NewComputeHandler(
	repository repository.RepositoryInterface,
	factInterpreter factinterpreter.FactInterpreterInterface,
) *ComputeHandler {
	return &ComputeHandler{repository: repository, factInterpreter: factInterpreter}
}

func (h *ComputeHandler) Compute(componentType, componentName, metricName string) string {

	stateMetricSource, errMSState := yaml.Parse[dtos.MetricSourceDTO](yaml.State, dtos.GetMetricSourceUniqueKey)
	if errMSState != nil {
		log.Fatalf("error: %v", errMSState)
	}

	componentIdentifier := utils.GetMetricSourceItentifier(metricName, componentName, componentType)
	var metricSource *dtos.MetricSourceDTO
	for _, metricSourceDTO := range stateMetricSource {
		if metricSourceDTO.Spec.Name == componentIdentifier {
			metricSource = metricSourceDTO
			break
		}
	}

	metricValue, processErr := h.factInterpreter.ProcessFacts(metricSource.Metadata.Facts)
	if processErr != nil {
		log.Fatalf("error: %v", processErr)
	}

	pushErr := h.repository.Push(context.Background(), *metricSource.Spec.ID, metricValue, time.Now())
	if pushErr != nil {
		log.Printf("metric source id: %s", *metricSource.Spec.ID)
		log.Fatalf("error: %v", pushErr)
	}

	return ""
}
