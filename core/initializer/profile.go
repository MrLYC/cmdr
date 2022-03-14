package initializer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

type ProfileInjector struct {
	scriptPath  string
	profilePath string
}

func (p *ProfileInjector) makeProfileStatement() string {
	return fmt.Sprintf(`source '%s'`, p.scriptPath)
}

func (p *ProfileInjector) makeProfileScript() (io.Reader, error) {
	scriptName := filepath.Base(p.scriptPath)
	scriptPath := fmt.Sprintf(`'%s'`, p.scriptPath)

	re := regexp.MustCompile(fmt.Sprintf(
		`^([^#]*?(?:^|\s|\||&|;)?(?:source|\.)\s+)(['"]?.*?%s['"]?)(.*?)$`,
		scriptName,
	))

	file, err := os.Open(p.profilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open %s", p.profilePath)
	}
	defer utils.CallClose(file)

	buffer := bytes.NewBuffer(nil)

	found := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := re.FindStringSubmatch(line)
		if len(parts) != 4 {
			_, _ = fmt.Fprintln(buffer, line)
			continue
		}

		if found {
			continue
		}

		found = true
		_, _ = fmt.Fprintln(buffer, strings.Join([]string{parts[1], scriptPath, parts[3]}, ""))
	}

	if !found {
		_, _ = fmt.Fprintln(buffer, p.makeProfileStatement())
	}

	err = scanner.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %s", p.profilePath)
	}

	return buffer, nil
}

func (p *ProfileInjector) Init() error {
	var script io.Reader

	_, err := os.Stat(p.profilePath)
	if err == nil {
		script, err = p.makeProfileScript()
		if err != nil {
			return errors.Wrapf(err, "failed to make profile script")
		}

	} else {
		script = strings.NewReader(fmt.Sprintf("%s\n", p.makeProfileStatement()))
	}

	file, err := os.OpenFile(p.profilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", p.profilePath)
	}
	defer utils.CallClose(file)

	_, err = io.Copy(file, script)
	if err != nil {
		return errors.Wrapf(err, "failed to write %s", p.profilePath)
	}

	return nil
}

func NewProfileInjector(scriptPath, profilePath string) *ProfileInjector {
	return &ProfileInjector{
		scriptPath:  scriptPath,
		profilePath: profilePath,
	}
}

func getProfilePathByShell(path string) (string, error) {
	shell := filepath.Base(path)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get user home dir")
	}

	switch shell {
	case "bash", "ash":
		return filepath.Join(homeDir, ".bashrc"), nil
	case "zsh":
		return filepath.Join(homeDir, ".zshrc"), nil
	case "sh":
		return filepath.Join(homeDir, ".profile"), nil
	default:
		return "", errors.Wrapf(core.ErrShellNotSupported, shell)
	}
}

func init() {
	core.RegisterInitializerFactory("profile-injector", func(cfg core.Configuration) (core.Initializer, error) {
		profilePath := cfg.GetString(core.CfgKeyCmdrProfilePath)
		if profilePath == "" {
			path, err := getProfilePathByShell(cfg.GetString(core.CfgKeyCmdrShell))
			if err != nil {
				return nil, err
			}

			profilePath = path
		}

		return NewProfileInjector(filepath.Join(
			cfg.GetString(core.CfgKeyCmdrProfileDir),
			"cmdr_initializer.sh",
		), profilePath), nil
	})
}
