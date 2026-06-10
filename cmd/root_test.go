package cmd

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Root", func() {
	It("should set default configuration values", func() {
		cfg := viper.New()
		preInitConfigWith(cfg, "/home/test/.cmdr")

		Expect(cfg.GetString(core.CfgKeyCmdrRootDir)).To(Equal("/home/test/.cmdr"))
		Expect(cfg.GetString(core.CfgKeyCmdrBinDir)).To(Equal("bin"))
		Expect(cfg.GetString(core.CfgKeyCmdrShimsDir)).To(Equal("shims"))
		Expect(cfg.GetString(core.CfgKeyCmdrProfileDir)).To(Equal("profile"))
		Expect(cfg.GetString(core.CfgKeyCmdrDatabasePath)).To(Equal("cmdr.db"))
		Expect(cfg.GetString(core.CfgKeyLogLevel)).To(Equal("info"))
		Expect(cfg.GetString(core.CfgKeyLogOutput)).To(Equal("stderr"))
	})

	It("should read config only when file exists", func() {
		dir, err := os.MkdirTemp("", "cmdr-root")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)

		cfg := viper.New()
		cfg.Set(core.CfgKeyCmdrConfigPath, filepath.Join(dir, "config.yaml"))
		Expect(os.WriteFile(cfg.GetString(core.CfgKeyCmdrConfigPath), []byte("core:\n  root_dir: /custom\n"), 0644)).To(Succeed())
		initConfigWith(cfg, os.Stat)
		Expect(cfg.GetString(core.CfgKeyCmdrRootDir)).To(Equal("/custom"))

		cfg = viper.New()
		cfg.Set(core.CfgKeyCmdrConfigPath, filepath.Join(dir, "missing.yaml"))
		initConfigWith(cfg, func(string) (os.FileInfo, error) {
			return nil, errors.New("missing")
		})
		Expect(cfg.ConfigFileUsed()).To(Equal(filepath.Join(dir, "missing.yaml")))
	})

	It("should convert relative paths under root dir", func() {
		cfg := viper.New()
		cfg.Set(core.CfgKeyCmdrRootDir, "/cmdr")
		cfg.Set(core.CfgKeyCmdrBinDir, "bin")
		cfg.Set(core.CfgKeyCmdrShimsDir, "/custom/shims")
		cfg.Set(core.CfgKeyCmdrProfileDir, "profile")
		cfg.Set(core.CfgKeyCmdrDatabasePath, "cmdr.db")

		postInitConfigWith(cfg)

		Expect(cfg.GetString(core.CfgKeyCmdrBinDir)).To(Equal("/cmdr/bin"))
		Expect(cfg.GetString(core.CfgKeyCmdrShimsDir)).To(Equal("/custom/shims"))
		Expect(cfg.GetString(core.CfgKeyCmdrProfileDir)).To(Equal("/cmdr/profile"))
		Expect(cfg.GetString(core.CfgKeyCmdrDatabasePath)).To(Equal("/cmdr/cmdr.db"))
	})

	It("should set only configured proxy env values", func() {
		cfg := viper.New()
		cfg.Set(core.CfgKeyProxyGo, "https://goproxy.example")
		cfg.Set(core.CfgKeyProxyHTTP, "")
		cfg.Set(core.CfgKeyProxyHTTPS, "https://proxy.example")
		values := map[string]string{}

		initProxyWith(cfg, func(key, value string) error {
			values[key] = value
			return nil
		})

		Expect(values).To(Equal(map[string]string{
			"GOPROXY":     "https://goproxy.example",
			"HTTPS_PROXY": "https://proxy.example",
		}))
	})

	It("should execute a simple root command", func() {
		previousCfg := core.GetConfiguration()
		previousFactory := core.GetDatabaseFactory()
		previousLogger := core.GetLogger()
		defer func() {
			core.SetConfiguration(previousCfg)
			core.SetDatabaseFactory(previousFactory)
			core.SetLogger(previousLogger)
			rootCmd.SetArgs(nil)
		}()

		cfg := viper.New()
		cfg.Set(core.CfgKeyCmdrConfigPath, filepath.Join(os.TempDir(), "cmdr-missing-config.yaml"))
		core.SetConfiguration(cfg)
		rootCmd.SetArgs([]string{"version"})

		Expect(func() { ExecuteContext(context.Background()) }).NotTo(Panic())
		Expect(core.GetDatabaseFactory()).NotTo(BeNil())
	})

	It("should create database factories", func() {
		dir, err := os.MkdirTemp("", "cmdr-root-db")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(dir)

		cfg := viper.New()
		cfg.Set(core.CfgKeyCmdrDatabasePath, filepath.Join(dir, "cmdr.db"))
		db, err := makeDatabaseFactory(cfg)()
		Expect(err).NotTo(HaveOccurred())
		Expect(db.Close()).To(Succeed())

		cfg.Set(core.CfgKeyCmdrDatabasePath, filepath.Join(dir, "missing", "cmdr.db"))
		_, err = makeDatabaseFactory(cfg)()
		Expect(err).To(HaveOccurred())
	})
})
