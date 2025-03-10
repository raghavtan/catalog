package factcollectors

import (
	"fmt"

	fsdtos "github.com/motain/of-catalog/internal/services/factsystem/dtos"
	"github.com/spyzhov/ajson"
)

func inspectJson(jsonData []byte, fact fsdtos.Fact) (string, error) {
	root, unmarshalErr := ajson.Unmarshal(jsonData)
	if unmarshalErr != nil {
		return "", unmarshalErr
	}

	node, evalErr := ajson.Eval(root, fact.JSONPath)
	if evalErr != nil {
		return "", fmt.Errorf("error parsing jsonpath: %v", evalErr)
	}

	nodeValue, vslueErr := node.Value()
	if vslueErr != nil {
		return "", fmt.Errorf("error parsing jsonpath value: %v", vslueErr)
	}

	return fmt.Sprintf("%v", nodeValue), nil
}
