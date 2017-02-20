package srcinfo_test

import (
	"github.com/ooesili/aurgo/internal/cache"
	. "github.com/ooesili/aurgo/internal/srcinfo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Srcinfo", func() {
	Context("given a SRCINFO with a single package", func() {
		var pkg cache.Package

		BeforeEach(func() {
			srcinfo := New("x86_64")

			var err error
			pkg, err = srcinfo.Parse([]byte(fixtureSinglePackge))
			Expect(err).ToNot(HaveOccurred())
		})

		It("it parses the depends field", func() {
			Expect(pkg.Depends).To(Equal([]string{"libdope", "leftpad"}))
		})

		It("parses the checkdepends field", func() {
			Expect(pkg.Checkdepends).To(Equal([]string{"checktool", "testlib"}))
		})

		It("parses the makedepends field", func() {
			Expect(pkg.Makedepends).To(Equal([]string{"cmake", "maven"}))
		})
	})

	Describe("parsing the arch field", func() {
		Context("given a package that supports x86_64 and i686", func() {
			It("succeds when looking for x86_64 packages", func() {
				srcinfo := New("x86_64")
				_, err := srcinfo.Parse([]byte(fixtureMultiArch))
				Expect(err).ToNot(HaveOccurred())
			})

			It("succeds when looking for i686 packages", func() {
				srcinfo := New("i686")
				_, err := srcinfo.Parse([]byte(fixtureMultiArch))
				Expect(err).ToNot(HaveOccurred())
			})

			It("fails when looking for arm packages", func() {
				srcinfo := New("arm")
				_, err := srcinfo.Parse([]byte(fixtureMultiArch))
				Expect(err).To(MatchError("package does not support arm"))
			})
		})

		Context("given a package that supports any architecture", func() {
			It("succeds when looking for x86_64 packages", func() {
				srcinfo := New("x86_64")
				_, err := srcinfo.Parse([]byte(fixtureAnyArch))
				Expect(err).ToNot(HaveOccurred())
			})

			It("succeds when looking for i686 packages", func() {
				srcinfo := New("i686")
				_, err := srcinfo.Parse([]byte(fixtureAnyArch))
				Expect(err).ToNot(HaveOccurred())
			})

			It("succeds when looking for arm packages", func() {
				srcinfo := New("arm")
				_, err := srcinfo.Parse([]byte(fixtureAnyArch))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("architecture specific dependencies", func() {
		var (
			targetArch string
			pkg        cache.Package
		)

		JustBeforeEach(func() {
			srcinfo := New(targetArch)
			var err error
			pkg, err = srcinfo.Parse([]byte(fixtureArchSpecificDeps))
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the architecture has specific dependencies", func() {
			BeforeEach(func() {
				targetArch = "x86_64"
			})

			It("adds architecture-specific depends", func() {
				Expect(pkg.Depends).To(Equal([]string{"leftpad", "libnice64"}))
			})

			It("adds architecture-specific checkdepends", func() {
				Expect(pkg.Checkdepends).To(Equal([]string{"checktool", "testlib"}))
			})

			It("adds architecture-specific makedepends", func() {
				Expect(pkg.Makedepends).To(Equal([]string{"cmake", "maven"}))
			})
		})

		Context("when the architecture does not have specific dependencies", func() {
			BeforeEach(func() {
				targetArch = "i686"
			})

			It("only adds the normally specified depends", func() {
				Expect(pkg.Depends).To(Equal([]string{"leftpad"}))
			})

			It("only adds the normally specified depends", func() {
				Expect(pkg.Checkdepends).To(Equal([]string{"checktool"}))
			})

			It("only adds the normally specified depends", func() {
				Expect(pkg.Makedepends).To(Equal([]string{"cmake"}))
			})
		})
	})

	Context("when given dependencies with version constraints", func() {
		var pkg cache.Package

		BeforeEach(func() {
			srcinfo := New("x86_64")
			var err error
			pkg, err = srcinfo.Parse([]byte(fixtureVersionConstraints))
			Expect(err).ToNot(HaveOccurred())
		})

		It("strips versions on depends", func() {
			Expect(pkg.Depends).To(Equal([]string{
				"dep1", "dep2", "dep3", "dep4", "dep5",
			}))
		})

		It("strips versions on makedepends", func() {
			Expect(pkg.Makedepends).To(Equal([]string{
				"makedep1", "makedep2", "makedep3", "makedep4", "makedep5",
			}))
		})

		It("strips versions on checkdepends", func() {
			Expect(pkg.Checkdepends).To(Equal([]string{
				"checkdep1", "checkdep2", "checkdep3", "checkdep4", "checkdep5",
			}))
		})
	})

	Context("when given a split package", func() {
		It("returns an error", func() {
			srcinfo := New("x86_64")
			_, err := srcinfo.Parse([]byte(fixtureSplitPackage))
			Expect(err).To(MatchError("cannot handle split packages"))
		})
	})
})
