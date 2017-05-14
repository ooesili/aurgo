package sys_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/ooesili/aurgo/internal/sys"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filesystem", func() {
	var (
		filesystem Filesystem
		tempDir    string
	)

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "aurgo-test-filesystem-")
		Expect(err).ToNot(HaveOccurred())

		filesystem = NewFilesystem()
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("Exists", func() {
		var (
			leadingDir string
			filename   string
			exists     bool
			err        error
		)

		BeforeEach(func() {
			leadingDir = filepath.Join(tempDir, "dir")
			Expect(os.Mkdir(leadingDir, 0755)).To(Succeed())

			filename = filepath.Join(leadingDir, "file")
			Expect(ioutil.WriteFile(filename, nil, 0644)).To(Succeed())
		})

		JustBeforeEach(func() {
			exists, err = filesystem.Exists(filename)
		})

		Context("when the file exists", func() {
			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns true", func() {
				Expect(exists).To(BeTrue())
			})
		})

		Context("when the file does not exist", func() {
			BeforeEach(func() {
				Expect(os.Remove(filename)).To(Succeed())
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns false", func() {
				Expect(exists).To(BeFalse())
			})
		})

		Context("when there is an error checking the file", func() {
			BeforeEach(func() {
				Expect(os.RemoveAll(leadingDir)).To(Succeed())
				Expect(ioutil.WriteFile(leadingDir, nil, 0644)).To(Succeed())
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ListFiles", func() {
		var (
			dirname string
			err     error
			files   []string
		)

		BeforeEach(func() {
			dirname = ""
		})

		JustBeforeEach(func() {
			files, err = filesystem.ListFiles(dirname)
		})

		Context("with no files", func() {
			BeforeEach(func() {
				dirname = tempDir
			})

			It("returns an empty result", func() {
				Expect(files).To(BeEmpty())
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("with multiple files", func() {
			BeforeEach(func() {
				dirname = tempDir
			})

			BeforeEach(func() {
				path1 := filepath.Join(tempDir, "hello")
				Expect(ioutil.WriteFile(path1, nil, 0600)).To(Succeed())

				path2 := filepath.Join(tempDir, "friend")
				Expect(ioutil.WriteFile(path2, nil, 0600)).To(Succeed())
			})

			It("returns the names of those files", func() {
				Expect(files).To(ConsistOf("hello", "friend"))
			})
		})

		Context("when the directory cannot be read", func() {
			BeforeEach(func() {
				path := filepath.Join(tempDir, "notadir")
				Expect(ioutil.WriteFile(path, nil, 0600)).To(Succeed())
				dirname = path
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
