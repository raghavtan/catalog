package utils_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/modules/component/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetSlug(t *testing.T) {
	tests := []struct {
		name          string
		componentType string
		expected      string
	}{
		{"my-service", "service", "svc-my-service"},
		{"my-cloud-resource", "cloud-resource", "unknown-my-cloud-resource"},
		{"my-website", "website", "web-my-website"},
		{"my-application", "application", "app-my-application"},
		{"my-unknown", "unknown-type", "unknown-my-unknown"},
		{"my-service", "SERVICE", "svc-my-service"},
		{"my-cloud-resource-upper_case", "CLOUD_RESOURCE", "cr-my-cloud-resource-upper_case"},
		{"my-website", "WEBSITE", "web-my-website"},
		{"my-application", "APPLICATION", "app-my-application"},
		{"my-unknown", "UNKNOWN-TYPE", "unknown-my-unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.GetSlug(tt.name, tt.componentType)
			assert.Equal(t, tt.expected, result)
		})
	}
}
