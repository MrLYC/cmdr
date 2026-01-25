package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/homedepot/flop"
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
)

type DoctorManager struct {
	binaryMgr   core.CommandManager
	databaseMgr core.CommandManager
}

func (m *DoctorManager) each(fn func(mgr core.CommandManager) error) error {
	for _, mgr := range []core.CommandManager{m.binaryMgr, m.databaseMgr} {
		err := fn(mgr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *DoctorManager) reverseEach(fn func(mgr core.CommandManager) error) error {
	for _, mgr := range []core.CommandManager{m.databaseMgr, m.binaryMgr} {
		err := fn(mgr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *DoctorManager) all(fn func(mgr core.CommandManager) error) error {
	var errs error
	for _, mgr := range []core.CommandManager{m.binaryMgr, m.databaseMgr} {
		err := fn(mgr)
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	return errs
}

func (m *DoctorManager) eachStopOnError(fn func(mgr core.CommandManager) error) error {
	for _, mgr := range []core.CommandManager{m.binaryMgr, m.databaseMgr} {
		err := fn(mgr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DoctorManager) Define(name, version, location string) (core.Command, error) {
	var command core.Command
	err := d.eachStopOnError(func(mgr core.CommandManager) error {
		cmd, err := mgr.Define(name, version, location)
		if err != nil {
			return err
		}
		command = cmd
		return nil
	})
	return command, err
}

func (d *DoctorManager) Undefine(name, version string) error {
	return d.eachStopOnError(func(mgr core.CommandManager) error {
		return mgr.Undefine(name, version)
	})
}

func (d *DoctorManager) Activate(name, version string) error {
	return d.eachStopOnError(func(mgr core.CommandManager) error {
		return mgr.Activate(name, version)
	})
}

func (d *DoctorManager) Deactivate(name string) error {
	return d.eachStopOnError(func(mgr core.CommandManager) error {
		return mgr.Deactivate(name)
	})
}

func (d *DoctorManager) Close() error {
	return d.all(func(mgr core.CommandManager) error {
		return mgr.Close()
	})
}

func (d *DoctorManager) Provider() core.CommandProvider {
	return core.CommandProviderDoctor
}

func (d *DoctorManager) Query() (core.CommandQuery, error) {
	mainQuery, mainErr := d.binaryMgr.Query()
	if mainErr != nil {
		return d.databaseMgr.Query()
	}

	recorderQuery, recorderErr := d.databaseMgr.Query()
	if recorderErr != nil {
		return mainQuery, nil
	}

	var queriedCommands []core.Command

	// merge two queries
	commands, mainErr := mainQuery.All()
	if mainErr == nil {
		queriedCommands = append(queriedCommands, commands...)
	}

	recorderCommands, recorderErr := recorderQuery.All()
	if recorderErr == nil {
		queriedCommands = append(queriedCommands, recorderCommands...)
	}

	indexes := make(map[string]int, len(queriedCommands))
	merged := make([]*Command, 0, len(queriedCommands))
	for i, cmd := range queriedCommands {
		name := cmd.GetName()
		version := cmd.GetVersion()

		key := fmt.Sprintf("%s-%s", name, version)
		index, ok := indexes[key]
		if ok {
			// update by recorder
			merged[index].Activated = cmd.GetActivated()
			continue
		}

		indexes[key] = i
		merged = append(merged, &Command{
			Name:      name,
			Version:   version,
			Activated: cmd.GetActivated(),
			Location:  cmd.GetLocation(),
		})
	}

	return NewCommandFilter(merged), nil
}

func NewDoctorManager(binaryMgr core.CommandManager, databaseMgr core.CommandManager) *DoctorManager {
	return &DoctorManager{
		binaryMgr:   binaryMgr,
		databaseMgr: databaseMgr,
	}
}

type CommandDoctor struct {
	core.CommandManager
	rootDir string
}

func checkFileAccessible(location string) (bool, error) {
	info, err := os.Stat(location)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "stat file failed: %s", location)
	}

	if info.IsDir() {
		return false, nil
	}

	absPath, err := filepath.Abs(location)
	if err != nil {
		return false, errors.Wrapf(err, "get abs path failed: %s", location)
	}
	_ = absPath

	if info.Mode()&0111 == 0 {
		return false, nil
	}

	return true, nil
}

func (d *CommandDoctor) backup() (string, error) {
	if d.rootDir == "" {
		return "", nil
	}

	info, err := os.Stat(d.rootDir)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", errors.Wrapf(err, "stat root dir failed: %s", d.rootDir)
	}
	if !info.IsDir() {
		return "", errors.Errorf("root dir is not a directory: %s", d.rootDir)
	}

	tmpDir := os.TempDir()
	rootDirBase := filepath.Base(d.rootDir)
	backupDir := filepath.Join(tmpDir, fmt.Sprintf("%s.backup.%s", rootDirBase, time.Now().Format("20060102-150405")))

	err = flop.Copy(d.rootDir, backupDir, flop.Options{
		MkdirAll:  true,
		Recursive: true,
	})
	if err != nil {
		return "", errors.Wrapf(err, "backup failed: %s -> %s", d.rootDir, backupDir)
	}

	return backupDir, nil
}

func (d *CommandDoctor) Fix(dryRun bool) error {
	return d.FixWithOptions(dryRun, true)
}

func (d *CommandDoctor) FixWithOptions(dryRun bool, backup bool) error {
	logger := core.GetLogger()

	if dryRun {
		logger.Info("running in dry-run mode, no changes will be made", nil)
	}

	if backup && !dryRun {
		backupDir, err := d.backup()
		if err != nil {
			return errors.Wrapf(err, "backup failed")
		}
		if backupDir != "" {
			logger.Info("backup created", map[string]interface{}{
				"backup_dir": backupDir,
			})
		}
	}

	query, err := d.Query()
	if err != nil {
		return errors.Wrapf(err, "make query failed")
	}

	commands, err := query.All()
	if err != nil {
		return errors.Wrapf(err, "query commands failed")
	}

	var availableCommands []core.Command
	for _, cmd := range commands {
		name := cmd.GetName()
		version := cmd.GetVersion()
		location := cmd.GetLocation()
		activated := cmd.GetActivated()

		logger.Debug("checking command", map[string]interface{}{
			"name":     name,
			"version":  version,
			"location": location,
		})

		accessible, err := checkFileAccessible(location)
		if err != nil {
			logger.Warn("check command accessibility failed, treat as unavailable", map[string]interface{}{
				"name":     name,
				"version":  version,
				"location": location,
				"error":    err.Error(),
			})
		}

		if accessible {
			logger.Debug("command is available", map[string]interface{}{
				"name":    name,
				"version": version,
			})

			availableCommands = append(availableCommands, cmd)
			continue
		}

		logger.Warn("command is not available", map[string]interface{}{
			"name":    name,
			"version": version,
		})

		if activated {
			if dryRun {
				logger.Info("[DRY-RUN] would deactivate command", map[string]interface{}{
					"name": name,
				})
			} else {
				logger.Info("deactivating command", map[string]interface{}{
					"name": name,
				})
				err = d.Deactivate(name)
				if err != nil {
					logger.Warn("deactivate command failed, try to remove it", map[string]interface{}{
						"name":  name,
						"error": err,
					})
				}
			}
		}

		if dryRun {
			logger.Info("[DRY-RUN] would remove command", map[string]interface{}{
				"name":    name,
				"version": version,
			})
		} else {
			logger.Info("removing command", map[string]interface{}{
				"name":    name,
				"version": version,
			})
			err = d.Undefine(name, version)
			if err != nil {
				logger.Error("remove command failed, continue", map[string]interface{}{
					"name":    name,
					"version": version,
					"error":   err,
				})
			}
		}
	}

	for _, cmd := range availableCommands {
		name := cmd.GetName()
		version := cmd.GetVersion()
		location := cmd.GetLocation()
		activated := cmd.GetActivated()

		// Skip re-define if the shim file already exists at the expected location.
		// The command is already available (checkFileAccessible passed), so calling
		// Define with the shim path as source would be a no-op at best, or could
		// fail when using copy mode (CopyFile deletes target before copying from source,
		// but source == target in this case).
		logger.Debug("skipping re-define for available command", map[string]interface{}{
			"name":     name,
			"version":  version,
			"location": location,
		})

		if activated {
			if dryRun {
				logger.Info("[DRY-RUN] would re-activate command", map[string]interface{}{
					"name":    name,
					"version": version,
				})
			} else {
				err = d.Activate(name, version)
				if err != nil {
					logger.Warn("re-activate command failed, continue", map[string]interface{}{
						"name":    name,
						"version": version,
					})
				}
			}
		}
	}

	if dryRun {
		logger.Info("dry-run completed, no changes were made", nil)
	}

	return nil
}

func NewCommandDoctor(manager core.CommandManager, rootDir string) *CommandDoctor {
	return &CommandDoctor{
		CommandManager: manager,
		rootDir:        rootDir,
	}
}

func init() {
	var _ core.CommandManager = (*DoctorManager)(nil)

	core.RegisterCommandManagerFactory(core.CommandProviderDoctor, func(cfg core.Configuration) (core.CommandManager, error) {
		mainMgr, err := core.NewCommandManager(core.CommandProviderBinary, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "new main command manager failed")
		}

		recorderMgr, err := core.NewCommandManager(core.CommandProviderDatabase, cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "new recorder command manager failed")
		}

		return NewDoctorManager(mainMgr, recorderMgr), nil
	})
}
