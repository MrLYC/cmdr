package integration_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

func TestInitCommand(t *testing.T) {
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	containerManager.Run(func(manager *ContainerManager, container *dockertest.Resource) {
		code, err := container.Exec([]string{"/tmp/cmdr", "init"}, dockertest.ExecOptions{
			Env:    []string{"SHELL=zsh"},
			StdErr: stderr,
		})

		assert.NoError(t, err)
		assert.Equal(t, code, 0)

		code, err = container.Exec([]string{"cmdr", "version"}, dockertest.ExecOptions{
			StdOut: stdout,
		})

		assert.NoError(t, err)
		assert.Equal(t, code, 0)
	})

	assert.Equal(t, stderr.String(), "")
	assert.True(t, strings.Contains(stdout.String(), "0.0.0"))
}
