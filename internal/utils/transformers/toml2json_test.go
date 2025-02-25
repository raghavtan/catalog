package transformers

import (
	"testing"
)

func TestToml2json(t *testing.T) {
	tests := []struct {
		name      string
		tomlData  string
		wantJSON  string
		wantError bool
	}{
		{
			name:     "valid TOML",
			tomlData: `key = "value"`,
			wantJSON: `{"key":"value"}`,
		},
		{
			name:      "invalid TOML",
			tomlData:  `key = value`,
			wantError: true,
		},
		{
			name:     "nested TOML",
			tomlData: `parent = { child = "value" }`,
			wantJSON: `{"parent":{"child":"value"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJSON, err := Toml2json(tt.tomlData)
			if (err != nil) != tt.wantError {
				t.Errorf("Toml2json() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && string(gotJSON) != tt.wantJSON {
				t.Errorf("Toml2json() = %v, want %v", string(gotJSON), tt.wantJSON)
			}
		})
	}
}
