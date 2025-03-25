package utils_test

import (
	"fmt"
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/utils"
)

func TestGetMetricSourceIdentifier(t *testing.T) {
	tests := []struct {
		metricName    string
		componentName string
		componentType string
		expected      string
	}{
		{"cpu_usage", "webserver-foo", "SERVICE", "cpu_usage-svc-webserver-foo"},
		{"memory_usage", "db-rds", "CLOUD_RESOURCE", "memory_usage-cr-db-rds"},
		{"disk_io", "ios-app", "APPLICATION", "disk_io-app-ios-app"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%s_%s", tt.metricName, tt.componentName, tt.componentType), func(t *testing.T) {
			result := utils.GetMetricSourceIdentifier(tt.metricName, tt.componentName, tt.componentType)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
