package transformers

import (
	"encoding/json"

	"github.com/pelletier/go-toml/v2"
)

func Toml2json(tomlData string) ([]byte, error) {
	// Parse TOML into a map
	var data map[string]interface{}
	errUnmarshal := toml.Unmarshal([]byte(tomlData), &data)
	if errUnmarshal != nil {
		return nil, errUnmarshal
	}

	return json.Marshal(data)
}
