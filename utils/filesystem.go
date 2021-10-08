package utils

import (
	"io"
	"os"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
)

func EnsureNotExists(path string) error {
	fs := define.FS
	_, err := fs.Stat(path)
	if err == nil {
		return fs.Remove(path)
	}

	return nil
}

func CopyFile(src, dst string) error {
	fs := define.FS
	stat, err := os.Stat(src)
	if err != nil {
		return errors.Wrapf(err, "stat src file %s failed", src)
	}

	in, err := fs.Open(src)
	if err != nil {
		return errors.Wrapf(err, "open src file %s failed", src)
	}

	defer CallClose(in)

	err = EnsureNotExists(dst)
	if err != nil {
		return errors.Wrapf(err, "delete dst file %s failed", dst)
	}

	out, err := fs.OpenFile(dst, os.O_WRONLY|os.O_CREATE, stat.Mode().Perm())
	if err != nil {
		return errors.Wrapf(err, "create dst file %s failed", dst)
	}

	defer CallClose(out)

	_, err = io.Copy(out, in)

	if err != nil {
		return errors.Wrapf(err, "copy file failed")
	}

	return nil
}
