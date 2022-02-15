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
})
