package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func ExtraTGZ(target, dir string) error {
	src, err := os.Open(target)
	if err != nil {
		return errors.Wrapf(err, "open %s failed", target)
	}

	defer src.Close()

	gzipReader, err := gzip.NewReader(src)
	if err != nil {
		return errors.Wrapf(err, "create gzip reader failed")
	}

	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return errors.Wrapf(err, "read failed")
			}
		}

		filename := filepath.Join(dir, header.Name)
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, header.FileInfo().Mode().Perm())
		if err != nil {
			return errors.Wrapf(err, "create file %s failed", filename)
		}
		io.Copy(file, tarReader)
	}
	return nil
}
