package factcollectors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"

	componentdtos "github.com/motain/fact-collector/internal/modules/component/dtos"
	"github.com/motain/fact-collector/internal/modules/metric/dtos"
	"github.com/motain/fact-collector/internal/utils/eval"
	"github.com/motain/fact-collector/internal/utils/yaml"
)

type ComponentFactCollectorInterface interface {
	Check(fact dtos.Fact) (bool, error)
	Inspect(fact dtos.Fact) (float64, error)
}

type ComponentFactCollector struct{}

func NewComponentFactCollector() *ComponentFactCollector {
	return &ComponentFactCollector{}
}

func (fc *ComponentFactCollector) Check(fact dtos.Fact) (bool, error) {
	if fact.FactType != dtos.JSONPathFact {
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

func (fc *ComponentFactCollector) Inspect(fact dtos.Fact) (float64, error) {
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

func (fc *ComponentFactCollector) extractData(fact dtos.Fact) ([]byte, error) {
	stateComponents, errState := yaml.Parse[componentdtos.ComponentDTO](yaml.StateLocation, false, componentdtos.GetComponentUniqueKey)
	if errState != nil {
		log.Fatalf("error: %v", errState)
	}

	component, found := stateComponents[fact.ComponentName]
	if !found {
		return nil, errors.New("component not found")
	}
	jsonData, marshalErr := json.Marshal(component)
	if marshalErr != nil {
		return nil, marshalErr
	}

	return jsonData, nil
}
