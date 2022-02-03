package define

import storm "github.com/asdine/storm/v3"

//go:generate mockgen -destination=mock/storm.go -package=mock github.com/asdine/storm/v3 Query
//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock DBClient

type DBClient interface {
	storm.TypeStore
	Close() error
}
