package core

import (
	"context"

	"github.com/hashicorp/go-multierror"

	"github.com/mrlyc/cmdr/define"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock Steper

type Steper interface {
	String() string
	Run(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
}

type BaseStep struct{}

func (s *BaseStep) String() string {
	return ""
}

func (s *BaseStep) Run(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (s *BaseStep) Commit(ctx context.Context) error {
	return nil
}

func (s *BaseStep) Rollback(ctx context.Context) {
}

type StepRunner struct {
	steps []Steper
}

func (r *StepRunner) Add(step ...Steper) *StepRunner {
	r.steps = append(r.steps, step...)
	return r
}

func (r *StepRunner) Layout() []string {
	var layout []string
	for _, step := range r.steps {
		layout = append(layout, step.String())
	}
	return layout
}

func (r *StepRunner) Run(ctx context.Context) (errs error) {
	logger := define.Logger
	failed := false
	var err error

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
			failed = true
			errs = multierror.Append(errs, err)
			break
		}

		defer func(step Steper) {
			if failed {
				logger.Warn("step rollback", map[string]interface{}{
					"step": step,
				})
				step.Rollback(ctx)
				return
			}

			logger.Debug("step finished", map[string]interface{}{
				"step": step,
			})
			err := step.Commit(ctx)
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
