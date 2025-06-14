// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/motain/of-catalog/internal/services/githubservice (interfaces: GitHubServiceInterface)

// Package githubservice is a generated GoMock package.
package githubservice

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	github "github.com/google/go-github/v58/github"
)

// MockGitHubServiceInterface is a mock of GitHubServiceInterface interface.
type MockGitHubServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockGitHubServiceInterfaceMockRecorder
}

// MockGitHubServiceInterfaceMockRecorder is the mock recorder for MockGitHubServiceInterface.
type MockGitHubServiceInterfaceMockRecorder struct {
	mock *MockGitHubServiceInterface
}

// NewMockGitHubServiceInterface creates a new mock instance.
func NewMockGitHubServiceInterface(ctrl *gomock.Controller) *MockGitHubServiceInterface {
	mock := &MockGitHubServiceInterface{ctrl: ctrl}
	mock.recorder = &MockGitHubServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitHubServiceInterface) EXPECT() *MockGitHubServiceInterfaceMockRecorder {
	return m.recorder
}

// GetFileContent mocks base method.
func (m *MockGitHubServiceInterface) GetFileContent(arg0, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileContent", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileContent indicates an expected call of GetFileContent.
func (mr *MockGitHubServiceInterfaceMockRecorder) GetFileContent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileContent", reflect.TypeOf((*MockGitHubServiceInterface)(nil).GetFileContent), arg0, arg1)
}

// GetFileExists mocks base method.
func (m *MockGitHubServiceInterface) GetFileExists(arg0, arg1 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileExists", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileExists indicates an expected call of GetFileExists.
func (mr *MockGitHubServiceInterfaceMockRecorder) GetFileExists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileExists", reflect.TypeOf((*MockGitHubServiceInterface)(nil).GetFileExists), arg0, arg1)
}

// GetRepo mocks base method.
func (m *MockGitHubServiceInterface) GetRepo(arg0 string) (*github.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepo", arg0)
	ret0, _ := ret[0].(*github.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepo indicates an expected call of GetRepo.
func (mr *MockGitHubServiceInterfaceMockRecorder) GetRepo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepo", reflect.TypeOf((*MockGitHubServiceInterface)(nil).GetRepo), arg0)
}

// GetRepoProperties mocks base method.
func (m *MockGitHubServiceInterface) GetRepoProperties(arg0 string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepoProperties", arg0)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepoProperties indicates an expected call of GetRepoProperties.
func (mr *MockGitHubServiceInterfaceMockRecorder) GetRepoProperties(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepoProperties", reflect.TypeOf((*MockGitHubServiceInterface)(nil).GetRepoProperties), arg0)
}

// GetRepoURL mocks base method.
func (m *MockGitHubServiceInterface) GetRepoURL(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepoURL", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetRepoURL indicates an expected call of GetRepoURL.
func (mr *MockGitHubServiceInterfaceMockRecorder) GetRepoURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepoURL", reflect.TypeOf((*MockGitHubServiceInterface)(nil).GetRepoURL), arg0)
}

// Search mocks base method.
func (m *MockGitHubServiceInterface) Search(arg0, arg1 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockGitHubServiceInterfaceMockRecorder) Search(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockGitHubServiceInterface)(nil).Search), arg0, arg1)
}
