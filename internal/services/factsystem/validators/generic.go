package validators

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/motain/of-catalog/internal/services/factsystem/dtos"
	"github.com/motain/of-catalog/internal/utils/eval"
)

type ValidatorInterface interface {
	Check(task *dtos.Task, deps []*dtos.Task) error
}

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (fc *Validator) Check(task *dtos.Task, deps []*dtos.Task) error {
	if dtos.ValidateType != dtos.TaskType(task.Type) {
		return nil
	}

	if deps == nil {
		return errors.New("too few dependencies provided for validate task")
	}

	if len(deps) > 1 {
		return fc.validateDepsRelations(task, deps)
	}

	dep := deps[0]
	if dep.Result == nil {
		// Should I fail or ignore and set task.Result = false ?
		return errors.New("dependency result not provided")
	}

	if values, isStringSlice := dep.Result.([]interface{}); isStringSlice {
		err := fc.validateList(task, values)
		if err != nil {
			return err
		}

		return nil
	}

	isValid, err := fc.validate(task, dep.Result)
	if err != nil {
		return err
	}

	task.Result = isValid

	return nil
}

func (fc *Validator) validateDepsRelations(task *dtos.Task, deps []*dtos.Task) error {
	switch dtos.TaskRule(task.Rule) {
	case dtos.DepsMatchRule:
		return fc.validateDependeciesMatch(task, deps)
	default:
		return nil
	}
}

func (fc *Validator) validateList(task *dtos.Task, list []interface{}) error {
	switch dtos.TaskRule(task.Rule) {
	case dtos.UniqueRule:
		return fc.validateUnique(task, list)
	default:
		return fc.validateEach(task, list)
	}
}

func (fc *Validator) validateUnique(task *dtos.Task, list []interface{}) error {
	uniqueMap := make(map[interface{}]bool)
	for _, v := range list {
		if _, ok := uniqueMap[v]; ok {
			task.Result = false
			return nil
		}
		uniqueMap[v] = true
	}

	task.Result = true
	return nil
}

func (fc *Validator) validateEach(task *dtos.Task, list []interface{}) error {
	res := make([]bool, len(list))
	for i, v := range list {
		ok, err := fc.validate(task, v)
		if err != nil {
			return err
		}
		res[i] = ok
	}

	task.Result = res
	return nil
}

func (fc *Validator) validate(task *dtos.Task, value interface{}) (bool, error) {
	strValue := fmt.Sprintf("%v", value)
	switch dtos.TaskRule(task.Rule) {
	case dtos.RegexMatchRule:
		return fc.validateRegex(task, strValue)
	case dtos.FormulaRule:
		return fc.validateFormula(task, strValue)
	default:
		return false, errors.New("unknown validation rule")
	}
}

func (fc *Validator) validateRegex(task *dtos.Task, value string) (bool, error) {
	regexPattern, regexErr := regexp.Compile(task.Pattern)
	if regexErr != nil {
		return false, regexErr
	}

	return regexPattern.MatchString(value), nil
}

func (fc *Validator) validateFormula(task *dtos.Task, value string) (bool, error) {
	return eval.Expression(fmt.Sprintf("%s %s", value, task.Pattern))
}

func (fc *Validator) validateDependeciesMatch(task *dtos.Task, deps []*dtos.Task) error {
	mappedDepResultsTypes := make(map[string][]*dtos.Task)
	for _, dep := range deps {
		resType := fmt.Sprintf("%T", dep.Result)
		mappedDepResultsTypes[resType] = append(mappedDepResultsTypes[resType], dep)
	}

	if len(mappedDepResultsTypes) > 1 {
		task.Result = false
		return nil
	}

	if _, isSLice := deps[0].Result.([]interface{}); isSLice {
		// Implement slice comparison
		// mappedDepResultsLength := make(map[int][]*dtos.Task)
		// for _, dep := range deps {
		// 	length := len(dep.Result.([]interface{}))
		// 	mappedDepResultsLength[length] = append(mappedDepResultsLength[length], dep)
		// }
		// if len(mappedDepResultsLength) > 1 {
		// 	task.Result = false
		// 	return nil
		// }
		return errors.New("slice comparison not implemented")
	}

	resultsMatch := true
	for _, dep := range deps {
		resultsMatch = resultsMatch && dep.Result != deps[0].Result
	}

	task.Result = resultsMatch
	return nil
}
