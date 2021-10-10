package main

import (
	"context"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mrlyc/cmdr/cmd"
)

// go:generate ent generate ./model/schema
func main() {
	ctx := context.Background()
	cmd.ExecuteContext(ctx)
}
