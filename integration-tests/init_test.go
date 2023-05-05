package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mrlyc/cmdr/core/utils"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type ContainerManager struct {
	Pool       *dockertest.Pool
	Repository string
	Tag        string
	Name       string
	TargetPath string
	BinaryPath string
}

func (c *ContainerManager) Run(callback func(manager *ContainerManager, container *dockertest.Resource)) {
	err := c.Pool.Client.Ping()
	utils.CheckError(err)

	cmdrContainer, err := c.Pool.RunWithOptions(&dockertest.RunOptions{
		Repository: c.Repository,
		Tag:        c.Tag,
		Name:       c.Name,
		Cmd:        []string{"sleep", "10"},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.Mounts = []docker.HostMount{
			{
				Target: c.TargetPath,
				Source: c.BinaryPath,
				Type:   "bind",
			},
		}
	})
	utils.CheckError(err)
	defer func() {
		err := c.Pool.Purge(cmdrContainer)
		utils.CheckError(err)
	}()

	callback(c, cmdrContainer)

}

var containerManager *ContainerManager

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	utils.CheckError(err)

	pwd, err := os.Getwd()
	utils.CheckError(err)

	containerManager = &ContainerManager{
		Pool:       pool,
		Repository: "mcr.microsoft.com/devcontainers/go",
		Tag:        "0-1.19",
		Name:       "cmdr-integration-test",
		TargetPath: "/tmp/cmdr",
		BinaryPath: filepath.Join(pwd, "..", "bin", "cmdr"),
	}

	m.Run()
}
