package utils

import (
	"sort"

	"github.com/mrlyc/cmdr/core"
)

// sort.Interface implementation for sorting a slice of maps by a given key
type sortedCommands []core.Command

// Len returns the length of the slice
func (c sortedCommands) Len() int {
	return len(c)
}

// Swap swaps the elements at the given indices
func (c sortedCommands) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less returns true if the element at index i is less than the element at index j
func (c sortedCommands) Less(i, j int) bool {
	name1 := c[i].GetName()
	name2 := c[j].GetName()
	if name1 != name2 {
		return name1 < name2
	}

	activated1 := c[i].GetActivated()
	activated2 := c[j].GetActivated()
	if activated1 != activated2 {
		return activated1
	}

	return c[i].GetVersion() < c[j].GetVersion()
}

func SortCommands(commands []core.Command) {
	sort.Sort(sortedCommands(commands))
}
