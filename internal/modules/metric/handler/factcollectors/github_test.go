package factcollectors_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/motain/fact-collector/internal/modules/metric/dtos"
	"github.com/motain/fact-collector/internal/modules/metric/handler/factcollectors"
	"github.com/motain/fact-collector/internal/services/githubservice"
)

func TestGithubFactCollector_checkFileJSONPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGithubService := githubservice.NewMockGitHubServiceInterface(ctrl)
	collector := factcollectors.NewGithubFactCollector(mockGithubService)

	tests := []struct {
		name          string
		fact          dtos.Fact
		setupMock     func()
		expected      bool
		expectedError error
	}{
		{
			name: "JSONPathFact valid JSON path",
			fact: dtos.Fact{
				FactType:      dtos.JSONPathFact,
				Repo:          "test-repo",
				FilePath:      "test.json",
				JSONPath:      "name",
				ExpectedValue: "test",
			},
			setupMock: func() {
				mockGithubService.EXPECT().GetFileContent("test-repo", "test.json").Return(`{"name": "test"}`, nil)
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name: "JSONPathFact invalid JSON path",
			fact: dtos.Fact{
				FactType:      dtos.JSONPathFact,
				Repo:          "test-repo",
				FilePath:      "test.json",
				JSONPath:      "invalid",
				ExpectedValue: "test",
			},
			setupMock: func() {
				mockGithubService.EXPECT().GetFileContent("test-repo", "test.json").Return(`{"name": "test"}`, nil)
			},
			expected:      false,
			expectedError: errors.New("jsonpath does not exist"),
		},
		{
			name: "JSONPathFact unsupported file extension",
			fact: dtos.Fact{
				FactType: dtos.JSONPathFact,
				Repo:     "test-repo",
				FilePath: "test.txt",
				JSONPath: "name",
			},
			setupMock:     func() {},
			expected:      false,
			expectedError: errors.New("unsupported file extension: .txt"),
		},
		{
			name: "JSONPathFact valid TOML path",
			fact: dtos.Fact{
				FactType:      dtos.JSONPathFact,
				Repo:          "test-repo",
				FilePath:      "test.toml",
				JSONPath:      "name",
				ExpectedValue: "test",
			},
			setupMock: func() {
				mockGithubService.EXPECT().GetFileContent("test-repo", "test.toml").Return(`name = "test"`, nil)
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name: "JSONPathFact TOML to JSON transformation error",
			fact: dtos.Fact{
				FactType: dtos.JSONPathFact,
				Repo:     "test-repo",
				FilePath: "test.toml",
				JSONPath: "name",
			},
			setupMock: func() {
				mockGithubService.EXPECT().GetFileContent("test-repo", "test.toml").Return(`name = "test"`, nil)
			},
			expected:      false,
			expectedError: errors.New("expected value or formula not provided"),
		},
		{
			name: "FileExistsFact existing file",
			fact: dtos.Fact{
				FactType:      dtos.FileExistsFact,
				Repo:          "test-repo",
				FilePath:      "test.json",
				JSONPath:      "",
				ExpectedValue: "",
			},
			setupMock: func() {
				mockGithubService.EXPECT().GetFileExists("test-repo", "test.json").Return(true, nil)
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name: "FileExistsFact not existing file",
			fact: dtos.Fact{
				FactType:      dtos.FileExistsFact,
				Repo:          "test-repo",
				FilePath:      "test.json",
				JSONPath:      "",
				ExpectedValue: "",
			},
			setupMock: func() {
				mockGithubService.EXPECT().GetFileExists("test-repo", "test.json").Return(false, nil)
			},
			expected:      false,
			expectedError: nil,
		},
		{
			name: "RepoPropertiesFact existing property matching expected value",
			fact: dtos.Fact{
				FactType:      dtos.RepoPropertiesFact,
				Repo:          "test-repo",
				RepoProperty:  "name",
				FilePath:      "",
				JSONPath:      "",
				ExpectedValue: "test",
			},
			setupMock: func() {
				mockGithubService.EXPECT().GetRepoProperties("test-repo").Return(map[string]string{"name": "test"}, nil)
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name: "RepoPropertiesFact existing property not matching expected value",
			fact: dtos.Fact{
				FactType:      dtos.RepoPropertiesFact,
				Repo:          "test-repo",
				RepoProperty:  "name",
				FilePath:      "",
				JSONPath:      "",
				ExpectedValue: "testRepo",
			},
			setupMock: func() {
				mockGithubService.EXPECT().GetRepoProperties("test-repo").Return(map[string]string{"name": "test"}, nil)
			},
			expected:      false,
			expectedError: nil,
		},
		{
			name: "RepoPropertiesFact not existing property",
			fact: dtos.Fact{
				FactType:      dtos.RepoPropertiesFact,
				Repo:          "test-repo",
				RepoProperty:  "foo",
				FilePath:      "",
				JSONPath:      "",
				ExpectedValue: "test",
			},
			setupMock: func() {
				mockGithubService.EXPECT().GetRepoProperties("test-repo").Return(map[string]string{"name": "test"}, nil)
			},
			expected:      false,
			expectedError: errors.New("repo property does not exist"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			result, err := collector.Check(tt.fact)
			if result != tt.expected || (err != nil && err.Error() != tt.expectedError.Error()) {
				t.Errorf("expected %v, got %v, expected error %v, got error %v", tt.expected, result, tt.expectedError, err)
			}
		})
	}
}
