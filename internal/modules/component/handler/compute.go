package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/motain/of-catalog/internal/modules/component/dtos"
	"github.com/motain/of-catalog/internal/modules/component/repository"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	"github.com/motain/of-catalog/internal/services/factsystem/processor"
	"github.com/motain/of-catalog/internal/utils/yaml"
)

type ComputeHandler struct {
	repository    repository.RepositoryInterface
	factProcessor processor.ProcessorInterface
}

func NewComputeHandler(
	repository repository.RepositoryInterface,
	factProcessor processor.ProcessorInterface,
) *ComputeHandler {
	return &ComputeHandler{repository: repository, factProcessor: factProcessor}
}

func (h *ComputeHandler) Compute(ctx context.Context, componentName string, all bool, metricName string, stateRootLocation string) {
	components, errCState := yaml.Parse(yaml.GetStateInput(stateRootLocation), dtos.GetComponentUniqueKey)
	if errCState != nil {
		log.Fatalf("error: %v", errCState)
	}

	component, componentExists := components[componentName]
	if !componentExists {
		log.Fatalf("compute: error: component not found for name %s", componentName)
	}

	if !all {
		fmt.Printf("Tracking metric '%s' for component '%s'\n", metricName, componentName)
		computeErr := h.computeMetric(ctx, component, metricName)
		if computeErr != nil {
			log.Fatalf("compute: %v", computeErr)
		}
		return
	}

	for metricName := range component.Spec.MetricSources {
		fmt.Printf("Tracking metric '%s' for component '%s'\n", metricName, componentName)
		computeErr := h.computeMetric(ctx, component, metricName)
		if computeErr != nil {
			log.Printf("compute metric %s: %v", metricName, computeErr)
		}
	}
}

func (h *ComputeHandler) computeMetric(ctx context.Context, component *dtos.ComponentDTO, metricName string) error {
	metricSource, msExists := component.Spec.MetricSources[metricName]
	if !msExists {
		return fmt.Errorf("error: metric source not found for metric %s", metricName)
	}

	metricValue, processErr := h.factProcessor.Process(ctx, metricSource.Facts)
	if processErr != nil {
		return fmt.Errorf("%v", processErr)
	}

	pushErr := h.repository.Push(ctx, MetricSourceDTOToResource(metricSource), metricValue, time.Now())
	if pushErr != nil {
		return fmt.Errorf("error: %v", pushErr)
	}

	return nil
}

func MetricSourceDTOToResource(metricSource *dtos.MetricSourceDTO) resources.MetricSource {
	return resources.MetricSource{
		ID:     metricSource.ID,
		Name:   metricSource.Name,
		Metric: metricSource.Metric,
		Facts:  metricSource.Facts,
	}
}
