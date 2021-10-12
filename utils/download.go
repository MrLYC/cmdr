package utils

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"

	"github.com/mrlyc/cmdr/define"
)

func DownloadToFile(ctx context.Context, url, output string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return errors.Wrapf(err, "create request failed")
	}

	f, err := define.FS.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrapf(err, "create output file failed")
	}
	defer f.Close()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "request failed")
	}
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(-1, "Downloading")
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		return errors.Wrapf(err, "download failed")
	}

	return nil
}
