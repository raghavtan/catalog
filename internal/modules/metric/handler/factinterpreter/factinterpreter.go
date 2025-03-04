package factinterpreter

import (
	"fmt"

	"github.com/motain/fact-collector/internal/modules/metric/dtos"
	"github.com/motain/fact-collector/internal/modules/metric/handler/factcollectors"
	"github.com/motain/fact-collector/internal/utils/transformers"
)

type FactInterpreterInterface interface {
	ProcessFacts(factOperations dtos.FactOperations) (float64, error)
}

type FactInterpreter struct {
	githhubFC   factcollectors.GithubFactCollectorInterface
	jsonAPIFC   factcollectors.JSONAPIFactCollectorInterface
	componentFC factcollectors.ComponentFactCollectorInterface
}

func NewFactInterpreter(
	ghfc factcollectors.GithubFactCollectorInterface,
	jsonAPIFC factcollectors.JSONAPIFactCollectorInterface,
	componentFC factcollectors.ComponentFactCollectorInterface,
) *FactInterpreter {
	return &FactInterpreter{githhubFC: ghfc, jsonAPIFC: jsonAPIFC, componentFC: componentFC}
}

func (f *FactInterpreter) ProcessFacts(factOperations dtos.FactOperations) (float64, error) {
	if len(factOperations.All) != 0 || len(factOperations.Any) != 0 {
		return f.processConditionalOperations(factOperations)
	}

	if factOperations.Inspect != nil {
		return f.processInspectOperation(*factOperations.Inspect)
	}

	return 0, nil
}

// processConditionalOperations processes a set of conditional operations and returns a float64 result.
// It first checks if all operations in the "All" field succeed. If not, it returns 0.
// Then, it checks if any operations in the "Any" field succeed.
// The final result is a float64 representation of the logical AND between the success of all "All" operations
// and the success of any "Any" operations.
//
// Parameters:
//   - factOperations: a dtos.FactOperations struct containing the operations to be processed.
//
// Returns:
//   - float64: 1 if all "All" operations succeed and any "Any" operations succeed, otherwise 0.
//   - error: an error if any occurs during the processing of operations.
func (f *FactInterpreter) processConditionalOperations(factOperations dtos.FactOperations) (float64, error) {
	operationsAllSucceed, allErr := f.ProcessOperationsAll(factOperations.All)
	if allErr != nil {
		return 0, fmt.Errorf("process operations: %v", allErr)
	}

	if !operationsAllSucceed {
		return 0, nil
	}

	operationsAnySucceed, anyErr := f.ProcessOperationsAny(factOperations.Any)
	if anyErr != nil {
		return 0, fmt.Errorf("process operations: %v", anyErr)
	}
	return transformers.Bool2Float64(operationsAllSucceed && operationsAnySucceed), nil
}

func (f *FactInterpreter) ProcessOperationsAll(facts []*dtos.Fact) (bool, error) {
	succeed := true
	for _, fact := range facts {
		if !isSourceEnabled(dtos.FactSource(fact.Source)) {
			continue
		}

		result, err := f.check(*fact)
		if err != nil {
			return false, fmt.Errorf("process operations all: %v", err)
		}

		if !result {
			succeed = false
			break
		}
	}

	return succeed, nil
}

func (f *FactInterpreter) ProcessOperationsAny(facts []*dtos.Fact) (bool, error) {
	succeed := len(facts) == 0
	for _, fact := range facts {
		if !isSourceEnabled(dtos.FactSource(fact.Source)) {
			continue
		}

		result, err := f.check(*fact)
		if err != nil {
			return false, fmt.Errorf("process operations any: %v", err)
		}

		if result {
			succeed = true
			break
		}
	}

	return succeed, nil
}

func (f *FactInterpreter) processInspectOperation(fact dtos.Fact) (float64, error) {
	if dtos.FactSource(fact.Source) == dtos.GitHubFactSource {
		return f.githhubFC.Inspect(fact)
	}

	if dtos.FactType(fact.Source) == dtos.FactType(dtos.JSONAPIFactSource) {
		return f.jsonAPIFC.Inspect(fact)
	}

	if dtos.FactType(fact.Source) == dtos.FactType(dtos.ComponentFactSource) {
		return f.componentFC.Inspect(fact)
	}

	return 0, fmt.Errorf("error: invalid fact source %s used for inspect", fact.Source)
}

func (f *FactInterpreter) check(fact dtos.Fact) (bool, error) {
	if dtos.FactSource(fact.Source) == dtos.GitHubFactSource {
		return f.githhubFC.Check(fact)
	}

	if dtos.FactType(fact.Source) == dtos.FactType(dtos.JSONAPIFactSource) {
		return f.jsonAPIFC.Check(fact)
	}

	if dtos.FactType(fact.Source) == dtos.FactType(dtos.ComponentFactSource) {
		return f.componentFC.Check(fact)
	}

	return false, fmt.Errorf("check: invalid fact source %s used for check", fact.Source)
}

func isSourceEnabled(item dtos.FactSource) bool {
	enabledSources := []dtos.FactSource{dtos.GitHubFactSource, dtos.JSONAPIFactSource, dtos.ComponentFactSource}
	for _, v := range enabledSources {
		if v == item {
			return true
		}
	}
	return false
}
