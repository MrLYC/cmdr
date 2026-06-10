package cmd

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/cmd/internal/testutils"
	"github.com/mrlyc/cmdr/core"
)

type initTestInitializer struct {
	step  string
	calls *[]string
	err   error
}

func (i initTestInitializer) Init(isUpgrade bool) error {
	*i.calls = append(*i.calls, i.step)
	return i.err
}

var _ = Describe("Init", func() {
	It("should check flags", func() {
		testutils.CheckCommandFlag(initCmd, "upgrade", "u", core.CfgKeyXInitUpgrade, "false", false)
	})

	It("should run initializers in order", func() {
		steps := []string{
			"profile-dir-backup",
			"binary",
			"database-migrator",
			"profile-dir-export",
			"profile-dir-render",
			"profile-injector",
			"cmdr-updater",
		}
		calls := []string{}
		for _, step := range steps {
			step := step
			core.RegisterInitializerFactory(step, func(core.Configuration) (core.Initializer, error) {
				return initTestInitializer{step: step, calls: &calls}, nil
			})
		}

		Expect(initCmd.Flags().Set("upgrade", "true")).To(Succeed())
		Expect(func() { initCmd.Run(initCmd, nil) }).NotTo(Panic())
		Expect(calls).To(Equal(steps))
		Expect(initCmd.Flags().Set("upgrade", "false")).To(Succeed())
	})

	It("should report initializer creation and init errors", func() {
		steps := []string{
			"profile-dir-backup",
			"binary",
			"database-migrator",
			"profile-dir-export",
			"profile-dir-render",
			"profile-injector",
			"cmdr-updater",
		}
		for _, step := range steps {
			step := step
			core.RegisterInitializerFactory(step, func(core.Configuration) (core.Initializer, error) {
				return initTestInitializer{step: step, calls: &[]string{}, err: errors.New("init failed")}, nil
			})
		}
		core.RegisterInitializerFactory("binary", func(core.Configuration) (core.Initializer, error) {
			return nil, errors.New("create failed")
		})

		Expect(func() { initCmd.Run(initCmd, nil) }).To(Panic())
	})
})
