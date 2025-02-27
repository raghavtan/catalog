package yaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type DefinitionType string

const (
	State  DefinitionType = "state"
	Config DefinitionType = "config"
)

func Parse[T any](defintionType DefinitionType, getKey func(def *T) string) (map[string]*T, error) {
	filePath, pathErr := getFilePath[T](defintionType)
	if pathErr != nil {
		return nil, pathErr
	}

	if filePath == "" {
		return make(map[string]*T), nil
	}

	defintions, parseErr := parse[T](filePath)
	if parseErr != nil {
		return nil, errors.Join(fmt.Errorf("failed to parse %s file: \"%s\"", defintionType, filePath), parseErr)
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
//   - defintionType: The type of the definitions to parse.
//   - getKey: A function that takes a pointer to a definition of type T and returns a string key for that definition.
//   - filter: A function that takes a pointer to a definition of type T and returns a boolean indicating whether the definition should be included in the result.
//
// Returns:
//   - A map where the keys are the results of the getKey function and the values are pointers to the filtered definitions of type T.
//   - An error if there was an issue getting the file path or parsing the YAML file.
func ParseFiltered[T any](defintionType DefinitionType, getKey func(def *T) string, filter func(def *T) bool) (map[string]*T, error) {
	filePath, pathErr := getFilePath[T](defintionType)
	if pathErr != nil {
		return nil, pathErr
	}

	if filePath == "" {
		return make(map[string]*T), nil
	}

	defintions, parseErr := parse[T](filePath)
	if parseErr != nil {
		return nil, errors.Join(fmt.Errorf("failed to parse %s file: \"%s\"", defintionType, filePath), parseErr)
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

func parse[T any](filePath string) ([]*T, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var results []*T
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	for {
		var result T
		err = decoder.Decode(&result)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		results = append(results, &result)
	}

	return results, nil
}

func WriteState[T any](data []*T) error {
	tKind, kindErr := GetKindFromGeneric(fmt.Sprintf("%T", new(T)))
	if kindErr != nil {
		return kindErr
	}

	stateFile := fmt.Sprintf(".state/%s.yaml", tKind)
	// Create the directory if it does not exist
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

	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	for _, item := range data {
		if err := encoder.Encode(item); err != nil {
			return err
		}
	}
	if err := encoder.Close(); err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf(".state/%s.yaml", tKind), buffer.Bytes(), 0644)
}

func getFilePath[T any](defintionType DefinitionType) (string, error) {
	var fileLocation string
	switch defintionType {
	case State:
		fileLocation = ".state"
	case Config:
		fileLocation = "config"
	default:
		log.Fatalf("unknown definition type: %s", defintionType)
	}

	tKind, kindErr := GetKindFromGeneric(fmt.Sprintf("%T", new(T)))
	if kindErr != nil {
		return "", kindErr
	}

	filePath := fmt.Sprintf("%s/%s.yaml", fileLocation, tKind)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", nil
	}

	return filePath, nil
}
