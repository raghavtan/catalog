package utils

import (
	"regexp"
)

const placeholderPattern = `:[a-zA-Z_]+`

func ReplacePlaceholder(target, value string) string {
	re := regexp.MustCompile(placeholderPattern)
	return re.ReplaceAllStringFunc(target, func(match string) string {
		return value
	})
}
