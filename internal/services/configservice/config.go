package configservice

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigServiceInterface interface {
	GetGithubOrg() string
	GetGithubUser() string
	GetCompassToken() string
	GetCompassHost() string
	GetCompassCloudId() string
}

type ConfigService struct{}

func NewConfigService() *ConfigService {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &ConfigService{}
}

func (c *ConfigService) GetGithubOrg() string {
	return os.Getenv("FC_GITHUB_ORG")
}

func (c *ConfigService) GetGithubUser() string {
	return os.Getenv("FC_GITHUB_USER")
}

func (c *ConfigService) GetCompassToken() string {
	return os.Getenv("FC_COMPASS_TOKEN")
}

func (c *ConfigService) GetCompassHost() string {
	return os.Getenv("FC_COMPASS_HOST")
}

func (c *ConfigService) GetCompassCloudId() string {
	return os.Getenv("FC_COMPASS_CLOUD_ID")
}
