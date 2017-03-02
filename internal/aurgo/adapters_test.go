package aurgo_test

import (
	"errors"

	. "github.com/ooesili/aurgo/internal/aurgo"
	"github.com/ooesili/aurgo/test/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RepoVisitor", func() {
	var (
		repo    *mocks.Repo
		visitor VisitorAdapater
		deps    []string
		err     error
	)

	BeforeEach(func() {
		repo = &mocks.Repo{}
		visitor = NewRepoVisitor(repo)
	})

	JustBeforeEach(func() {
		deps, err = visitor.Visit("dopepkg")
	})

	Context("when all repo operations succeed", func() {
		BeforeEach(func() {
			repo.GetDepsCall.Returns.Pkgs = []string{"libdope", "leftpad"}
		})

		It("succeeds", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("syncs the package", func() {
			Expect(repo.SyncCall.Received.Pkg).To(Equal("dopepkg"))
		})

		It("gets the dependencies from the package", func() {
			Expect(repo.GetDepsCall.Received.Pkg).To(Equal("dopepkg"))
		})

		It("returns the results for that package", func() {
			Expect(deps).To(Equal([]string{"libdope", "leftpad"}))
		})
	})

	Context("when the package cannot be synced", func() {
		BeforeEach(func() {
			repo.SyncCall.Returns.Err = errors.New("dang")
		})

		It("returns an error", func() {
			Expect(err).To(HaveOccurred())
		})

		It("does not try to get dependencies", func() {
			Expect(repo.GetDepsCall.Received.Pkg).To(BeZero())
		})
	})

	Context("when getting the dependencies fails", func() {
		BeforeEach(func() {
			repo.GetDepsCall.Returns.Err = errors.New("dang")
		})

		It("returns an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})
})
