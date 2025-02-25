package factinterpreter

import (
	"log"

	"github.com/motain/fact-collector/internal/modules/metric/dtos"
	"github.com/motain/fact-collector/internal/modules/metric/handler/factcollectors"
	"github.com/motain/fact-collector/internal/utils/transformers"
)

type FactInterpreterInterface interface {
	ProcessFacts(factOperations dtos.FactOperations) float64
}

type FactInterpreter struct {
	githhubFC factcollectors.GithubFactCollectorInterface
}

func NewFactInterpreter(
	ghfc factcollectors.GithubFactCollectorInterface,
) *FactInterpreter {
	return &FactInterpreter{githhubFC: ghfc}
}

func (f *FactInterpreter) ProcessFacts(factOperations dtos.FactOperations) float64 {
	if factOperations.All != nil || factOperations.Any != nil {
		return f.processConditionalOperations(factOperations)
	}

	return 0
}

// processConditionalOperations evaluates a set of conditional operations on facts.
// It verifies if all facts in the "All" list and any fact in the "Any" list have a source of "github"
// and satisfy certain conditions checked by the githhubFC.Check method.
// If both All and Any conditions are provided, the results of All and Any are combined using a logical AND operator.
//
// Parameters:
//   - factOperations: a dtos.FactOperations object containing lists of facts to be checked.
//
// Returns:
//   - float64: returns 1 if all conditions are met, otherwise returns 0.
func (f *FactInterpreter) processConditionalOperations(factOperations dtos.FactOperations) float64 {
	operationsAllSucceed := f.ProcessOperationsAll(factOperations.All)
	if !operationsAllSucceed {
		return 0
	}

	operationsAnySucceed := f.ProcessOperationsAny(factOperations.Any)
	return transformers.Bool2Float64(operationsAllSucceed && operationsAnySucceed)
}

func (f *FactInterpreter) ProcessOperationsAll(facts []dtos.Fact) bool {
	succeed := true
	for _, fact := range facts {
		if fact.Source != "github" {
			continue
		}

		result, err := f.githhubFC.Check(fact)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		if !result {
			succeed = false
			break
		}
	}

	return succeed
}

func (f *FactInterpreter) ProcessOperationsAny(facts []dtos.Fact) bool {
	succeed := facts == nil
	for _, fact := range facts {
		if fact.Source != "github" {
			continue
		}

		result, err := f.githhubFC.Check(fact)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		if result {
			succeed = true
			break
		}
	}

	return succeed
}
