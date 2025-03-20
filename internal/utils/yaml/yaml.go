package yaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

const (
	StateLocation  = ".state"
	Kind           = "Kind"
	DTO            = "DTO"
	FilePermission = 0644
)

type ParseInput struct {
	RootLocation string
	Recursive    bool
}

type KeyExtractor[T any] func(def *T) string
type Filter[T any] func(def *T) bool

func GetStateInput(stateRootLocation string) ParseInput {
	return ParseInput{
		RootLocation: StateLocation,
		Recursive:    false,
	}
}

func Parse[T any](parseInput ParseInput, getKey KeyExtractor[T]) (map[string]*T, error) {
	return ParseFiltered(parseInput, getKey, func(def *T) bool { return true })
}

func ParseFiltered[T any](parseInput ParseInput, getKey KeyExtractor[T], filter Filter[T]) (map[string]*T, error) {
	defintions, parseErr := getDefinitions[T](parseInput)
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
	end := strings.Index(typeName, DTO)

	err := errors.New("could not extract DTO name from literal type")
	if start == -1 {
		return "", err
	}

	if end == -1 {
		return "", err
	}

	if start >= end {
		return "", err
	}

	return strings.ToLower(typeName[start:end]), nil
}

func parse[T any](tKind, globString string) ([]*T, error) {
	basepath, pattern := doublestar.SplitPattern(globString)
	matches, globErr := doublestar.Glob(os.DirFS(basepath), pattern)
	if globErr != nil {
		return nil, globErr
	}

	var results []*T
	for _, match := range matches {
		decodedResults, decodeErr := decodeData[T](tKind, filepath.Join(basepath, match))
		if decodeErr != nil {
			return nil, decodeErr
		}

		results = append(results, decodedResults...)
	}
	return results, nil
}

func SortResults[T any](result []*T, getKey KeyExtractor[T]) []*T {
	componentsName := make([]string, 0, len(result))
	componentsMap := make(map[string]*T)
	for _, component := range result {
		key := getKey(component)
		componentsMap[key] = component
		componentsName = append(componentsName, key)
	}
	sort.Strings(componentsName)

	uniqueSortedComponentsName := make([]*T, 0, len(componentsName))
	for i, componentName := range componentsName {
		if i == 0 || componentName != componentsName[i-1] {
			uniqueSortedComponentsName = append(uniqueSortedComponentsName, componentsMap[componentName])
		}
	}
	return uniqueSortedComponentsName
}

func WriteState[T any](data []*T) error {
	tKind, kindErr := GetKindFromGeneric(fmt.Sprintf("%T", new(T)))
	if kindErr != nil {
		return kindErr
	}

	stateFileLocation := filepath.Join(StateLocation, getKindFileName(tKind))
	dir := filepath.Dir(stateFileLocation)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if len(data) == 0 {
		if err := os.Remove(stateFileLocation); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	buffer, encodeErr := encodeData(data)
	if encodeErr != nil {
		return encodeErr
	}

	return os.WriteFile(stateFileLocation, buffer, FilePermission)
}

func getDefinitions[T any](parseInput ParseInput) ([]*T, error) {
	tKind, kindErr := GetKindFromGeneric(fmt.Sprintf("%T", new(T)))
	if kindErr != nil {
		return nil, kindErr
	}

	filePath, pathErr := getFilePath[T](tKind, parseInput)
	if pathErr != nil {
		return nil, pathErr
	}

	if filePath == "" {
		return nil, nil
	}

	defintions, parseErr := parse[T](tKind, filePath)
	if parseErr != nil {
		return nil, errors.Join(fmt.Errorf("failed to parse files at %s: \"%s\"", parseInput.RootLocation, filePath), parseErr)
	}

	return defintions, nil
}

func getFilePath[T any](tKind string, parseInput ParseInput) (string, error) {
	directory := strings.TrimRight(parseInput.RootLocation, string(filepath.Separator))
	if parseInput.Recursive {
		directory = fmt.Sprintf("%s/**", directory)
	}

	fileString := fmt.Sprintf("%s*", tKind)
	return filepath.Join(directory, getKindFileName(fileString)), nil
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

func decodeData[T any](tKind, fileName string) ([]*T, error) {
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

		// Assuming the struct has a field named Kind
		kindField := reflect.ValueOf(result).FieldByName(Kind)
		if kindField.IsValid() && strings.EqualFold(kindField.String(), tKind) {
			results = append(results, &result)
		}
	}

	return results, nil
}

func getKindFileName(kind string) string {
	return fmt.Sprintf("%s.yaml", kind)
}
