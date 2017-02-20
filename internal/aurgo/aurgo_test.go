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
			pacman *mocks.Pacman
			err    error
		)

		BeforeEach(func() {
			config = &mocks.Config{}
			cache = &mocks.Cache{}
			pacman = &mocks.Pacman{}
		})

		JustBeforeEach(func() {
			aurgo = New(config, cache, pacman)
			err = aurgo.SyncAll()
		})

		Context("with a single package with no dependencies", func() {
			BeforeEach(func() {
				cache.GetDepsCall.DepMap = map[string][]string{
					"dopepkg": {},
				}
				config.PackagesCall.Returns.Packages = []string{"dopepkg"}
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("syncs the package in the cache", func() {
				Expect(cache.SyncCall.SyncedPackages).To(ConsistOf("dopepkg"))
			})
		})

		Context("with a package with a dependency", func() {
			BeforeEach(func() {
				cache.GetDepsCall.DepMap = map[string][]string{
					"dopepkg": {"libdope"},
					"libdope": {},
				}
				config.PackagesCall.Returns.Packages = []string{"dopepkg"}
			})

			It("syncs the dependencies", func() {
				Expect(cache.SyncCall.SyncedPackages).To(ConsistOf(
					"dopepkg", "libdope",
				))
			})
		})

		Context("with a package with a transitive dependency", func() {
			BeforeEach(func() {
				cache.GetDepsCall.DepMap = map[string][]string{
					"dopepkg": {"libdope"},
					"libdope": {"leftpad"},
					"leftpad": {},
				}
				config.PackagesCall.Returns.Packages = []string{"dopepkg"}
			})

			It("syncs the all dependencies", func() {
				Expect(cache.SyncCall.SyncedPackages).To(ConsistOf(
					"dopepkg", "libdope", "leftpad",
				))
			})
		})

		Context("when a diamond dependency exists", func() {
			BeforeEach(func() {
				cache.GetDepsCall.DepMap = map[string][]string{
					"dopepkg": {"libdope", "libcool"},
					"libdope": {"leftpad"},
					"libcool": {"leftpad"},
					"leftpad": {},
				}
				config.PackagesCall.Returns.Packages = []string{"dopepkg"}
			})

			It("syncs each dependency exactly once", func() {
				Expect(cache.SyncCall.SyncedPackages).To(ConsistOf(
					"dopepkg", "libdope", "libcool", "leftpad",
				))
			})
		})

		Context("when a transitive dependency is also explicitly dependended on", func() {
			BeforeEach(func() {
				cache.GetDepsCall.DepMap = map[string][]string{
					"dopepkg": {"libdope"},
					"libdope": {"leftpad"},
					"leftpad": {},
				}
				config.PackagesCall.Returns.Packages = []string{"dopepkg", "leftpad"}
			})

			It("syncs each dependency exactly once", func() {
				Expect(cache.SyncCall.SyncedPackages).To(ConsistOf(
					"dopepkg", "libdope", "leftpad",
				))
			})
		})

		Context("when some dependencies are already available", func() {
			BeforeEach(func() {
				cache.GetDepsCall.DepMap = map[string][]string{
					"dopepkg": {"libdope", "libcool"},
					"libdope": {"leftpad"},
					"leftpad": {"openssl"},
				}
				config.PackagesCall.Returns.Packages = []string{"dopepkg"}
				pacman.ListAvailableCall.Returns.Packages = []string{"libcool", "openssl"}
			})

			It("skips those packages", func() {
				Expect(cache.SyncCall.SyncedPackages).To(ConsistOf(
					"dopepkg",
					"libdope",
					"leftpad",
				))
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
				cache.SyncCall.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when getting package dependencies fails", func() {
			BeforeEach(func() {
				config.PackagesCall.Returns.Packages = []string{"dopepkg"}
				cache.GetDepsCall.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
