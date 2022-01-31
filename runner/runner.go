package runner

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
	op "github.com/mrlyc/cmdr/operator"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock Runner

type Runner interface {
	Run(ctx context.Context) (errs error)
}

type OperatorRunner struct {
	operators []op.Operator
}

func (r *OperatorRunner) Add(operator ...op.Operator) Runner {
	r.operators = append(r.operators, operator...)
	return r
}

func (r *OperatorRunner) Layout() []string {
	var layout []string
	for _, operator := range r.operators {
		layout = append(layout, operator.String())
	}
	return layout
}

func (r *OperatorRunner) Run(ctx context.Context) (errs error) {
	logger := define.Logger
	failed := false
	var err error

	for _, operator := range r.operators {
		logger.Debug("running operator", map[string]interface{}{
			"operator": operator,
		})

		ctx, err = operator.Run(ctx)
		if err != nil {
			logger.Debug("operator failed", map[string]interface{}{
				"operator": operator,
				"error":    err,
			})
			failed = true
			errs = multierror.Append(errs, errors.WithMessagef(err, "run on operator %s", operator))
			break
		}

		defer func(operator op.Operator) {
			if failed {
				logger.Warn("operator rollback", map[string]interface{}{
					"operator": operator,
				})
				operator.Rollback(ctx)
				return
			}

			logger.Debug("operator finished", map[string]interface{}{
				"operator": operator,
			})
			err := operator.Commit(ctx)
			if err != nil {
				logger.Debug("operator error", map[string]interface{}{
					"operator": operator,
					"error":    err,
				})
				errs = multierror.Append(errs, errors.WithMessagef(err, "commit on operator %s", operator))
			}
		}(operator)
	}

	return errs
}

func New(operators ...op.Operator) *OperatorRunner {
	return &OperatorRunner{
		operators: operators,
	}
}
