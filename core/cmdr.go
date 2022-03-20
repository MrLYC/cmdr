package core

import (
	"context"
	"fmt"
)

//go:generate stringer -type=CmdrSearcherProvider
//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock CmdrSearcher

type CmdrReleaseAsset struct {
	Name    string
	Version string
	Asset   string
	Url     string
}

type CmdrSearcher interface {
	GetReleaseAsset(ctx context.Context, releaseName, assetName string) (CmdrReleaseAsset, error)
}

type CmdrSearcherProvider int

const (
	CmdrSearcherProviderUnknown CmdrSearcherProvider = iota
	CmdrSearcherProviderDefault
	CmdrSearcherProviderApi
	CmdrSearcherProviderAtom
)

type factoryCmdrSearcher func(cfg Configuration) (CmdrSearcher, error)

var (
	ErrCmdrSearcherFactoryeNotFound = fmt.Errorf("factory not found")
	factoriesCmdrSearcher           map[CmdrSearcherProvider]factoryCmdrSearcher
)

func RegisterCmdrSearcherFactory(provider CmdrSearcherProvider, fn func(cfg Configuration) (CmdrSearcher, error)) {
	factoriesCmdrSearcher[provider] = fn
}

func NewCmdrSearcher(provider CmdrSearcherProvider, cfg Configuration) (CmdrSearcher, error) {
	fn, ok := factoriesCmdrSearcher[provider]

	if !ok {
		return nil, ErrCmdrSearcherFactoryeNotFound
	}

	return fn(cfg)
}

func init() {
	factoriesCmdrSearcher = make(map[CmdrSearcherProvider]factoryCmdrSearcher)
}
