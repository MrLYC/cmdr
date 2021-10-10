package core

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

type ModelTable struct {
	*tablewriter.Table
}

func NewModleTablePrinter(writer io.Writer) *ModelTable {
	table := tablewriter.NewWriter(writer)
	table.SetBorder(false)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")

	return &ModelTable{
		Table: table,
	}
}
