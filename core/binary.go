package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/homedepot/flop"
	"github.com/pkg/errors"
)

type dirManager struct {
	path string
}

func (m *dirManager) MkdirAll(mode os.FileMode) error {
	err := os.MkdirAll(m.path, mode)
	if err != nil {
		return errors.Wrapf(err, "create dir %s failed", m.path)
	}

	return nil
}

func (m *dirManager) Child(name string) *dirManager {
	return &dirManager{
		path: filepath.Join(m.path, name),
	}
}

func (m *dirManager) EnsureNotExists(name string) error {
	path := filepath.Join(m.path, name)

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

func (m *dirManager) Exists(name string) error {
	_, err := os.Lstat(filepath.Join(m.path, name))
	if err != nil {
		return errors.Wrapf(err, "stat dir %s failed", m.path)
	}

	return nil
}

func (m *dirManager) SymbolLink(name, target string, mode os.FileMode) error {
	err := m.EnsureNotExists(name)
	if err != nil {
		return err
	}

	path := filepath.Join(m.path, name)
	err = os.Symlink(target, path)
	if err != nil {
		return errors.Wrapf(err, "create symbol link failed")
	}

	err = os.Chmod(path, mode)
	if err != nil {
		return errors.Wrapf(err, "chmod %s failed", path)
	}

	return nil
}

func (m *dirManager) CopyFile(name, target string, mode os.FileMode) error {
	err := m.EnsureNotExists(name)
	if err != nil {
		return err
	}

	path := filepath.Join(m.path, name)
	err = flop.Copy(target, path, flop.Options{
		MkdirAll:  true,
		Recursive: true,
	})
	if err != nil {
		return errors.Wrapf(err, "copy file failed")
	}

	err = os.Chmod(path, mode)
	if err != nil {
		return errors.Wrapf(err, "chmod %s failed", path)
	}

	return nil
}

func (m *dirManager) RealPath(name string) (string, error) {
	path, err := m.AbsPath(name)
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

func (m *dirManager) AbsPath(name string) (string, error) {
	path := filepath.Join(m.path, name)

	err := m.Exists(name)
	if err != nil {
		return "", errors.Wrapf(err, "path %s not found", path)
	}

	realPath, err := filepath.Abs(path)
	if err != nil {
		return "", errors.Wrapf(err, "get real path %s failed", path)
	}

	return realPath, nil
}

func (m *dirManager) Path() string {
	return m.path
}

func newDirManager(path string) *dirManager {
	return &dirManager{
		path: path,
	}
}

type BinaryManager struct {
	BinManager   *dirManager
	ShimsManager *dirManager
	dirMode      os.FileMode
}

func (m *BinaryManager) Init() error {
	for _, mgr := range []*dirManager{m.BinManager, m.ShimsManager} {
		err := mgr.MkdirAll(m.dirMode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *BinaryManager) ShimsName(name, version string) string {
	return fmt.Sprintf("%s_%s", name, version)
}

func (m *BinaryManager) Install(name, version, location string, symlink bool) error {
	mgr := m.ShimsManager.Child(name)

	err := mgr.MkdirAll(m.dirMode)
	if err != nil {
		return errors.WithMessagef(err, "create dir %s failed", mgr.path)
	}

	shimsName := m.ShimsName(name, version)

	if symlink {
		err = mgr.SymbolLink(shimsName, location, 0755)
	} else {
		err = mgr.CopyFile(shimsName, location, 0755)
	}

	if err != nil {
		return errors.WithMessagef(err, "install %s failed", location)
	}

	return nil
}

func (m *BinaryManager) Uninstall(name, version string) error {
	mgr := m.ShimsManager.Child(name)
	shimsName := m.ShimsName(name, version)

	err := mgr.EnsureNotExists(shimsName)
	if err != nil {
		return err
	}

	return nil
}

func (m *BinaryManager) Activate(name, version string) error {
	mgr := m.ShimsManager.Child(name)
	shimsName := m.ShimsName(name, version)

	path, err := mgr.AbsPath(shimsName)
	if err != nil {
		return errors.WithMessagef(err, "get shims %s failed", shimsName)
	}

	err = m.BinManager.SymbolLink(name, path, 0755)
	if err != nil {
		return errors.WithMessagef(err, "symlink %s failed", path)
	}

	return nil
}

func (m *BinaryManager) Deactivate(name string) error {
	err := m.BinManager.EnsureNotExists(name)
	if err != nil {
		return errors.Wrapf(err, "remove %s failed", name)
	}

	return nil
}

func NewBinaryManager(root string) *BinaryManager {
	bin := newDirManager(filepath.Join(root, "bin"))
	shims := newDirManager(filepath.Join(root, "shims"))

	return &BinaryManager{
		BinManager:   bin,
		ShimsManager: shims,
		dirMode:      0755,
	}
}
