package initializer_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/initializer"
)

var _ = Describe("Profile", func() {
	Context("ProfileInjector", func() {
		var (
			rootDir    string
			scriptPath = "~/.cmdr/cmdr_profile"
		)

		BeforeEach(func() {
			var err error

			rootDir, err = os.MkdirTemp("", "")
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			Expect(os.RemoveAll(rootDir)).To(Succeed())
		})

		It("should inject when path not exists", func() {
			injector := initializer.NewProfileInjector(scriptPath, filepath.Join(rootDir, "not_exists"))
			Expect(injector.Init(false)).To(Succeed())
		})

		DescribeTable("should inject profile", func(content, excepted string) {
			profilePath := filepath.Join(rootDir, "profile")
			Expect(os.WriteFile(profilePath, []byte(content), 0755)).To(Succeed())

			injector := initializer.NewProfileInjector(scriptPath, profilePath)
			Expect(injector.Init(false)).To(Succeed())

			profile, err := os.ReadFile(profilePath)
			Expect(err).To(BeNil())
			Expect(string(profile)).To(Equal(excepted))
		},
			Entry("empty", "", "source '~/.cmdr/cmdr_profile'\n"),
			Entry("update path", "source ~/cmdr_profile", "source '~/.cmdr/cmdr_profile'\n"),
			Entry("update path with quotes", "source '~/cmdr_profile'", "source '~/.cmdr/cmdr_profile'\n"),
			Entry("update path with double quotes", "source \"~/cmdr_profile\"", "source '~/.cmdr/cmdr_profile'\n"),
			Entry("quotes", "source '~/.cmdr/cmdr_profile'", "source '~/.cmdr/cmdr_profile'\n"),
			Entry("double quotes", "source \"~/.cmdr/cmdr_profile\"", "source '~/.cmdr/cmdr_profile'\n"),
			Entry("no quotes", "source ~/.cmdr/cmdr_profile", "source '~/.cmdr/cmdr_profile'\n"),
			Entry("source with dot", ". ~/cmdr_profile", ". '~/.cmdr/cmdr_profile'\n"),
			Entry("with space", "  . ~/cmdr_profile", "  . '~/.cmdr/cmdr_profile'\n"),
			Entry("with tab", "\t. ~/cmdr_profile", "\t. '~/.cmdr/cmdr_profile'\n"),
			Entry("&& with space", "should_run && . ~/cmdr_profile", "should_run && . '~/.cmdr/cmdr_profile'\n"),
			Entry("&& without space", "should_run&&. ~/cmdr_profile", "should_run&&. '~/.cmdr/cmdr_profile'\n"),
			Entry("|| with space", "should_not_run || . ~/cmdr_profile", "should_not_run || . '~/.cmdr/cmdr_profile'\n"),
			Entry("|| without space", "should_not_run||. ~/cmdr_profile", "should_not_run||. '~/.cmdr/cmdr_profile'\n"),
			Entry("in if statement with space", "if should_run; then . ~/cmdr_profile; fi", "if should_run; then . '~/.cmdr/cmdr_profile'; fi\n"),
			Entry("in if statement without space", "if should_run;then . ~/cmdr_profile;fi", "if should_run;then . '~/.cmdr/cmdr_profile';fi\n"),
			Entry("with comment", ". ~/cmdr_profile  # load cmdr", ". '~/.cmdr/cmdr_profile'  # load cmdr\n"),
			Entry("comment only", "# . ~/cmdr_profile", "# . ~/cmdr_profile\nsource '~/.cmdr/cmdr_profile'\n"),
			Entry(
				"multiple line",
				"#!/bin/bash\n# a comment\ndo_first_thing\n. ~/cmdr_profile\ndo_second_thing\n",
				"#!/bin/bash\n# a comment\ndo_first_thing\n. '~/.cmdr/cmdr_profile'\ndo_second_thing\n",
			),
			Entry(
				"multiple profile",
				"#!/bin/bash\n# a comment\n. ~/cmdr_profile\ndo_some_thing\n. ~/cmdr_profile\n",
				"#!/bin/bash\n# a comment\n. '~/.cmdr/cmdr_profile'\ndo_some_thing\n",
			),
		)
	})

	Context("Injector", func() {
		homeDir, _ := os.UserHomeDir()

		DescribeTable("Profile path", func(shell, excepted string) {
			cfg := viper.New()
			cfg.Set(core.CfgKeyCmdrShell, shell)

			i, err := core.NewInitializer("profile-injector", cfg)
			Expect(err).To(BeNil())

			injector := i.(*initializer.ProfileInjector)

			Expect(injector.ProfilePath()).To(Equal(excepted))
		},
			Entry("bash", "bash", filepath.Join(homeDir, ".bashrc")),
			Entry("ash", "ash", filepath.Join(homeDir, ".profile")),
			Entry("zsh", "zsh", filepath.Join(homeDir, ".zshrc")),
			Entry("sh", "sh", filepath.Join(homeDir, ".profile")),
			Entry("fish", "fish", filepath.Join(homeDir, ".config", "fish", "config.fish")),
		)
	})
})
