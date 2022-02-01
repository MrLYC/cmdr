package operator

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type BrokenCommandsFixer struct {
	BaseOperator
}

func (s *BrokenCommandsFixer) String() string {
	return "broken-commands-fixer"
}

func (s *BrokenCommandsFixer) Run(ctx context.Context) (context.Context, error) {
	fs := define.FS
	logger := define.Logger
	shimsDir := GetShimsDir(ctx)

	var errs error
	client := GetDBClientFromContext(ctx)
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	for _, command := range commands {
		location := command.Location
		if command.Managed {
			location = utils.GetCommandShimsPath(shimsDir, command.Name, command.Version)
		}

		_, err := fs.Stat(location)
		if err == nil {
			continue
		}

		logger.Debug("deleting command", map[string]interface{}{
			"name":     command.Name,
			"version":  command.Version,
			"location": command.Location,
			"err":      err,
		})

		err = client.DeleteStruct(command)
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

func NewBrokenCommandsFixer() *BrokenCommandsFixer {
	return &BrokenCommandsFixer{}
}
