package command

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/define"
)

func TestCommand(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Command Suite")
}

func checkCommandFlag(cmd *cobra.Command, name, shorthand, configKey string, value string, required bool) {
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
		Expect(define.Config.Get(configKey)).NotTo(BeNil())
	}
}
