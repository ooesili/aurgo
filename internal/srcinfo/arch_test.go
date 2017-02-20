package srcinfo_test

import (
	. "github.com/ooesili/aurgo/internal/srcinfo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("ArchString", func() {
	DescribeTable("valid architectures",
		func(goarch, expectedArch string) {
			arch, err := ArchString(goarch)
			Expect(err).ToNot(HaveOccurred())
			Expect(arch).To(Equal(expectedArch))
		},
		Entry("amd64", "amd64", "x86_64"),
	)

	Context("when given an unsupported architecture", func() {
		It("returns an error", func() {
			_, err := ArchString("foobar")
			Expect(err).To(MatchError("unsupported architecture: foobar"))
		})
	})
})
