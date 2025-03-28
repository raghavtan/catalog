package ownerservice_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	githubservice "github.com/motain/of-catalog/internal/services/githubservice/mocks"
	"github.com/motain/of-catalog/internal/services/ownerservice"
	"github.com/motain/of-catalog/internal/services/ownerservice/dtos"
	"github.com/stretchr/testify/assert"
)

var ofOrgMainYAML = `
apiVersion: backstage.io/v1alpha1
kind: Group
metadata:
  name: squad1
  description: ""
  annotations:
    jiraTeamID: "ari:cloud:identity::team/squad1"
  links:
    - url: "https://onefootball.slack.com/archives/FOOBAR"
      title: squad1-chat
      type: slack
      icon: chat
    - url: "https://onefootball.atlassian.net/jira/servicedesk/projects/FOO"
      title: "IT Helpdesk"
      type: project
      icon: jira
    - url: "https://onefootball.atlassian.net/jira/software/projects/ITP/boards/1278"
      title: "IT Projects"
      type: project
      icon: jira
spec:
  profile:
    displayName: SQUAD ONE
  type: squad
  parent: TRIBE FOOBARBZ42
  children: []
`

func TestOwnerService_extractData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGitHubService := githubservice.NewMockGitHubServiceInterface(ctrl)

	tests := []struct {
		name          string
		tribe         string
		squad         string
		mockSetup     func()
		extraCal      func(os *ownerservice.OwnerService)
		expectedError error
		expectedData  *dtos.Owner
	}{
		{
			name:  "successfully extracts data",
			tribe: "TRIBE FOOBARBZ42",
			squad: "squad1",
			mockSetup: func() {
				mockGitHubService.EXPECT().
					GetFileContent("of-org", "main.yaml").
					Return(ofOrgMainYAML, nil)
			},
			expectedError: nil,
			expectedData: &dtos.Owner{
				OwnerID: "ari:cloud:identity::team/squad1",
				SlackChannels: map[string]string{
					"squad1-chat": "https://onefootball.slack.com/archives/FOOBAR",
				},
				Projects: map[string]string{
					"IT Helpdesk": "https://onefootball.atlassian.net/jira/servicedesk/projects/FOO",
					"IT Projects": "https://onefootball.atlassian.net/jira/software/projects/ITP/boards/1278",
				},
				DisplayName: "SQUAD ONE",
			},
		},
		{
			name:  "returns error when GetFileContent fails",
			tribe: "TRIBE FOOBARBZ42",
			squad: "squad1",
			mockSetup: func() {
				mockGitHubService.EXPECT().
					GetFileContent("of-org", "main.yaml").
					Return("", errors.New("failed to fetch file content"))
			},
			expectedError: errors.New("failed to fetch file content"),
			expectedData:  nil,
		},
		{
			name: "returns error when YAML decoding fails",
			mockSetup: func() {
				mockGitHubService.EXPECT().
					GetFileContent("of-org", "main.yaml").
					Return("invalid yaml content", nil)
			},
			expectedError: errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `invalid...` into dtos.Group"),
			expectedData:  nil,
		},
		{
			name:  "do not call twice GetFileContent but uses the cached data",
			tribe: "TRIBE FOOBARBZ42",
			squad: "squad1",
			mockSetup: func() {
				mockGitHubService.EXPECT().
					GetFileContent("of-org", "main.yaml").
					Return(ofOrgMainYAML, nil)
			},
			extraCal: func(os *ownerservice.OwnerService) {
				os.GetOwnerByTribeAndSquad("TRIBE FOOBARBZ42", "squad1")
			},
			expectedError: nil,
			expectedData: &dtos.Owner{
				OwnerID: "ari:cloud:identity::team/squad1",
				SlackChannels: map[string]string{
					"squad1-chat": "https://onefootball.slack.com/archives/FOOBAR",
				},
				Projects: map[string]string{
					"IT Helpdesk": "https://onefootball.atlassian.net/jira/servicedesk/projects/FOO",
					"IT Projects": "https://onefootball.atlassian.net/jira/software/projects/ITP/boards/1278",
				},
				DisplayName: "SQUAD ONE",
			},
		},
		{
			name:  "returns nil when tribe does not match",
			tribe: "TRIBE WRONG",
			squad: "squad1",
			mockSetup: func() {
				mockGitHubService.EXPECT().
					GetFileContent("of-org", "main.yaml").
					Return(ofOrgMainYAML, nil)
			},
			expectedError: errors.New("no matching group found"),
			expectedData:  nil,
		},
		{
			name:  "returns nil when squad does not match",
			tribe: "TRIBE FOOBARBZ42",
			squad: "wrong-squad",
			mockSetup: func() {
				mockGitHubService.EXPECT().
					GetFileContent("of-org", "main.yaml").
					Return(ofOrgMainYAML, nil)
			},
			expectedError: errors.New("no matching group found"),
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			service := ownerservice.NewOwnerService(mockGitHubService)
			data, err := service.GetOwnerByTribeAndSquad(tt.tribe, tt.squad)

			if tt.extraCal != nil {
				tt.extraCal(service)
			}

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, data)
			}
		})
	}
}
