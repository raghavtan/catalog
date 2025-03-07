package handler

import (
	"context"
	"fmt"
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

func (h *ComputeHandler) Compute(componentName string, all bool, metricName string, stateRootLocation string) {
	components, errCState := yaml.Parse(stateRootLocation, false, dtos.GetComponentUniqueKey)
	if errCState != nil {
		log.Fatalf("error: %v", errCState)
	}

	component, componentExists := components[componentName]
	if !componentExists {
		log.Fatalf("compute: error: component not found for name %s", componentName)
	}

	if !all {
		fmt.Printf("Tracking metric '%s' for component '%s'\n", metricName, componentName)
		computeErr := h.computeMetric(component, metricName)
		if computeErr != nil {
			log.Fatalf("compute: %v", computeErr)
		}
		return
	}

	for metricName := range component.Spec.MetricSources {
		fmt.Printf("Tracking metric '%s' for component '%s'\n", metricName, componentName)
		computeErr := h.computeMetric(component, metricName)
		if computeErr != nil {
			log.Printf("compute metric %s: %v", metricName, computeErr)
		}
	}
}

func (h *ComputeHandler) computeMetric(component *dtos.ComponentDTO, metricName string) error {
	metricSource, msExists := component.Spec.MetricSources[metricName]
	if !msExists {
		return fmt.Errorf("error: metric source not found for metric %s", metricName)
	}

	metricValue, processErr := h.factInterpreter.ProcessFacts(metricSource.Facts)
	if processErr != nil {
		return fmt.Errorf("%v", processErr)
	}

	pushErr := h.repository.Push(context.Background(), metricSource.ID, metricValue, time.Now())
	if pushErr != nil {
		return fmt.Errorf("error: %v", pushErr)
	}

	return nil
}
