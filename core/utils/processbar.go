package utils

import (
	"fmt"
	"io"
	"time"

	"github.com/schollz/progressbar/v3"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock progressBar
//go:generate mockgen -destination=mock/io.go -package=mock io ReadCloser

type progressBar interface {
	Add(num int) error
	Finish() error
}

type ProgressBar struct {
	io.ReadCloser
	bar progressBar
}

func (p *ProgressBar) Read(b []byte) (n int, err error) {
	n, err = p.ReadCloser.Read(b)
	_ = p.bar.Add(n)
	return n, err
}

func (p *ProgressBar) Close() error {
	_ = p.bar.Finish()
	return p.ReadCloser.Close()
}

func NewProgressBar(r io.ReadCloser, bar progressBar) *ProgressBar {
	return &ProgressBar{
		ReadCloser: r,
		bar:        bar,
	}
}

type ProgressBarTracker struct {
	description string
	stream      io.Writer
}

func (t *ProgressBarTracker) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) (body io.ReadCloser) {
	bar := progressbar.NewOptions64(
		-1,
		progressbar.OptionSetDescription(t.description),
		progressbar.OptionSetWriter(t.stream),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprintln(t.stream, "")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)

	_ = bar.RenderBlank()
	_ = bar.Set(int(currentSize))

	return NewProgressBar(stream, bar)
}

func NewProgressBarTracker(description string, stream io.Writer) *ProgressBarTracker {
	return &ProgressBarTracker{
		description: description,
		stream:      stream,
	}
}
