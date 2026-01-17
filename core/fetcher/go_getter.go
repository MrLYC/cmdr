package fetcher

import (
	"io"
	"os"
	"sync"

	"github.com/hashicorp/go-getter"
	"github.com/mrlyc/cmdr/core/utils"
	"github.com/pkg/errors"
)

//go:generate mockgen -destination=mock/go_getter.go -package=mock github.com/hashicorp/go-getter Detector,ProgressTracker

type GoGetter struct {
	progressListener getter.ProgressTracker
	detectors        []getter.Detector
	options          []getter.ClientOption
	optionsMutex     sync.RWMutex
}

func (d *GoGetter) IsSupport(uri string) bool {
	_, err := getter.Detect(uri, os.TempDir(), d.detectors)
	return err == nil
}

func (d *GoGetter) Fetch(name, version, uri, dst string) error {
	d.optionsMutex.RLock()
	options := d.options
	d.optionsMutex.RUnlock()

	client := getter.Client{
		Src:              uri,
		Dst:              dst,
		Pwd:              os.TempDir(),
		Mode:             getter.ClientModeAny,
		Detectors:        d.detectors,
		Options:          options,
		ProgressListener: d.progressListener,
	}

	err := client.Get()
	if err != nil {
		return errors.Wrapf(err, "download failed")
	}

	return nil
}

func (d *GoGetter) SetOptions(options []getter.ClientOption) {
	d.optionsMutex.Lock()
	defer d.optionsMutex.Unlock()
	d.options = options
}

func NewGoGetter(progressListener getter.ProgressTracker, detectors []getter.Detector, options []getter.ClientOption) *GoGetter {
	return &GoGetter{
		progressListener: progressListener,
		detectors:        detectors,
		options:          options,
	}
}

func NewDefaultGoGetter(stream io.Writer, options ...getter.ClientOption) *GoGetter {
	tracker := utils.NewProgressBarTracker("downloading", stream)
	detectors := make([]getter.Detector, 0, len(getter.Detectors))
	for _, d := range getter.Detectors {
		switch d.(type) {
		case *getter.FileDetector:
			// it is no need to download a local file
			continue
		}

		detectors = append(detectors, d)
	}

	return NewGoGetter(tracker, detectors, options)
}
