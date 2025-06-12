package ownerservice_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/services/ownerservice"
	"github.com/motain/of-catalog/internal/services/ownerservice/dtos"
	"github.com/stretchr/testify/assert"
)

// createMockSquads creates a map of mock Group structs for testing
func createMockSquads() map[string]ownerservice.Group {
	return map[string]ownerservice.Group{
		"squad1": {
			APIVersion: "backstage.io/v1alpha1",
			Kind:       "Group",
			Metadata: struct {
				Name        string            `yaml:"name"`
				Description string            `yaml:"description"`
				Annotations map[string]string `yaml:"annotations"`
				Links       []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				} `yaml:"links"`
			}{
				Name:        "squad1",
				Description: "Test squad 1",
				Annotations: map[string]string{
					"jiraTeamID": "ari:cloud:identity::team/squad1",
				},
				Links: []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				}{
					{
						Title: "squad1-chat",
						URL:   "https://onefootball.slack.com/archives/FOOBAR",
						Type:  "slack",
					},
					{
						Title: "IT Projects",
						URL:   "https://onefootball.atlassian.net/jira/software/projects/ITP/boards/1278",
						Icon:  "jira",
						Type:  "project",
					},
				},
			},
			Spec: struct {
				Profile struct {
					DisplayName string `yaml:"displayName"`
				} `yaml:"profile"`
				Type     string   `yaml:"type"`
				Parent   string   `yaml:"parent"`
				Children []string `yaml:"children"`
			}{
				Type:   "squad",
				Parent: "TRIBE FOOBARBZ42",
			},
		},
		"squad-no-slack": {
			APIVersion: "backstage.io/v1alpha1",
			Kind:       "Group",
			Metadata: struct {
				Name        string            `yaml:"name"`
				Description string            `yaml:"description"`
				Annotations map[string]string `yaml:"annotations"`
				Links       []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				} `yaml:"links"`
			}{
				Name:        "squad-no-slack",
				Description: "Test squad without slack",
				Annotations: map[string]string{
					"jiraTeamID": "ari:cloud:identity::team/squad-no-slack",
				},
				Links: []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				}{
					{
						Title: "Some Project",
						URL:   "https://onefootball.atlassian.net/jira/software/projects/SP/boards/123",
						Icon:  "jira",
						Type:  "project",
					},
				},
			},
			Spec: struct {
				Profile struct {
					DisplayName string `yaml:"displayName"`
				} `yaml:"profile"`
				Type     string   `yaml:"type"`
				Parent   string   `yaml:"parent"`
				Children []string `yaml:"children"`
			}{
				Type:   "squad",
				Parent: "TRIBE FOOBARBZ42",
			},
		},
		"squad-no-projects": {
			APIVersion: "backstage.io/v1alpha1",
			Kind:       "Group",
			Metadata: struct {
				Name        string            `yaml:"name"`
				Description string            `yaml:"description"`
				Annotations map[string]string `yaml:"annotations"`
				Links       []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				} `yaml:"links"`
			}{
				Name:        "squad-no-projects",
				Description: "Test squad without projects",
				Annotations: map[string]string{
					"jiraTeamID": "ari:cloud:identity::team/squad-no-projects",
				},
				Links: []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				}{
					{
						Title: "squad-chat",
						URL:   "https://onefootball.slack.com/archives/CHANNEL",
						Type:  "slack",
					},
				},
			},
			Spec: struct {
				Profile struct {
					DisplayName string `yaml:"displayName"`
				} `yaml:"profile"`
				Type     string   `yaml:"type"`
				Parent   string   `yaml:"parent"`
				Children []string `yaml:"children"`
			}{
				Type:   "squad",
				Parent: "TRIBE FOOBARBZ42",
			},
		},
		"minimal-squad": {
			APIVersion: "backstage.io/v1alpha1",
			Kind:       "Group",
			Metadata: struct {
				Name        string            `yaml:"name"`
				Description string            `yaml:"description"`
				Annotations map[string]string `yaml:"annotations"`
				Links       []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				} `yaml:"links"`
			}{
				Name:        "minimal-squad",
				Description: "Minimal test squad",
				Annotations: map[string]string{
					"jiraTeamID": "ari:cloud:identity::team/minimal-squad",
				},
			},
			Spec: struct {
				Profile struct {
					DisplayName string `yaml:"displayName"`
				} `yaml:"profile"`
				Type     string   `yaml:"type"`
				Parent   string   `yaml:"parent"`
				Children []string `yaml:"children"`
			}{
				Type:   "squad",
				Parent: "TRIBE FOOBARBZ42",
			},
		},
		"no-jira-team": {
			APIVersion: "backstage.io/v1alpha1",
			Kind:       "Group",
			Metadata: struct {
				Name        string            `yaml:"name"`
				Description string            `yaml:"description"`
				Annotations map[string]string `yaml:"annotations"`
				Links       []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				} `yaml:"links"`
			}{
				Name:        "no-jira-team",
				Description: "Squad without jira team",
				Links: []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				}{
					{
						Title: "no-jira-chat",
						URL:   "https://onefootball.slack.com/archives/NOJIRA",
						Type:  "slack",
					},
				},
			},
			Spec: struct {
				Profile struct {
					DisplayName string `yaml:"displayName"`
				} `yaml:"profile"`
				Type     string   `yaml:"type"`
				Parent   string   `yaml:"parent"`
				Children []string `yaml:"children"`
			}{
				Type:   "squad",
				Parent: "TRIBE FOOBARBZ42",
			},
		},
		"different-tribe-squad": {
			APIVersion: "backstage.io/v1alpha1",
			Kind:       "Group",
			Metadata: struct {
				Name        string            `yaml:"name"`
				Description string            `yaml:"description"`
				Annotations map[string]string `yaml:"annotations"`
				Links       []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				} `yaml:"links"`
			}{
				Name:        "different-tribe-squad",
				Description: "Squad from different tribe",
				Annotations: map[string]string{
					"jiraTeamID": "ari:cloud:identity::team/different-tribe-squad",
				},
				Links: []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				}{
					{
						Title: "different-squad-chat",
						URL:   "https://onefootball.slack.com/archives/DIFFERENT",
						Type:  "slack",
					},
					{
						Title: "Different Project",
						URL:   "https://onefootball.atlassian.net/jira/software/projects/DIFF/boards/999",
						Icon:  "jira",
						Type:  "project",
					},
				},
			},
			Spec: struct {
				Profile struct {
					DisplayName string `yaml:"displayName"`
				} `yaml:"profile"`
				Type     string   `yaml:"type"`
				Parent   string   `yaml:"parent"`
				Children []string `yaml:"children"`
			}{
				Type:   "squad",
				Parent: "DIFFERENT TRIBE",
			},
		},
	}
}

func TestOwnerService_GetOwnerByTribeAndSquad(t *testing.T) {
	tests := []struct {
		name          string
		tribe         string
		squad         string
		expectedError string
		expectedData  *dtos.Owner
	}{
		{
			name:  "successfully gets owner data with all fields",
			tribe: "TRIBE FOOBARBZ42",
			squad: "squad1",
			expectedData: &dtos.Owner{
				OwnerID: "ari:cloud:identity::team/squad1",
				SlackChannels: map[string]string{
					"squad1-chat": "https://onefootball.slack.com/archives/FOOBAR",
				},
				Projects: map[string]string{
					"IT Projects": "https://onefootball.atlassian.net/jira/software/projects/ITP/boards/1278",
				},
				DisplayName: "squad1",
			},
		},
		{
			name:          "returns error when squad not found",
			tribe:         "TRIBE FOOBARBZ42",
			squad:         "nonexistent-squad",
			expectedError: "squad 'nonexistent-squad' not found",
		},
		{
			name:  "handles squad with no Slack channel",
			tribe: "TRIBE FOOBARBZ42",
			squad: "squad-no-slack",
			expectedData: &dtos.Owner{
				OwnerID:       "ari:cloud:identity::team/squad-no-slack",
				SlackChannels: map[string]string{},
				Projects: map[string]string{
					"Some Project": "https://onefootball.atlassian.net/jira/software/projects/SP/boards/123",
				},
				DisplayName: "squad-no-slack",
			},
		},
		{
			name:  "handles squad with no projects",
			tribe: "TRIBE FOOBARBZ42",
			squad: "squad-no-projects",
			expectedData: &dtos.Owner{
				OwnerID: "ari:cloud:identity::team/squad-no-projects",
				SlackChannels: map[string]string{
					"squad-chat": "https://onefootball.slack.com/archives/CHANNEL",
				},
				Projects:    map[string]string{},
				DisplayName: "squad-no-projects",
			},
		},
		{
			name:  "handles squad with minimal data",
			tribe: "TRIBE FOOBARBZ42",
			squad: "minimal-squad",
			expectedData: &dtos.Owner{
				OwnerID:       "ari:cloud:identity::team/minimal-squad",
				SlackChannels: map[string]string{},
				Projects:      map[string]string{},
				DisplayName:   "minimal-squad",
			},
		},
		{
			name:  "handles squad with empty Jira team ID",
			tribe: "TRIBE FOOBARBZ42",
			squad: "no-jira-team",
			expectedData: &dtos.Owner{
				OwnerID: "",
				SlackChannels: map[string]string{
					"no-jira-chat": "https://onefootball.slack.com/archives/NOJIRA",
				},
				Projects:    map[string]string{},
				DisplayName: "no-jira-team",
			},
		},
		{
			name:  "successfully gets owner data for different tribe",
			tribe: "DIFFERENT TRIBE",
			squad: "different-tribe-squad",
			expectedData: &dtos.Owner{
				OwnerID: "ari:cloud:identity::team/different-tribe-squad",
				SlackChannels: map[string]string{
					"different-squad-chat": "https://onefootball.slack.com/archives/DIFFERENT",
				},
				Projects: map[string]string{
					"Different Project": "https://onefootball.atlassian.net/jira/software/projects/DIFF/boards/999",
				},
				DisplayName: "different-tribe-squad",
			},
		},
		{
			name:          "handles empty squad parameter",
			tribe:         "TRIBE FOOBARBZ42",
			squad:         "",
			expectedError: "squad '' not found",
		},
	}

	// Create service with mock squads
	mockSquads := createMockSquads()
	service := ownerservice.NewOwnerServiceWithSquads(mockSquads)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := service.GetOwnerByTribeAndSquad(tt.tribe, tt.squad)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, data)
			}
		})
	}
}

func TestOwnerService_GetOwnerByTribeAndSquad_Interface(t *testing.T) {
	// Test that OwnerService implements OwnerServiceInterface
	var _ ownerservice.OwnerServiceInterface = &ownerservice.OwnerService{}

	mockSquads := createMockSquads()
	service := ownerservice.NewOwnerServiceWithSquads(mockSquads)

	// Verify the interface method works
	owner, err := service.GetOwnerByTribeAndSquad("TRIBE FOOBARBZ42", "squad1")

	assert.NoError(t, err)
	assert.NotNil(t, owner)
	assert.Equal(t, "ari:cloud:identity::team/squad1", owner.OwnerID)
}

func TestNewOwnerService(t *testing.T) {
	// Test that NewOwnerService creates a service
	service := ownerservice.NewOwnerService()
	assert.NotNil(t, service)

	// The service should be ready to use (though squads map might be empty if YAML file is not found)
	// This test ensures the constructor works correctly
}

func TestNewOwnerServiceWithSquads(t *testing.T) {
	// Test that NewOwnerServiceWithSquads works with provided squads
	mockSquads := map[string]ownerservice.Group{
		"test-squad": {
			APIVersion: "backstage.io/v1alpha1",
			Kind:       "Group",
			Metadata: struct {
				Name        string            `yaml:"name"`
				Description string            `yaml:"description"`
				Annotations map[string]string `yaml:"annotations"`
				Links       []struct {
					Title string `yaml:"title"`
					URL   string `yaml:"url"`
					Icon  string `yaml:"icon"`
					Type  string `yaml:"type"`
				} `yaml:"links"`
			}{
				Name:        "test-squad",
				Annotations: map[string]string{"jiraTeamID": "test-jira-id"},
			},
			Spec: struct {
				Profile struct {
					DisplayName string `yaml:"displayName"`
				} `yaml:"profile"`
				Type     string   `yaml:"type"`
				Parent   string   `yaml:"parent"`
				Children []string `yaml:"children"`
			}{
				Type:   "squad",
				Parent: "TEST TRIBE",
			},
		},
	}

	service := ownerservice.NewOwnerServiceWithSquads(mockSquads)
	assert.NotNil(t, service)

	owner, err := service.GetOwnerByTribeAndSquad("TEST TRIBE", "test-squad")
	assert.NoError(t, err)
	assert.Equal(t, "test-jira-id", owner.OwnerID)
	assert.Equal(t, "test-squad", owner.DisplayName)
}
