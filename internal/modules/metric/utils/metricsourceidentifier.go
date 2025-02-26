package utils

import (
	"fmt"

	componentutils "github.com/motain/fact-collector/internal/modules/component/utils"
)

func GetMetricSourceItentifier(metricName, componentName, componentType string) string {
	componentSlug := componentutils.GetSlug(componentName, componentType)
	return fmt.Sprintf("%s-%s", metricName, componentSlug)
}
