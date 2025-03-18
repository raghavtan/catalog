package configservice_test

import (
	"os"
	"testing"

	"github.com/motain/of-catalog/internal/services/configservice"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	os.Setenv("FOO", "BAR")
	cfg := configservice.NewConfigService()
	assert.Equal(t, "BAR", cfg.Get("FOO"))
}

func TestGetDefaultGithubOrg(t *testing.T) {
	cfg := configservice.NewConfigService()
	assert.Equal(t, "motain", cfg.GetGithubOrg())
}

func TestGetGithubOrg(t *testing.T) {
	os.Setenv("GITHUB_ORG", "my-org")
	cfg := configservice.NewConfigService()
	assert.Equal(t, "my-org", cfg.GetGithubOrg())
}

func TestGetGithubUser(t *testing.T) {
	os.Setenv("GITHUB_USER", "foo.bar@baz.42")
	cfg := configservice.NewConfigService()
	assert.Equal(t, "foo.bar@baz.42", cfg.GetGithubUser())
}

func GetCompassToken(t *testing.T) {
	os.Setenv("COMPASS_TOKEN", "Zm9vLWJhci1iYXotNDIK")
	cfg := configservice.NewConfigService()
	assert.Equal(t, "Zm9vLWJhci1iYXotNDIK", cfg.GetCompassToken())
}

func TestGetCompassHost(t *testing.T) {
	os.Setenv("COMPASS_HOST", "https://compass.example.com")
	cfg := configservice.NewConfigService()
	assert.Equal(t, "https://compass.example.com", cfg.GetCompassHost())
}

func TestGetCompassCloudId(t *testing.T) {
	os.Setenv("COMPASS_CLOUD_ID", "123456")
	cfg := configservice.NewConfigService()
	assert.Equal(t, "123456", cfg.GetCompassCloudId())
}

func TestGetGithubToken(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "123456")
	cfg := configservice.NewConfigService()
	assert.Equal(t, "123456", cfg.GetGithubToken())
}
