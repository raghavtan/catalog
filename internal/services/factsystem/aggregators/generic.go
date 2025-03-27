package aggregators

import (
	"context"
	"errors"
	"fmt"

	"github.com/motain/of-catalog/internal/services/factsystem/dtos"
	"github.com/motain/of-catalog/internal/services/factsystem/utils"
)

type AggregatorInterface interface {
	Combine(ctx context.Context, task *dtos.Task, deps []*dtos.Task) error
}

type Aggregator struct{}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (ag *Aggregator) Combine(ctx context.Context, task *dtos.Task, deps []*dtos.Task) error {
	if dtos.TaskType(task.Type) != dtos.AggregateType {
		return nil
	}

	partials := make([]interface{}, len(deps))
	for i, dep := range deps {
		if dep.Result == nil {
			return errors.New("dependency result not provided")
		}
		combinedDepResult, depErr := ag.combineResults(task, dep.Result)
		if depErr != nil {
			return depErr
		}

		partials[i] = combinedDepResult
	}

	combinedPartialsResult, depErr := ag.combineResults(task, partials)
	if depErr != nil {
		return depErr
	}

	task.Result = combinedPartialsResult
	return nil
}

func (ag *Aggregator) combineResults(task *dtos.Task, results interface{}) (interface{}, error) {
	switch dtos.TaskMethod(task.Method) {
	case dtos.CountMethod:
		return ag.count(results)
	case dtos.SumMethod:
		return ag.sum(results)
	case dtos.AndMethod:
		return ag.and(results)
	case dtos.OrMethod:
		return ag.or(results)
	default:
		return nil, errors.New("unknown method")
	}
}

func (ag *Aggregator) count(results interface{}) (float64, error) {
	switch v := results.(type) {
	case []interface{}:
		return float64(len(v)), nil
	case []string:
		return float64(len(v)), nil
	case []bool:
		return float64(len(v)), nil
	case []int:
		return float64(len(v)), nil
	default:
		return 0, fmt.Errorf("unsupported type for count: %T", results)
	}
}

func (ag *Aggregator) sum(results interface{}) (float64, error) {
	values, castErr := utils.ToSlice[float64](results)
	if castErr != nil {
		return 0.0, fmt.Errorf("combineResult error for method \"sum\": %s", castErr)
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum, nil
}

func (ag *Aggregator) and(results interface{}) (bool, error) {
	if value, ok := results.(bool); ok {
		return value, nil
	}

	values, castErr := utils.ToSlice[bool](results)
	if castErr != nil {
		return false, fmt.Errorf("combineResult error for method \"and\": %s", castErr)
	}

	for _, v := range values {
		if !v {
			return false, nil
		}
	}
	return true, nil
}

func (ag *Aggregator) or(results interface{}) (bool, error) {
	if value, ok := results.(bool); ok {
		return value, nil
	}

	values, castErr := utils.ToSlice[bool](results)
	if castErr != nil {
		return false, fmt.Errorf("combineResult error for method \"or\": %s", castErr)
	}
	for _, v := range values {
		if v {
			return true, nil
		}
	}
	return false, nil
}
