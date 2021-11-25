package core

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
)

type commandRunner struct {
	command string
	args    []string
}

func (r *commandRunner) run(ctx context.Context) error {
	logger := define.Logger

	process := exec.CommandContext(ctx, r.command, r.args...)
	logger.Debug("running process", map[string]interface{}{
		"process": process,
	})

	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	err := process.Run()
	if err != nil {
		return errors.Wrapf(err, "failed to run subprocess %v", process)
	}

	return nil
}

func newCommandRunner(command string, args ...string) *commandRunner {
	return &commandRunner{command: command, args: args}
}

type FinalizeCommandRunner struct {
	BaseStep
	*commandRunner
}

func (s *FinalizeCommandRunner) String() string {
	return "finalize-command-runner"
}

func (s *FinalizeCommandRunner) Finish(ctx context.Context) error {
	logger := define.Logger
	logger.Info("executing cmdr setup command")
	return s.run(ctx)
}

func NewUpgradeSetupRunner(args ...string) *FinalizeCommandRunner {
	return &FinalizeCommandRunner{commandRunner: newCommandRunner(filepath.Join(GetBinDir(), define.Name), args...)}
}
