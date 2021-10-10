package main

import (
	"context"

	_ "modernc.org/sqlite"

	"github.com/mrlyc/cmdr/cmd"
)

// go:generate ent generate ./model/schema
func main() {
	ctx := context.Background()
	cmd.ExecuteContext(ctx)
}
