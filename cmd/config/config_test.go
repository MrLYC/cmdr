package config

import (
	"io"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"logur.dev/logur"

	"github.com/mrlyc/cmdr/core"
)

var _ = Describe("Config", func() {
	captureStdout := func(fn func()) string {
		old := os.Stdout
		r, w, err := os.Pipe()
		Expect(err).NotTo(HaveOccurred())
		os.Stdout = w
		fn()
		Expect(w.Close()).To(Succeed())
		os.Stdout = old
		out, err := io.ReadAll(r)
		Expect(err).NotTo(HaveOccurred())
		return string(out)
	}

	It("should render a config value as yaml", func() {
		cfg := core.NewConfiguration()
		cfg.Set("download.replace", []map[string]string{
			{"match": "https://(.*)", "template": "mirror://$1"},
		})

		out, err := renderConfigValue(cfg, "download.replace")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(out)).To(ContainSubstring("match:"))
		Expect(string(out)).To(ContainSubstring("template:"))
	})

	It("should filter private settings", func() {
		settings := publicSettings(map[string]interface{}{
			"core": map[string]interface{}{"root_dir": "/tmp/cmdr"},
			"_":    map[string]interface{}{"config": "hidden"},
		})

		Expect(settings).To(HaveKey("core"))
		Expect(settings).NotTo(HaveKey("_"))
	})

	It("should create a config file with yaml values", func() {
		root, err := os.MkdirTemp("", "cmdr-config")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(root)

		configFile := filepath.Join(root, "nested", "config.yaml")
		Expect(setConfigValue(configFile, "core.root_dir", "/opt/cmdr", logur.NoopLogger{})).To(Succeed())

		info, err := os.Stat(filepath.Dir(configFile))
		Expect(err).NotTo(HaveOccurred())
		Expect(info.IsDir()).To(BeTrue())
		Expect(info.Mode().Perm() & 0100).NotTo(Equal(os.FileMode(0)))

		cfg := core.NewConfiguration()
		cfg.SetConfigFile(configFile)
		Expect(cfg.ReadInConfig()).To(Succeed())
		Expect(cfg.GetString("core.root_dir")).To(Equal("/opt/cmdr"))
	})

	It("should return invalid yaml errors", func() {
		root, err := os.MkdirTemp("", "cmdr-config")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(root)

		configFile := filepath.Join(root, "config.yaml")
		Expect(setConfigValue(configFile, "core.root_dir", "[", logur.NoopLogger{})).To(HaveOccurred())
	})

	It("should run get and list commands", func() {
		previous := core.GetConfiguration()
		defer core.SetConfiguration(previous)

		cfg := core.NewConfiguration()
		cfg.Set("core.root_dir", "/tmp/cmdr")
		cfg.Set("_", map[string]interface{}{"hidden": true})
		cfg.Set(core.CfgKeyXConfigGetKey, "core.root_dir")
		core.SetConfiguration(cfg)

		Expect(captureStdout(func() { getCmd.Run(getCmd, nil) })).To(Equal("/tmp/cmdr\n"))
		Expect(captureStdout(func() { listCmd.Run(listCmd, nil) })).To(ContainSubstring("core:"))
	})

	It("should run set command", func() {
		previous := core.GetConfiguration()
		defer core.SetConfiguration(previous)

		root, err := os.MkdirTemp("", "cmdr-config")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(root)

		cfg := core.NewConfiguration()
		configFile := filepath.Join(root, "config.yaml")
		cfg.SetConfigFile(configFile)
		cfg.Set(core.CfgKeyXConfigSetKey, "core.root_dir")
		cfg.Set(core.CfgKeyXConfigSetValue, "/tmp/cmdr")
		core.SetConfiguration(cfg)

		Expect(func() { setCmd.Run(setCmd, nil) }).NotTo(Panic())
		Expect(configFile).To(BeAnExistingFile())
	})
})
