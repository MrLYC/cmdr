package utils

import (
	"io"
	"os"

	"github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
)

type Downloader struct {
	progressListener getter.ProgressTracker
	detectors        []getter.Detector
	options          []getter.ClientOption
}

func (d *Downloader) IsSupport(uri string) bool {
	_, err := getter.Detect(uri, os.TempDir(), d.detectors)
	return err == nil
}

func (d *Downloader) Fetch(uri, dst string) error {
	client := getter.Client{
		Src:              uri,
		Dst:              dst,
		Pwd:              os.TempDir(),
		Mode:             getter.ClientModeAny,
		Detectors:        d.detectors,
		Options:          d.options,
		ProgressListener: d.progressListener,
	}

	err := client.Get()
	if err != nil {
		return errors.Wrapf(err, "download failed")
	}

	return nil
}

func NewDownloader(progressListener getter.ProgressTracker, detectors []getter.Detector, options []getter.ClientOption) *Downloader {
	return &Downloader{
		progressListener: progressListener,
		detectors:        detectors,
		options:          options,
	}
}

func NewDefaultDownloader(stream io.Writer, options ...getter.ClientOption) *Downloader {
	tracker := NewProgressBarTracker("downloading", stream)
	detectors := make([]getter.Detector, 0, len(getter.Detectors))
	for _, d := range getter.Detectors {
		switch d.(type) {
		case *getter.FileDetector:
			// it is no need to download a local file
			continue
		}

		detectors = append(detectors, d)
	}

	return NewDownloader(tracker, detectors, options)
}
