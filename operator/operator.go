package operator

import "context"

type BaseOperator struct{}

func (s *BaseOperator) String() string {
	return ""
}

func (s *BaseOperator) Run(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (s *BaseOperator) Commit(ctx context.Context) error {
	return nil
}

func (s *BaseOperator) Rollback(ctx context.Context) {
}
