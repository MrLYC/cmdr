package testutils

import (
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/mrlyc/cmdr/core"
)

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
