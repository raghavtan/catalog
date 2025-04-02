package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/motain/of-catalog/internal/modules/component/utils"
	metricdtos "github.com/motain/of-catalog/internal/modules/metric/dtos"
	fsdtos "github.com/motain/of-catalog/internal/services/factsystem/dtos"

	"github.com/motain/of-catalog/internal/modules/component/dtos"
	"github.com/motain/of-catalog/internal/modules/component/repository"
	"github.com/motain/of-catalog/internal/services/githubservice"
	"github.com/motain/of-catalog/internal/utils/yaml"
)

type BindHandler struct {
	github     githubservice.GitHubServiceInterface
	repository repository.RepositoryInterface
}

func NewBindHandler(
	gh githubservice.GitHubServiceInterface,
	repository repository.RepositoryInterface,
) *BindHandler {
	return &BindHandler{github: gh, repository: repository}
}

func (h *BindHandler) Bind(ctx context.Context, stateRootLocation string) {
	components, errCState := yaml.Parse(yaml.GetStateInput(stateRootLocation), dtos.GetComponentUniqueKey)
	if errCState != nil {
		log.Fatalf("error: %v", errCState)
	}

	metricsMap := h.getMetricsGroupedByCompoentType(stateRootLocation)

	for _, component := range components {
		for metricName, metricSource := range component.Spec.MetricSources {
			if _, exists := metricsMap[component.Metadata.ComponentType][metricName]; !exists {
				errDelete := h.repository.UnbindMetric(ctx, MetricSourceDTOToResource(metricSource))
				if errDelete != nil {
					fmt.Printf("Failed to delete metric source %s: %v\n", metricSource.Name, errDelete)
				}
			}
		}

		for metricName, metric := range metricsMap[component.Metadata.ComponentType] {
			bindErr := h.handleBind(ctx, component, metric)
			if bindErr != nil {
				fmt.Printf("Failed to bind metric %s to component %s: %v\n", metricName, component.Metadata.Name, bindErr)
			}
		}
	}

	state := make([]*dtos.ComponentDTO, len(components))
	i := 0
	for _, component := range components {
		state[i] = component
		i += 1
	}

	err := yaml.WriteState(state)
	if err != nil {
		log.Fatalf("error writing metrics to file: %v", err)
	}
}

func (*BindHandler) getMetricsGroupedByCompoentType(
	stateRootLocation string,
) map[string]map[string]*metricdtos.MetricDTO {
	metrics, errMState := yaml.Parse(yaml.GetStateInput(stateRootLocation), metricdtos.GetMetricUniqueKey)
	if errMState != nil {
		log.Fatalf("error: %v", errMState)
	}

	metricsMap := make(map[string]map[string]*metricdtos.MetricDTO)
	for _, metric := range metrics {
		for _, componentType := range metric.Metadata.ComponentType {
			if _, exists := metricsMap[componentType]; !exists {
				metricsMap[componentType] = make(map[string]*metricdtos.MetricDTO)
			}
			metricsMap[componentType][metric.Metadata.Name] = metric
		}
	}

	return metricsMap
}

func (h *BindHandler) handleBind(ctx context.Context, component *dtos.ComponentDTO, metric *metricdtos.MetricDTO) error {
	fmt.Printf("Binding component %s to metric %s\n", component.Metadata.Name, metric.Metadata.Name)

	metricName := metric.Metadata.Name
	componentName := component.Metadata.Name
	identifier := utils.GetMetricSourceIdentifier(metricName, componentName, component.Metadata.ComponentType)
	tasks := h.prepareSourceMetricFacts(metric.Metadata.Facts, *component)

	if _, exists := component.Spec.MetricSources[metricName]; exists {
		component.Spec.MetricSources[metricName].Facts = tasks
		component.Spec.MetricSources[metricName].Name = identifier
		return nil
	}

	id, errBind := h.repository.BindMetric(ctx, componentDTOToResource(component), metric.Spec.ID, identifier)
	if errBind != nil {
		return fmt.Errorf("failed to create metric source for %s/%s (component/metric): %v", componentName, metricName, errBind)
	}

	if component.Spec.MetricSources == nil {
		component.Spec.MetricSources = make(map[string]*dtos.MetricSourceDTO)
	}

	component.Spec.MetricSources[metricName] = &dtos.MetricSourceDTO{
		ID:     id,
		Name:   identifier,
		Metric: metric.Spec.ID,
		Facts:  tasks,
	}

	return nil
}

func (h *BindHandler) prepareSourceMetricFacts(tasks []*fsdtos.Task, component dtos.ComponentDTO) []*fsdtos.Task {
	processedFacts := make([]*fsdtos.Task, len(tasks))
	for i, task := range tasks {
		processedFacts[i] = h.prepareSourceMetricFact(task, component)
	}
	return processedFacts
}

func (h *BindHandler) prepareSourceMetricFact(task *fsdtos.Task, component dtos.ComponentDTO) *fsdtos.Task {
	if task == nil {
		return nil
	}

	// if fact.URI != "" {
	// 	fmt.Printf("Processing fact %s for component %s\n", fact.Name, component.Metadata.Name)
	// 	fmt.Printf("Fact URI: %s\n", fact.URI)
	// 	parsed := utils.ReplaceMetricFactPlaceholders(fact.URI, component)
	// 	fmt.Printf("Fact URI: %s\n", parsed)
	// }

	processedFact := fsdtos.Task{
		ID:     task.ID,
		Name:   task.Name,
		Source: task.Source,
		URI:    utils.ReplaceMetricFactPlaceholders(task.URI, component),
		Auth:   task.Auth,
		// ComponentName: utils.ReplaceMetricFactPlaceholders(task.ComponentName, component),
		Repo:            utils.ReplaceMetricFactPlaceholders(task.Repo, component),
		Type:            task.Type,
		FilePath:        task.FilePath,
		JSONPath:        task.JSONPath,
		Rule:            task.Rule,
		Pattern:         utils.ReplaceMetricFactPlaceholders(task.Pattern, component),
		DependsOn:       task.DependsOn,
		Method:          task.Method,
		SearchString:    task.SearchString,
		PrometheusQuery: utils.ReplaceMetricFactPlaceholders(task.PrometheusQuery, component),

		// Are these still worth it?
		// RegexPattern:     task.RegexPattern,
		// RepoProperty:     task.RepoProperty,
		// ReposSearchQuery: task.ReposSearchQuery,

		// I need to reintegrate this !!
		// ExpectedFormula:  task.ExpectedFormula,
	}

	return &processedFact
}
