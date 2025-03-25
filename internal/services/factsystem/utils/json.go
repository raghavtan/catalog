package utils

import (
	"fmt"

	"github.com/spyzhov/ajson"
)

func InspectExtractedData(JSONPath string, jsonData []byte) (interface{}, error) {
	root, unmarshalErr := ajson.Unmarshal(jsonData)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	node, evalErr := ajson.Eval(root, JSONPath)
	if evalErr != nil {
		return nil, fmt.Errorf("error parsing jsonpath: %v", evalErr)
	}

	nodeValue, vslueErr := node.Value()
	if vslueErr != nil {
		return nil, fmt.Errorf("error parsing jsonpath value: %v", vslueErr)
	}

	if resSlice, ok := nodeValue.([]*ajson.Node); ok {
		result := make([]string, len(resSlice))
		for i, node := range resSlice {
			result[i] = node.String()
		}

		return result, nil
	}

	return nodeValue, nil
}
