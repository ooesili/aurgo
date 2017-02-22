package logging_test

import (
	"bytes"
	"errors"

	"github.com/ooesili/aurgo/internal/aurgo"
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

	Describe("Sync", func() {
		var err error

		BeforeEach(func() {
			realCache.SyncCall.Returns.Err = errors.New("darn")

			err = cache.Sync("dopepkg")
		})

		It("forwards its arguments to the real cache", func() {
			Expect(realCache.SyncCall.Received.Pkg).To(Equal("dopepkg"))
		})

		It("forwards the return values back to the caller", func() {
			Expect(err).To(MatchError("darn"))
		})

		It("logs the package its about to sync", func() {
			expectedMessage := "---> syncing package: dopepkg\n"
			Expect(out.String()).To(Equal(expectedMessage))
		})
	})

	Describe("Remove", func() {
		var err error

		BeforeEach(func() {
			realCache.RemoveCall.Returns.Err = errors.New("darn")

			err = cache.Remove("somepkg")
		})

		It("forwards the arguments to the real cache", func() {
			Expect(realCache.RemoveCall.Received.Pkg).To(Equal("somepkg"))
		})

		It("forwards the return values back to the caller", func() {
			Expect(err).To(MatchError("darn"))
		})

		It("logs the package its about to remove", func() {
			expectedMessage := "---> removing package: somepkg\n"
			Expect(out.String()).To(Equal(expectedMessage))
		})
	})
})

type MockCache struct {
	aurgo.Cache

	SyncCall struct {
		Received struct {
			Pkg string
		}
		Returns struct {
			Err error
		}
	}

	RemoveCall struct {
		Received struct {
			Pkg string
		}
		Returns struct {
			Err error
		}
	}
}

func (c *MockCache) Sync(pkg string) error {
	c.SyncCall.Received.Pkg = pkg
	return c.SyncCall.Returns.Err
}

func (c *MockCache) Remove(pkg string) error {
	c.RemoveCall.Received.Pkg = pkg
	return c.RemoveCall.Returns.Err
}
