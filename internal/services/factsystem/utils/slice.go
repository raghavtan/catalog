package utils

import (
	"fmt"
	"reflect"
)

func ToSlice[T any](list interface{}) ([]T, error) {
	if values, isT := list.([]T); isT {
		return values, nil
	}

	listValue := reflect.ValueOf(list)
	if listValue.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected a slice, got %T", list)
	}

	values := make([]T, listValue.Len())
	for i := 0; i < listValue.Len(); i++ {
		item := listValue.Index(i).Interface()
		if casted, ok := item.(T); ok {
			values[i] = casted
			continue
		}
		return nil, fmt.Errorf("invalid type, expected %T, got %T", reflect.TypeOf((*T)(nil)).Elem(), item)
	}

	return values, nil
}
