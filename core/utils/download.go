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

func NewProgressBarDownloader(stream io.Writer, options ...getter.ClientOption) *Downloader {
	tracker := NewProgressBarTracker("downloading", stream)

	return NewDownloader(
		tracker,
		[]getter.Detector{
			new(getter.GitHubDetector),
			new(getter.GitLabDetector),
			new(getter.GitDetector),
			new(getter.BitBucketDetector),
			new(getter.S3Detector),
			new(getter.GCSDetector),
		},
		options,
	)
}
