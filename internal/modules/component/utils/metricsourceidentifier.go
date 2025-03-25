package utils

import (
	"fmt"
)

func GetMetricSourceIdentifier(metricName, componentName, componentType string) string {
	componentSlug := GetSlug(componentName, componentType)
	return fmt.Sprintf("%s-%s", metricName, componentSlug)
}
