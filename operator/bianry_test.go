package operator_test

import (
	"context"
	"fmt"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/model"
	"github.com/mrlyc/cmdr/operator"
	"github.com/mrlyc/cmdr/utils"
)

var _ = Describe("Bianry", func() {
	var (
		command1, command2 *model.Command
		ctx                context.Context
		helper             *utils.CmdrHelper
	)

	BeforeEach(func() {
		tempDir, err := afero.TempDir(define.FS, "", "")
		Expect(err).To(BeNil())
		helper = utils.NewCmdrHelper(tempDir)

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
		Expect(define.FS.RemoveAll(helper.GetRootDir())).To(Succeed())
	})

	Context("BinariesInstaller", func() {
		var installer *operator.BinariesInstaller

		BeforeEach(func() {
			installer = operator.NewBinariesInstaller(helper)
		})

		It("context not found", func() {
			_, err := installer.Run(context.Background())
			Expect(errors.Cause(err)).To(Equal(operator.ErrContextValueNotFound))
		})

		It("install binaries", func() {
			_, err := installer.Run(ctx)
			Expect(err).To(BeNil())

			Expect(
				afero.Exists(define.FS, helper.GetCommandShimsPath(command1.Name, command1.Version)),
			).To(BeTrue())
			Expect(
				afero.Exists(define.FS, helper.GetCommandShimsPath(command2.Name, command2.Version)),
			).To(BeTrue())
		})

		It("install binaries partital success", func() {
			command1.Location = fmt.Sprintf("%s_no_exists", command1.Location)

			_, err := installer.Run(ctx)
			Expect(err).NotTo(BeNil())

			Expect(afero.Exists(
				define.FS, helper.GetCommandShimsPath(command1.Name, command1.Version),
			)).To(BeFalse())
			Expect(afero.Exists(
				define.FS, helper.GetCommandShimsPath(command2.Name, command2.Version),
			)).To(BeTrue())
		})

		It("install binaries with not managed", func() {
			command2.Managed = false

			_, err := installer.Run(ctx)
			Expect(err).To(BeNil())

			Expect(afero.Exists(
				define.FS, helper.GetCommandShimsPath(command1.Name, command1.Version),
			)).To(BeTrue())
			Expect(afero.Exists(
				define.FS, helper.GetCommandShimsPath(command2.Name, command2.Version),
			)).To(BeFalse())
		})
	})

	Context("BinariesUninstaller", func() {
		var uninstaller *operator.BinariesUninstaller

		BeforeEach(func() {
			uninstaller = operator.NewBinariesUninstaller()
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
			Expect(errors.Cause(err)).To(Equal(operator.ErrContextValueNotFound))
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

	Context("BinariesActivator", func() {
		var activator *operator.BinariesActivator

		BeforeEach(func() {
			command1.Managed = true
			command2.Managed = false

			Expect(define.FS.MkdirAll(helper.GetCommandShimsDir(command1.Name), 0755)).To(Succeed())
			Expect(define.FS.Rename(command1.Location, helper.GetCommandShimsPath(command1.Name, command1.Version))).To(Succeed())
			activator = operator.NewBinariesActivator(helper)
		})

		It("context not found", func() {
			_, err := activator.Run(context.Background())
			Expect(errors.Cause(err)).To(Equal(operator.ErrContextValueNotFound))
		})

		It("activate binaries", func() {
			_, err := activator.Run(ctx)
			Expect(err).To(BeNil())

			Expect(
				afero.Exists(define.FS, helper.GetCommandBinPath(command1.Name)),
			).To(BeTrue())
			Expect(
				afero.Exists(define.FS, helper.GetCommandBinPath(command2.Name)),
			).To(BeTrue())
		})

		It("activate binaries that already exists", func() {
			Expect(define.FS.MkdirAll(helper.GetBinDir(), 0755)).To(BeNil())
			Expect(afero.WriteFile(define.FS, helper.GetCommandBinPath(command1.Name), []byte(""), 0755)).To(BeNil())
			Expect(define.FS.MkdirAll(helper.GetCommandBinPath(command2.Name), 0755)).To(BeNil())

			_, err := activator.Run(ctx)
			Expect(err).To(BeNil())

			Expect(
				afero.Exists(define.FS, helper.GetCommandBinPath(command1.Name)),
			).To(BeTrue())
			Expect(
				afero.Exists(define.FS, helper.GetCommandBinPath(command2.Name)),
			).To(BeTrue())
		})

	})
})
