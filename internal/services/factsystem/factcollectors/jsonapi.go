package factcollectors

//go:generate mockgen -destination=./mock_jsonapiFact_fact_collector.go -package=factcollectors github.com/motain/of-catalog/internal/modules/metric/handler/factcollectors JSONAPIFactCollectorInterface

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/motain/of-catalog/internal/services/configservice"
	fsdtos "github.com/motain/of-catalog/internal/services/factsystem/dtos"
	"github.com/motain/of-catalog/internal/services/jsonservice"
	"github.com/motain/of-catalog/internal/utils/eval"
)

type JSONAPIFactCollectorInterface interface {
	Check(fact fsdtos.Fact) (bool, error)
	Inspect(fact fsdtos.Fact) (float64, error)
}

type JSONAPIFactCollector struct {
	config      configservice.ConfigServiceInterface
	jsonService jsonservice.JSONServiceInterface
}

func NewJSONAPIFactCollector(
	config configservice.ConfigServiceInterface,
	jsonService jsonservice.JSONServiceInterface,
) *JSONAPIFactCollector {
	return &JSONAPIFactCollector{config: config, jsonService: jsonService}
}

func (fc *JSONAPIFactCollector) Check(fact fsdtos.Fact) (bool, error) {
	if fact.FactType != fsdtos.JSONPathFact {
		return false, nil
	}

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

func (fc *JSONAPIFactCollector) Inspect(fact fsdtos.Fact) (float64, error) {
	if fact.FactType != fsdtos.JSONPathFact {
		return 0, nil
	}

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

func (fc *JSONAPIFactCollector) extractData(fact fsdtos.Fact) ([]byte, error) {
	var jsonData []byte

	req, err := http.NewRequest("GET", fact.URI, nil)
	if err != nil {
		return jsonData, fmt.Errorf("failed to create request: %v", err)
	}

	if fact.Auth != nil {
		token := fc.config.Get(fact.Auth.TokenEnvVariable)
		req.Header.Set(fact.Auth.Header, token)
	}

	resp, doErr := fc.jsonService.Do(req)
	if doErr != nil {
		return jsonData, fmt.Errorf("failed to do request: %v", doErr)
	}

	defer resp.Body.Close()
	jsonData, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return jsonData, fmt.Errorf("failed to read response body: %v", readErr)
	}

	return jsonData, nil
}
