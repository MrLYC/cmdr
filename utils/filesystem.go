package utils

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

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
		return errors.Wrapf(err, "stat source file %s failed", src)
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

func GetSymbolLinker() afero.Linker {
	linker, ok := define.FS.(afero.Linker)
	if !ok {
		return nil
	}

	return linker
}

func GetSymbolLinkReader() afero.LinkReader {
	reader, ok := define.FS.(afero.LinkReader)
	if !ok {
		return nil
	}

	return reader
}

func GetFsLstater() afero.Lstater {
	lister, ok := define.FS.(afero.Lstater)
	if !ok {
		return nil
	}

	return lister
}

func GetRealPath(path string) string {
	linkReader := GetSymbolLinkReader()

	if linkReader == nil {
		return path
	}

	realPath, err := linkReader.ReadlinkIfPossible(path)
	if err != nil {
		return path
	}

	return realPath
}
