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
