package yaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

const StateLocation = ".state"

type DefinitionType string

const (
	State  DefinitionType = "state"
	Config DefinitionType = "config"
)

func Parse[T any](rootLocation string, recursive bool, getKey func(def *T) string) (map[string]*T, error) {
	defintions, parseErr := getDefinitions[T](rootLocation, recursive)
	if parseErr != nil {
		return nil, parseErr
	}

	mappedDefinition := make(map[string]*T, len(defintions))
	for _, defintion := range defintions {
		key := getKey(defintion)
		mappedDefinition[key] = defintion
	}

	return mappedDefinition, nil
}

// ParseFiltered parses a YAML file based on the provided DefinitionType, filters the parsed definitions using the provided filter function,
// and returns a map of the filtered definitions keyed by the result of the getKey function.
//
// T is a generic type parameter representing the type of the definitions.
//
// Parameters:
//   - getKey: A function that takes a pointer to a definition of type T and returns a string key for that definition.
//   - filter: A function that takes a pointer to a definition of type T and returns a boolean indicating whether the definition should be included in the result.
//
// Returns:
//   - A map where the keys are the results of the getKey function and the values are pointers to the filtered definitions of type T.
//   - An error if there was an issue getting the file path or parsing the YAML file.
func ParseFiltered[T any](rootLocation string, recursive bool, getKey func(def *T) string, filter func(def *T) bool) (map[string]*T, error) {
	defintions, parseErr := getDefinitions[T](rootLocation, recursive)
	if parseErr != nil {
		return nil, parseErr
	}

	mappedDefinition := make(map[string]*T)
	for _, defintion := range defintions {
		if filter(defintion) {
			key := getKey(defintion)
			mappedDefinition[key] = defintion
		}
	}

	return mappedDefinition, nil
}

func GetKindFromGeneric(typeName string) (string, error) {
	start := strings.LastIndex(typeName, ".") + 1
	end := strings.Index(typeName, "DTO")

	if start == -1 {
		return "", errors.New("could not extract part from type name")
	}

	if end == -1 {
		return "", errors.New("could not extract part from type name")
	}

	if start >= end {
		return "", errors.New("could not extract part from type name")
	}

	return strings.ToLower(typeName[start:end]), nil
}

func parse[T any](globString string) ([]*T, error) {
	basepath, pattern := doublestar.SplitPattern(globString)
	matches, globErr := doublestar.Glob(os.DirFS(basepath), pattern)
	if globErr != nil {
		return nil, globErr
	}

	var results []*T
	for _, match := range matches {
		decodedResults, decodeErr := decodeData[T](basepath + "/" + match)
		if decodeErr != nil {
			return nil, decodeErr
		}

		results = append(results, decodedResults...)
	}
	return results, nil
}

func WriteState[T any](data []*T) error {
	tKind, kindErr := GetKindFromGeneric(fmt.Sprintf("%T", new(T)))
	if kindErr != nil {
		return kindErr
	}

	stateFile := fmt.Sprintf("%s/%s.yaml", StateLocation, tKind)
	dir := filepath.Dir(stateFile)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if len(data) == 0 {
		if err := os.Remove(stateFile); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	buffer, encodeErr := encodeData(data)
	if encodeErr != nil {
		return encodeErr
	}

	return os.WriteFile(stateFile, buffer, 0644)
}

func getDefinitions[T any](rootLocation string, recursive bool) ([]*T, error) {
	filePath, pathErr := getFilePath[T](rootLocation, recursive)
	if pathErr != nil {
		return nil, pathErr
	}

	if filePath == "" {
		return nil, nil
	}

	defintions, parseErr := parse[T](filePath)
	if parseErr != nil {
		return nil, errors.Join(fmt.Errorf("failed to parse files at %s: \"%s\"", rootLocation, filePath), parseErr)
	}

	return defintions, nil
}

func getFilePath[T any](rootLocation string, recursive bool) (string, error) {
	directory := strings.TrimRight(rootLocation, "/")
	if recursive {
		directory = fmt.Sprintf("%s/**", directory)
	}

	tKind, kindErr := GetKindFromGeneric(fmt.Sprintf("%T", new(T)))
	if kindErr != nil {
		return "", kindErr
	}

	filePath := fmt.Sprintf("%s/%s*.yaml", directory, tKind)

	return filePath, nil
}

func encodeData[T any](data []*T) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	for _, item := range data {
		if encodeErr := encoder.Encode(item); encodeErr != nil {
			return nil, encodeErr
		}
	}
	if closeErr := encoder.Close(); closeErr != nil {
		return nil, closeErr
	}

	return buffer.Bytes(), nil
}

func decodeData[T any](fileName string) ([]*T, error) {
	data, readErr := os.ReadFile(fileName)
	if readErr != nil {
		return nil, readErr
	}

	var results []*T
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	for {
		var result T
		decodeErr := decoder.Decode(&result)
		if decodeErr != nil {
			if decodeErr == io.EOF {
				break
			}
			return nil, decodeErr
		}

		results = append(results, &result)
	}

	return results, nil
}
