package operator

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
)

type BrokenCommandsFixer struct {
	*CmdrOperator
}

func (i *BrokenCommandsFixer) String() string {
	return "broken-commands-fixer"
}

func (i *BrokenCommandsFixer) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger

	var errs error
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	for _, command := range commands {
		shimsName := i.cmdr.BinaryManager.ShimsName(command.Name, command.Version)
		if i.cmdr.BinaryManager.ShimsManager.Exists(shimsName) == nil {
			continue
		}

		logger.Debug("deleting command", map[string]interface{}{
			"name":     command.Name,
			"version":  command.Version,
			"location": command.Location,
			"err":      err,
		})

		err = i.cmdr.CommandManager.Delete(command)
		if err != nil {
			errs = multierror.Append(errs, errors.Wrapf(err, "remove command %s(%s) failed", command.Name, command.Version))
			continue
		}
		logger.Info("command deleted", map[string]interface{}{
			"name":    command.Name,
			"version": command.Version,
		})

	}

	return ctx, errs
}

func NewBrokenCommandsFixer(cmdr *core.Cmdr) *BrokenCommandsFixer {
	return &BrokenCommandsFixer{
		CmdrOperator: NewCmdrOperator(cmdr),
	}
}
