package githubservice_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-github/v58/github"
	"github.com/motain/of-catalog/internal/services/githubservice"
	"github.com/stretchr/testify/assert"
)

func TestGetRepoProperties(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := githubservice.NewMockGitHubRepositoriesInterface(ctrl)
	service := githubservice.NewGitHubService(mockClient)

	tests := []struct {
		name          string
		repo          string
		mockSetup     func()
		expectedProps map[string]string
		expectedErr   error
	}{
		{
			name: "successful fetch",
			repo: "test-repo",
			mockSetup: func() {
				repo := &github.Repository{
					Name:            github.String("test-repo"),
					Description:     github.String("A test repository"),
					DefaultBranch:   github.String("main"),
					Visibility:      github.String("public"),
					OpenIssuesCount: github.Int(5),
					License: &github.License{
						Name: github.String("MIT"),
					},
				}
				mockClient.EXPECT().Get(gomock.Any(), "motain", "test-repo").Return(repo, nil, nil)
			},
			expectedProps: map[string]string{
				"Name":          "test-repo",
				"Description":   "A test repository",
				"DefaultBranch": "main",
				"Visibility":    "public",
				"OpenIssues":    "5",
				"License":       "MIT",
			},
			expectedErr: nil,
		},
		{
			name: "repository not found",
			repo: "non-existent-repo",
			mockSetup: func() {
				mockClient.EXPECT().Get(gomock.Any(), "motain", "non-existent-repo").Return(nil, nil, errors.New("repository not found"))
			},
			expectedProps: nil,
			expectedErr:   errors.New("repository not found"),
		},
		{
			name: "repository with no license",
			repo: "no-license-repo",
			mockSetup: func() {
				repo := &github.Repository{
					Name:            github.String("no-license-repo"),
					Description:     github.String("A repository with no license"),
					DefaultBranch:   github.String("main"),
					Visibility:      github.String("public"),
					OpenIssuesCount: github.Int(3),
					License:         nil,
				}
				mockClient.EXPECT().Get(gomock.Any(), "motain", "no-license-repo").Return(repo, nil, nil)
			},
			expectedProps: map[string]string{
				"Name":          "no-license-repo",
				"Description":   "A repository with no license",
				"DefaultBranch": "main",
				"Visibility":    "public",
				"OpenIssues":    "3",
				"License":       "",
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			props, err := service.GetRepoProperties(tt.repo)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedProps, props)
		})
	}
}
func TestGetFileContent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := githubservice.NewMockGitHubRepositoriesInterface(ctrl)
	service := githubservice.NewGitHubService(mockClient)

	tests := []struct {
		name        string
		repo        string
		path        string
		mockSetup   func()
		expected    string
		expectedErr error
	}{
		{
			name: "successful fetch",
			repo: "test-repo",
			path: "README.md",
			mockSetup: func() {
				fileContent := &github.RepositoryContent{
					Content: github.String("This is a test file content"),
				}
				mockClient.EXPECT().GetContents(gomock.Any(), "motain", "test-repo", "README.md", nil).Return(fileContent, nil, nil, nil)
			},
			expected:    "This is a test file content",
			expectedErr: nil,
		},
		{
			name: "file not found",
			repo: "test-repo",
			path: "NON_EXISTENT.md",
			mockSetup: func() {
				mockClient.EXPECT().GetContents(gomock.Any(), "motain", "test-repo", "NON_EXISTENT.md", nil).Return(nil, nil, nil, &github.ErrorResponse{Response: &http.Response{StatusCode: 404}})
			},
			expected:    "",
			expectedErr: fmt.Errorf("failed to fetch file: %w", &github.ErrorResponse{Response: &http.Response{StatusCode: 404}}),
		},
		{
			name: "failed to decode file content",
			repo: "test-repo",
			path: "README.md",
			mockSetup: func() {
				mockClient.EXPECT().GetContents(gomock.Any(), "motain", "test-repo", "README.md", nil).Return(nil, nil, nil, errors.New("illegal base64 data at input byte 7"))
			},
			expected:    "",
			expectedErr: fmt.Errorf("failed to fetch file: illegal base64 data at input byte 7"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			content, err := service.GetFileContent(tt.repo, tt.path)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, content)
		})
	}
}
func TestGetFileExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := githubservice.NewMockGitHubRepositoriesInterface(ctrl)
	service := githubservice.NewGitHubService(mockClient)

	tests := []struct {
		name        string
		repo        string
		path        string
		mockSetup   func()
		expected    bool
		expectedErr error
	}{
		{
			name: "file exists",
			repo: "test-repo",
			path: "README.md",
			mockSetup: func() {
				fileContent := &github.RepositoryContent{}
				mockClient.EXPECT().GetContents(gomock.Any(), "motain", "test-repo", "README.md", nil).Return(fileContent, nil, nil, nil)
			},
			expected:    true,
			expectedErr: nil,
		},
		{
			name: "file does not exist",
			repo: "test-repo",
			path: "NON_EXISTENT.md",
			mockSetup: func() {
				mockClient.EXPECT().GetContents(gomock.Any(), "motain", "test-repo", "NON_EXISTENT.md", nil).Return(nil, nil, nil, &github.ErrorResponse{Response: &http.Response{StatusCode: 404}})
			},
			expected:    false,
			expectedErr: nil,
		},
		{
			name: "error fetching file",
			repo: "test-repo",
			path: "README.md",
			mockSetup: func() {
				mockClient.EXPECT().GetContents(gomock.Any(), "motain", "test-repo", "README.md", nil).Return(nil, nil, nil, errors.New("failed to fetch file"))
			},
			expected:    false,
			expectedErr: fmt.Errorf("failed to fetch file: %w", errors.New("failed to fetch file")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			exists, err := service.GetFileExists(tt.repo, tt.path)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, exists)
		})
	}
}
func TestGetRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := githubservice.NewMockGitHubRepositoriesInterface(ctrl)
	service := githubservice.NewGitHubService(mockClient)

	tests := []struct {
		name        string
		repo        string
		mockSetup   func()
		expected    *github.Repository
		expectedErr error
	}{
		{
			name: "successful fetch",
			repo: "test-repo",
			mockSetup: func() {
				repo := &github.Repository{
					Name:          github.String("test-repo"),
					Description:   github.String("A test repository"),
					DefaultBranch: github.String("main"),
					Visibility:    github.String("public"),
				}
				mockClient.EXPECT().Get(gomock.Any(), "motain", "test-repo").Return(repo, nil, nil)
			},
			expected: &github.Repository{
				Name:          github.String("test-repo"),
				Description:   github.String("A test repository"),
				DefaultBranch: github.String("main"),
				Visibility:    github.String("public"),
			},
			expectedErr: nil,
		},
		{
			name: "repository not found",
			repo: "non-existent-repo",
			mockSetup: func() {
				mockClient.EXPECT().Get(gomock.Any(), "motain", "non-existent-repo").Return(nil, nil, errors.New("repository not found"))
			},
			expected:    nil,
			expectedErr: fmt.Errorf("failed to fetch repo: %w", errors.New("repository not found")),
		},
		{
			name: "error fetching repository",
			repo: "error-repo",
			mockSetup: func() {
				mockClient.EXPECT().Get(gomock.Any(), "motain", "error-repo").Return(nil, nil, errors.New("some error"))
			},
			expected:    nil,
			expectedErr: fmt.Errorf("failed to fetch repo: %w", errors.New("some error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			repo, err := service.GetRepo(tt.repo)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expected, repo)
		})
	}
}
