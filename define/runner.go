package define

import "context"

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock Runner

type Runner interface {
	Run(ctx context.Context) (errs error)
}
