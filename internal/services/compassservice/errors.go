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
	matched, _ := regexp.MatchString("^.*already exists\\.", err)

	return matched
}
