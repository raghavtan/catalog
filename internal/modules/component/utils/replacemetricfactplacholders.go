package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/motain/of-catalog/internal/modules/component/dtos"
)

func ReplaceMetricFactPlaceholders(placeholder string, component dtos.ComponentDTO) string {
	re := regexp.MustCompile(`\$\{(.*?)\}`)
	replaceFromComponent := func(path string) string {
		value, err := getFieldByPath(component, path)
		if err != nil {
			return ""
		}

		return fmt.Sprintf("%v", value)
	}

	return re.ReplaceAllStringFunc(placeholder, func(match string) string {
		path := re.FindStringSubmatch(match)[1]
		return replaceFromComponent(path)
	})
}

// getFieldByPath fetches a nested field value using dot notation
func getFieldByPath(obj interface{}, path string) (interface{}, error) {
	fields := strings.Split(path, ".")
	val := reflect.ValueOf(obj)

	// Traverse fields
	for _, field := range fields {
		// Dereference pointer if necessary
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		// Ensure it's a struct
		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("invalid path: %s", path)
		}

		// Get field by name
		val = val.FieldByName(field)

		// If field is invalid, return error
		if !val.IsValid() {
			return nil, fmt.Errorf("field not found: %s", field)
		}
	}

	return val.Interface(), nil
}
