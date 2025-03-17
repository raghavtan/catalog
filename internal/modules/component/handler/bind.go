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

func (h *BindHandler) Bind(stateRootLocation string) {
	components, errCState := yaml.Parse(stateRootLocation, false, dtos.GetComponentUniqueKey)
	if errCState != nil {
		log.Fatalf("error: %v", errCState)
	}

	metricsMap := h.getMetricsGroupedByCompoentType(stateRootLocation)

	for _, component := range components {
		for metricName, metricSource := range component.Spec.MetricSources {
			if _, exists := metricsMap[component.Metadata.ComponentType][metricName]; !exists {
				errDelete := h.repository.UnBindMetric(context.Background(), metricSource.ID)
				if errDelete != nil {
					fmt.Printf("Failed to delete metric source %s: %v\n", metricSource.Name, errDelete)
				}
			}
		}

		for metricName, metric := range metricsMap[component.Metadata.ComponentType] {
			bindErr := h.handleBind(component, metric)
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

	err := yaml.WriteState[dtos.ComponentDTO](state)
	if err != nil {
		log.Fatalf("error writing metrics to file: %v", err)
	}
}

func (*BindHandler) getMetricsGroupedByCompoentType(
	stateRootLocation string,
) map[string]map[string]*metricdtos.MetricDTO {
	metrics, errMState := yaml.Parse(stateRootLocation, false, metricdtos.GetMetricUniqueKey)
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

func (h *BindHandler) handleBind(component *dtos.ComponentDTO, metric *metricdtos.MetricDTO) error {
	fmt.Printf("Binding component %s to metric %s\n", component.Metadata.Name, metric.Metadata.Name)

	metricName := metric.Metadata.Name
	componentName := component.Metadata.Name
	identifier := utils.GetMetricSourceItentifier(metricName, componentName, component.Metadata.ComponentType)
	msFactOperations, msFactsErr := h.prepareSourceMetricFactOperations(metric.Metadata.Facts, *component)
	if msFactsErr != nil {
		fmt.Printf("Failed to prepare facts for %s/%s (component/metric): %v\n", componentName, metricName, msFactsErr)
	}

	if _, exists := component.Spec.MetricSources[metricName]; exists {
		component.Spec.MetricSources[metricName].Facts = msFactOperations
		component.Spec.MetricSources[metricName].Name = identifier
		return nil
	}

	id, errBind := h.repository.BindMetric(context.Background(), component.Spec.ID, metric.Spec.ID, identifier)
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
		Facts:  msFactOperations,
	}

	return nil
}

func (h *BindHandler) prepareSourceMetricFactOperations(
	factOperations fsdtos.FactOperations,
	component dtos.ComponentDTO,
) (fsdtos.FactOperations, error) {
	operatorAll, errAll := h.prepareSourceMetricFacts(factOperations.All, component)
	if errAll != nil {
		return fsdtos.FactOperations{}, errAll
	}

	operatorAny, errAny := h.prepareSourceMetricFacts(factOperations.Any, component)
	if errAny != nil {
		return fsdtos.FactOperations{}, errAny
	}

	return fsdtos.FactOperations{
		All:     operatorAll,
		Any:     operatorAny,
		Inspect: h.prepareSourceMetricFact(factOperations.Inspect, component),
	}, nil
}

func (h *BindHandler) prepareSourceMetricFacts(facts []*fsdtos.Fact, component dtos.ComponentDTO) ([]*fsdtos.Fact, error) {
	msFacts := make([]*fsdtos.Fact, len(facts))
	for i, fact := range facts {
		msFacts[i] = h.prepareSourceMetricFact(fact, component)
	}

	return msFacts, nil
}

func (h *BindHandler) prepareSourceMetricFact(fact *fsdtos.Fact, component dtos.ComponentDTO) *fsdtos.Fact {
	if fact == nil {
		return nil
	}

	// if fact.URI != "" {
	// 	fmt.Printf("Processing fact %s for component %s\n", fact.Name, component.Metadata.Name)
	// 	fmt.Printf("Fact URI: %s\n", fact.URI)
	// 	parsed := utils.ReplaceMetricFactPlaceholders(fact.URI, component)
	// 	fmt.Printf("Fact URI: %s\n", parsed)
	// }

	processedFact := fsdtos.Fact{
		Name:             fact.Name,
		Source:           fact.Source,
		URI:              utils.ReplaceMetricFactPlaceholders(fact.URI, component),
		ComponentName:    utils.ReplaceMetricFactPlaceholders(fact.ComponentName, component),
		Repo:             utils.ReplaceMetricFactPlaceholders(fact.Repo, component),
		FactType:         fact.FactType,
		FilePath:         fact.FilePath,
		RegexPattern:     fact.RegexPattern,
		JSONPath:         fact.JSONPath,
		RepoProperty:     fact.RepoProperty,
		ReposSearchQuery: fact.ReposSearchQuery,
		ExpectedValue:    utils.ReplaceMetricFactPlaceholders(fact.ExpectedValue, component),
		ExpectedFormula:  fact.ExpectedFormula,
		Auth:             fact.Auth,
	}

	return &processedFact
}
