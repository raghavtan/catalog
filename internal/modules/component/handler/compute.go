package handler

import (
	"context"
	"log"
	"time"

	"github.com/motain/fact-collector/internal/modules/component/dtos"
	"github.com/motain/fact-collector/internal/modules/component/repository"
	"github.com/motain/fact-collector/internal/services/factsystem/factinterpreter"
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

func (h *ComputeHandler) Compute(componentName, metricName string, stateRootLocation string) {
	components, errCState := yaml.Parse(stateRootLocation, false, dtos.GetComponentUniqueKey)
	if errCState != nil {
		log.Fatalf("error: %v", errCState)
	}

	component, componentExists := components[componentName]
	if !componentExists {
		log.Fatalf("compute: error: component not found for name %s", componentName)
	}

	metricSource, msExists := component.Spec.MetricSources[metricName]
	if !msExists {
		log.Fatalf("compute: error: metric source not found for metric %s", metricName)
	}

	metricValue, processErr := h.factInterpreter.ProcessFacts(metricSource.Facts)
	if processErr != nil {
		log.Fatalf("compute: %v", processErr)
	}

	pushErr := h.repository.Push(context.Background(), metricSource.ID, metricValue, time.Now())
	if pushErr != nil {
		log.Printf("metric source id: %s", metricSource.ID)
		log.Fatalf("compute: error: %v", pushErr)
	}
}
