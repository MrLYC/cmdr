package core

import (
	"encoding/json"
	"io"

	"entgo.io/ent"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type ModelTable struct {
	fields []ent.Field
	*tablewriter.Table
}

func (p *ModelTable) Append(model interface{}) error {
	data, err := json.Marshal(model)
	if err != nil {
		return errors.Wrapf(err, "dump failed")
	}

	mappings := make(map[string]interface{}, len(p.fields))
	err = json.Unmarshal(data, &mappings)
	if err != nil {
		return errors.Wrapf(err, "load failed")
	}

	row := make([]string, 0, len(p.fields))
	for _, f := range p.fields {
		name := f.Descriptor().Name
		var content string
		value, ok := mappings[name]
		if ok {
			content = cast.ToString(value)
		}
		row = append(row, content)
	}

	p.Table.Append(row)

	return nil
}

func NewModleTablePrinter(schema ent.Interface, writer io.Writer) *ModelTable {
	table := tablewriter.NewWriter(writer)
	fields := schema.Fields()
	headers := make([]string, 0, len(fields))

	for _, f := range fields {
		headers = append(headers, f.Descriptor().Name)
	}

	table.SetHeader(headers)

	return &ModelTable{
		fields: fields,
		Table:  table,
	}
}
