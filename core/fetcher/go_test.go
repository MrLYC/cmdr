package fetcher_test

import (
	"os"
	"path/filepath"

	. "github.com/mrlyc/cmdr/core/fetcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoGetter", func() {
	var (
		installer *GoInstaller
	)

	BeforeEach(func() {
		installer = NewDefaultGoInstaller()
	})

	Context("IsSupport", func() {
		It("should return true", func() {
			Expect(installer.IsSupport("go://github.com/mrlyc/cmdr")).To(BeTrue())
		})
	})

	Context("Fetch", func() {
		var (
			root   string
			dst    string
			goPath string
		)

		BeforeEach(func() {
			var err error
			root, err = os.MkdirTemp("", "cmdr-go-installer")
			Expect(err).NotTo(HaveOccurred())
			dst = filepath.Join(root, "bin")
			Expect(os.MkdirAll(dst, 0755)).To(Succeed())
			goPath = filepath.Join(root, "go")
		})

		AfterEach(func() {
			os.RemoveAll(root)
		})

		writeGo := func(script string) {
			Expect(os.WriteFile(goPath, []byte(script), 0755)).To(Succeed())
			installer = NewGoInstaller(goPath, "go://")
		}

		It("should install direct locations with suffix", func() {
			writeGo("#!/bin/sh\nprintf '%s\\n' \"$3\" >> \"$PWD/calls\"\ntouch \"$PWD/cmdr\"\n")

			Expect(installer.Fetch("cmdr", "1.0.0", "go://github.com/mrlyc/cmdr@latest", dst)).To(Succeed())
			Expect(filepath.Join(dst, "cmdr")).To(BeAnExistingFile())
			content, err := os.ReadFile(filepath.Join(dst, "calls"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal("github.com/mrlyc/cmdr@latest\n"))
		})

		It("should detect a working version suffix", func() {
			writeGo("#!/bin/sh\nprintf '%s\\n' \"$3\" >> \"$PWD/calls\"\nif [ \"$3\" = 'github.com/mrlyc/cmdr@v1.2.3' ]; then touch \"$PWD/cmdr\"; exit 0; fi\nexit 1\n")

			Expect(installer.Fetch("cmdr", "1.2.3", "go://github.com/mrlyc/cmdr", dst)).To(Succeed())
			content, err := os.ReadFile(filepath.Join(dst, "calls"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(content)).To(Equal("github.com/mrlyc/cmdr@1.2.3\ngithub.com/mrlyc/cmdr@v1.2.3\n"))
		})

		It("should return install errors", func() {
			writeGo("#!/bin/sh\nexit 1\n")

			Expect(installer.Fetch("cmdr", "1.0.0", "go://github.com/mrlyc/cmdr@latest", dst)).To(HaveOccurred())
		})
	})
})
