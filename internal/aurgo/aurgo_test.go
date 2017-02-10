package aurgo_test

import (
	"errors"

	. "github.com/ooesili/aurgo/internal/aurgo"
	"github.com/ooesili/aurgo/test/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Aurgo", func() {
	Describe("SyncAll", func() {
		var (
			aurgo  Aurgo
			config *mocks.Config
			cache  *mocks.Cache
			err    error
		)

		BeforeEach(func() {
			config = &mocks.Config{}
			cache = &mocks.Cache{}
		})

		JustBeforeEach(func() {
			aurgo = New(config, cache)
			err = aurgo.SyncAll()
		})

		Context("when all dependencies succeed", func() {
			BeforeEach(func() {
				config.PackagesCall.Returns.Packages = []string{"dopepkg"}
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("syncs the package in the cache", func() {
				Expect(cache.SyncCall.Received.Package).To(Equal("dopepkg"))
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

		Context("when syncing a package fails", func() {
			BeforeEach(func() {
				config.PackagesCall.Returns.Packages = []string{"dopepkg"}
				cache.SyncCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
