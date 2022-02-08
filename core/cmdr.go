package core

import (
	"github.com/pkg/errors"

	"github.com/mrlyc/cmdr/define"
)

type Cmdr struct {
	BinaryManager  *BinaryManager
	CommandManager *CommandManager
}

func (c *Cmdr) Init() error {
	for _, mgr := range []define.Initializer{
		c.BinaryManager,
		c.CommandManager,
	} {
		err := mgr.Init()
		if err != nil {
			return errors.Wrapf(err, "initialize %s failed", mgr)
		}
	}

	return nil
}

func NewCmdr(root string) (*Cmdr, error) {
	mgr, err := NewCommandManager(root)
	if err != nil {
		return nil, err
	}

	return &Cmdr{
		BinaryManager:  NewBinaryManager(root),
		CommandManager: mgr,
	}, nil
}
