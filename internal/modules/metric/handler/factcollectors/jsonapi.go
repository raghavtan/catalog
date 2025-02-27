package factcollectors

//go:generate mockgen -destination=./mock_jsonapiFact_fact_collector.go -package=factcollectors github.com/motain/fact-collector/internal/modules/metric/handler/factcollectors JSONAPIFactCollectorInterface

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/motain/fact-collector/internal/modules/metric/dtos"
	"github.com/motain/fact-collector/internal/services/jsonservice"
	"github.com/motain/fact-collector/internal/utils/eval"
	"github.com/tidwall/gjson"
)

type JSONAPIFactCollectorInterface interface {
	Check(fact dtos.Fact) (bool, error)
	Inspect(fact dtos.Fact) (float64, error)
}

type JSONAPIFactCollector struct {
	jsonService jsonservice.JSONServiceInterface
}

func NewJSONAPIFactCollector(jsonService jsonservice.JSONServiceInterface) *JSONAPIFactCollector {
	return &JSONAPIFactCollector{jsonService: jsonService}
}

func (fc *JSONAPIFactCollector) Check(fact dtos.Fact) (bool, error) {
	if fact.FactType != dtos.FileJSONPathFact {
		return false, nil
	}

	jsonData, extractionErr := fc.extractData(fact)
	if extractionErr != nil {
		return false, extractionErr
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

func (fc *JSONAPIFactCollector) Inspect(fact dtos.Fact) (float64, error) {
	if fact.FactType != dtos.FileJSONPathFact {
		return 0, nil
	}

	jsonData, extractionErr := fc.extractData(fact)
	if extractionErr != nil {
		return 0, extractionErr
	}

	value := gjson.GetBytes(jsonData, fact.JSONPath)
	if !value.Exists() {
		return 0, fmt.Errorf("jsonpath does not exist")
	}

	return value.Float(), nil
}

func (fc *JSONAPIFactCollector) extractData(fact dtos.Fact) ([]byte, error) {
	var jsonData []byte

	req, err := http.NewRequest("GET", fact.URI, nil)
	if err != nil {
		return jsonData, fmt.Errorf("failed to create request: %v", err)
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
