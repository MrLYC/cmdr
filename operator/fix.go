package operator

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
)

type BrokenCommandsFixer struct {
	BaseOperator
	shimsDir string
}

func (s *BrokenCommandsFixer) String() string {
	return "broken-commands-fixer"
}

func (s *BrokenCommandsFixer) Run(ctx context.Context) (context.Context, error) {
	fs := define.FS
	logger := define.Logger
	var errs error
	client := GetDBClientFromContext(ctx)
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, err
	}

	availableCommands := make([]*model.Command, 0, len(commands))
	for _, command := range commands {
		location := command.Location
		if command.Managed {
			location = GetCommandShimsPath(s.shimsDir, command.Name, command.Version)
		}

		_, err := fs.Stat(location)
		if err == nil {
			availableCommands = append(availableCommands, command)
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

func NewBrokenCommandsFixer(shimsDir string) *BrokenCommandsFixer {
	return &BrokenCommandsFixer{
		shimsDir: shimsDir,
	}
}
