package core_test

import (
	"context"
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
)

var _ = Describe("Bianry", func() {
	var (
		command1, command2 *model.Command
		ctx                context.Context
		shimsDir           string
	)

	BeforeEach(func() {
		tempDir, err := afero.TempDir(define.FS, "", "")
		Expect(err).To(BeNil())
		shimsDir = filepath.Join(tempDir, "shims")

		location1 := filepath.Join(tempDir, "test1.sh")
		err = afero.WriteFile(define.FS, location1, []byte(`#!/bin/sh\necho $@`), 0755)
		Expect(err).To(BeNil())

		command1 = &model.Command{
			Name:     "test1",
			Version:  "1.0.0",
			Location: location1,
			Managed:  true,
		}

		location2 := filepath.Join(tempDir, "test2.sh")
		err = afero.WriteFile(define.FS, location2, []byte(`#!/bin/sh\necho $@`), 0755)
		Expect(err).To(BeNil())

		command2 = &model.Command{
			Name:     "test2",
			Version:  "1.0.0",
			Location: location2,
			Managed:  true,
		}

		ctx = context.WithValue(context.Background(), define.ContextKeyCommands, []*model.Command{command1, command2})
	})

	AfterEach(func() {
		Expect(define.FS.RemoveAll(shimsDir)).To(Succeed())
	})

	Context("BinariesInstaller", func() {
		var installer *core.BinariesInstaller

		BeforeEach(func() {
			installer = core.NewBinariesInstaller(shimsDir)
		})

		It("context not found", func() {
			_, err := installer.Run(context.Background())
			Expect(errors.Cause(err)).To(Equal(core.ErrContextValueNotFound))
		})

		It("install binaries", func() {
			_, err := installer.Run(ctx)
			Expect(err).To(BeNil())

			Expect(
				afero.Exists(define.FS, core.GetCommandPath(shimsDir, command1.Name, command1.Version)),
			).To(BeTrue())
			Expect(
				afero.Exists(define.FS, core.GetCommandPath(shimsDir, command2.Name, command2.Version)),
			).To(BeTrue())
		})

		It("install binaries partital success", func() {
			command1.Location = fmt.Sprintf("%s_no_exists", command1.Location)

			_, err := installer.Run(ctx)
			Expect(err).NotTo(BeNil())

			Expect(afero.Exists(
				define.FS, core.GetCommandPath(shimsDir, command1.Name, command1.Version),
			)).To(BeFalse())
			Expect(afero.Exists(
				define.FS, core.GetCommandPath(shimsDir, command2.Name, command2.Version),
			)).To(BeTrue())
		})

		It("install binaries with not managed", func() {
			command2.Managed = false

			_, err := installer.Run(ctx)
			Expect(err).To(BeNil())

			Expect(afero.Exists(
				define.FS, core.GetCommandPath(shimsDir, command1.Name, command1.Version),
			)).To(BeTrue())
			Expect(afero.Exists(
				define.FS, core.GetCommandPath(shimsDir, command2.Name, command2.Version),
			)).To(BeFalse())
		})
	})

	Context("BinariesUninstaller", func() {
		var uninstaller *core.BinariesUninstaller

		BeforeEach(func() {
			uninstaller = core.NewBinariesUninstaller()
		})

		It("uninstall binaries", func() {
			_, err := uninstaller.Run(ctx)
			Expect(err).To(BeNil())

			Expect(
				afero.Exists(define.FS, command1.Location),
			).To(BeFalse())
			Expect(
				afero.Exists(define.FS, command2.Location),
			).To(BeFalse())
		})

		It("context not found", func() {
			_, err := uninstaller.Run(context.Background())
			Expect(errors.Cause(err)).To(Equal(core.ErrContextValueNotFound))
		})

		It("uninstall binaries with not managed", func() {
			command1.Managed = false

			_, err := uninstaller.Run(ctx)
			Expect(err).To(BeNil())

			Expect(
				afero.Exists(define.FS, command1.Location),
			).To(BeTrue())
			Expect(
				afero.Exists(define.FS, command2.Location),
			).To(BeFalse())
		})

		It("uninstall binaries that not exists", func() {
			define.FS.Remove(command1.Location)

			_, err := uninstaller.Run(ctx)
			Expect(err).To(BeNil())

			Expect(
				afero.Exists(define.FS, command1.Location),
			).To(BeFalse())
			Expect(
				afero.Exists(define.FS, command2.Location),
			).To(BeFalse())
		})
	})
})
