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
			depWalker *MockDepWalker
			cleaner   *MockCleaner
			config    *mocks.Config
			aurgo     Aurgo
			err       error
		)

		BeforeEach(func() {
			depWalker = &MockDepWalker{}
			cleaner = &MockCleaner{}
			config = &mocks.Config{}
			aurgo = New(depWalker, cleaner, config)

			config.PackagesCall.Returns.Packages = []string{
				"dopepkg", "otherpkg",
			}
		})

		JustBeforeEach(func() {
			err = aurgo.SyncAll()
		})

		Context("when all subsystems succeed", func() {
			BeforeEach(func() {
				depWalker.WalkCall.Returns.Pkgs = []string{
					"dopepkg", "libdope", "otherpkg",
				}
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("walks the packages from the config for dependencies", func() {
				Expect(depWalker.WalkCall.Received.Pkgs).To(Equal([]string{
					"dopepkg", "otherpkg",
				}))
			})

			It("cleans up unused packages", func() {
				Expect(cleaner.CleanCall.Received.Pkgs).To(Equal([]string{
					"dopepkg", "libdope", "otherpkg",
				}))
			})
		})

		Context("when the dependencies cannot be resolved", func() {
			BeforeEach(func() {
				depWalker.WalkCall.Returns.Err = errors.New("darn")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when unused packages cannot be cleaned", func() {
			BeforeEach(func() {
				depWalker.WalkCall.Returns.Pkgs = []string{
					"dopepkg", "libdope", "otherpkg",
				}
				cleaner.CleanCall.Returns.Err = errors.New("darn")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})

type MockDepWalker struct {
	WalkCall struct {
		Received struct {
			Pkgs []string
		}
		Returns struct {
			Pkgs []string
			Err  error
		}
	}
}

func (m *MockDepWalker) Walk(pkgs []string) ([]string, error) {
	m.WalkCall.Received.Pkgs = pkgs
	returns := m.WalkCall.Returns
	return returns.Pkgs, returns.Err
}

type MockCleaner struct {
	CleanCall struct {
		Received struct {
			Pkgs []string
		}
		Returns struct {
			Err error
		}
	}
}

func (m *MockCleaner) Clean(pkgs []string) error {
	m.CleanCall.Received.Pkgs = pkgs
	return m.CleanCall.Returns.Err
}
