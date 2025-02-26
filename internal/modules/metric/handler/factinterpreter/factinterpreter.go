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
	githhubFC factcollectors.GithubFactCollectorInterface
}

func NewFactInterpreter(
	ghfc factcollectors.GithubFactCollectorInterface,
) *FactInterpreter {
	return &FactInterpreter{githhubFC: ghfc}
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
		return 0, fmt.Errorf("error: %v", allErr)
	}

	if !operationsAllSucceed {
		return 0, nil
	}

	operationsAnySucceed, anyErr := f.ProcessOperationsAny(factOperations.Any)
	if anyErr != nil {
		return 0, fmt.Errorf("error: %v", anyErr)
	}
	return transformers.Bool2Float64(operationsAllSucceed && operationsAnySucceed), nil
}

func (f *FactInterpreter) ProcessOperationsAll(facts []*dtos.Fact) (bool, error) {
	succeed := true
	for _, fact := range facts {
		if fact.Source != "github" {
			continue
		}

		result, err := f.githhubFC.Check(*fact)
		if err != nil {
			return false, fmt.Errorf("error: %v", err)
		}

		if !result {
			succeed = false
			break
		}
	}

	return succeed, nil
}

func (f *FactInterpreter) ProcessOperationsAny(facts []*dtos.Fact) (bool, error) {
	succeed := facts == nil
	for _, fact := range facts {
		if fact.Source != "github" {
			continue
		}

		result, err := f.githhubFC.Check(*fact)
		if err != nil {
			return false, fmt.Errorf("error: %v", err)
		}

		if result {
			succeed = true
			break
		}
	}

	return succeed, nil
}

func (f *FactInterpreter) processInspectOperation(fact dtos.Fact) (float64, error) {
	return f.githhubFC.Inspect(fact)
}
