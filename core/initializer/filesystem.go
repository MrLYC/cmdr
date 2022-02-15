package initializer

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/homedepot/flop"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

type FSBackup struct {
	path   string
	target string
}

func (b *FSBackup) Init() error {
	dir, err := os.MkdirTemp("", fmt.Sprintf("%s-backup-%s-*", core.Name, filepath.Base(b.path)))
	if err != nil {
		return errors.Wrapf(err, "failed to create backup directory for %s", b.path)
	}

	err = flop.Copy(b.path, dir, flop.Options{
		Recursive:        true,
		AppendNameToPath: true,
	})
	switch errors.Cause(err) {
	case nil:
	case flop.ErrFileNotExist:
		return nil
	default:
		return errors.WithMessagef(err, "failed to backup %s", b.path)
	}

	b.target = dir
	return nil
}

func (b *FSBackup) Target() string {
	return b.target
}

func NewFSBackup(path string) *FSBackup {
	return &FSBackup{
		path: path,
	}
}

type EmbedFSExporter struct {
	embedFS fs.FS
	srcPath string
	dstPath string
}

func (e *EmbedFSExporter) copyDir(dstPath string, perm os.FileMode) error {
	err := os.MkdirAll(dstPath, perm)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", dstPath)
	}

	return nil
}

func (e *EmbedFSExporter) copyFile(srcPath, dstPath string, perm os.FileMode) error {
	srcFile, err := e.embedFS.Open(srcPath)
	if err != nil {
		return errors.Wrap(err, "failed to open source file")
	}
	defer utils.CallClose(srcFile)

	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE, perm)
	if err != nil {
		return errors.Wrap(err, "failed to open destination file")
	}
	defer utils.CallClose(dstFile)

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.Wrap(err, "failed to copy file")
	}

	return nil
}

func (e *EmbedFSExporter) exportDir(srcPath string, d fs.DirEntry, err error) error {
	if err != nil {
		return errors.Wrap(err, "failed to walk directory")
	}

	path, err := filepath.Rel(e.srcPath, srcPath)
	if err != nil {
		return errors.Wrap(err, "failed to get relative path")
	}

	dstPath := filepath.Join(e.dstPath, path)

	fileinfo, err := d.Info()
	if err != nil {
		return errors.Wrap(err, "failed to get file info")
	}

	perm := fileinfo.Mode().Perm()
	if d.IsDir() {
		return e.copyDir(dstPath, perm)
	}

	return e.copyFile(srcPath, dstPath, perm)
}

func (e *EmbedFSExporter) Init() error {
	info, err := fs.Stat(e.embedFS, e.srcPath)
	if err != nil {
		return errors.Wrap(err, "failed to stat source path")
	}

	err = os.MkdirAll(e.dstPath, info.Mode().Perm())
	if err != nil {
		return errors.Wrap(err, "failed to create destination directory")
	}

	err = fs.WalkDir(e.embedFS, e.srcPath, e.exportDir)
	if err != nil {
		return errors.Wrap(err, "failed to walk source directory")
	}

	return nil
}

func NewEmbedFSExporter(embedFS fs.FS, srcPath, dstPath string) *EmbedFSExporter {
	return &EmbedFSExporter{
		embedFS: embedFS,
		srcPath: srcPath,
		dstPath: dstPath,
	}
}
