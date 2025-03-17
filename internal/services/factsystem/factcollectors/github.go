package factcollectors

//go:generate mockgen -destination=./mock/mock_github_fact_collector.go -package=factcollectors github.com/motain/of-catalog/internal/services/factsystem/factcollectors GithubFactCollectorInterface

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	fsdtos "github.com/motain/of-catalog/internal/services/factsystem/dtos"
	"github.com/motain/of-catalog/internal/services/githubservice"
	"github.com/motain/of-catalog/internal/utils/eval"
	"github.com/motain/of-catalog/internal/utils/transformers"
)

type GithubFactCollectorInterface interface {
	Check(fact fsdtos.Fact) (bool, error)
	Inspect(fact fsdtos.Fact) (float64, error)
}

type GithubFactCollector struct {
	github githubservice.GitHubServiceInterface
}

func NewGithubFactCollector(github githubservice.GitHubServiceInterface) *GithubFactCollector {
	return &GithubFactCollector{github: github}
}

func (fc *GithubFactCollector) Check(fact fsdtos.Fact) (bool, error) {
	if fact.FactType == fsdtos.FileExistsFact {
		return fc.checkFileExists(fact)
	}

	if fact.FactType == fsdtos.FileRegexFact {
		return fc.checkFileRegex(fact)
	}

	if fact.FactType == fsdtos.JSONPathFact {
		return fc.checkFileJSONPath(fact)
	}

	if fact.FactType == fsdtos.RepoPropertiesFact {
		return fc.checkRepoProperties(fact)
	}

	if fact.FactType == fsdtos.RepoSearchFact {
		return fc.checkRepoSearch(fact)
	}

	return false, errors.New("unsupported fact type")
}

func (fc *GithubFactCollector) Inspect(fact fsdtos.Fact) (float64, error) {
	jsonData, extractionErr := fc.extractData(fact)
	if extractionErr != nil {
		return 0, extractionErr
	}

	value, inspectErr := inspectJson(jsonData, fact)
	if inspectErr != nil {
		return 0, inspectErr
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return floatValue, nil
}

func (fc *GithubFactCollector) checkFileExists(fact fsdtos.Fact) (bool, error) {
	exists, fileErr := fc.github.GetFileExists(fact.Repo, fact.FilePath)
	if fileErr != nil {
		return false, fileErr
	}

	return exists, nil
}

func (fc *GithubFactCollector) checkFileRegex(fact fsdtos.Fact) (bool, error) {
	fileContent, fileErr := fc.github.GetFileContent(fact.Repo, fact.FilePath)
	if fileErr != nil {
		return false, fileErr
	}

	regexPattern, regexErr := regexp.Compile(fact.RegexPattern)
	if regexErr != nil {
		return false, regexErr
	}
	matched := regexPattern.MatchString(fileContent)
	if !matched {
		return false, nil
	}

	return true, nil
}

func (fc *GithubFactCollector) checkFileJSONPath(fact fsdtos.Fact) (bool, error) {
	jsonData, extractionErr := fc.extractData(fact)
	if extractionErr != nil {
		return false, extractionErr
	}

	value, inspectErr := inspectJson(jsonData, fact)
	if inspectErr != nil {
		return false, inspectErr
	}

	if fact.ExpectedValue == "" && fact.ExpectedFormula == "" {
		return false, errors.New("expected value or formula not provided")
	}

	if fact.ExpectedFormula != "" {
		return eval.Expression(fmt.Sprintf("%s %s", value, fact.ExpectedFormula))
	}

	regexPattern, regexErr := regexp.Compile(fact.ExpectedValue)
	if regexErr != nil {
		return false, regexErr
	}

	return regexPattern.MatchString(value), nil
}

func (fc *GithubFactCollector) checkRepoProperties(fact fsdtos.Fact) (bool, error) {
	repoProperties, repoErr := fc.github.GetRepoProperties(fact.Repo)
	if repoErr != nil {
		return false, repoErr
	}

	value, ok := repoProperties[fact.RepoProperty]
	if !ok {
		return false, fmt.Errorf("repo property does not exist")
	}

	if fact.ExpectedValue != "" {
		return value == fact.ExpectedValue, nil
	}

	return eval.Expression(fmt.Sprintf("%s %s", value, fact.ExpectedFormula))
}

func (fc *GithubFactCollector) checkRepoSearch(fact fsdtos.Fact) (bool, error) {
	repoSearchResults, repoErr := fc.github.Search(fact.Repo, fact.ReposSearchQuery)
	if repoErr != nil {
		return false, repoErr
	}

	return len(repoSearchResults) != 0, nil
}

func (fc *GithubFactCollector) extractData(fact fsdtos.Fact) ([]byte, error) {
	var jsonData []byte

	fileExtension := filepath.Ext(fact.FilePath)
	if fileExtension != ".json" && fileExtension != ".toml" {
		return jsonData, fmt.Errorf("unsupported file extension: %s", fileExtension)
	}

	fileContent, fileErr := fc.github.GetFileContent(fact.Repo, fact.FilePath)
	if fileErr != nil {
		return jsonData, fileErr
	}

	if fileExtension == ".toml" {
		return transformers.Toml2json(fileContent)
	}

	if fileExtension == ".json" {
		return []byte(fileContent), nil
	}

	return jsonData, fmt.Errorf("jsonpath %s does not exist in %s", fact.JSONPath, fact.FilePath)
}
