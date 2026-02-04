package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

type cleanCandidate struct {
	name     string
	version  string
	location string
	addedAt  time.Time
}

func defaultCleanTrashDir() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", errors.Wrap(err, "get user home dir")
		}
		return filepath.Join(home, ".Trash", "cmdr-cleaned"), nil
	case "linux":
		return filepath.Join(string(os.PathSeparator), "tmp", "cmdr-cleaned"), nil
	default:
		return filepath.Join(os.TempDir(), "cmdr-cleaned"), nil
	}
}

func ensureDir(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return errors.Wrapf(err, "create dir %s failed", dir)
	}
	return nil
}

func uniquePath(dir, base string) (string, error) {
	dst := filepath.Join(dir, base)
	_, statErr := os.Lstat(dst)
	if os.IsNotExist(statErr) {
		return dst, nil
	}
	if statErr != nil {
		return "", errors.Wrapf(statErr, "stat %s failed", dst)
	}

	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	for i := 1; ; i++ {
		candidate := filepath.Join(dir, fmt.Sprintf("%s-%d%s", name, i, ext))
		_, candidateErr := os.Lstat(candidate)
		if os.IsNotExist(candidateErr) {
			return candidate, nil
		}
		if candidateErr != nil {
			return "", errors.Wrapf(candidateErr, "stat %s failed", candidate)
		}
	}
}

func moveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}
	if linkErr, ok := err.(*os.LinkError); !ok || linkErr.Err != syscall.EXDEV {
		return errors.Wrapf(err, "rename %s to %s failed", src, dst)
	}

	// Cross-device rename. Copy then remove.
	in, err := os.Open(src)
	if err != nil {
		return errors.Wrapf(err, "open %s failed", src)
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return errors.Wrapf(err, "stat %s failed", src)
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return errors.Wrapf(err, "create %s failed", dst)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return errors.Wrapf(err, "copy %s to %s failed", src, dst)
	}

	if err := os.Remove(src); err != nil {
		return errors.Wrapf(err, "remove %s failed", src)
	}

	return nil
}

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean old inactive command versions",
	Run: utils.RunCobraCommandWith(core.CommandProviderDefault, func(cfg core.Configuration, manager core.CommandManager) error {
		logger := core.GetLogger()

		ageDays := cfg.GetInt(core.CfgKeyXCleanAgeDays)
		keep := cfg.GetInt(core.CfgKeyXCleanKeep)
		wantedNames := cfg.GetStringSlice(core.CfgKeyXCleanName)
		if ageDays < 0 {
			return errors.Errorf("age must be >= 0")
		}
		if keep < 0 {
			return errors.Errorf("keep must be >= 0")
		}

		trashRoot, err := defaultCleanTrashDir()
		if err != nil {
			return err
		}
		if err := ensureDir(trashRoot); err != nil {
			return err
		}

		binDir := cfg.GetString(core.CfgKeyCmdrBinDir)
		binHelper := utils.NewPathHelper(binDir)

		threshold := time.Now().Add(-time.Duration(ageDays) * 24 * time.Hour)

		query, err := manager.Query()
		if err != nil {
			return err
		}
		commands, err := query.All()
		if err != nil {
			return err
		}

		// If names are specified, only clean those commands.
		wantedNameSet := map[string]struct{}{}
		if len(wantedNames) > 0 {
			for _, n := range wantedNames {
				n = strings.TrimSpace(n)
				if n == "" {
					continue
				}
				wantedNameSet[n] = struct{}{}
			}

			if len(wantedNameSet) == 0 {
				return errors.Errorf("name must not be empty")
			}

			seenNameSet := map[string]struct{}{}
			for _, cmd := range commands {
				seenNameSet[cmd.GetName()] = struct{}{}
			}
			missing := []string{}
			for n := range wantedNameSet {
				if _, ok := seenNameSet[n]; !ok {
					missing = append(missing, n)
				}
			}
			if len(missing) > 0 {
				sort.Strings(missing)
				return errors.Errorf("command(s) not found: %s", strings.Join(missing, ", "))
			}
		}

		activeLocationByName := map[string]string{}
		for _, cmd := range commands {
			name := cmd.GetName()
			if len(wantedNameSet) > 0 {
				if _, ok := wantedNameSet[name]; !ok {
					continue
				}
			}
			if _, ok := activeLocationByName[name]; ok {
				continue
			}
			location, err := binHelper.RealPath(name)
			if err == nil {
				activeLocationByName[name] = filepath.Clean(location)
			}
		}

		inactiveByName := map[string][]cleanCandidate{}
		for _, cmd := range commands {
			if len(wantedNameSet) > 0 {
				if _, ok := wantedNameSet[cmd.GetName()]; !ok {
					continue
				}
			}
			if cmd.GetActivated() {
				continue
			}

			src := filepath.Clean(cmd.GetLocation())
			if active, ok := activeLocationByName[cmd.GetName()]; ok && filepath.Clean(active) == src {
				logger.Warn("skip cleaning activated version detected by shim", map[string]interface{}{
					"name":     cmd.GetName(),
					"version":  cmd.GetVersion(),
					"location": src,
				})
				continue
			}

			info, err := os.Stat(src)
			if err != nil {
				logger.Warn("skip cleaning version with missing shim", map[string]interface{}{
					"name":     cmd.GetName(),
					"version":  cmd.GetVersion(),
					"location": src,
					"error":    err,
				})
				continue
			}

			inactiveByName[cmd.GetName()] = append(inactiveByName[cmd.GetName()], cleanCandidate{
				name:     cmd.GetName(),
				version:  cmd.GetVersion(),
				location: src,
				addedAt:  info.ModTime(),
			})
		}

		var resultErr error
		cleaned := 0
		for name, candidates := range inactiveByName {
			sort.Slice(candidates, func(i, j int) bool {
				if candidates[i].addedAt.Equal(candidates[j].addedAt) {
					return candidates[i].version > candidates[j].version
				}
				return candidates[i].addedAt.After(candidates[j].addedAt)
			})

			for idx, c := range candidates {
				if idx < keep {
					continue
				}
				if c.addedAt.After(threshold) {
					continue
				}

				dstDir := filepath.Join(trashRoot, name)
				if err := ensureDir(dstDir); err != nil {
					resultErr = multierror.Append(resultErr, err)
					continue
				}

				dst, err := uniquePath(dstDir, filepath.Base(c.location))
				if err != nil {
					resultErr = multierror.Append(resultErr, err)
					continue
				}

				if err := moveFile(c.location, dst); err != nil {
					resultErr = multierror.Append(resultErr, err)
					continue
				}

				if err := manager.Undefine(name, c.version); err != nil {
					// Best-effort rollback: restore shim back to original place.
					if rbErr := moveFile(dst, c.location); rbErr != nil {
						logger.Warn("failed to rollback shim after undefine failure", map[string]interface{}{
							"name":    name,
							"version": c.version,
							"error":   rbErr,
						})
					}
					resultErr = multierror.Append(resultErr, errors.Wrapf(err, "undefine %s:%s failed", name, c.version))
					continue
				}

				cleaned++
				logger.Info("cleaned inactive version", map[string]interface{}{
					"name":       name,
					"version":    c.version,
					"added_at":   c.addedAt.Format(time.RFC3339),
					"trashed_to": dst,
				})
			}
		}

		logger.Info("clean finished", map[string]interface{}{
			"cleaned":   cleaned,
			"age_days":  ageDays,
			"keep":      keep,
			"trash_dir": trashRoot,
		})

		return resultErr
	}),
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	cfg := core.GetConfiguration()
	flags := cleanCmd.Flags()

	flags.Int("age", 100, "only clean versions older than this many days")
	flags.Int("keep", 3, "keep this many newest inactive versions for each command")
	flags.StringArrayP("name", "n", []string{}, "only clean specified command name(s) (can be repeated)")

	utils.PanicOnError("binding flags",
		cfg.BindPFlag(core.CfgKeyXCleanAgeDays, flags.Lookup("age")),
		cfg.BindPFlag(core.CfgKeyXCleanKeep, flags.Lookup("keep")),
		cfg.BindPFlag(core.CfgKeyXCleanName, flags.Lookup("name")),
	)
}
