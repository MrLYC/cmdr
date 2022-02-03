package define

import "context"

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock Operator

type Operator interface {
	String() string
	Run(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context)
}
