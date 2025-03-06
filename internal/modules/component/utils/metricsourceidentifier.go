package utils

import (
	"fmt"
)

func GetMetricSourceItentifier(metricName, componentName, componentType string) string {
	componentSlug := GetSlug(componentName, componentType)
	return fmt.Sprintf("%s-%s", metricName, componentSlug)
}
