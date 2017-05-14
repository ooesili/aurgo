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
})
