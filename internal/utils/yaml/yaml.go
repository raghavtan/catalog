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

func ParseState[T any]() ([]*T, error) {
	tKind, kindErr := GetKindFromGeneric(fmt.Sprintf("%T", new(T)))
	if kindErr != nil {
		return nil, kindErr
	}

	filePath := fmt.Sprintf(".state/%s.yaml", tKind)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []*T{}, nil
	}

	res, parseErr := parse[T](filePath)
	if parseErr != nil {
		return nil, errors.Join(errors.New("failed to parse state file"), parseErr)
	}

	return res, nil
}

func ParseConfig[T any]() ([]*T, error) {
	tKind, kindErr := GetKindFromGeneric(fmt.Sprintf("%T", new(T)))
	if kindErr != nil {
		return nil, kindErr
	}

	filePath := fmt.Sprintf("config/%s.yaml", tKind)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return []*T{}, nil
	}

	res, parseErr := parse[T](filePath)
	if parseErr != nil {
		log.Printf("Failed to parse state: %v", parseErr)
		return nil, errors.Join(errors.New("failed to parse config file"), parseErr)
	}

	return res, nil
}

func GetKindFromGeneric(typeName string) (string, error) {
	start := strings.LastIndex(typeName, ".") + 1
	end := strings.Index(typeName, "DTO")

	if start == -1 || end == -1 || start >= end {
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
