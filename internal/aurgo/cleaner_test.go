package aurgo_test

import (
	"errors"

	. "github.com/ooesili/aurgo/internal/aurgo"
	"github.com/ooesili/aurgo/test/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("RepoCleaner", func() {
	var (
		repo    *mocks.Repo
		cleaner RepoCleaner
	)

	BeforeEach(func() {
		repo = &mocks.Repo{}
		cleaner = NewRepoCleaner(repo)
	})

	DescribeTable("cleaning unused packages",
		func(used, existing, expectedCleaned []string) {
			repo.ListCall.Returns.Pkgs = existing

			Expect(cleaner.Clean(used)).To(Succeed())
			Expect(repo.RemoveCall.Removed).To(Equal(expectedCleaned))
		},

		Entry("no packages",
			[]string{},
			[]string{},
			nil,
		),

		Entry("one used package",
			[]string{"dopepkg"},
			[]string{"dopepkg"},
			nil,
		),

		Entry("one unused existing package",
			[]string{},
			[]string{"dopepkg"},
			[]string{"dopepkg"},
		),

		Entry("many used packages",
			[]string{"dopepkg", "libdope", "leftpad"},
			[]string{"dopepkg", "libdope", "leftpad"},
			nil,
		),

		Entry("one unused in many packages",
			[]string{"dopepkg", "libdope"},
			[]string{"dopepkg", "libdope", "leftpad"},
			[]string{"leftpad"},
		),
	)

	Describe("failure", func() {
		Context("when existing packages cannot be listed", func() {
			It("erturns an error", func() {
				repo.ListCall.Returns.Err = errors.New("gosh")

				err := cleaner.Clean([]string{"dopepkg"})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when a package cannot be removed", func() {
			It("returns an error", func() {
				repo.ListCall.Returns.Pkgs = []string{"dopepkg", "leftpad"}
				repo.RemoveCall.Err = errors.New("darn")

				err := cleaner.Clean([]string{"dopepkg"})
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
