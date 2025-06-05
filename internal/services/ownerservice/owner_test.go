package ownerservice_test

import (
	"testing"

	"github.com/motain/of-catalog/internal/services/ownerservice"
	"github.com/motain/of-catalog/internal/services/ownerservice/dtos"
	"github.com/stretchr/testify/assert"
)

func TestOwnerService_GetOwnerByTribeAndSquad(t *testing.T) {
	tests := []struct {
		name          string
		tribe         string
		squad         string
		expectedError string
		expectedData  *dtos.Owner
	}{
		{
			name:  "successfully gets owner data",
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
			name:          "returns error when tribe doesn't match",
			tribe:         "WRONG TRIBE",
			squad:         "squad1",
			expectedError: "squad 'squad1' belongs to tribe 'TRIBE FOOBARBZ42', not 'WRONG TRIBE'",
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := ownerservice.NewOwnerService()
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
