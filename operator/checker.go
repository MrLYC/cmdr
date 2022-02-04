package operator

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/define"
)

type BinariesChecker struct {
	BaseOperator
}

func (c *BinariesChecker) String() string {
	return "binaries-checker"
}

func (c *BinariesChecker) Run(ctx context.Context) (context.Context, error) {
	fs := define.FS
	commands, err := GetCommandsFromContext(ctx)
	if err != nil {
		return ctx, errors.Wrapf(err, "get commands from context failed")
	}

	var errs error
	for _, command := range commands {
		exists, err := afero.Exists(fs, command.Location)
		if err != nil {
			errs = multierror.Append(errs, errors.WithMessage(err, command.Location))
			continue
		}

		if !exists {
			errs = multierror.Append(errs, errors.WithMessagef(ErrLocationNotExists, "%s not exists", command.Location))
			continue
		}
	}

	return ctx, errs
}

func NewBinariesChecker() *BinariesChecker {
	return &BinariesChecker{}
}
