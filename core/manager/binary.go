package manager

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	. "github.com/ahmetb/go-linq/v3"
	ver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

type Binary struct {
	binDir    string
	shimsDir  string
	name      string
	version   string
	shimsName string
}

func (b *Binary) String() string {
	return b.GetLocation()
}

func (b *Binary) GetName() string {
	return b.name
}

func (b *Binary) GetVersion() string {
	return b.version
}

func (b *Binary) GetActivated() bool {
	binHelper := utils.NewPathHelper(b.binDir)
	binPath, err := binHelper.RealPath(b.name)
	if err != nil {
		return false
	}

	return binPath == b.GetLocation()
}

func (b *Binary) GetLocation() string {
	return utils.NewPathHelper(b.shimsDir).
		Child(b.name).
		Child(b.shimsName).
		Path()
}

func NewBinary(binDir, shimsDir, name, version, shimsName string) *Binary {
	return &Binary{
		binDir:    binDir,
		shimsDir:  shimsDir,
		name:      name,
		version:   version,
		shimsName: shimsName,
	}
}

type BinariesFilter struct {
	binaries []*Binary
}

func (f *BinariesFilter) Filter(fn func(b interface{}) bool) *BinariesFilter {
	From(f.binaries).Where(fn).ToSlice(&f.binaries)
	return f
}

func (f *BinariesFilter) WithName(name string) core.CommandQuery {
	return f.Filter(func(b interface{}) bool {
		return b.(*Binary).GetName() == name
	})
}

func (f *BinariesFilter) WithVersion(version string) core.CommandQuery {
	semver := ver.Must(ver.NewVersion(version))

	return f.Filter(func(b interface{}) bool {
		binVer := ver.Must(ver.NewVersion(b.(*Binary).GetVersion()))
		return semver.Equal(binVer)
	})
}

func (f *BinariesFilter) WithActivated(activated bool) core.CommandQuery {
	return f.Filter(func(b interface{}) bool {
		return b.(*Binary).GetActivated() == activated
	})
}

func (f *BinariesFilter) WithLocation(location string) core.CommandQuery {
	return f.Filter(func(b interface{}) bool {
		return b.(*Binary).GetLocation() == location
	})
}

func (f *BinariesFilter) All() ([]core.Command, error) {
	commands := make([]core.Command, 0, len(f.binaries))
	for _, b := range f.binaries {
		commands = append(commands, b)
	}

	return commands, nil
}

func (f *BinariesFilter) One() (core.Command, error) {
	if len(f.binaries) == 0 {
		return nil, errors.Wrapf(core.ErrBinaryNotFound, "binaries not found")
	}

	return f.binaries[0], nil
}

func (f *BinariesFilter) Count() (int, error) {
	return len(f.binaries), nil
}

func NewBinariesFilter(binaries []*Binary) *BinariesFilter {
	return &BinariesFilter{binaries}
}

type BinaryManager struct {
	binDir   string
	shimsDir string
	dirMode  os.FileMode
	linkFn   func(shimsHelper *utils.PathHelper, source, shimsName string, mode os.FileMode) error
}

func (m *BinaryManager) Init(isUpgrade bool) error {
	logger := core.GetLogger()
	logger.Debug("creating directory", map[string]interface{}{
		"bin_dir":   m.binDir,
		"shims_dir": m.shimsDir,
	})

	for _, path := range []string{m.binDir, m.shimsDir} {

		helper := utils.NewPathHelper(path)
		err := helper.MkdirAll(m.dirMode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *BinaryManager) Close() error {
	return nil
}

func (m *BinaryManager) Provider() core.CommandProvider {
	return core.CommandProviderBinary
}

func (m *BinaryManager) ShimsName(name, version string) string {
	return fmt.Sprintf("%s_%s", name, version)
}

func (m *BinaryManager) GetNormalizedVersion(version string) string {
	semver := ver.Must(ver.NewVersion(version))
	return semver.String()
}

func (m *BinaryManager) getAllPossibleShimsNames(name, version string) []string {
	shimsDir := filepath.Join(m.shimsDir, name)
	normalizedVersion := m.GetNormalizedVersion(version)

	return []string{
		filepath.Join(shimsDir, fmt.Sprintf("%s_%s", name, version)),
		filepath.Join(shimsDir, fmt.Sprintf("%s_%s", name, normalizedVersion)),
	}
}

func (m *BinaryManager) Query() (core.CommandQuery, error) {
	var binaries []*Binary

	err := filepath.Walk(m.shimsDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to walk %s", path)
		}

		if info.IsDir() {
			return nil
		}

		dir, filename := filepath.Split(path)
		name := filepath.Base(dir)
		if !strings.HasPrefix(filename, name) {
			return nil
		}

		version := strings.TrimPrefix(filename, name+"_")

		bin := NewBinary(m.binDir, m.shimsDir, name, version, filename)
		binaries = append(binaries, bin)

		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to query binaries")
	}

	return NewBinariesFilter(binaries), nil
}

func (m *BinaryManager) linkBinary(name string, version string, location string, shimsName string) error {
	logger := core.GetLogger()
	shimsHelper := utils.NewPathHelper(m.shimsDir).Child(name)

	err := shimsHelper.MkdirAll(m.dirMode)
	if err != nil {
		return errors.WithMessagef(err, "create dir %s failed", shimsHelper.Path())
	}

	logger.Debug("defining binary", map[string]interface{}{
		"name":     name,
		"version":  version,
		"location": location,
	})

	srcLocation, err := filepath.Abs(location)
	if err != nil {
		return errors.Wrapf(err, "get source path of %s failed", location)
	}

	err = m.linkFn(shimsHelper, srcLocation, shimsName, 0755)
	if err != nil {
		return errors.WithMessagef(err, "link %s to %s failed", location, shimsName)
	}

	return nil
}

func (m *BinaryManager) getNormalizedShimsName(name, version string) string {
	semver := ver.Must(ver.NewVersion(version))
	return fmt.Sprintf("%s_%s", name, semver.String())
}

func (m *BinaryManager) Define(name string, version string, location string) (core.Command, error) {
	normalizedShimsName := m.getNormalizedShimsName(name, version)
	oldShimsName := m.ShimsName(name, version)

	shimsHelper := utils.NewPathHelper(m.shimsDir).Child(name)

	oldPath := shimsHelper.Child(oldShimsName).Path()
	normalizedPath := shimsHelper.Child(normalizedShimsName).Path()

	var shimsName string
	_, err := os.Stat(normalizedPath)
	if err == nil {
		shimsName = normalizedShimsName
	} else if !os.IsNotExist(err) {
		return nil, err
	} else {
		_, err = os.Stat(oldPath)
		if err == nil {
			if oldShimsName != normalizedShimsName {
				core.GetLogger().Info("migrating old version format to new", map[string]interface{}{
					"name":        name,
					"old_version": version,
					"new_version": normalizedShimsName[len(name)+1:],
				})
				if err := os.Rename(oldPath, normalizedPath); err != nil {
					core.GetLogger().Warn("failed to migrate file, using old format", map[string]interface{}{
						"old_path": oldPath,
						"new_path": normalizedPath,
						"error":    err.Error(),
					})
					shimsName = oldShimsName
				} else {
					shimsName = normalizedShimsName
				}
			} else {
				shimsName = oldShimsName
			}
		} else if !os.IsNotExist(err) {
			return nil, errors.Wrapf(err, "stat old path %s failed", oldPath)
		} else {
			shimsName = normalizedShimsName
		}
	}

	err = m.linkBinary(name, version, location, shimsName)
	if err != nil {
		return nil, err
	}

	return NewBinary(m.binDir, m.shimsDir, name, version, shimsName), nil
}

func (m *BinaryManager) Undefine(name string, version string) error {
	helper := utils.NewPathHelper(m.shimsDir).Child(name)
	normalizedShimsName := m.getNormalizedShimsName(name, version)
	oldShimsName := m.ShimsName(name, version)

	core.GetLogger().Debug("undefining binary", map[string]interface{}{
		"name":    name,
		"version": version,
	})

	for _, shimsName := range []string{normalizedShimsName, oldShimsName} {
		if shimsName != "" {
			err := helper.EnsureNotExists(shimsName)
			if err != nil {
				core.GetLogger().Warn("failed to remove shims file", map[string]interface{}{
					"shims_name": shimsName,
					"error":      err.Error(),
				})
			}
		}
	}

	return nil
}

func (m *BinaryManager) Activate(name, version string) error {
	shimsHelper := utils.NewPathHelper(m.shimsDir).Child(name)
	normalizedShimsName := m.getNormalizedShimsName(name, version)
	oldShimsName := m.ShimsName(name, version)

	var path string
	var err error

	for _, shimsName := range []string{normalizedShimsName, oldShimsName} {
		if shimsName == "" {
			continue
		}
		path, err = shimsHelper.AbsPath(shimsName)
		if err == nil {
			break
		}
	}

	if err != nil {
		return errors.WithMessagef(err, "get shims failed")
	}

	binHelper := utils.NewPathHelper(m.binDir)

	core.GetLogger().Debug("activating binary", map[string]interface{}{
		"name":    name,
		"version": version,
	})

	err = binHelper.SymbolLink(name, path, 0755)
	if err != nil {
		return errors.WithMessagef(err, "symlink %s failed", path)
	}

	return nil
}

func (m *BinaryManager) Deactivate(name string) error {
	binHelper := utils.NewPathHelper(m.binDir)

	core.GetLogger().Debug("deactivating binary", map[string]interface{}{
		"name": name,
	})

	err := binHelper.EnsureNotExists(name)
	if err != nil {
		return errors.Wrapf(err, "remove %s failed", name)
	}

	return nil
}

func (m *BinaryManager) GetShimsDir() string {
	return m.shimsDir
}

func GetNormalizedVersion(version string) string {
	semver := ver.Must(ver.NewVersion(version))
	return semver.String()
}

func NewBinaryManager(
	binDir, shimsDir string,
	dirMode os.FileMode,
	linkFn func(shimsHelper *utils.PathHelper, source, shimsName string, mode os.FileMode) error,
) *BinaryManager {
	return &BinaryManager{binDir, shimsDir, dirMode, linkFn}
}

func NewBinaryManagerWithCopy(
	binDir, shimsDir string,
	dirMode os.FileMode,
) *BinaryManager {
	return NewBinaryManager(
		binDir, shimsDir, dirMode,
		func(shimsHelper *utils.PathHelper, source, shimsName string, mode os.FileMode) error {
			return shimsHelper.CopyFile(shimsName, source, mode)
		},
	)
}

func NewBinaryManagerWithLink(
	binDir, shimsDir string,
	dirMode os.FileMode,
) *BinaryManager {
	return NewBinaryManager(
		binDir, shimsDir, dirMode,
		func(shimsHelper *utils.PathHelper, source, shimsName string, mode os.FileMode) error {
			return shimsHelper.SymbolLink(shimsName, source, mode)
		},
	)
}

func newBinaryManagerByConfiguration(cfg core.Configuration) *BinaryManager {
	binDir := cfg.GetString(core.CfgKeyCmdrBinDir)
	shimsDir := cfg.GetString(core.CfgKeyCmdrShimsDir)

	switch cfg.GetString(core.CfgKeyCmdrLinkMode) {
	case "link":
		return NewBinaryManagerWithLink(binDir, shimsDir, 0755)
	default:
		return NewBinaryManagerWithCopy(binDir, shimsDir, 0755)
	}
}

func init() {
	var (
		_ core.Command        = (*Binary)(nil)
		_ core.CommandQuery   = (*BinariesFilter)(nil)
		_ core.CommandManager = (*BinaryManager)(nil)
		_ core.Initializer    = (*BinaryManager)(nil)
	)

	core.RegisterCommandManagerFactory(core.CommandProviderBinary, func(cfg core.Configuration) (core.CommandManager, error) {
		return newBinaryManagerByConfiguration(cfg), nil
	})

	core.RegisterInitializerFactory("binary", func(cfg core.Configuration) (core.Initializer, error) {
		return newBinaryManagerByConfiguration(cfg), nil
	})
}
