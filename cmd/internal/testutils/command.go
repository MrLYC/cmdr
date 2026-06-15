package testutils

import (
	"github.com/golang/mock/gomock"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/mock"
)

type CommandFlagSpec struct {
	Name       string
	Shorthand  string
	ConfigKey  string
	Default    string
	IsRequired bool
}

func CheckCommandFlag(cmd *cobra.Command, name, shorthand, configKey string, value string, required bool) {
	flag := cmd.Flag(name)
	Expect(flag).NotTo(BeNil())

	Expect(flag.Name).To(Equal(name))
	if shorthand != "" {
		Expect(flag.Shorthand).To(Equal(shorthand))
	}

	if value != "" {
		Expect(flag.DefValue).To(Equal(value))
	}

	if required {
		Expect(flag.Annotations).To(HaveKey(cobra.BashCompOneRequiredFlag))
	}

	if configKey != "" {
		Expect(core.GetConfiguration().Get(configKey)).NotTo(BeNil())
	}
}

func CheckCommandFlags(cmd *cobra.Command, specs ...CommandFlagSpec) {
	for _, spec := range specs {
		CheckCommandFlag(cmd, spec.Name, spec.Shorthand, spec.ConfigKey, spec.Default, spec.IsRequired)
	}
}

type CommandHarness struct {
	Ctrl    *gomock.Controller
	Cfg     core.Configuration
	Manager *mock.MockCommandManager

	provider        core.CommandProvider
	previousCfg     core.Configuration
	previousFactory func(core.Configuration) (core.CommandManager, error)
}

func NewCommandHarness(provider core.CommandProvider) *CommandHarness {
	harness := &CommandHarness{
		provider:        provider,
		previousCfg:     core.GetConfiguration(),
		previousFactory: core.GetCommandManagerFactory(provider),
		Ctrl:            gomock.NewController(ginkgo.GinkgoT()),
		Cfg:             viper.New(),
	}
	harness.Manager = mock.NewMockCommandManager(harness.Ctrl)
	core.RegisterCommandManagerFactory(provider, func(core.Configuration) (core.CommandManager, error) {
		return harness.Manager, nil
	})
	core.SetConfiguration(harness.Cfg)
	return harness
}

func (h *CommandHarness) Finish() {
	if h.Ctrl != nil {
		h.Ctrl.Finish()
	}
	core.RegisterCommandManagerFactory(h.provider, h.previousFactory)
	core.SetConfiguration(h.previousCfg)
}

func (h *CommandHarness) ExpectClose() {
	h.Manager.EXPECT().Close().Return(nil)
}

func (h *CommandHarness) QueryReturning(commands *[]core.Command) *mock.MockCommandQuery {
	query := mock.NewMockCommandQuery(h.Ctrl)
	query.EXPECT().All().DoAndReturn(func() ([]core.Command, error) {
		if commands == nil {
			return nil, nil
		}
		return *commands, nil
	}).AnyTimes()
	h.Manager.EXPECT().Query().Return(query, nil)
	return query
}

func (h *CommandHarness) Command(name, version, location string, activated bool) core.Command {
	command := mock.NewMockCommand(h.Ctrl)
	command.EXPECT().GetActivated().Return(activated).AnyTimes()
	command.EXPECT().GetName().Return(name).AnyTimes()
	command.EXPECT().GetVersion().Return(version).AnyTimes()
	command.EXPECT().GetLocation().Return(location).AnyTimes()
	return command
}
