package core_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
)

var _ = Describe("Sorter", func() {
	It("Run", func() {
		commands := []*model.Command{
			{Name: "a", Version: "1.0.0", Activated: true},
			{Name: "b", Version: "1.0.1", Activated: true},
			{Name: "b", Version: "1.0.0", Activated: false},
			{Name: "a", Version: "1.0.1", Activated: false},
			{Name: "c", Version: "1.0.0", Activated: false},
		}

		sorter := NewCommandSorter()
		_, err := sorter.Run(context.WithValue(context.Background(), define.ContextKeyCommands, commands))
		Expect(err).To(Succeed())
		Expect(commands).To(Equal([]*model.Command{
			{Name: "a", Version: "1.0.0", Activated: true},
			{Name: "b", Version: "1.0.1", Activated: true},
			{Name: "a", Version: "1.0.1", Activated: false},
			{Name: "b", Version: "1.0.0", Activated: false},
			{Name: "c", Version: "1.0.0", Activated: false},
		}))
	})
})
