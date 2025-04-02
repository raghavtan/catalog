package configservice

//go:generate mockgen -destination=./mocks/mock_config_service.go -package=configservice github.com/motain/of-catalog/internal/services/configservice ConfigServiceInterface

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigServiceInterface interface {
	Get(envVar string) string
	GetGithubOrg() string
	GetGithubToken() string
	GetGithubUser() string
	GetCompassToken() string
	GetCompassHost() string
	GetCompassCloudId() string
	GetPrometheusURL() string
	GetAWSRegion() string
	GetAWSRole() string
}

type ConfigService struct{}

func NewConfigService() *ConfigService {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, defaulting to environment variables")
	}

	return &ConfigService{}
}

func (c *ConfigService) Get(envVar string) string {
	return os.Getenv(envVar)
}

func (c *ConfigService) GetGithubOrg() string {
	githubOrg := os.Getenv("GITHUB_ORG")
	if githubOrg == "" {
		return "motain"
	}
	return githubOrg
}

func (c *ConfigService) GetGithubToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

func (c *ConfigService) GetGithubUser() string {
	return os.Getenv("GITHUB_USER")
}

func (c *ConfigService) GetCompassToken() string {
	return os.Getenv("COMPASS_TOKEN")
}

func (c *ConfigService) GetCompassHost() string {
	return os.Getenv("COMPASS_HOST")
}

func (c *ConfigService) GetCompassCloudId() string {
	return os.Getenv("COMPASS_CLOUD_ID")
}

func (c *ConfigService) GetPrometheusURL() string {
	return os.Getenv("PROMETHEUS_URL")
}

func (c *ConfigService) GetAWSRegion() string {
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		return "eu-west-1"
	}
	return awsRegion
}

func (c *ConfigService) GetAWSRole() string { return os.Getenv("AWS_ROLE") }
