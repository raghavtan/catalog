//go:build unit
// +build unit

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

func TestGetGithubToken(t *testing.T) {
	os.Setenv("FC_GITHUB_TOKEN", "foo-bar-baz-42")
	cfg := configservice.NewConfigService()

	org := cfg.GetGithubToken()
	assert.Equal(t, "foo-bar-baz-42", org)
}

func TestGetGithubUser(t *testing.T) {
	os.Setenv("FC_GITHUB_USER", "foo.bar@baz.42")
	cfg := configservice.NewConfigService()

	org := cfg.GetGithubUser()
	assert.Equal(t, "foo.bar@baz.42", org)
}
