package extractors

//go:generate mockgen -destination=./mocks/mock_extractor.go -package=extractors github.com/motain/of-catalog/internal/services/factsystem/extractors ExtractorInterface

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/motain/of-catalog/internal/services/configservice"
	"github.com/motain/of-catalog/internal/services/factsystem/dtos"
	"github.com/motain/of-catalog/internal/services/factsystem/utils"
	"github.com/motain/of-catalog/internal/services/githubservice"
	"github.com/motain/of-catalog/internal/services/jsonservice"
	"github.com/motain/of-catalog/internal/services/prometheusservice"
	"github.com/motain/of-catalog/internal/utils/transformers"
)

type ExtractorInterface interface {
	Extract(ctx context.Context, task *dtos.Task, deps []*dtos.Task) error
}

type Extractor struct {
	config            configservice.ConfigServiceInterface
	jsonService       jsonservice.JSONServiceInterface
	github            githubservice.GitHubServiceInterface
	prometheusService prometheusservice.PrometheusServiceInterface
}

func NewExtractor(
	config configservice.ConfigServiceInterface,
	jsonService jsonservice.JSONServiceInterface,
	github githubservice.GitHubServiceInterface,
	prometheusService prometheusservice.PrometheusServiceInterface,
) *Extractor {
	return &Extractor{config: config, jsonService: jsonService, github: github, prometheusService: prometheusService}
}

func (ex *Extractor) Extract(ctx context.Context, task *dtos.Task, deps []*dtos.Task) error {
	if len(deps) > 1 {
		return errors.New("too many dependencies provided in extract context")
	}

	if len(deps) == 0 || deps == nil {
		return ex.handleSingleResult(ctx, task, "")
	}

	if deps[0].Result == nil {
		return errors.New("dependency result not provided")
	}

	if values, ok := deps[0].Result.([]interface{}); ok {
		if len(values) == 0 {
			return errors.New("dependency result not provided")
		}

		stringValues := make([]string, len(values))
		for i, value := range values {
			stringValues[i] = fmt.Sprintf("%v", value)
		}
		return ex.handleMultipleResults(ctx, task, stringValues)
	}

	if values, ok := deps[0].Result.([]string); ok {
		if len(values) == 0 {
			return errors.New("dependency result not provided")
		}

		return ex.handleMultipleResults(ctx, task, values)
	}

	return ex.handleSingleResult(ctx, task, fmt.Sprintf("%v", deps[0].Result))
}

func (ex *Extractor) handleSingleResult(ctx context.Context, task *dtos.Task, dependencyResult string) error {
	result, processErr := ex.processData(ctx, task, dependencyResult)
	if processErr != nil {
		return fmt.Errorf("single result handler failed to process request: %v", processErr)
	}
	task.Result = result
	return nil
}

func (ex *Extractor) handleMultipleResults(ctx context.Context, task *dtos.Task, dependencyResults []string) error {
	results := make([]string, 0)
	for _, value := range dependencyResults {
		result, processErr := ex.processData(ctx, task, value)
		if processErr != nil {
			return fmt.Errorf("multiple results handler failed to process request: %v", processErr)
		}

		if values, ok := result.([]interface{}); ok {
			for _, v := range values {
				results = append(results, fmt.Sprintf("%v", v))
			}
			continue
		}

		if values, ok := result.([]string); ok {
			results = append(results, values...)
			continue
		}

		if value, ok := result.(string); ok {
			results = append(results, value)
			continue
		}
	}

	task.Result = results
	return nil
}

func (ex *Extractor) processData(ctx context.Context, task *dtos.Task, dependencyResult string) (interface{}, error) {
	var jsonData []byte
	var dataErr error
	switch dtos.TaskSource(task.Source) {
	case dtos.GitHubTaskSource:
		if task.Rule == string(dtos.SearchRule) {
			searchListResult, searchErr := ex.github.Search(task.Repo, task.SearchString)
			if searchErr != nil {
				return nil, fmt.Errorf("failed to process github Search request for source for string %s %s: %v", task.SearchString, task.Source, searchErr)
			}
			return len(searchListResult) != 0, nil
		}
		jsonData, dataErr = ex.processGithub(task, unquoted(dependencyResult))
	case dtos.JSONAPITaskSource:
		jsonData, dataErr = ex.processJSONAPI(ctx, task, unquoted(dependencyResult))
	case dtos.PrometheusTaskSource:
		jsonData, dataErr = ex.queryPrometheus(task, unquoted(dependencyResult))
	default:
		return nil, fmt.Errorf("no data extracted, unknown source %s", task.Source)
	}
	if dataErr != nil {
		return nil, fmt.Errorf("failed to process request for source %s: %v", task.Source, dataErr)
	}

	switch dtos.TaskRule(task.Rule) {
	case dtos.JSONPathRule:
		return utils.InspectExtractedData(task.JSONPath, jsonData)
	case dtos.NotEmptyRule:
		return jsonData != nil, nil
	default:
		return jsonData, nil
	}
}

func (ex *Extractor) processGithub(task *dtos.Task, result string) ([]byte, error) {
	extractFilePath := utils.ReplacePlaceholder(task.FilePath, result)
	fileContent, fileErr := ex.github.GetFileContent(task.Repo, extractFilePath)
	if fileErr != nil {
		re := regexp.MustCompile(`404 Not Found`)
		if re.MatchString(fileErr.Error()) {
			return nil, nil
		}
		return nil, fileErr
	}

	if dtos.TaskRule(task.Rule) != dtos.JSONPathRule {
		return []byte(fileContent), nil
	}

	fileExtension := filepath.Ext(task.FilePath)
	if fileExtension != ".json" && fileExtension != ".toml" {
		return nil, fmt.Errorf("unsupported file extension: %s", fileExtension)
	}

	if fileExtension == ".toml" {
		jsonData, transformErr := transformers.Toml2json(fileContent)
		if transformErr != nil {
			return nil, fmt.Errorf("failed to transform toml file to json: %v", transformErr)
		}
		return jsonData, nil
	}

	return []byte(fileContent), nil
}

func (ex *Extractor) processJSONAPI(ctx context.Context, task *dtos.Task, result string) ([]byte, error) {
	extractURI := utils.ReplacePlaceholder(task.URI, result)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, extractURI, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	if task.Auth != nil {
		token := ex.config.Get(task.Auth.TokenVar)
		req.Header.Set(task.Auth.Header, token)
	}

	resp, fileErr := ex.jsonService.Do(req)
	if fileErr != nil {
		return nil, fileErr
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)
	jsonData, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read response body: %v", readErr)
	}

	return jsonData, nil
}

func unquoted(toUnquote string) string {
	unquoted, unQuoteErr := strconv.Unquote(toUnquote) //nolint: errcheck
	if unQuoteErr != nil {
		unquoted = toUnquote
	}

	return unquoted
}

func (ex *Extractor) queryPrometheus(task *dtos.Task, result string) ([]byte, error) {
	prometheusQuery := utils.ReplacePlaceholder(task.PrometheusQuery, result)
	response, err := ex.prometheusService.InstantQuery(prometheusQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to query prometheus: %v", err)
	}

	return json.Marshal(response)
}
