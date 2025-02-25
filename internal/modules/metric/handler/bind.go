package handler

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

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
		msFacts, msFactsErr := h.prepareSourceMetricFactOperations(metric.Metadata.Facts, *component)
		if msFactsErr != nil {
			fmt.Printf("Failed to prepare facts for metric source %s: %v\n", identifier, msFactsErr)
		}
		metricSourceMap[identifier].Metadata.Facts = msFacts
		return append(result, metricSourceMap[identifier])
	}

	id, errBind := h.repository.CreateMetricSource(context.Background(), *metric.Spec.ID, *component.Spec.ID, identifier)
	if errBind != nil {
		panic(errBind)
	}

	msFacts, msFactsErr := h.prepareSourceMetricFactOperations(metric.Metadata.Facts, *component)
	if msFactsErr != nil {
		fmt.Printf("Failed to prepare facts for metric source %s: %v\n", identifier, msFactsErr)
	}

	metricSourceDTO := dtos.MetricSourceDTO{
		APIVersion: "v1",
		Kind:       "MetricSource",
		Metadata: dtos.MetricSourceMetadataDTO{
			Name:          identifier,
			ComponentType: []string{component.Metadata.ComponentType},
			Status:        "active",
			Facts:         msFacts,
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

// getFieldByPath fetches a nested field value using dot notation
func getFieldByPath(obj interface{}, path string) (interface{}, error) {
	fields := strings.Split(path, ".")
	val := reflect.ValueOf(obj)

	// Traverse fields
	for _, field := range fields {
		// Dereference pointer if necessary
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		// Ensure it's a struct
		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("invalid path: %s", path)
		}

		// Get field by name
		val = val.FieldByName(field)

		// If field is invalid, return error
		if !val.IsValid() {
			return nil, fmt.Errorf("field not found: %s", field)
		}
	}

	return val.Interface(), nil
}

func (h *BindHandler) replaceMetricFactPlaceholder(placeholder string, component componentdtos.ComponentDTO) (string, error) {
	re := regexp.MustCompile(`\$\{(.*)\}`)
	matches := re.FindStringSubmatch(placeholder)
	if len(matches) != 2 {
		return placeholder, nil
	}
	capturedGroup := matches[1]

	value, err := getFieldByPath(component, capturedGroup)
	if err != nil {
		return "", fmt.Errorf("error fetching '%s': %v", capturedGroup, err)
	}

	return fmt.Sprintf("%v", value), nil
}

func (h *BindHandler) prepareSourceMetricFactOperations(factOperations dtos.FactOperations, component componentdtos.ComponentDTO) (dtos.FactOperations, error) {
	operatorAll, errAll := h.prepareSourceMetricFacts(factOperations.All, component)
	if errAll != nil {
		return dtos.FactOperations{}, errAll
	}

	operatorAny, errAny := h.prepareSourceMetricFacts(factOperations.Any, component)
	if errAny != nil {
		return dtos.FactOperations{}, errAny
	}

	operatorReport, errReport := h.prepareSourceMetricFacts(factOperations.Report, component)
	if errReport != nil {
		return dtos.FactOperations{}, errReport
	}

	return dtos.FactOperations{
		All:    operatorAll,
		Any:    operatorAny,
		Report: operatorReport,
	}, nil
}

func (h *BindHandler) prepareSourceMetricFacts(facts []dtos.Fact, component componentdtos.ComponentDTO) ([]dtos.Fact, error) {
	msFacts := make([]dtos.Fact, len(facts))
	for i, fact := range facts {
		repo, errRepo := h.replaceMetricFactPlaceholder(fact.Repo, component)
		if errRepo != nil {
			fmt.Printf("Failed to replace placeholder for role in fact %s: %v\n", fact.Name, errRepo)
			return nil, errRepo
		}
		expectedValue, errExpValue := h.replaceMetricFactPlaceholder(fact.ExpectedValue, component)
		if errExpValue != nil {
			fmt.Printf("Failed to replace placeholder for expectedValue in fact %s: %v\n", fact.Name, errExpValue)
			return nil, errExpValue
		}

		msFacts[i] = dtos.Fact{
			Name:            fact.Name,
			Source:          fact.Source,
			URI:             fact.URI,
			Repo:            repo,
			FactType:        fact.FactType,
			FilePath:        fact.FilePath,
			RegexPattern:    fact.RegexPattern,
			JSONPath:        fact.JSONPath,
			RepoProperty:    fact.RepoProperty,
			ExpectedValue:   expectedValue,
			ExpectedFormula: fact.ExpectedFormula,
		}
	}

	return msFacts, nil
}
