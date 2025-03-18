package compassservice

import "regexp"

type CompassError struct {
	Message string
}

func HasAlreadyExistsError(errs []CompassError) bool {
	for _, err := range errs {
		if isAlreadyExistsError(err.Message) {
			return true
		}
	}

	return false
}

func isAlreadyExistsError(err string) bool {
	matched, _ := regexp.MatchString("^.*already exists", err)

	return matched
}

func HasNotFoundError(errs []CompassError) bool {
	for _, err := range errs {
		if isNotFoundError(err.Message) {
			return true
		}
	}

	return false
}

func isNotFoundError(err string) bool {
	matched, _ := regexp.MatchString(".*not found", err)

	return matched
}
