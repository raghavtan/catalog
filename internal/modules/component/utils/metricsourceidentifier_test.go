package utils_test

import (
	"fmt"
	"testing"

	"github.com/motain/fact-collector/internal/modules/component/utils"
)

func TestGetMetricSourceItentifier(t *testing.T) {
	tests := []struct {
		metricName    string
		componentName string
		componentType string
		expected      string
	}{
		{"cpu_usage", "webserver-foo", "service", "cpu_usage-svc-webserver-foo"},
		{"memory_usage", "db-rds", "cloud-resource", "memory_usage-cr-db-rds"},
		{"disk_io", "ios-app", "application", "disk_io-app-ios-app"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%s_%s", tt.metricName, tt.componentName, tt.componentType), func(t *testing.T) {
			result := utils.GetMetricSourceItentifier(tt.metricName, tt.componentName, tt.componentType)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
