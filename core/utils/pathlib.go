package utils

import (
	"os"
	"path/filepath"

	"github.com/homedepot/flop"
	"github.com/pkg/errors"
)

type PathHelper struct {
	path string
}

func (p *PathHelper) MkdirAll(mode os.FileMode) error {
	err := os.MkdirAll(p.path, mode)
	if err != nil {
		return errors.Wrapf(err, "create dir %s failed", p.path)
	}

	return nil
}

func (p *PathHelper) Child(name string) *PathHelper {
	return &PathHelper{
		path: filepath.Join(p.path, name),
	}
}

func (p *PathHelper) Parent() *PathHelper {
	return &PathHelper{
		path: filepath.Dir(p.path),
	}
}

func (p *PathHelper) EnsureNotExists(name string) error {
	path := filepath.Join(p.path, name)

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return errors.Wrapf(err, "stat dir %s failed", path)
	}

	if info.IsDir() {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}

	if err != nil {
		return errors.Wrapf(err, "remove dir %s failed", path)
	}

	return nil
}

func (p *PathHelper) Exists(name string) error {
	_, err := os.Lstat(filepath.Join(p.path, name))
	if err != nil {
		return errors.Wrapf(err, "stat dir %s failed", p.path)
	}

	return nil
}

func (p *PathHelper) Chmod(name string, mode os.FileMode) error {
	path := filepath.Join(p.path, name)

	err := os.Chmod(path, mode)
	if err != nil {
		return errors.Wrapf(err, "chmod %s failed", path)
	}

	return nil
}

func (p *PathHelper) SymbolLink(name, target string, mode os.FileMode) error {
	err := p.EnsureNotExists(name)
	if err != nil {
		return err
	}

	path := filepath.Join(p.path, name)
	err = os.Symlink(target, path)
	if err != nil {
		return errors.Wrapf(err, "create symbol link failed")
	}

	return p.Chmod(name, mode)
}

func (p *PathHelper) CopyFile(name, target string, mode os.FileMode) error {
	err := p.EnsureNotExists(name)
	if err != nil {
		return err
	}

	path := filepath.Join(p.path, name)
	err = flop.Copy(target, path, flop.Options{
		MkdirAll:  true,
		Recursive: true,
	})
	if err != nil {
		return errors.WithMessagef(err, "copy file failed")
	}

	return p.Chmod(name, mode)
}

func (p *PathHelper) RealPath(name string) (string, error) {
	path, err := p.AbsPath(name)
	if err != nil {
		return "", errors.Wrapf(err, "get real path %s failed", name)
	}

	info, err := os.Lstat(path)
	if err != nil {
		return "", errors.Wrapf(err, "stat dir %s failed", path)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return path, nil
	}

	location, err := os.Readlink(path)
	if err != nil {
		return "", errors.Wrapf(err, "read link %s failed", path)
	}

	realPath, err := filepath.Abs(location)
	if err != nil {
		return "", errors.Wrapf(err, "get real path %s failed", path)
	}

	return realPath, nil
}

func (p *PathHelper) AbsPath(name string) (string, error) {
	path := filepath.Join(p.path, name)

	err := p.Exists(name)
	if err != nil {
		return "", errors.Wrapf(err, "path %s not found", path)
	}

	realPath, err := filepath.Abs(path)
	if err != nil {
		return "", errors.Wrapf(err, "get real path %s failed", path)
	}

	return realPath, nil
}

func (p *PathHelper) Path() string {
	return p.path
}

func NewPathHelper(path string) *PathHelper {
	return &PathHelper{
		path: path,
	}
}
