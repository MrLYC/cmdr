package operator

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
	BaseOperator
	script string
}

func (s *ShellProfiler) String() string {
	return "shell-profiler"
}

func (s *ShellProfiler) isContainsProfile(path string) bool {
	file, err := os.Open(path)
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

func (s *ShellProfiler) getProfilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get user home dir")
	}

	profilePath := os.Getenv("CMDR_PROFILE_PATH")
	if profilePath != "" {
		return profilePath, nil
	}

	shell := filepath.Base(os.Getenv("SHELL"))

	switch shell {
	case "bash":
		profilePath = filepath.Join(homeDir, ".bashrc")
	case "zsh":
		profilePath = filepath.Join(homeDir, ".zshrc")
	default:
		return "", errors.Wrapf(ErrNotSupported, shell)
	}

	return profilePath, nil
}

func (s *ShellProfiler) Run(ctx context.Context) (context.Context, error) {
	script := s.script
	profile, err := s.getProfilePath()
	if err != nil {
		return ctx, err
	}

	if s.isContainsProfile(profile) {
		return ctx, nil
	}

	file, err := os.OpenFile(profile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to open profile file")
	}
	defer utils.CallClose(file)

	_, err = fmt.Fprintf(file, "\n%s", script)
	if err != nil {
		return ctx, errors.Wrapf(err, "failed to write to profile file")
	}

	return ctx, nil
}

func NewShellProfiler(helper *utils.CmdrHelper) *ShellProfiler {
	return &ShellProfiler{
		script: fmt.Sprintf(`eval "$(%s init)"`, helper.GetCommandBinPath(define.Name)),
	}
}
