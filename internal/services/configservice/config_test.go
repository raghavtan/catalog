package configservice_test

import (
	"os"
	"testing"

	"github.com/motain/fact-collector/internal/services/configservice"
	"github.com/stretchr/testify/assert"
)

func TestGetGithubOrg(t *testing.T) {
	os.Setenv("FC_GITHUB_ORG", "my-org")
	cfg := configservice.NewConfigService()

	org := cfg.GetGithubOrg()
	assert.Equal(t, "my-org", org)
}

func TestGetGithubUser(t *testing.T) {
	os.Setenv("FC_GITHUB_USER", "foo.bar@baz.42")
	cfg := configservice.NewConfigService()

	org := cfg.GetGithubUser()
	assert.Equal(t, "foo.bar@baz.42", org)
}

func GetCompassToken(t *testing.T) {
	os.Setenv("FC_COMPASS_TOKEN", "Zm9vLWJhci1iYXotNDIK")
	cfg := configservice.NewConfigService()

	org := cfg.GetCompassToken()
	assert.Equal(t, "Zm9vLWJhci1iYXotNDIK", org)
}

func TestGetCompassHost(t *testing.T) {
	os.Setenv("FC_COMPASS_HOST", "https://compass.example.com")
	cfg := configservice.NewConfigService()

	org := cfg.GetCompassHost()
	assert.Equal(t, "https://compass.example.com", org)
}

func TestGetCompassCloudId(t *testing.T) {
	os.Setenv("FC_COMPASS_CLOUD_ID", "123456")
	cfg := configservice.NewConfigService()

	org := cfg.GetCompassCloudId()
	assert.Equal(t, "123456", org)
}
