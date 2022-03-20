package main

import (
	"context"
	"os"

	"github.com/mrlyc/cmdr/cmd"
	"github.com/mrlyc/cmdr/core"
	_ "github.com/mrlyc/cmdr/core/initializer"
	_ "github.com/mrlyc/cmdr/core/manager"
)

// go:generate ent generate ./model/schema
func main() {
	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}

		switch err := recovered.(type) {
		case core.ExitError:
			os.Exit(err.Code())
		default:
			panic(err)
		}

	}()

	ctx := context.Background()
	cmd.ExecuteContext(ctx)
}
