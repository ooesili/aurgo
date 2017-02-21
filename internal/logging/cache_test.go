package logging_test

import (
	"bytes"
	"errors"

	. "github.com/ooesili/aurgo/internal/logging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	var (
		out       *bytes.Buffer
		realCache *MockCache
		cache     Cache
	)

	BeforeEach(func() {
		out = &bytes.Buffer{}
		realCache = &MockCache{}
		cache = NewCache(realCache, out)
	})

	Describe("GetDeps", func() {
		var (
			pkgs []string
			err  error
		)

		BeforeEach(func() {
			realCache.GetDepsCall.Returns.Pkgs = []string{"pkg1", "pkg2"}
			realCache.GetDepsCall.Returns.Err = errors.New("shoot")

			pkgs, err = cache.GetDeps("somepkg")
		})

		It("forwards its arguments to the real cache", func() {
			Expect(realCache.GetDepsCall.Received.Pkg).To(Equal("somepkg"))
		})

		It("forwards the return values back to the caller", func() {
			Expect(pkgs).To(Equal([]string{"pkg1", "pkg2"}))
			Expect(err).To(MatchError("shoot"))
		})
	})

	Describe("Sync", func() {
		var err error

		BeforeEach(func() {
			realCache.SyncCall.Returns.Err = errors.New("darn")

			err = cache.Sync("otherpkg")
		})

		It("forwards its arguments to the real cache", func() {
			Expect(realCache.SyncCall.Received.Pkg).To(Equal("otherpkg"))
		})

		It("forwards the return values back to the caller", func() {
			Expect(err).To(MatchError("darn"))
		})

		It("logs the package its about to sync", func() {
			expectedMessage := "---> syncing package: otherpkg\n"
			Expect(out.String()).To(Equal(expectedMessage))
		})
	})
})

type MockCache struct {
	GetDepsCall struct {
		Received struct {
			Pkg string
		}
		Returns struct {
			Pkgs []string
			Err  error
		}
	}

	SyncCall struct {
		Received struct {
			Pkg string
		}
		Returns struct {
			Err error
		}
	}
}

func (c *MockCache) GetDeps(pkg string) ([]string, error) {
	c.GetDepsCall.Received.Pkg = pkg
	returns := c.GetDepsCall.Returns
	return returns.Pkgs, returns.Err
}

func (c *MockCache) Sync(pkg string) error {
	c.SyncCall.Received.Pkg = pkg
	return c.SyncCall.Returns.Err
}
