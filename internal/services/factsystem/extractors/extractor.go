package extractors

//go:generate mockgen -destination=./mocks/mock_extractor.go -package=extractors github.com/motain/of-catalog/internal/services/factsystem/extractors ExtractorInterface

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/motain/of-catalog/internal/services/configservice"
	"github.com/motain/of-catalog/internal/services/factsystem/dtos"
	"github.com/motain/of-catalog/internal/services/factsystem/utils"
	"github.com/motain/of-catalog/internal/services/githubservice"
	"github.com/motain/of-catalog/internal/services/jsonservice"
	"github.com/motain/of-catalog/internal/utils/transformers"
)

type ExtractorInterface interface {
	Extract(ctx context.Context, task *dtos.Task, deps []*dtos.Task) error
}

type Extractor struct {
	config      configservice.ConfigServiceInterface
	jsonService jsonservice.JSONServiceInterface
	github      githubservice.GitHubServiceInterface
}

func NewExtractor(
	config configservice.ConfigServiceInterface,
	jsonService jsonservice.JSONServiceInterface,
	github githubservice.GitHubServiceInterface,
) *Extractor {
	return &Extractor{config: config, jsonService: jsonService, github: github}
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

	if value, ok := deps[0].Result.(string); ok {
		return ex.handleSingleResult(ctx, task, value)
	}

	if values, ok := deps[0].Result.([]string); ok {
		return ex.handleMultipleResults(ctx, task, values)
	}

	return nil
}

func (ex *Extractor) handleSingleResult(ctx context.Context, task *dtos.Task, dependencyResult string) error {
	result, processErr := ex.processData(ctx, task, dependencyResult)
	if processErr != nil {
		return fmt.Errorf("failed to process request: %v", processErr)
	}

	task.Result = result
	return nil
}

func (ex *Extractor) handleMultipleResults(ctx context.Context, task *dtos.Task, dependencyResults []string) error {
	results := make([]interface{}, len(dependencyResults))
	for i, value := range dependencyResults {
		result, processErr := ex.processData(ctx, task, value)
		if processErr != nil {
			return fmt.Errorf("failed to process request: %v", processErr)
		}
		results[i] = result
	}

	task.Result = results
	return nil
}

func (ex *Extractor) processData(ctx context.Context, task *dtos.Task, dependencyResult string) (interface{}, error) {
	unquotedResult, _ := strconv.Unquote(dependencyResult) //nolint: errcheck

	var jsonData []byte
	var dataErr error
	switch dtos.TaskSource(task.Source) {
	case dtos.GitHubTaskSource:
		jsonData, dataErr = ex.processGithub(task, unquotedResult)
	case dtos.JSONAPITaskSource:
		jsonData, dataErr = ex.processJSONAPI(ctx, task, unquotedResult)
	default:
		return nil, fmt.Errorf("no data extracted, unknown source %s", task.Source)
	}
	if dataErr != nil {
		return nil, fmt.Errorf("failed to process request for source %s: %v", task.Source, dataErr)
	}

	return utils.InspectExtractedData(task.JSONPath, jsonData)
}

func (fe *Extractor) processGithub(task *dtos.Task, result string) ([]byte, error) {
	extractFilePath := utils.ReplacePlaceholder(task.URI, result)
	fileContent, fileErr := fe.github.GetFileContent(task.Repo, extractFilePath)
	if fileErr != nil {
		return nil, fileErr
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

	if fileExtension == ".json" {
		return []byte(fileContent), nil
	}

	return nil, fmt.Errorf("unsupported file extension: %s", fileExtension)
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

	defer resp.Body.Close()
	jsonData, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("failed to read response body: %v", readErr)
	}

	return jsonData, nil
}
