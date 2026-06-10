package utils

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pathlib", func() {
	var root string

	BeforeEach(func() {
		var err error
		root, err = os.MkdirTemp("", "cmdr-path")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(root)
	})

	It("should create, derive and expose paths", func() {
		helper := NewPathHelper(root)
		Expect(helper.Path()).To(Equal(root))
		Expect(helper.Child("child").Path()).To(Equal(filepath.Join(root, "child")))
		Expect(helper.Child("child").Parent().Path()).To(Equal(root))
		Expect(helper.Child("dir").MkdirAll(0750)).To(Succeed())
		Expect(filepath.Join(root, "dir")).To(BeADirectory())
	})

	It("should detect existing and missing paths", func() {
		helper := NewPathHelper(root)
		Expect(os.WriteFile(filepath.Join(root, "file"), []byte("x"), 0644)).To(Succeed())

		Expect(helper.Exists("file")).To(Succeed())
		Expect(helper.Exists("missing")).To(HaveOccurred())
		absPath, err := helper.AbsPath("file")
		Expect(err).NotTo(HaveOccurred())
		Expect(absPath).To(Equal(filepath.Join(root, "file")))
		_, err = helper.AbsPath("missing")
		Expect(err).To(HaveOccurred())
	})

	It("should remove files, directories and broken symlinks", func() {
		helper := NewPathHelper(root)
		Expect(os.WriteFile(filepath.Join(root, "file"), []byte("x"), 0644)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(root, "dir", "nested"), 0755)).To(Succeed())
		Expect(os.Symlink(filepath.Join(root, "missing"), filepath.Join(root, "broken"))).To(Succeed())

		Expect(helper.EnsureNotExists("file")).To(Succeed())
		Expect(filepath.Join(root, "file")).NotTo(BeAnExistingFile())
		Expect(helper.EnsureNotExists("dir")).To(Succeed())
		Expect(filepath.Join(root, "dir")).NotTo(BeADirectory())
		Expect(helper.EnsureNotExists("broken")).To(Succeed())
		Expect(helper.EnsureNotExists("missing")).To(Succeed())
	})

	It("should chmod, copy and symlink files", func() {
		helper := NewPathHelper(root)
		src := filepath.Join(root, "source")
		Expect(os.WriteFile(src, []byte("data"), 0644)).To(Succeed())

		Expect(helper.CopyFile("copies/target", src, 0755)).To(Succeed())
		copied := filepath.Join(root, "copies", "target")
		Expect(copied).To(BeAnExistingFile())
		info, err := os.Stat(copied)
		Expect(err).NotTo(HaveOccurred())
		Expect(info.Mode().Perm()).To(Equal(os.FileMode(0755)))

		Expect(helper.Chmod("copies/target", 0644)).To(Succeed())
		info, err = os.Stat(copied)
		Expect(err).NotTo(HaveOccurred())
		Expect(info.Mode().Perm()).To(Equal(os.FileMode(0644)))

		Expect(helper.SymbolLink("link", src, 0755)).To(Succeed())
		realPath, err := helper.RealPath("link")
		Expect(err).NotTo(HaveOccurred())
		Expect(realPath).To(Equal(src))
	})

	It("should return filesystem operation errors", func() {
		helper := NewPathHelper(root)
		blocker := filepath.Join(root, "blocker")
		Expect(os.WriteFile(blocker, []byte("x"), 0644)).To(Succeed())
		Expect(NewPathHelper(blocker).MkdirAll(0755)).To(HaveOccurred())
		Expect(helper.Chmod("missing", 0644)).To(HaveOccurred())
		Expect(helper.CopyFile("copy", filepath.Join(root, "missing"), 0644)).To(HaveOccurred())
		Expect(helper.SymbolLink("link", filepath.Join(root, "missing"), 0755)).To(HaveOccurred())
	})

	It("should return real paths for normal files and errors for bad links", func() {
		helper := NewPathHelper(root)
		path := filepath.Join(root, "file")
		Expect(os.WriteFile(path, []byte("x"), 0644)).To(Succeed())

		realPath, err := helper.RealPath("file")
		Expect(err).NotTo(HaveOccurred())
		Expect(realPath).To(Equal(path))

		_, err = helper.RealPath("missing")
		Expect(err).To(HaveOccurred())
	})
})
