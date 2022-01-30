package operator

import (
	"context"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/utils"
)

type DirectoryMaker struct {
	BaseOperator
	dirs map[string]string
}

func (m *DirectoryMaker) String() string {
	return "directory-maker"
}

func (m *DirectoryMaker) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	for n, p := range m.dirs {
		logger.Info("creating dir", map[string]interface{}{
			"name": n,
			"dir":  p,
		})
		utils.ExitWithError(
			define.FS.MkdirAll(p, 0755),
			"making dir %s failed", n,
		)
	}

	return ctx, nil
}

func NewDirectoryMaker(dirs map[string]string) *DirectoryMaker {
	return &DirectoryMaker{
		dirs: dirs,
	}
}

type DirectoryRemover struct {
	BaseOperator
	dirs map[string]string
}

func (r *DirectoryRemover) String() string {
	return "directory-remover"
}

func (r *DirectoryRemover) Run(ctx context.Context) (context.Context, error) {
	logger := define.Logger
	for n, p := range r.dirs {
		logger.Info("removing dir", map[string]interface{}{
			"name": n,
			"dir":  p,
		})
		utils.ExitWithError(
			define.FS.RemoveAll(p),
			"removing dir %s failed", n,
		)
	}

	return ctx, nil
}

func NewDirectoryRemover(dirs map[string]string) *DirectoryRemover {
	return &DirectoryRemover{
		dirs: dirs,
	}
}
