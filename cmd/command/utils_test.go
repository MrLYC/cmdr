package command

import (
	"fmt"
	"os"

	"github.com/agiledragon/gomonkey"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/config"
	"github.com/mrlyc/cmdr/define"
	"github.com/mrlyc/cmdr/runner"
	"github.com/mrlyc/cmdr/runner/mock"
)

var _ = Describe("Utils", func() {
	var (
		rawConfig, cfg define.Configuration
	)

	BeforeEach(func() {
		rawConfig = config.Global
		cfg = viper.New()
		config.Global = cfg
	})

	AfterEach(func() {
		config.Global = rawConfig
	})

	Describe("executeRunner", func() {
		var (
			ctrl       *gomock.Controller
			mockRunner *mock.MockRunner
			cmd        cobra.Command
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockRunner = mock.NewMockRunner(ctrl)
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		mockFactory := func(define.Configuration) runner.Runner {
			return mockRunner
		}

		It("should call runner.Run", func() {
			mockRunner.EXPECT().Run(cmd.Context()).Return(nil)

			fn := executeRunner(mockFactory)
			fn(&cmd, []string{})
		})

		It("should exit when return error", func() {
			var exitCode int
			patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
				exitCode = code
			})
			defer patches.Reset()

			mockRunner.EXPECT().Run(cmd.Context()).Return(fmt.Errorf("testing"))

			fn := executeRunner(mockFactory)
			fn(&cmd, []string{})
			Expect(exitCode).To(Equal(-1))
		})
	})
})
