//go:build unit
// +build unit

package yaml_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	thisyaml "github.com/motain/fact-collector/internal/utils/yaml"
)

type TestDTO struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
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
				{Name: "John", Age: 30},
				{Name: "Jane", Age: 25},
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
func TestParseState(t *testing.T) {
	tests := []struct {
		name      string
		setup     func()
		teardown  func()
		expected  []*TestDTO
		expectErr bool
	}{
		{
			name: "successful_parse",
			setup: func() {
				os.MkdirAll(".state", os.ModePerm)
				data := `
---
name: John
age: 30
---
name: Jane
age: 25
`
				os.WriteFile(".state/test.yaml", []byte(data), 0644)
			},
			teardown: func() {
				os.RemoveAll(".state")
			},
			expected: []*TestDTO{
				{Name: "John", Age: 30},
				{Name: "Jane", Age: 25},
			},
			expectErr: false,
		},
		{
			name: "file_not_exist",
			setup: func() {
				os.RemoveAll(".state")
			},
			teardown:  func() {},
			expected:  []*TestDTO{},
			expectErr: false,
		},
		{
			name: "invalid_yaml",
			setup: func() {
				os.MkdirAll(".state", os.ModePerm)
				data := `
		- name: John
			age: thirty
		`
				os.WriteFile(".state/test.yaml", []byte(data), 0644)
			},
			teardown: func() {
				os.RemoveAll(".state")
			},
			expected:  nil,
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

			result, err := thisyaml.ParseState[TestDTO]()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func TestParseConfig(t *testing.T) {
	tests := []struct {
		name      string
		setup     func()
		teardown  func()
		expected  []*TestDTO
		expectErr bool
	}{
		{
			name: "successful_parse",
			setup: func() {
				os.MkdirAll("config", os.ModePerm)
				data := `
---
name: John
age: 30
---
name: Jane
age: 25
`
				os.WriteFile("config/test.yaml", []byte(data), 0644)
			},
			teardown: func() {
				os.RemoveAll("config")
			},
			expected: []*TestDTO{
				{Name: "John", Age: 30},
				{Name: "Jane", Age: 25},
			},
			expectErr: false,
		},
		{
			name: "file_not_exist",
			setup: func() {
				os.RemoveAll("config")
			},
			teardown:  func() {},
			expected:  []*TestDTO{},
			expectErr: false,
		},
		{
			name: "invalid_yaml",
			setup: func() {
				os.MkdirAll("config", os.ModePerm)
				data := `
		- name: John
			age: thirty
		`
				os.WriteFile("config/test.yaml", []byte(data), 0644)
			},
			teardown: func() {
				os.RemoveAll("config")
			},
			expected:  nil,
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

			result, err := thisyaml.ParseConfig[TestDTO]()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write error")
}
