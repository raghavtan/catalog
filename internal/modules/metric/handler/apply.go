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

func (h *ApplyHandler) Apply(configRootLocation string, stateRootLocation string, recursive bool) {
	stateMetrics, errState := yaml.Parse(stateRootLocation, false, dtos.GetMetricUniqueKey)
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	configMetrics, errConfig := yaml.Parse[dtos.MetricDTO](configRootLocation, recursive, dtos.GetMetricUniqueKey)
	if errConfig != nil {
		log.Fatalf("error: %v", errConfig)
	}

	created, updated, deleted, unchanged := drift.Detect(
		stateMetrics,
		configMetrics,
		dtos.FromStateToConfig,
		dtos.IsEqualMetric,
	)
	h.handleDeleted(deleted)

	var result []*dtos.MetricDTO
	result = h.handleUnchanged(result, unchanged)
	result = h.handleCreated(result, created)
	result = h.handleUpdated(result, updated)

	err := yaml.WriteState(result)
	if err != nil {
		log.Fatalf("error writing metrics to file: %v", err)
	}
}

func (h *ApplyHandler) handleDeleted(metrics map[string]*dtos.MetricDTO) {
	for _, metricDTO := range metrics {
		err := h.repository.Delete(context.Background(), metricDTO.Spec.ID)
		if err != nil {
			panic(err)
		}
	}
}

func (h *ApplyHandler) handleUnchanged(result []*dtos.MetricDTO, metrics map[string]*dtos.MetricDTO) []*dtos.MetricDTO {
	for _, metricDTO := range metrics {
		result = append(result, metricDTO)
	}
	return result
}

func (h *ApplyHandler) handleCreated(result []*dtos.MetricDTO, metrics map[string]*dtos.MetricDTO) []*dtos.MetricDTO {
	for _, metricDTO := range metrics {
		metric := metricDTOToResource(metricDTO)

		id, err := h.repository.Create(context.Background(), metric)
		if err != nil {
			panic(err)
		}

		metricDTO.Spec.ID = id
		result = append(result, metricDTO)
	}

	return result
}

func (h *ApplyHandler) handleUpdated(result []*dtos.MetricDTO, metrics map[string]*dtos.MetricDTO) []*dtos.MetricDTO {
	for _, metricDTO := range metrics {
		metric := metricDTOToResource(metricDTO)
		err := h.repository.Update(context.Background(), metric)
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
