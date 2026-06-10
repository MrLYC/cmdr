package cmd

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	pkgerrors "github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/utils"
)

type upgradeTestSearcher struct {
	err error
}

func (s upgradeTestSearcher) GetReleaseAsset(context.Context, string, string) (core.CmdrReleaseAsset, error) {
	if s.err != nil {
		return core.CmdrReleaseAsset{}, s.err
	}
	return core.CmdrReleaseAsset{
		Name:    "v1.2.3",
		Version: "1.2.3",
		Asset:   "cmdr",
		Url:     "https://example.com/cmdr",
	}, nil
}

var _ = Describe("Upgrade", func() {
	var cfg *viper.Viper

	BeforeEach(func() {
		cfg = viper.New()
		cfg.Set(core.CfgKeyXUpgradeRelease, "latest")
		cfg.Set(core.CfgKeyXUpgradeAsset, "cmdr")
		cfg.Set(core.CfgKeyXUpgradeArgs, []string{"init", "--upgrade"})
		core.RegisterCmdrSearcherFactory(core.CmdrSearcherProviderDefault, func(core.Configuration) (core.CmdrSearcher, error) {
			return upgradeTestSearcher{}, nil
		})
	})

	AfterEach(func() {
		upgradeCmdrFn = utils.UpgradeCmdr
	})

	It("should set init upgrade args in pre-run", func() {
		previous := core.GetConfiguration()
		defer core.SetConfiguration(previous)
		core.SetConfiguration(cfg)

		upgradeCmd.PreRun(upgradeCmd, nil)
		Expect(cfg.GetStringSlice(core.CfgKeyXUpgradeArgs)).To(Equal([]string{"init", "--upgrade"}))
	})

	It("should run upgrade success and non-fatal states", func() {
		var receivedArgs []string
		upgradeCmdrFn = func(ctx context.Context, cfg core.Configuration, url, version string, args []string) error {
			Expect(url).To(Equal("https://example.com/cmdr"))
			Expect(version).To(Equal("1.2.3"))
			receivedArgs = args
			return nil
		}
		Expect(func() { runUpgrade(context.Background(), cfg, []string{"--shell", "bash"}) }).NotTo(Panic())
		Expect(receivedArgs).To(Equal([]string{"init", "--upgrade", "--shell", "bash"}))

		upgradeCmdrFn = func(context.Context, core.Configuration, string, string, []string) error {
			return utils.ErrCmdrAlreadyLatestVersion
		}
		Expect(func() { runUpgrade(context.Background(), cfg, nil) }).NotTo(Panic())

		upgradeCmdrFn = func(context.Context, core.Configuration, string, string, []string) error {
			return pkgerrors.Wrap(utils.ErrCmdrCommandAlreadyDefined, "wrapped")
		}
		Expect(func() { runUpgrade(context.Background(), cfg, nil) }).NotTo(Panic())
	})

	It("should panic on searcher and upgrade errors", func() {
		core.RegisterCmdrSearcherFactory(core.CmdrSearcherProviderDefault, func(core.Configuration) (core.CmdrSearcher, error) {
			return nil, errors.New("factory failed")
		})
		Expect(func() { runUpgrade(context.Background(), cfg, nil) }).To(Panic())

		core.RegisterCmdrSearcherFactory(core.CmdrSearcherProviderDefault, func(core.Configuration) (core.CmdrSearcher, error) {
			return upgradeTestSearcher{err: errors.New("search failed")}, nil
		})
		Expect(func() { runUpgrade(context.Background(), cfg, nil) }).To(Panic())

		core.RegisterCmdrSearcherFactory(core.CmdrSearcherProviderDefault, func(core.Configuration) (core.CmdrSearcher, error) {
			return upgradeTestSearcher{}, nil
		})
		upgradeCmdrFn = func(context.Context, core.Configuration, string, string, []string) error {
			return errors.New("upgrade failed")
		}
		Expect(func() { runUpgrade(context.Background(), cfg, nil) }).To(Panic())
	})
})
