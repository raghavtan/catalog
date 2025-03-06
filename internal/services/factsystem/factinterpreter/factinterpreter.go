package factinterpreter

import (
	"fmt"

	fsdtos "github.com/motain/fact-collector/internal/services/factsystem/dtos"
	"github.com/motain/fact-collector/internal/services/factsystem/factcollectors"
	"github.com/motain/fact-collector/internal/utils/transformers"
)

type FactInterpreterInterface interface {
	ProcessFacts(factOperations fsdtos.FactOperations) (float64, error)
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

func (f *FactInterpreter) ProcessFacts(factOperations fsdtos.FactOperations) (float64, error) {
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
//   - factOperations: a fsdtos.FactOperations struct containing the operations to be processed.
//
// Returns:
//   - float64: 1 if all "All" operations succeed and any "Any" operations succeed, otherwise 0.
//   - error: an error if any occurs during the processing of operations.
func (f *FactInterpreter) processConditionalOperations(factOperations fsdtos.FactOperations) (float64, error) {
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

func (f *FactInterpreter) ProcessOperationsAll(facts []*fsdtos.Fact) (bool, error) {
	succeed := true
	for _, fact := range facts {
		if !isSourceEnabled(fsdtos.FactSource(fact.Source)) {
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

func (f *FactInterpreter) ProcessOperationsAny(facts []*fsdtos.Fact) (bool, error) {
	succeed := len(facts) == 0
	for _, fact := range facts {
		if !isSourceEnabled(fsdtos.FactSource(fact.Source)) {
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

func (f *FactInterpreter) processInspectOperation(fact fsdtos.Fact) (float64, error) {
	if fsdtos.FactSource(fact.Source) == fsdtos.GitHubFactSource {
		return f.githhubFC.Inspect(fact)
	}

	if fsdtos.FactType(fact.Source) == fsdtos.FactType(fsdtos.JSONAPIFactSource) {
		return f.jsonAPIFC.Inspect(fact)
	}

	if fsdtos.FactType(fact.Source) == fsdtos.FactType(fsdtos.ComponentFactSource) {
		return f.componentFC.Inspect(fact)
	}

	return 0, fmt.Errorf("error: invalid fact source %s used for inspect", fact.Source)
}

func (f *FactInterpreter) check(fact fsdtos.Fact) (bool, error) {
	if fsdtos.FactSource(fact.Source) == fsdtos.GitHubFactSource {
		return f.githhubFC.Check(fact)
	}

	if fsdtos.FactType(fact.Source) == fsdtos.FactType(fsdtos.JSONAPIFactSource) {
		return f.jsonAPIFC.Check(fact)
	}

	if fsdtos.FactType(fact.Source) == fsdtos.FactType(fsdtos.ComponentFactSource) {
		return f.componentFC.Check(fact)
	}

	return false, fmt.Errorf("check: invalid fact source %s used for check", fact.Source)
}

func isSourceEnabled(item fsdtos.FactSource) bool {
	enabledSources := []fsdtos.FactSource{fsdtos.GitHubFactSource, fsdtos.JSONAPIFactSource, fsdtos.ComponentFactSource}
	for _, v := range enabledSources {
		if v == item {
			return true
		}
	}
	return false
}
