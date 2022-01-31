package operator

import "context"

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock Operator

type Operator interface {
	String() string
	Run(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
}

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
