package handler

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/motain/of-catalog/internal/modules/metric/dtos"
	"github.com/motain/of-catalog/internal/modules/metric/repository"
	"github.com/motain/of-catalog/internal/modules/metric/resources"
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
	fmt.Println("DEBUG: Checking .state/metric/ directory:")
	entries, err := os.ReadDir(".state/metric")
	if err != nil {
		fmt.Printf("DEBUG: Error reading .state/metric directory: %v\n", err)
	} else {
		fmt.Printf("DEBUG: Found %d files in .state/metric/\n", len(entries))
		for _, entry := range entries {
			fmt.Printf("  - %s\n", entry.Name())
		}
	}

	// DEBUG: Try to read split state with verbose logging
	fmt.Println("DEBUG: Attempting to read split metric state...")
	stateInput := yaml.GetMetricStateInput()
	fmt.Printf("DEBUG: State input: RootLocation=%s, Recursive=%t\n", stateInput.RootLocation, stateInput.Recursive)

	stateMetrics, errState := yaml.Parse(stateInput, dtos.GetMetricUniqueKey)
	if errState != nil {
		fmt.Printf("DEBUG: Error reading split state: %v\n", errState)
		log.Fatalf("error: %v", errState)
	}

	parseInput := yaml.ParseInput{
		RootLocation: configRootLocation,
		Recursive:    recursive,
	}
	configMetrics, errConfig := yaml.Parse(parseInput, dtos.GetMetricUniqueKey)
	if errConfig != nil {
		log.Fatalf("error: %v", errConfig)
	}

	fmt.Printf("DEBUG: Config metrics count: %d\n", len(configMetrics))
	for name := range configMetrics {
		fmt.Printf("DEBUG: Config metric: %s\n", name)
	}

	// If split state is empty but we have split files, there's a parsing issue
	if len(stateMetrics) == 0 && len(entries) > 0 {
		if len(entries) > 0 {
			firstFile := fmt.Sprintf(".state/metric/%s", entries[0].Name())
			content, readErr := os.ReadFile(firstFile)
			if readErr != nil {
				fmt.Printf("DEBUG: Error reading %s: %v\n", firstFile, readErr)
			} else {
				fmt.Printf("DEBUG: Content of %s:\n%s\n", firstFile, string(content))
			}
		}
	}

	created, updated, deleted, unchanged := drift.Detect(
		stateMetrics,
		configMetrics,
		dtos.FromStateToConfig,
		dtos.IsEqualMetric,
	)

	// DEBUG: Print drift detection results
	fmt.Printf("DEBUG: Drift detection results:\n")
	fmt.Printf("  Created: %d\n", len(created))
	fmt.Printf("  Updated: %d\n", len(updated))
	fmt.Printf("  Deleted: %d\n", len(deleted))
	fmt.Printf("  Unchanged: %d\n", len(unchanged))

	for name := range created {
		fmt.Printf("  - Created: %s\n", name)
	}

	var result []*dtos.MetricDTO
	//h.handleDeleted(ctx, deleted)
	result = h.handleUnchanged(ctx, result, unchanged)
	result = h.handleCreated(ctx, result, created)
	result = h.handleUpdated(ctx, result, updated)

	err = yaml.WriteMetricStates(result, dtos.GetMetricUniqueKey)
	if err != nil {
		log.Fatalf("error writing metrics to file: %v", err)
	}
}

func (h *ApplyHandler) handleDeleted(ctx context.Context, metrics map[string]*dtos.MetricDTO) {
	for _, metricDTO := range metrics {
		err := h.repository.Delete(ctx, metricDTO.Spec.ID)
		if err != nil {
			panic(err)
		}
	}
}

func (h *ApplyHandler) handleUnchanged(ctx context.Context, result []*dtos.MetricDTO, metrics map[string]*dtos.MetricDTO) []*dtos.MetricDTO {
	for _, metricDTO := range metrics {
		result = append(result, metricDTO)
	}
	return result
}

func (h *ApplyHandler) handleCreated(ctx context.Context, result []*dtos.MetricDTO, metrics map[string]*dtos.MetricDTO) []*dtos.MetricDTO {
	for name, metricDTO := range metrics {
		fmt.Printf("DEBUG: Creating metric: %s\n", name)

		metric := metricDTOToResource(metricDTO)

		id, err := h.repository.Create(ctx, metric)
		if err != nil {
			fmt.Printf("DEBUG: Failed to create metric %s: %v\n", name, err)
			panic(err)
		}

		metricDTO.Spec.ID = id
		result = append(result, metricDTO)
	}

	return result
}

func (h *ApplyHandler) handleUpdated(ctx context.Context, result []*dtos.MetricDTO, metrics map[string]*dtos.MetricDTO) []*dtos.MetricDTO {
	for _, metricDTO := range metrics {
		metric := metricDTOToResource(metricDTO)
		err := h.repository.Update(ctx, metric)
		if err != nil {
			panic(err)
		}

		result = append(result, metricDTO)
	}

	return result
}

func metricDTOToResource(metricDTO *dtos.MetricDTO) resources.Metric {
	return resources.Metric{
		ID:          metricDTO.Spec.ID,
		Name:        metricDTO.Spec.Name,
		Description: metricDTO.Spec.Description,
		Format: resources.MetricFormat{
			Unit: metricDTO.Spec.Format.Unit,
		},
	}
}
