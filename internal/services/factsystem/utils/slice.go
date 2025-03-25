package utils

import (
	"fmt"
	"reflect"
)

func ToSlice[T any](list interface{}) ([]T, error) {
	if values, isT := list.([]T); isT {
		return values, nil
	}

	values := make([]T, len(list.([]interface{})))
	for i, v := range list.([]interface{}) {
		if casted, ok := v.(T); ok {
			values[i] = casted
			continue
		}
		return nil, fmt.Errorf("invalid type, expectd %T, got %T", reflect.TypeOf((*T)(nil)).Elem(), v)
	}

	return values, nil
}
