package logging_test

import (
	"bytes"
	"errors"

	. "github.com/ooesili/aurgo/internal/logging"
	"github.com/ooesili/aurgo/test/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	var (
		out      *bytes.Buffer
		realRepo *mocks.Repo
		repo     Repo
	)

	BeforeEach(func() {
		out = &bytes.Buffer{}
		realRepo = &mocks.Repo{}
		repo = NewRepo(realRepo, out)
	})

	Describe("Sync", func() {
		var err error

		BeforeEach(func() {
			realRepo.SyncCall.Returns.Err = errors.New("darn")

			err = repo.Sync("dopepkg")
		})

		It("forwards its arguments to the real repo", func() {
			Expect(realRepo.SyncCall.Received.Pkg).To(Equal("dopepkg"))
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
			realRepo.RemoveCall.Err = errors.New("darn")

			err = repo.Remove("somepkg")
		})

		It("forwards the arguments to the real repo", func() {
			Expect(realRepo.RemoveCall.Removed).To(Equal([]string{"somepkg"}))
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
