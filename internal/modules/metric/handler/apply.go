package handler

import (
	"context"
	"log"

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
	stateMetrics, errState := yaml.Parse(yaml.GetStateInput(stateRootLocation), dtos.GetMetricUniqueKey)
	if errState != nil {
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
	// jsonData, err := json.MarshalIndent(configMetrics, "", "  ")
	// if err != nil {
	// 	log.Fatalf("error converting configMetrics to JSON: %v", err)
	// }
	// fmt.Println(string(jsonData))
	// os.Exit(1)
	created, updated, deleted, unchanged := drift.Detect(
		stateMetrics,
		configMetrics,
		dtos.FromStateToConfig,
		dtos.IsEqualMetric,
	)

	var result []*dtos.MetricDTO
	h.handleDeleted(ctx, deleted)
	result = h.handleUnchanged(ctx, result, unchanged)
	result = h.handleCreated(ctx, result, created)
	result = h.handleUpdated(ctx, result, updated)

	err := yaml.WriteState(result)
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
	for _, metricDTO := range metrics {
		metric := metricDTOToResource(metricDTO)

		id, err := h.repository.Create(ctx, metric)
		if err != nil {
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
