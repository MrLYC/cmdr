package core

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock Fetcher

type Fetcher interface {
	IsSupport(uri string) bool
	Fetch(uri, dir string) error
}
