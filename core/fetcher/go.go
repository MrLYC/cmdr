package fetcher

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mrlyc/cmdr/core"
	"github.com/pkg/errors"
)

type GoInstaller struct {
	goPath string
	scheme string
}

func (g *GoInstaller) IsSupport(uri string) bool {
	return strings.HasPrefix(uri, g.scheme)
}

func (g *GoInstaller) install(location, dst string) error {
	cmd := exec.Command(g.goPath, "install", "-v", location)
	cmd.Dir = dst
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	envs := []string{fmt.Sprintf("GOBIN=%s", dst)}
	cmd.Env = append(envs, os.Environ()...)

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "install command %s fail", location)
	}

	return nil
}

func (g *GoInstaller) Fetch(name, version, uri, dst string) error {
	logger := core.GetLogger()
	location := strings.TrimPrefix(uri, g.scheme)

	if strings.Contains(location, "@") {
		return g.install(location, dst)
	}

	var err error
	for _, detected := range []string{
		fmt.Sprintf("%s@%s", location, version),
		fmt.Sprintf("%s@v%s", location, version),
		fmt.Sprintf("%s@latest", location),
		location,
	} {
		logger.Warn("version suffix not set, detecting", map[string]interface{}{
			"suffix": detected,
		})
		err := g.install(detected, dst)
		if err == nil {
			break
		}
	}

	return err
}

func NewGoInstaller(goPath string, schema string) *GoInstaller {
	return &GoInstaller{
		goPath: goPath,
		scheme: schema,
	}
}

func NewDefaultGoInstaller() *GoInstaller {
	return NewGoInstaller("go", "go://")
}
