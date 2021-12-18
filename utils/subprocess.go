package utils

import (
	"context"
	"os"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
)

func WaitProcess(ctx context.Context, command string, args []string) error {
	logger := define.Logger
	process := exec.CommandContext(ctx, command, args...)
	logger.Debug("running process", map[string]interface{}{
		"process": process,
	})

	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	err := process.Run()
	if err != nil {
		return errors.Wrapf(err, "run process failed")
	}

	return nil
}
