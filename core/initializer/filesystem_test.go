package initializer_test

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"github.com/mrlyc/cmdr/core"
	"github.com/mrlyc/cmdr/core/initializer"
)

var _ = Describe("Filesystem", func() {
	Context("EmbedFSExporter", func() {
		var (
			embedFS fs.FS
			rootDir string
			dstDir  string
		)

		BeforeEach(func() {
			var err error

			rootDir, err = os.MkdirTemp("", "")
			Expect(err).To(BeNil())

			dstDir, err = os.MkdirTemp("", "")
			Expect(err).To(BeNil())

			embedFS = os.DirFS(rootDir)
		})

		AfterEach(func() {
			Expect(os.RemoveAll(rootDir)).To(Succeed())
			Expect(os.RemoveAll(dstDir)).To(Succeed())
		})

		It("should export dir structure", func() {
			dirs := []string{
				filepath.Join("root", "empty_dir"),
				filepath.Join("root", "dir"),
			}
			files := []string{
				filepath.Join("root", "dir", "a.txt"),
				filepath.Join("root", "b.txt"),
			}

			for _, path := range dirs {
				Expect(os.MkdirAll(filepath.Join(rootDir, path), 0755)).To(Succeed())
			}

			for i, path := range files {
				Expect(os.WriteFile(
					filepath.Join(rootDir, path),
					[]byte(fmt.Sprintf("file %d", i)),
					0644,
				)).To(Succeed())
			}

			exporter := initializer.NewEmbedFSExporter(embedFS, "root", filepath.Join(dstDir, "root"), 0644)
			Expect(exporter.Init()).To(Succeed())

			for _, path := range dirs {
				info, err := os.Stat(filepath.Join(dstDir, path))
				Expect(err).To(BeNil())
				Expect(info.IsDir()).To(BeTrue())
			}

			for i, path := range files {
				info, err := os.Stat(filepath.Join(dstDir, path))
				Expect(err).To(BeNil())
				Expect(info.IsDir()).To(BeFalse())

				content, err := os.ReadFile(filepath.Join(dstDir, path))
				Expect(err).To(BeNil())
				Expect(string(content)).To(Equal(fmt.Sprintf("file %d", i)))
			}
		})
	})

	Context("FSBackup", func() {
		var (
			rootDir string
			path    string
		)

		BeforeEach(func() {
			var err error

			rootDir, err = os.MkdirTemp("", "")
			Expect(err).To(BeNil())

			path = filepath.Join(rootDir, "dir", "a.txt")
			Expect(os.MkdirAll(filepath.Dir(path), 0755)).To(Succeed())
			Expect(os.WriteFile(path, []byte("hello"), 0644)).To(Succeed())
		})

		AfterEach(func() {
			Expect(os.RemoveAll(rootDir)).To(Succeed())
		})

		It("should backup a dir", func() {
			backup := initializer.NewFSBackup(rootDir)
			Expect(backup.Init()).To(Succeed())

			target := backup.Target()
			Expect(filepath.Join(target, "dir", "a.txt")).To(BeAnExistingFile())
		})

		It("should backup a file", func() {
			backup := initializer.NewFSBackup(path)
			Expect(backup.Init()).To(Succeed())

			target := backup.Target()
			Expect(filepath.Join(target, "a.txt")).To(BeAnExistingFile())
		})

		It("should ok when path not exists", func() {
			backup := initializer.NewFSBackup(filepath.Join(rootDir, "not_exists"))
			Expect(backup.Init()).To(Succeed())
		})
	})

	Context("DirRender", func() {
		var (
			rootDir string
			cfg     core.Configuration
		)

		BeforeEach(func() {
			var err error

			rootDir, err = os.MkdirTemp("", "")
			Expect(err).To(BeNil())

			cfg = viper.New()
			cfg.Set("config.key", "hello")
		})

		AfterEach(func() {
			Expect(os.RemoveAll(rootDir)).To(Succeed())
		})

		It("should render a file", func() {
			Expect(os.WriteFile(filepath.Join(rootDir, `{{.GetString "config.key"}}.txt.gotmpl`), []byte(`{{.GetString "config.key"}}`), 0644)).To(Succeed())

			render := initializer.NewDirRender(rootDir, ".gotmpl", cfg)
			Expect(render.Init()).To(Succeed())

			Expect(filepath.Join(rootDir, "hello.txt")).To(BeAnExistingFile())
			Expect(filepath.Join(rootDir, `{{.GetString "config.key"}}.txt.gotmpl`)).NotTo(BeARegularFile())

			content, err := ioutil.ReadFile(filepath.Join(rootDir, "hello.txt"))
			Expect(err).To(BeNil())
			Expect(string(content)).To(Equal("hello"))
		})

		It("should render a dir", func() {
			Expect(os.MkdirAll(filepath.Join(rootDir, `{{.GetString "config.key"}}.gotmpl`), 0755)).To(Succeed())

			render := initializer.NewDirRender(rootDir, ".gotmpl", cfg)
			Expect(render.Init()).To(Succeed())

			Expect(filepath.Join(rootDir, "hello")).To(BeADirectory())
		})

		It("should render nothing", func() {
			Expect(os.WriteFile(filepath.Join(rootDir, "a.txt"), []byte(`{{.GetString "config.key"}}`), 0644)).To(Succeed())

			render := initializer.NewDirRender(rootDir, ".gotmpl", cfg)
			Expect(render.Init()).To(Succeed())

			Expect(filepath.Join(rootDir, "a.txt")).To(BeAnExistingFile())
			Expect(filepath.Join(rootDir, "a.txt.gotmpl")).NotTo(BeARegularFile())
		})

		It("should render deep structure", func() {
			dir := filepath.Join(rootDir, "dir")
			Expect(os.MkdirAll(dir, 0755)).To(Succeed())

			Expect(os.WriteFile(filepath.Join(dir, `{{.GetString "config.key"}}.txt.gotmpl`), []byte(`{{.GetString "config.key"}}`), 0644)).To(Succeed())
			Expect(os.MkdirAll(filepath.Join(dir, `{{.GetString "config.key"}}.gotmpl`), 0755)).To(Succeed())

			render := initializer.NewDirRender(rootDir, ".gotmpl", cfg)
			Expect(render.Init()).To(Succeed())

			Expect(filepath.Join(dir, "hello.txt")).To(BeAnExistingFile())
			Expect(filepath.Join(dir, "hello")).To(BeADirectory())
		})
	})
})
