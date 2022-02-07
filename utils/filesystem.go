package utils

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

func EnsureNotExists(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return os.Remove(path)
	}

	return nil
}

func CopyFile(src, dst string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return errors.Wrapf(err, "stat source file %s failed", src)
	}

	in, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "open src file %s failed", src)
	}

	defer CallClose(in)

	err = EnsureNotExists(dst)
	if err != nil {
		return errors.Wrapf(err, "delete dst file %s failed", dst)
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, stat.Mode().Perm())
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

func GetRealPath(path string) string {
	realPath, err := os.Readlink(path)
	if err != nil {
		return path
	}

	return realPath
}
