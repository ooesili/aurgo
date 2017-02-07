package aurgo_test

import (
	"errors"

	. "github.com/ooesili/aurgo/internal/aurgo"
	"github.com/ooesili/aurgo/test/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aurgo", func() {
	Describe("Sync", func() {
		var (
			aurgo  Aurgo
			config *mocks.Config
			git    *mocks.Git
			err    error
		)

		BeforeEach(func() {
			config = &mocks.Config{}
			git = &mocks.Git{}
		})

		JustBeforeEach(func() {
			aurgo = New(config, git)
			err = aurgo.Sync()
		})

		Context("when all dependencies succeed", func() {
			BeforeEach(func() {
				config.PackagesCall.Returns.Packages = []string{"package1"}
				config.SourcePathCall.Returns.Path = "/path/to/package1"
				config.AurRepoURLCall.Returns.URL = "https://example.com/package1.git"
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("looks up the source path to each package", func() {
				Expect(config.SourcePathCall.Received.Package).To(Equal("package1"))
			})

			It("looks up the aur repo URL to each package", func() {
				Expect(config.AurRepoURLCall.Received.Package).To(Equal("package1"))
			})

			It("clones the git repo from the AUR", func() {
				Expect(git.CloneCall.Received.Path).To(Equal("/path/to/package1"))
				Expect(git.CloneCall.Received.URL).To(Equal("https://example.com/package1.git"))
			})
		})

		Context("when listing the packages fails", func() {
			BeforeEach(func() {
				config.PackagesCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when looking up the source path for a package fails", func() {
			BeforeEach(func() {
				config.PackagesCall.Returns.Packages = []string{"package1"}
				config.SourcePathCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when cloning a package fails", func() {
			BeforeEach(func() {
				config.PackagesCall.Returns.Packages = []string{"package1"}
				config.SourcePathCall.Returns.Path = "/path/to/package1"
				config.AurRepoURLCall.Returns.URL = "https://example.com/package1.git"
				git.CloneCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
