package aurgo_test

import (
	"fmt"

	. "github.com/ooesili/aurgo/internal/aurgo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("VisitingWalker", func() {
	DescribeTable("successful dependency walking",
		func(pkgs []string, depMap map[string][]string, expected []string) {
			visitor := &MockVisitor{}
			visitor.VisitCall.DepMap = depMap

			visitingWalker := NewVisitingDepWalker(visitor)
			visited, err := visitingWalker.Walk(pkgs)

			Expect(err).ToNot(HaveOccurred())
			Expect(visited).To(Equal(expected))
		},

		Entry("no packages",
			[]string{},
			map[string][]string{},
			nil,
		),

		Entry("no dependencies",
			[]string{"dopepkg"},
			map[string][]string{
				"dopepkg": {},
			},
			[]string{"dopepkg"},
		),

		Entry("multiple packages",
			[]string{"dopepkg", "otherpkg"},
			map[string][]string{
				"dopepkg":  {},
				"otherpkg": {},
			},
			[]string{"dopepkg", "otherpkg"},
		),

		Entry("single dependency",
			[]string{"dopepkg"},
			map[string][]string{
				"dopepkg": {"libdope"},
				"libdope": {},
			},
			[]string{"dopepkg", "libdope"},
		),

		Entry("multiple dependencies",
			[]string{"dopepkg"},
			map[string][]string{
				"dopepkg": {"libdope", "leftpad"},
				"libdope": {},
				"leftpad": {},
			},
			[]string{"dopepkg", "libdope", "leftpad"},
		),

		Entry("transitive dependencies",
			[]string{"dopepkg"},
			map[string][]string{
				"dopepkg": {"libdope"},
				"libdope": {"leftpad"},
				"leftpad": {},
			},
			[]string{"dopepkg", "libdope", "leftpad"},
		),

		Entry("diamond dependencies",
			[]string{"dopepkg"},
			map[string][]string{
				"dopepkg": {"libdope", "libcool"},
				"libdope": {"leftpad"},
				"libcool": {"leftpad"},
				"leftpad": {},
			},
			[]string{"dopepkg", "libdope", "leftpad", "libcool"},
		),

		Entry("transitive dependencies also explicitly dependended on",
			[]string{"dopepkg", "leftpad"},
			map[string][]string{
				"dopepkg": {"libdope"},
				"libdope": {"leftpad"},
				"leftpad": {},
			},
			[]string{"dopepkg", "libdope", "leftpad"},
		),

		Entry("multiple dependencies with a diamond dependency",
			[]string{"dopepkg", "otherpkg"},
			map[string][]string{
				"dopepkg":  {"dopelib"},
				"otherpkg": {"dopelib", "otherlib"},
				"dopelib":  {},
				"otherlib": {},
			},
			[]string{"dopepkg", "dopelib", "otherpkg", "otherlib"},
		),
	)

	DescribeTable("failure",
		func(depMap map[string][]string) {
			visitor := &MockVisitor{}
			visitor.VisitCall.DepMap = depMap

			visitingWalker := NewVisitingDepWalker(visitor)

			_, err := visitingWalker.Walk([]string{"dopepkg"})
			Expect(err).To(HaveOccurred())
		},

		Entry("immediate failure",
			map[string][]string{},
		),

		Entry("failure visiting a dependency",
			map[string][]string{
				"dopepkg": {"notathing"},
			},
		),

		Entry("failure visiting a transitive dependency",
			map[string][]string{
				"dopepkg": {"libdope"},
				"libdope": {"leftpad"},
			},
		),
	)
})

var _ = Describe("FilteringVisitor", func() {
	DescribeTable("dependency filtering",
		func(deps []string, seen []string, expected []string) {
			visitor := &MockVisitor{}
			visitor.VisitCall.DepMap = map[string][]string{"dopepkg": deps}

			filteringVisitor := NewFilteringVisitor(visitor, seen)
			deps, err := filteringVisitor.Visit("dopepkg")

			Expect(err).ToNot(HaveOccurred())
			Expect(deps).To(Equal(expected))
		},

		Entry("no dependencies",
			[]string{},
			[]string{},
			[]string{},
		),

		Entry("single unseen dependency",
			[]string{"libdope"},
			[]string{},
			[]string{"libdope"},
		),

		Entry("multiple unseen dependencies",
			[]string{"libdope", "leftpad"},
			[]string{},
			[]string{"libdope", "leftpad"},
		),

		Entry("single seen dependency",
			[]string{"libdope"},
			[]string{"libdope"},
			[]string{},
		),

		Entry("multiple seen dependencies",
			[]string{"libdope", "leftpad"},
			[]string{"libdope", "leftpad"},
			[]string{},
		),

		Entry("some seen and unseen dependencies",
			[]string{"libdope", "leftpad"},
			[]string{"leftpad"},
			[]string{"libdope"},
		),

		Entry("no seen dependencies",
			[]string{"libdope", "leftpad"},
			[]string{"somethingelse"},
			[]string{"libdope", "leftpad"},
		),
	)

	Describe("failure", func() {
		It("returns errors from the wrapped visitor", func() {
			visitor := &MockVisitor{}
			visitor.VisitCall.DepMap = map[string][]string{}

			seen := []string{}

			filteringVisitor := NewFilteringVisitor(visitor, seen)
			_, err := filteringVisitor.Visit("dopepkg")

			Expect(err).To(HaveOccurred())
		})
	})
})

type MockVisitor struct {
	VisitCall struct {
		DepMap map[string][]string
	}
}

func (v *MockVisitor) Visit(name string) ([]string, error) {
	deps, ok := v.VisitCall.DepMap[name]
	if !ok {
		return nil, fmt.Errorf("not found: %s", name)
	}

	return deps, nil
}
