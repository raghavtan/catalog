package configservice

import "os"

type ConfigServiceInterface interface {
	GetGithubOrg() string
	GetGithubToken() string
	GetGithubUser() string
}

type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (c *ConfigService) GetGithubOrg() string {
	return os.Getenv("FC_GITHUB_ORG")
}

func (c *ConfigService) GetGithubToken() string {
	return os.Getenv("FC_GITHUB_TOKEN")
}

func (c *ConfigService) GetGithubUser() string {
	return os.Getenv("FC_GITHUB_USER")
}
