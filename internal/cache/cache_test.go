package cache_test

import (
	"errors"

	. "github.com/ooesili/aurgo/internal/cache"
	"github.com/ooesili/aurgo/test/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	var (
		cache  Cache
		config *mocks.Config
		git    *mocks.Git
		err    error
	)

	BeforeEach(func() {
		config = &mocks.Config{}
		git = &mocks.Git{}
	})

	JustBeforeEach(func() {
		cache = New(config, git)
	})

	Describe("Sync", func() {
		JustBeforeEach(func() {
			err = cache.Sync("coolpkg")
		})

		Context("when all dependencies succeed", func() {
			BeforeEach(func() {
				config.SourcePathCall.Returns.Path = "/path/to/coolpkg"
				config.AurRepoURLCall.Returns.URL = "https://example.com/coolpkg.git"
			})

			It("looks up the source path for the package", func() {
				Expect(config.SourcePathCall.Received.Package).To(Equal("coolpkg"))
			})

			It("looks up the aur repo URL for the package", func() {
				Expect(config.AurRepoURLCall.Received.Package).To(Equal("coolpkg"))
			})

			It("cones the git repo from the URL", func() {
				Expect(git.CloneCall.Received.Path).To(Equal("/path/to/coolpkg"))
				Expect(git.CloneCall.Received.URL).To(Equal("https://example.com/coolpkg.git"))
			})
		})

		Context("when looking up the source path for the package fails", func() {
			BeforeEach(func() {
				config.SourcePathCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when cloning the package fails", func() {
			BeforeEach(func() {
				config.SourcePathCall.Returns.Path = "/path/to/coolpkg"
				config.AurRepoURLCall.Returns.URL = "https://example.com/coolpkg.git"
				git.CloneCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
