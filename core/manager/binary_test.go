package manager_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/manager"
	"github.com/mrlyc/cmdr/core/utils"
)

var _ = Describe("Binary", func() {
	Context("Binary", func() {
		var (
			binDir      string
			shimsDir    string
			err         error
			commandName = "command"
			version     = "1.0.0"
			shimsName   = "command_1.0.0"
		)

		BeforeEach(func() {
			binDir, err = os.MkdirTemp("", "")
			Expect(err).To(BeNil())

			shimsDir, err = os.MkdirTemp("", "")
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			Expect(os.RemoveAll(binDir)).To(Succeed())
			Expect(os.RemoveAll(shimsDir)).To(Succeed())
		})

		It("should return name", func() {
			b := manager.NewBinary(binDir, shimsDir, commandName, version, shimsName)
			Expect(b.GetName()).To(Equal(commandName))
		})

		It("should return version", func() {
			b := manager.NewBinary(binDir, shimsDir, commandName, version, shimsName)
			Expect(b.GetVersion()).To(Equal(version))
		})

		It("should return location", func() {
			b := manager.NewBinary(binDir, shimsDir, commandName, version, shimsName)
			Expect(b.GetLocation()).To(Equal(filepath.Join(shimsDir, commandName, shimsName)))
		})

		Context("Activate", func() {
			var (
				binHelper   *utils.PathHelper
				shimsHelper *utils.PathHelper
			)

			BeforeEach(func() {
				binHelper = utils.NewPathHelper(binDir)
				shimsHelper = utils.NewPathHelper(filepath.Join(shimsDir)).Child(commandName)

				Expect(binHelper.MkdirAll(0755)).To(Succeed())
				Expect(shimsHelper.MkdirAll(0755)).To(Succeed())
				Expect(os.WriteFile(shimsHelper.Child(shimsName).Path(), []byte(""), 0755)).To(Succeed())
			})

			It("should not return activated", func() {
				b := manager.NewBinary(binDir, shimsDir, commandName, version, shimsName)
				Expect(b.GetActivated()).To(BeFalse())
			})

			It("should return activated", func() {
				Expect(binHelper.SymbolLink(commandName, shimsHelper.Child(shimsName).Path(), 0755)).To(Succeed())

				b := manager.NewBinary(binDir, shimsDir, commandName, version, shimsName)
				Expect(b.GetActivated()).To(BeTrue())
			})
		})
	})

	Context("BinariesFilter", func() {
		var (
			binaries []*manager.Binary
			filter   *manager.BinariesFilter
		)

		Context("Binaries is empty", func() {
			BeforeEach(func() {
				binaries = []*manager.Binary{}
				filter = manager.NewBinariesFilter(binaries)
			})

			It("should return empty", func() {
				result, err := filter.All()
				Expect(err).To(BeNil())
				Expect(result).To(BeEmpty())
			})

			It("should return error", func() {
				_, err := filter.One()
				Expect(errors.Cause(err)).To(Equal(core.ErrBinaryNotFound))
			})

			It("should return 0", func() {
				result, err := filter.Count()
				Expect(err).To(BeNil())
				Expect(result).To(Equal(0))
			})
		})

		Context("Binaries filter", func() {
			var (
				binaryA, binaryB *manager.Binary
			)

			BeforeEach(func() {
				binaryA = manager.NewBinary(
					"bin",
					"shims",
					"command-a",
					"1.0.0",
					"location-a",
				)
				binaryB = manager.NewBinary(
					"bin",
					"shims",
					"command-b",
					"1.0.1",
					"location-b",
				)
				binaries = []*manager.Binary{binaryA, binaryB}
				filter = manager.NewBinariesFilter(binaries)
			})

			It("should return array", func() {
				result, err := filter.All()
				Expect(err).To(BeNil())
				Expect(result).To(Equal([]core.Command{binaryA, binaryB}))
			})

			It("should return first command", func() {
				result, err := filter.One()
				Expect(err).To(BeNil())
				Expect(result).To(Equal(binaryA))
			})

			It("should return 1", func() {
				count, err := filter.Count()
				Expect(err).To(BeNil())
				Expect(count).To(Equal(2))
			})
		})
	})

	Context("BinaryManager", func() {
		var (
			binDir      string
			shimsDir    string
			version     = "1.0.0"
			commandName = "exists"
			err         error
			mgr         *manager.BinaryManager
		)

		getShimsPath := func(name string) string {
			return filepath.Join(
				shimsDir,
				name,
				fmt.Sprintf("%s_%s", name, version),
			)
		}

		getBinPath := func(name string) string {
			return filepath.Join(binDir, name)
		}

		BeforeEach(func() {
			binDir, err = os.MkdirTemp("", "bin")
			Expect(err).To(BeNil())

			shimsDir, err = os.MkdirTemp("", "shims")
			Expect(err).To(BeNil())

			mgr = manager.NewBinaryManagerWithLink(binDir, shimsDir, 0755)
		})

		JustBeforeEach(func() {
			cmdShimsDir := filepath.Join(shimsDir, commandName)
			Expect(os.MkdirAll(cmdShimsDir, 0755)).To(Succeed())
			Expect(os.WriteFile(
				filepath.Join(cmdShimsDir, fmt.Sprintf("%s_%s", commandName, version)),
				[]byte(""),
				0644,
			)).To(Succeed())
		})

		AfterEach(func() {
			Expect(os.RemoveAll(binDir)).To(Succeed())
			Expect(os.RemoveAll(shimsDir)).To(Succeed())
		})

		It("should init a manager", func() {
			Expect(os.RemoveAll(binDir)).To(Succeed())
			Expect(os.RemoveAll(shimsDir)).To(Succeed())

			mgr = manager.NewBinaryManagerWithLink(binDir, shimsDir, 0755)
			Expect(mgr.Init(false)).To(Succeed())
			Expect(binDir).To(BeADirectory())
			Expect(shimsDir).To(BeADirectory())
		})

		It("should close a manager", func() {
			Expect(mgr.Close()).To(Succeed())
		})

		It("should return provider", func() {
			Expect(mgr.Provider()).To(Equal(core.CommandProviderBinary))
		})

		Context("Define", func() {
			var (
				tempDir  string
				location string
			)

			BeforeEach(func() {
				tempDir, err = os.MkdirTemp("", "")
				Expect(err).To(BeNil())

				location = filepath.Join(tempDir, "location")
				Expect(os.WriteFile(location, []byte(""), 0755)).To(Succeed())
			})

			AfterEach(func() {
				Expect(os.RemoveAll(tempDir)).To(Succeed())
			})

			checkDefineResult := func(name string) {
				shimsPath := getShimsPath(name)
				Expect(shimsPath).To(BeARegularFile())

				info, err := os.Lstat(shimsPath)
				Expect(err).To(BeNil())
				Expect(info.Mode() & os.FileMode(0755)).To(Equal(os.FileMode(0755)))
			}

			It("should define a command", func() {
				nonexistsCommand := "nonexists"

				_, err := mgr.Define(nonexistsCommand, version, location)
				Expect(err).To(BeNil())
				checkDefineResult(nonexistsCommand)
			})

			It("should redefine a command", func() {
				_, err := mgr.Define(commandName, version, location)
				Expect(err).To(BeNil())
				checkDefineResult(commandName)
			})

			checkUndefineResult := func(name string) {
				shimsPath := getShimsPath(name)
				Expect(shimsPath).NotTo(BeAnExistingFile())
			}

			It("should undefine a command", func() {
				Expect(mgr.Undefine(commandName, version)).To(Succeed())
				checkUndefineResult(commandName)
			})

			It("should undefine a non-exists command", func() {
				nonexistsCommand := "nonexists"

				Expect(mgr.Undefine(nonexistsCommand, version)).To(Succeed())
				checkUndefineResult(nonexistsCommand)
			})
		})

		Context("Activate", func() {
			checkActivateResult := func(name string) {
				binPath := getBinPath(name)
				Expect(binPath).To(BeAnExistingFile())

				path, err := os.Readlink(binPath)
				Expect(err).To(BeNil())
				Expect(path).To(Equal(getShimsPath(name)))
			}

			It("should activate a command", func() {
				Expect(mgr.Activate(commandName, version)).To(Succeed())
				checkActivateResult(commandName)
			})

			It("should reactivate a command", func() {
				Expect(mgr.Activate(commandName, version)).To(Succeed())
				Expect(mgr.Activate(commandName, version)).To(Succeed())
				checkActivateResult(commandName)
			})

			It("should not activate a non-exists command", func() {
				nonexistsCommand := "nonexists"

				Expect(mgr.Activate(nonexistsCommand, version)).NotTo(Succeed())
				Expect(getBinPath(nonexistsCommand)).NotTo(BeAnExistingFile())
			})

			checkDeactivateResult := func(name string) {
				binPath := getBinPath(name)
				Expect(binPath).NotTo(BeAnExistingFile())
			}

			It("should deactivate a command", func() {
				Expect(os.Symlink(getShimsPath(commandName), getBinPath(commandName))).To(Succeed())
				Expect(mgr.Deactivate(commandName)).To(Succeed())
				checkDeactivateResult(commandName)
			})

			It("should deactivate a non-exists command", func() {
				nonexistsCommand := "nonexists"
				Expect(mgr.Deactivate(nonexistsCommand)).To(Succeed())
				checkDeactivateResult(nonexistsCommand)
			})
		})

		Context("Query", func() {
			It("should return a query", func() {
				query, err := mgr.Query()
				Expect(err).To(BeNil())

				command, err := query.One()
				Expect(err).To(BeNil())

				Expect(command.GetName()).To(Equal(commandName))
				Expect(command.GetVersion()).To(Equal(version))
				Expect(command.GetLocation()).To(Equal(getShimsPath(command.GetName())))
			})
		})

	})

	Context("Factory", func() {
		var cfg core.Configuration

		BeforeEach(func() {
			cfg = viper.New()
		})

		It("should new command manager", func() {
			mgr, err := core.NewCommandManager(core.CommandProviderBinary, cfg)
			Expect(err).To(BeNil())

			_, ok := mgr.(*manager.BinaryManager)
			Expect(ok).To(BeTrue())
		})

		It("should new initializer", func() {
			initializer, err := core.NewInitializer("binary", cfg)
			Expect(err).To(BeNil())

			_, ok := initializer.(*manager.BinaryManager)
			Expect(ok).To(BeTrue())
		})
	})
})
