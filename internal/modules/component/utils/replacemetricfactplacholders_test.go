package utils

import (
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/dtos"
)

func TestGetFieldByPath(t *testing.T) {
	type Nested struct {
		Field string
	}

	type TestStruct struct {
		SimpleField string
		NestedField Nested
	}

	tests := []struct {
		name      string
		obj       interface{}
		path      string
		want      interface{}
		expectErr bool
	}{
		{
			name: "simple field",
			obj: TestStruct{
				SimpleField: "simpleValue",
			},
			path:      "SimpleField",
			want:      "simpleValue",
			expectErr: false,
		},
		{
			name: "nested field",
			obj: TestStruct{
				NestedField: Nested{
					Field: "nestedValue",
				},
			},
			path:      "NestedField.Field",
			want:      "nestedValue",
			expectErr: false,
		},
		{
			name: "invalid path",
			obj: TestStruct{
				SimpleField: "simpleValue",
			},
			path:      "InvalidField",
			want:      nil,
			expectErr: true,
		},
		{
			name: "invalid nested path",
			obj: TestStruct{
				NestedField: Nested{
					Field: "nestedValue",
				},
			},
			path:      "NestedField.InvalidField",
			want:      nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFieldByPath(tt.obj, tt.path)
			if (err != nil) != tt.expectErr {
				t.Errorf("getFieldByPath() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFieldByPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestReplaceMetricFactPlaceholders(t *testing.T) {
	tests := []struct {
		name        string
		placeholder string
		want        string
	}{
		{
			name:        "simple placeholder",
			placeholder: "Value: ${Spec.Name}",
			want:        "Value: simpleValue",
		},
		{
			name:        "nested placeholder",
			placeholder: "Value: ${Spec.Name}",
			want:        "Value: simpleValue",
		},
		{
			name:        "multiple placeholders",
			placeholder: "Simple: ${Metadata.Name}, Nested: ${Spec.Name}",
			want:        "Simple: simpleValue, Nested: simpleValue",
		},
		{
			name:        "invalid placeholder",
			placeholder: "Value: ${Meta}",
			want:        "Value: ",
		},
		{
			name:        "invalid nested placeholder",
			placeholder: "Value: ${Meta.InvalidField}",
			want:        "Value: ",
		},
	}

	comp := dtos.ComponentDTO{
		Metadata: dtos.Metadata{Name: "simpleValue"},
		Spec:     dtos.Spec{Name: "simpleValue"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceMetricFactPlaceholders(tt.placeholder, comp)
			if got != tt.want {
				t.Errorf("ReplaceMetricFactPlaceholders() = %v, want %v", got, tt.want)
			}
		})
	}
}
