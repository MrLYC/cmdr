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

	"github.com/mrlyc/cmdr/core/utils"
)

type ProfileInjector struct {
	scriptPath  string
	profilePath string
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
		statement := fmt.Sprintf(`source %s`, scriptPath)
		_, _ = fmt.Fprintln(buffer, statement)
	}

	err = scanner.Err()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %s", p.profilePath)
	}

	return buffer, nil
}

func (p *ProfileInjector) Init() error {
	stat, err := os.Stat(p.profilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to stat %s", p.profilePath)
	}

	script, err := p.makeProfileScript()
	if err != nil {
		return errors.Wrapf(err, "failed to make profile script")
	}

	file, err := os.OpenFile(p.profilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, stat.Mode().Perm())
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
