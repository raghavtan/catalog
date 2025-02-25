package factcollectors

//go:generate mockgen -destination=./mock_github_fact_collector.go -package=factcollectors github.com/motain/fact-collector/internal/modules/metric/handler/factcollectors GithubFactCollectorInterface

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/motain/fact-collector/internal/modules/metric/dtos"
	"github.com/motain/fact-collector/internal/services/githubservice"
	"github.com/motain/fact-collector/internal/utils/eval"
	"github.com/motain/fact-collector/internal/utils/transformers"
	"github.com/tidwall/gjson"
)

type GithubFactCollectorInterface interface {
	Check(fact dtos.Fact) (bool, error)
	// Collect(fact dtos.Fact) (string, error)
}

type GithubFactCollector struct {
	github githubservice.GitHubRepositoriesServiceInterface
}

func NewGithubFactCollector(github githubservice.GitHubRepositoriesServiceInterface) *GithubFactCollector {
	return &GithubFactCollector{github: github}
}

func (fc *GithubFactCollector) Check(fact dtos.Fact) (bool, error) {
	if fact.FactType == dtos.FileExistsFact {
		return fc.checkFileExists(fact)
	}

	if fact.FactType == dtos.FileRegexFact {
		return fc.checkFileRegex(fact)
	}

	if fact.FactType == dtos.FileJSONPathFact {
		return fc.checkFileJSONPath(fact)
	}

	if fact.FactType == dtos.RepoPropertiesFact {
		return fc.checkRepoProperties(fact)
	}

	return false, nil
}

func (fc *GithubFactCollector) checkFileExists(fact dtos.Fact) (bool, error) {
	exists, fileErr := fc.github.GetFileExists(fact.Repo, fact.FilePath)
	if fileErr != nil {
		return false, fileErr
	}

	return exists, nil
}

func (fc *GithubFactCollector) checkFileRegex(fact dtos.Fact) (bool, error) {
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

func (fc *GithubFactCollector) checkFileJSONPath(fact dtos.Fact) (bool, error) {
	fileExtension := filepath.Ext(fact.FilePath)
	if fileExtension != ".json" && fileExtension != ".toml" {
		return false, fmt.Errorf("unsupported file extension: %s", fileExtension)
	}

	fileContent, fileErr := fc.github.GetFileContent(fact.Repo, fact.FilePath)
	if fileErr != nil {
		return false, fileErr
	}

	var jsonData []byte
	var fcransformationError error

	if fileExtension == ".toml" {
		jsonData, fcransformationError = transformers.Toml2json(fileContent)
		if fcransformationError != nil {
			return false, fcransformationError
		}
	} else {
		jsonData = []byte(fileContent)
	}

	if fileExtension == ".json" {
		jsonData = []byte(fileContent)
	}

	value := gjson.GetBytes(jsonData, fact.JSONPath)
	if !value.Exists() {
		return false, fmt.Errorf("jsonpath does not exist")
	}

	if fact.ExpectedValue == "" && fact.ExpectedFormula == "" {
		return false, errors.New("expected value or formula not provided")
	}

	if fact.ExpectedFormula != "" {
		return eval.Expression(fmt.Sprintf("%s %s", value.String(), fact.ExpectedFormula))
	}

	return value.String() == fact.ExpectedValue, nil
}

func (fc *GithubFactCollector) checkRepoProperties(fact dtos.Fact) (bool, error) {
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
