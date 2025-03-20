package yaml_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	thisyaml "github.com/motain/of-catalog/internal/utils/yaml"
)

type TestDTO struct {
	APIVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
	Spec       Spec   `yaml:"spec" json:"spec"`
}

type Spec struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}

func getTestDTO(name string, age int) *TestDTO {
	return &TestDTO{
		APIVersion: "v1",
		Kind:       "test",
		Spec: Spec{
			Name: name,
			Age:  age,
		},
	}
}

func TestWriteState(t *testing.T) {
	tests := []struct {
		name      string
		data      []*TestDTO
		setup     func()
		teardown  func()
		expectErr bool
	}{
		{
			name: "successful_write",
			data: []*TestDTO{
				getTestDTO("John", 30),
				getTestDTO("Jane", 25),
			},
			setup: func() {
				os.MkdirAll(".state", os.ModePerm)
			},
			teardown: func() {
				os.RemoveAll(".state")
			},
			expectErr: false,
		},
		// {
		// 	name: "error_creating_directory",
		// 	data: []*TestDTO{
		// 		{Name: "John", Age: 30},
		// 	},
		// 	setup: func() {
		// 		os.MkdirAll(".state", os.ModePerm)
		// 		os.Chmod(".", 0444)
		// 	},
		// 	teardown: func() {
		// 		err := os.Chmod(".", 0755)
		// 		require.NoError(t, err)
		// 		os.RemoveAll(".state")
		// 	},
		// 	expectErr: true,
		// },
		// {
		// 	name: "error encoding data",
		// 	data: []*TestDTO{
		// 		{Name: "John", Age: 30},
		// 	},
		// 	setup: func() {
		// 		thisyaml.NewEncoder = func(w io.Writer) *yaml.Encoder {
		// 			return &yaml.Encoder{Encoder: &errorWriter{}}
		// 		}
		// 	},
		// 	teardown: func() {
		// 		yaml.NewEncoder = yaml.NewEncoder
		// 	},
		// 	expectErr: true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.teardown != nil {
				defer tt.teardown()
			}

			err := thisyaml.WriteState(tt.data)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				filePath := filepath.Join(".state", "test.yaml")
				data, readErr := os.ReadFile(filePath)
				require.NoError(t, readErr)

				var results []*TestDTO
				decoder := yaml.NewDecoder(bytes.NewReader(data))
				for {
					var result TestDTO
					err = decoder.Decode(&result)
					if err != nil {
						if err == io.EOF {
							break
						}
						require.NoError(t, err)
					}
					results = append(results, &result)
				}

				assert.Equal(t, tt.data, results)
			}
		})
	}
}
func TestGetKindFromGeneric(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		expectErr bool
	}{
		{
			name:      "valid_type",
			input:     fmt.Sprintf("%T", new(TestDTO)),
			expected:  "test",
			expectErr: false,
		},
		{
			name:      "invalid_type_no_dto",
			input:     fmt.Sprintf("%T", new(struct{ Name string })),
			expected:  "",
			expectErr: true,
		},
		{
			name:      "invalid_type_empty",
			input:     fmt.Sprintf("%T", new(struct{})),
			expected:  "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := thisyaml.GetKindFromGeneric(tt.input)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		parseInput thisyaml.ParseInput
		getKey     func(def *TestDTO) string
		setup      func()
		teardown   func()
		expected   map[string]*TestDTO
		expectErr  bool
	}{
		{
			name:       "successful_parse_non_recursive",
			parseInput: thisyaml.ParseInput{RootLocation: ".state", Recursive: false},
			getKey: func(def *TestDTO) string {
				return def.Spec.Name
			},
			setup: func() {
				os.MkdirAll(".state", os.ModePerm)
				data := []*TestDTO{
					getTestDTO("John", 30),
					getTestDTO("Jane", 25),
				}
				thisyaml.WriteState(data)
			},
			teardown: func() {
				os.RemoveAll(".state")
			},
			expected: map[string]*TestDTO{
				"John": getTestDTO("John", 30),
				"Jane": getTestDTO("Jane", 25),
			},
			expectErr: false,
		},
		{
			name:       "successful_parse_recursive",
			parseInput: thisyaml.ParseInput{RootLocation: ".state", Recursive: true},
			getKey: func(def *TestDTO) string {
				return def.Spec.Name
			},
			setup: func() {
				os.MkdirAll(".state/nested", os.ModePerm)
				data := []*TestDTO{
					getTestDTO("Alice", 40),
					getTestDTO("Bob", 35),
				}
				thisyaml.WriteState(data)
			},
			teardown: func() {
				os.RemoveAll(".state")
			},
			expected: map[string]*TestDTO{
				"Alice": getTestDTO("Alice", 40),
				"Bob":   getTestDTO("Bob", 35),
			},
			expectErr: false,
		},
		{
			name:       "no_files_found",
			parseInput: thisyaml.ParseInput{RootLocation: ".state", Recursive: false},
			getKey: func(def *TestDTO) string {
				return def.Spec.Name
			},
			setup: func() {
				os.MkdirAll(".state", os.ModePerm)
			},
			teardown: func() {
				os.RemoveAll(".state")
			},
			expected:  map[string]*TestDTO{},
			expectErr: false,
		},
		{
			name:       "invalid_yaml_file",
			parseInput: thisyaml.ParseInput{RootLocation: ".state", Recursive: false},
			getKey: func(def *TestDTO) string {
				return def.Spec.Name
			},
			setup: func() {
				os.MkdirAll(".state", os.ModePerm)
				os.WriteFile(".state/testdto.yaml", []byte("invalid_yaml"), 0644)
			},
			teardown: func() {
				os.RemoveAll(".state")
			},
			expected:  map[string]*TestDTO{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.teardown != nil {
				defer tt.teardown()
			}

			result, err := thisyaml.Parse(tt.parseInput, tt.getKey)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func TestSortResults(t *testing.T) {
	tests := []struct {
		name     string
		input    []*TestDTO
		getKey   func(def *TestDTO) string
		expected []*TestDTO
	}{
		{
			name: "sort_by_name",
			input: []*TestDTO{
				getTestDTO("Charlie", 40),
				getTestDTO("Alice", 30),
				getTestDTO("Bob", 25),
			},
			getKey: func(def *TestDTO) string {
				return def.Spec.Name
			},
			expected: []*TestDTO{
				getTestDTO("Alice", 30),
				getTestDTO("Bob", 25),
				getTestDTO("Charlie", 40),
			},
		},
		{
			name: "already_sorted",
			input: []*TestDTO{
				getTestDTO("Alice", 30),
				getTestDTO("Bob", 25),
				getTestDTO("Charlie", 40),
			},
			getKey: func(def *TestDTO) string {
				return def.Spec.Name
			},
			expected: []*TestDTO{
				getTestDTO("Alice", 30),
				getTestDTO("Bob", 25),
				getTestDTO("Charlie", 40),
			},
		},
		{
			name:  "empty_input",
			input: []*TestDTO{},
			getKey: func(def *TestDTO) string {
				return def.Spec.Name
			},
			expected: []*TestDTO{},
		},
		{
			name: "single_element",
			input: []*TestDTO{
				getTestDTO("Alice", 30),
			},
			getKey: func(def *TestDTO) string {
				return def.Spec.Name
			},
			expected: []*TestDTO{
				getTestDTO("Alice", 30),
			},
		},
		{
			name: "duplicate_keys",
			input: []*TestDTO{
				getTestDTO("Alice", 25),
				getTestDTO("Alice", 30),
				getTestDTO("Bob", 40),
			},
			getKey: func(def *TestDTO) string {
				return def.Spec.Name
			},
			expected: []*TestDTO{
				getTestDTO("Alice", 30),
				getTestDTO("Bob", 40),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := thisyaml.SortResults(tt.input, tt.getKey)
			assert.Equal(t, tt.expected, result)
		})
	}
}
