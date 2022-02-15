package initializer_test

import (
	"io/fs"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrlyc/cmdr/core/initializer"
)

var _ = Describe("Filesystem", func() {
	Context("DirRemover", func() {
		It("should succeed when dir exists", func() {
			dir, err := os.MkdirTemp("", "")
			Expect(err).To(BeNil())

			dirRemover := initializer.NewDirRemover(dir)
			Expect(dirRemover.Init()).To(BeNil())

			Expect(dir).NotTo(BeADirectory())
		})
	})

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
			files := map[string]os.FileMode{
				filepath.Join("root", "dir", "a.txt"): 0644,
				filepath.Join("root", "b.txt"):        0600,
			}

			for _, path := range dirs {
				Expect(os.MkdirAll(filepath.Join(rootDir, path), 0755)).To(Succeed())
			}

			for path, perm := range files {
				Expect(os.WriteFile(
					filepath.Join(rootDir, path),
					[]byte(perm.String()),
					perm,
				)).To(Succeed())
			}

			exporter := initializer.NewEmbedFSExporter(embedFS, "root", filepath.Join(dstDir, "root"))
			Expect(exporter.Init()).To(Succeed())

			for _, path := range dirs {
				info, err := os.Stat(filepath.Join(dstDir, path))
				Expect(err).To(BeNil())
				Expect(info.IsDir()).To(BeTrue())
			}

			for path, perm := range files {
				info, err := os.Stat(filepath.Join(dstDir, path))
				Expect(err).To(BeNil())
				Expect(info.Mode().Perm()).To(Equal(perm))
				Expect(info.IsDir()).To(BeFalse())

				content, err := os.ReadFile(filepath.Join(dstDir, path))
				Expect(err).To(BeNil())
				Expect(string(content)).To(Equal(perm.String()))
			}
		})
	})
})
