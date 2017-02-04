package acceptance_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Acceptance", func() {
	It("do", func() {
		cmd := exec.Command(aurgoBinary, "info", "xcape")
		output, err := cmd.CombinedOutput()
		Expect(err).ToNot(HaveOccurred())

		Expect(string(output)).To(Equal("aur/xcape 1.2-1\n"))
	})
})
