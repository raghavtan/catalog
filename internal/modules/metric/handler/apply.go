package handler

import (
	"context"
	"log"
	"reflect"

	"github.com/motain/fact-collector/internal/modules/metric/dtos"
	"github.com/motain/fact-collector/internal/modules/metric/repository"
	"github.com/motain/fact-collector/internal/modules/metric/resources"
	"github.com/motain/fact-collector/internal/services/githubservice"
	"github.com/motain/fact-collector/internal/utils/drift"
	"github.com/motain/fact-collector/internal/utils/yaml"
)

type ApplyHandler struct {
	github     githubservice.GitHubRepositoriesServiceInterface
	repository repository.RepositoryInterface
}

func NewApplyHandler(
	gh githubservice.GitHubRepositoriesServiceInterface,
	repository repository.RepositoryInterface,
) *ApplyHandler {
	return &ApplyHandler{github: gh, repository: repository}
}

func (h *ApplyHandler) Apply() string {
	configMetrics, errConfig := yaml.ParseConfig[dtos.MetricDTO]()
	if errConfig != nil {
		log.Fatalf("error: %v", errConfig)
	}

	stateMetrics, errState := yaml.ParseState[dtos.MetricDTO]()
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	getUniqueKey := func(m *dtos.MetricDTO) string {
		return m.Spec.Name
	}
	setID := func(m *dtos.MetricDTO, id string) {
		m.Spec.ID = &id
	}
	getID := func(m *dtos.MetricDTO) string {
		return *m.Spec.ID
	}
	isEqual := func(m1, m2 *dtos.MetricDTO) bool {
		return m1.Spec.Name == m2.Spec.Name && m1.Spec.Description == m2.Spec.Description && reflect.DeepEqual(m1.Spec.Format, m2.Spec.Format)
	}

	newMetrics, updatedMetrics, removedMetrics, unchangedMetrics := drift.Detect(stateMetrics, configMetrics, getUniqueKey, getID, setID, isEqual)
	for _, metricDTO := range removedMetrics {
		errMetric := h.repository.Delete(context.Background(), *metricDTO.Spec.ID)
		if errMetric != nil {
			panic(errMetric)
		}
	}

	var result = unchangedMetrics
	for _, metricDTO := range newMetrics {
		metric := resources.Metric{
			Name:        metricDTO.Spec.Name,
			Description: metricDTO.Spec.Description,
			Format: resources.MetricFormat{
				Unit: metricDTO.Spec.Format.Unit,
			},
		}

		id, errMetric := h.repository.Create(context.Background(), metric)
		if errMetric != nil {
			panic(errMetric)
		}

		metricDTO.Spec.ID = &id
		result = append(result, metricDTO)
	}

	for _, metricDTO := range updatedMetrics {
		metric := resources.Metric{
			ID:          metricDTO.Spec.ID,
			Name:        metricDTO.Spec.Name,
			Description: metricDTO.Spec.Description,
			Format: resources.MetricFormat{
				Unit: metricDTO.Spec.Format.Unit,
			},
		}

		errMetric := h.repository.Update(context.Background(), metric)
		if errMetric != nil {
			panic(errMetric)
		}

		result = append(result, metricDTO)
	}

	err := yaml.WriteState[dtos.MetricDTO](result)
	if err != nil {
		log.Fatalf("error writing metrics to file: %v", err)
	}

	return ""
}
