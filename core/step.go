package core

import (
	"context"

	"github.com/hashicorp/go-multierror"

	"github.com/mrlyc/cmdr/define"
)

type Steper interface {
	String() string
	Run(ctx context.Context) (context.Context, error)
	Finish(ctx context.Context) error
}

type BaseStep struct{}

func (s *BaseStep) String() string {
	return ""
}

func (s *BaseStep) Run(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (s *BaseStep) Finish(ctx context.Context) error {
	return nil
}

type StepRunner struct {
	steps []Steper
}

func (r *StepRunner) Add(step ...Steper) {
	r.steps = append(r.steps, step...)
}

func (r *StepRunner) Run(ctx context.Context) (err error) {
	logger := define.Logger
	var errs error

	for _, step := range r.steps {
		logger.Debug("running step", map[string]interface{}{
			"step": step,
		})

		ctx, err = step.Run(ctx)
		if err != nil {
			logger.Debug("step failed", map[string]interface{}{
				"step":  step,
				"error": err,
			})
			return err
		}

		defer func(step Steper) {
			logger.Info("step finished", map[string]interface{}{
				"step": step,
			})

			err = step.Finish(ctx)
			if err != nil {
				logger.Debug("step error", map[string]interface{}{
					"step":  step,
					"error": err,
				})
				errs = multierror.Append(errs, err)
			}
		}(step)
	}

	return errs
}

func NewStepRunner(steps ...Steper) *StepRunner {
	return &StepRunner{
		steps: steps,
	}
}
