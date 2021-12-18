package core

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type ShellProfiler struct {
	BaseStep
	shell  string
	script string
}

func (s *ShellProfiler) String() string {
	return "shell-profiler"
}

func (s *ShellProfiler) isContainsProfile(path string) bool {
	fs := define.FS

	file, err := fs.Open(path)
	if err != nil {
		return false
	}
	defer utils.CallClose(file)

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return false
	}

	return bytes.Contains(content, []byte(s.script))
}

func (s *ShellProfiler) Run(ctx context.Context) (context.Context, error) {
	fs := define.FS
	logger := define.Logger
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to get user home dir")
	}

	script := s.script
	var profile string
	switch s.shell {
	case "bash":
		profile = filepath.Join(homeDir, ".bashrc")
	case "zsh":
		profile = filepath.Join(homeDir, ".zshrc")
	default:
		logger.Warn("shell is not supported, please execute this script to init cmdr environment", map[string]interface{}{
			"shell":  s.shell,
			"script": script,
		})
	}

	if s.isContainsProfile(profile) {
		return ctx, nil
	}

	file, err := fs.OpenFile(profile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to open profile file")
	}
	defer utils.CallClose(file)

	_, err = fmt.Fprintf(file, "\n%s\n", script)
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to write to profile file")
	}

	return ctx, nil
}

func NewShellProfiler(binDir, shell string) *ShellProfiler {
	return &ShellProfiler{
		shell:  filepath.Base(shell),
		script: fmt.Sprintf(`eval "$(%s init)"`, GetCommandBinPath(binDir, define.Name)),
	}
}
