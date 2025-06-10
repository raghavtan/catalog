package ownerservice_test

import (
	"fmt"
	"testing"

	"github.com/motain/of-catalog/internal/services/ownerservice"
	"github.com/motain/of-catalog/internal/services/ownerservice/dtos"
	"github.com/stretchr/testify/assert"
)

func mockGetSquadDetails(squadName string) (ownerservice.SquadDetails, error) {
	mockData := map[string]ownerservice.SquadDetails{
		"squad1": {
			JiraTeamID:      "ari:cloud:identity::team/squad1",
			SlackURL:        "https://onefootball.slack.com/archives/FOOBAR",
			SlackTitle:      "squad1-chat",
			JiraProjectURL:  "https://onefootball.atlassian.net/jira/software/projects/ITP/boards/1278",
			JiraProjectName: "IT Projects",
			Tribe:           "TRIBE FOOBARBZ42",
		},
		"squad-no-slack": {
			JiraTeamID:      "ari:cloud:identity::team/squad-no-slack",
			SlackURL:        "",
			SlackTitle:      "",
			JiraProjectURL:  "https://onefootball.atlassian.net/jira/software/projects/SP/boards/123",
			JiraProjectName: "Some Project",
			Tribe:           "TRIBE FOOBARBZ42",
		},
		"squad-no-projects": {
			JiraTeamID:      "ari:cloud:identity::team/squad-no-projects",
			SlackURL:        "https://onefootball.slack.com/archives/CHANNEL",
			SlackTitle:      "squad-chat",
			JiraProjectURL:  "",
			JiraProjectName: "",
			Tribe:           "TRIBE FOOBARBZ42",
		},
		"minimal-squad": {
			JiraTeamID:      "ari:cloud:identity::team/minimal-squad",
			SlackURL:        "",
			SlackTitle:      "",
			JiraProjectURL:  "",
			JiraProjectName: "",
			Tribe:           "TRIBE FOOBARBZ42",
		},
		"no-jira-team": {
			JiraTeamID:      "",
			SlackURL:        "https://onefootball.slack.com/archives/NOJIRA",
			SlackTitle:      "no-jira-chat",
			JiraProjectURL:  "",
			JiraProjectName: "",
			Tribe:           "TRIBE FOOBARBZ42",
		},
		"different-tribe-squad": {
			JiraTeamID:      "ari:cloud:identity::team/different-tribe-squad",
			SlackURL:        "https://onefootball.slack.com/archives/DIFFERENT",
			SlackTitle:      "different-squad-chat",
			JiraProjectURL:  "https://onefootball.atlassian.net/jira/software/projects/DIFF/boards/999",
			JiraProjectName: "Different Project",
			Tribe:           "DIFFERENT TRIBE",
		},
	}

	if details, exists := mockData[squadName]; exists {
		return details, nil
	}
	return ownerservice.SquadDetails{}, fmt.Errorf("squad '%s' not found", squadName)
}

func mockGetSquadDetailsWithError(squadName string) (ownerservice.SquadDetails, error) {
	return ownerservice.SquadDetails{}, fmt.Errorf("failed to fetch squad details for '%s'", squadName)
}

func TestOwnerService_GetOwnerByTribeAndSquad(t *testing.T) {
	tests := []struct {
		name               string
		tribe              string
		squad              string
		mockFunc           func(string) (ownerservice.SquadDetails, error)
		expectedError      string
		expectedData       *dtos.Owner
		shouldContainError bool
	}{
		{
			name:     "successfully gets owner data with all fields",
			tribe:    "TRIBE FOOBARBZ42",
			squad:    "squad1",
			mockFunc: mockGetSquadDetails,
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
			mockFunc:      mockGetSquadDetails,
			expectedError: "squad 'nonexistent-squad' not found",
		},
		{
			name:          "returns error when tribe doesn't match",
			tribe:         "WRONG TRIBE",
			squad:         "squad1",
			mockFunc:      mockGetSquadDetails,
			expectedError: "squad 'squad1' belongs to tribe 'TRIBE FOOBARBZ42', not 'WRONG TRIBE'",
		},
		{
			name:     "handles squad with no Slack channel",
			tribe:    "TRIBE FOOBARBZ42",
			squad:    "squad-no-slack",
			mockFunc: mockGetSquadDetails,
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
			name:     "handles squad with no projects",
			tribe:    "TRIBE FOOBARBZ42",
			squad:    "squad-no-projects",
			mockFunc: mockGetSquadDetails,
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
			name:     "handles squad with minimal data",
			tribe:    "TRIBE FOOBARBZ42",
			squad:    "minimal-squad",
			mockFunc: mockGetSquadDetails,
			expectedData: &dtos.Owner{
				OwnerID:       "ari:cloud:identity::team/minimal-squad",
				SlackChannels: map[string]string{},
				Projects:      map[string]string{},
				DisplayName:   "minimal-squad",
			},
		},
		{
			name:     "handles squad with empty Jira team ID",
			tribe:    "TRIBE FOOBARBZ42",
			squad:    "no-jira-team",
			mockFunc: mockGetSquadDetails,
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
			name:     "successfully gets owner data for different tribe",
			tribe:    "DIFFERENT TRIBE",
			squad:    "different-tribe-squad",
			mockFunc: mockGetSquadDetails,
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
			name:          "handles GetSquadDetails error",
			tribe:         "TRIBE FOOBARBZ42",
			squad:         "any-squad",
			mockFunc:      mockGetSquadDetailsWithError,
			expectedError: "failed to fetch squad details for 'any-squad'",
		},
		{
			name:          "handles empty tribe parameter",
			tribe:         "",
			squad:         "squad1",
			mockFunc:      mockGetSquadDetails,
			expectedError: "squad 'squad1' belongs to tribe 'TRIBE FOOBARBZ42', not ''",
		},
		{
			name:          "handles empty squad parameter",
			tribe:         "TRIBE FOOBARBZ42",
			squad:         "",
			mockFunc:      mockGetSquadDetails,
			expectedError: "squad '' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use dependency injection to provide the mock function
			service := ownerservice.NewOwnerServiceWithDependencies(tt.mockFunc)
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

	service := ownerservice.NewOwnerServiceWithDependencies(mockGetSquadDetails)

	// Verify the interface method works
	owner, err := service.GetOwnerByTribeAndSquad("TRIBE FOOBARBZ42", "squad1")

	assert.NoError(t, err)
	assert.NotNil(t, owner)
	assert.Equal(t, "ari:cloud:identity::team/squad1", owner.OwnerID)
}

func TestNewOwnerService(t *testing.T) {
	// Test that NewOwnerService creates a service with default dependencies
	service := ownerservice.NewOwnerService()
	assert.NotNil(t, service)

	// The service should be ready to use (though we can't test GetSquadDetails without real data)
	// This test ensures the constructor works correctly
}

func TestNewOwnerServiceWithDependencies(t *testing.T) {
	// Test that dependency injection works
	mockFunc := func(squadName string) (ownerservice.SquadDetails, error) {
		return ownerservice.SquadDetails{
			JiraTeamID: "test-jira-id",
			Tribe:      "TEST TRIBE",
		}, nil
	}

	service := ownerservice.NewOwnerServiceWithDependencies(mockFunc)
	assert.NotNil(t, service)

	owner, err := service.GetOwnerByTribeAndSquad("TEST TRIBE", "test-squad")
	assert.NoError(t, err)
	assert.Equal(t, "test-jira-id", owner.OwnerID)
	assert.Equal(t, "test-squad", owner.DisplayName)
}
