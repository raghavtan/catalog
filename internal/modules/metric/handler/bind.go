package handler

import (
	"context"
	"fmt"
	"log"

	componentdtos "github.com/motain/fact-collector/internal/modules/component/dtos"

	"github.com/motain/fact-collector/internal/modules/metric/dtos"
	"github.com/motain/fact-collector/internal/modules/metric/repository"
	"github.com/motain/fact-collector/internal/services/githubservice"
	"github.com/motain/fact-collector/internal/utils/yaml"
)

type BindHandler struct {
	github     githubservice.GitHubRepositoriesServiceInterface
	repository repository.RepositoryInterface
}

func NewBindHandler(
	gh githubservice.GitHubRepositoriesServiceInterface,
	repository repository.RepositoryInterface,
) *BindHandler {
	return &BindHandler{github: gh, repository: repository}
}

func (h *BindHandler) Bind() string {
	stateMetrics, errState := yaml.ParseState[dtos.MetricDTO]()
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	stateComponents, errState := yaml.ParseState[componentdtos.ComponentDTO]()
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	metricSourceMap := h.getStateMetricSourceHashedByName()
	componentMap := h.getStateComponentsGroupedByType(stateComponents)

	result := make([]*dtos.MetricSourceDTO, 0)
	for _, metric := range stateMetrics {
		for _, componentType := range metric.Metadata.ComponentType {
			components, exists := componentMap[componentType]
			if !exists {
				continue
			}

			for _, component := range components {
				result = h.handleBind(result, metric, component, metricSourceMap)
			}
		}

		result = h.resolveDrifts(metricSourceMap, result)

		err := yaml.WriteState[dtos.MetricSourceDTO](result)
		if err != nil {
			log.Fatalf("error writing metrics to file: %v", err)
		}

	}

	return ""
}

func (*BindHandler) getStateComponentsGroupedByType(stateComponents []*componentdtos.ComponentDTO) map[string][]*componentdtos.ComponentDTO {
	componentMap := make(map[string][]*componentdtos.ComponentDTO)
	for _, component := range stateComponents {
		componentType := component.Metadata.ComponentType
		componentMap[componentType] = append(componentMap[componentType], component)
	}
	return componentMap
}

func (*BindHandler) getStateMetricSourceHashedByName() map[string]*dtos.MetricSourceDTO {
	stateMetricSource, errState := yaml.ParseState[dtos.MetricSourceDTO]()
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	metricSourceMap := make(map[string]*dtos.MetricSourceDTO)
	for _, metricSource := range stateMetricSource {
		if metricSource.Metadata.Status != "inactive" {
			metricSourceMap[metricSource.Spec.Name] = metricSource
		}
	}

	return metricSourceMap
}

func (h *BindHandler) handleBind(
	result []*dtos.MetricSourceDTO,
	metric *dtos.MetricDTO,
	component *componentdtos.ComponentDTO,
	metricSourceMap map[string]*dtos.MetricSourceDTO,
) []*dtos.MetricSourceDTO {
	fmt.Printf("Binding metric %s to component %s\n", metric.Spec.Name, component.Metadata.Name)

	identifier := fmt.Sprintf("%s-%s", metric.Spec.Name, component.Spec.Slug)
	if _, exists := metricSourceMap[identifier]; exists {
		return append(result, metricSourceMap[identifier])
	}

	id, errBind := h.repository.CreateMetricSource(context.Background(), *metric.Spec.ID, *component.Spec.ID, identifier)
	if errBind != nil {
		panic(errBind)
	}

	metricSourceDTO := dtos.MetricSourceDTO{
		APIVersion: "v1",
		Kind:       "MetricSource",
		Metadata: dtos.MetricSourceMetadataDTO{
			Name:          identifier,
			ComponentType: []string{component.Metadata.ComponentType},
			Status:        "active",
		},
		Spec: dtos.MetricSourceSpecDTO{
			ID:        &id,
			Name:      identifier,
			Metric:    *metric.Spec.ID,
			Component: *component.Spec.ID,
		},
	}

	result = append(result, &metricSourceDTO)

	return result
}

func (h *BindHandler) resolveDrifts(preBind map[string]*dtos.MetricSourceDTO, postBind []*dtos.MetricSourceDTO) []*dtos.MetricSourceDTO {
	postBindMap := make(map[string]*dtos.MetricSourceDTO)
	for _, metricSource := range postBind {
		postBindMap[metricSource.Spec.Name] = metricSource
	}

	for _, metricSource := range preBind {
		if _, exists := postBindMap[metricSource.Spec.Name]; !exists {
			errDelete := h.repository.DeleteMetricSource(context.Background(), *metricSource.Spec.ID)
			if errDelete != nil {
				fmt.Printf("Failed to delete metric source %s: %v\n", metricSource.Spec.Name, errDelete)
				continue
			}

			metricSource.Metadata.Status = "inactive"
			postBind = append(postBind, metricSource)
		}
	}

	return postBind
}
