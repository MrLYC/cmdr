package core

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock Fetcher

type Fetcher interface {
	IsSupport(uri string) bool
	Fetch(name, version, uri, dir string) error
}
