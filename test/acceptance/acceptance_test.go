package acceptance_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/exec"
)

var _ = Describe("Acceptance", func() {
	var aurgoPath string

	BeforeEach(func() {
		var err error
		aurgoPath, err = ioutil.TempDir("", "aurgo-")
		Expect(err).ToNot(HaveOccurred())

		repoYml := filepath.Join(aurgoPath, "repo.yml")
		err = ioutil.WriteFile(repoYml, []byte(fixtureRepoYaml), 0644)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(aurgoPath)
	})

	Describe("aurgo sync", func() {
		It("can download a PKGBUILD from the AUR", func() {
			cmd := exec.Command(aurgoBinary, "sync")
			err := cmd.Run(
				exec.Stdout(GinkgoWriter),
				exec.Stderr(GinkgoWriter),
				exec.Setenv("AURGOPATH", aurgoPath),
			)
			Expect(err).ToNot(HaveOccurred())

			pkgbuildPath := filepath.Join(aurgoPath, "src", "xcape", "PKGBUILD")
			Expect(pkgbuildPath).To(BeARegularFile())
		})

		It("can be run twice", func() {
			cmd := exec.Command(aurgoBinary, "sync")
			err := cmd.Run(
				exec.Stdout(GinkgoWriter),
				exec.Stderr(GinkgoWriter),
				exec.Setenv("AURGOPATH", aurgoPath),
			)
			Expect(err).ToNot(HaveOccurred())

			cmd = exec.Command(aurgoBinary, "sync")
			err = cmd.Run(
				exec.Stdout(GinkgoWriter),
				exec.Stderr(GinkgoWriter),
				exec.Setenv("AURGOPATH", aurgoPath),
			)
			Expect(err).ToNot(HaveOccurred())

			pkgbuildPath := filepath.Join(aurgoPath, "src", "xcape", "PKGBUILD")
			Expect(pkgbuildPath).To(BeARegularFile())
		})
	})

	It("can view the version of a package in the AUR", func() {
		cmd := exec.Command(aurgoBinary, "info", "xcape")
		output, err := cmd.CombinedOutput()
		Expect(err).ToNot(HaveOccurred())

		Expect(string(output)).To(Equal("aur/xcape 1.2-1\n"))
	})
})

var fixtureRepoYaml = `---
packages:
- xcape
`
