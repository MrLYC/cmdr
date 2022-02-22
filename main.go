package main

import (
	"context"

	"github.com/mrlyc/cmdr/cmd"
	_ "github.com/mrlyc/cmdr/core/initializer"
	_ "github.com/mrlyc/cmdr/core/manager"
)

// go:generate ent generate ./model/schema
func main() {
	ctx := context.Background()
	cmd.ExecuteContext(ctx)
}
