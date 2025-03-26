package utils

import (
	"encoding/json"
	"log"

	"github.com/itchyny/gojq"
)

func InspectExtractedData(JSONPath string, jsonData []byte) (interface{}, error) {
	query, parseQueryErr := gojq.Parse(JSONPath)
	if parseQueryErr != nil {
		return nil, parseQueryErr
	}

	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Fatal(err)
	}

	res := make([]string, 0)
	iter := query.Run(data)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if err, ok := v.(error); ok {
			return nil, err
		}

		res = append(res, v.(string))
	}

	return res, nil
}
