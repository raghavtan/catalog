//go:build functional
// +build functional

package tests

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLIOutput(t *testing.T) {
	cmd := exec.Command("go", "run", "../main.go")
	output, err := cmd.CombinedOutput()

	cmd.Env = append(cmd.Env, "FC_GITHUB_USER=github_username")

	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(output), "Data from GitHub Org:"))
}
