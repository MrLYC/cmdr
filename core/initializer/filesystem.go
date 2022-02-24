package initializer

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
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

	core.Logger.Debug("backup directory", map[string]interface{}{
		"path": dir,
	})

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
	filesystem fs.FS
	srcPath    string
	dstPath    string
	fileMode   os.FileMode
}

func (e *EmbedFSExporter) copyDir(dstPath string, perm os.FileMode) error {
	err := os.MkdirAll(dstPath, perm)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", dstPath)
	}

	return nil
}

func (e *EmbedFSExporter) copyFile(srcPath, dstPath string, perm os.FileMode) error {
	srcFile, err := e.filesystem.Open(srcPath)
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
		return errors.Wrap(err, "failed to export directory")
	}

	path, err := filepath.Rel(e.srcPath, srcPath)
	if err != nil {
		return errors.Wrap(err, "failed to get relative path")
	}

	dstPath := filepath.Join(e.dstPath, path)

	if d.IsDir() {
		return e.copyDir(dstPath, 0755)
	}

	return e.copyFile(srcPath, dstPath, e.fileMode)
}

func (e *EmbedFSExporter) Init() error {
	err := os.MkdirAll(e.dstPath, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create destination directory")
	}

	core.Logger.Debug("exporting embedded filesystem", map[string]interface{}{
		"src": e.srcPath,
		"dst": e.dstPath,
	})

	err = fs.WalkDir(e.filesystem, e.srcPath, e.exportDir)
	if err != nil {
		return errors.Wrap(err, "failed to walk source directory")
	}

	return nil
}

func NewEmbedFSExporter(embedFS fs.FS, srcPath, dstPath string, fileMode os.FileMode) *EmbedFSExporter {
	return &EmbedFSExporter{
		filesystem: embedFS,
		srcPath:    srcPath,
		dstPath:    dstPath,
		fileMode:   fileMode,
	}
}

type DirRender struct {
	data    interface{}
	srcPath string
	ext     string
}

func (r *DirRender) walkTemplates(fileHandler, dirHandler func(path string, info fs.FileInfo) error) error {
	var errs error

	err := filepath.Walk(r.srcPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to walk directory %s", r.srcPath)
		}

		if filepath.Ext(path) != r.ext {
			return nil
		}

		if info.IsDir() {
			return dirHandler(path, info)
		}

		return fileHandler(path, info)
	})

	if err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs
}

func (r *DirRender) renderTemplate(name, content string, target io.Writer) error {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return errors.Wrapf(err, "failed to parse template %s", name)
	}

	err = tmpl.Execute(target, r.data)
	if err != nil {
		return errors.Wrapf(err, "failed to execute template %s", name)
	}

	return nil
}

func (r *DirRender) renderPath(path string) (string, error) {
	dirPath := filepath.Dir(path)
	templateName := filepath.Base(path)

	var builder strings.Builder

	err := r.renderTemplate(path, templateName[:len(filepath.Base(path))-len(r.ext)], &builder)
	if err != nil {
		return "", err
	}

	return filepath.Join(dirPath, builder.String()), nil
}

func (r *DirRender) renderFile(path string, info os.FileInfo) error {
	targetPath, err := r.renderPath(path)
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE, info.Mode().Perm())
	if err != nil {
		return errors.Wrapf(err, "failed to create %s", path)
	}
	defer utils.CallClose(dstFile)

	templateContent, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", path)
	}
	defer os.Remove(path)

	core.Logger.Debug("rendering file", map[string]interface{}{
		"path": path,
	})
	return r.renderTemplate(path, string(templateContent), dstFile)
}

func (r *DirRender) renderDir(path string, info os.FileInfo) error {
	core.Logger.Debug("rendering dir", map[string]interface{}{
		"path": path,
	})

	targetPath, err := r.renderPath(path)
	if err != nil {
		return err
	}

	err = os.MkdirAll(targetPath, info.Mode().Perm())
	if err != nil {
		return errors.Wrapf(err, "failed to create %s", targetPath)
	}

	return nil
}

func (r *DirRender) Init() error {
	return r.walkTemplates(r.renderFile, r.renderDir)
}

func NewDirRender(srcPath, ext string, data interface{}) *DirRender {
	return &DirRender{
		data:    data,
		srcPath: srcPath,
		ext:     ext,
	}
}

func init() {
	core.RegisterInitializerFactory("profile-dir-backup", func(cfg core.Configuration) (core.Initializer, error) {
		return NewFSBackup(cfg.GetString(core.CfgKeyCmdrProfileDir)), nil
	})

	core.RegisterInitializerFactory("profile-dir-export", func(cfg core.Configuration) (core.Initializer, error) {
		return NewEmbedFSExporter(
			core.EmbedFS,
			filepath.Join("embed", "profile"),
			cfg.GetString(core.CfgKeyCmdrProfileDir),
			0644,
		), nil
	})

	core.RegisterInitializerFactory("profile-dir-render", func(cfg core.Configuration) (core.Initializer, error) {
		return NewDirRender(cfg.GetString(core.CfgKeyCmdrProfileDir), ".gotmpl", struct {
			Configuration core.Configuration
		}{cfg}), nil
	})
}
