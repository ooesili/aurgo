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
	var (
		aurgoPath string
		repoYml   string
		fixture   string
	)

	BeforeEach(func() {
		var err error
		aurgoPath, err = ioutil.TempDir("", "aurgo-")
		Expect(err).ToNot(HaveOccurred())

		repoYml = filepath.Join(aurgoPath, "repo.yml")

		fixture = ""
	})

	JustBeforeEach(func() {
		err := ioutil.WriteFile(repoYml, []byte(fixture), 0644)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(aurgoPath)
	})

	Describe("aurgo sync", func() {
		Context("with a package with no dependencies", func() {
			BeforeEach(func() {
				fixture = fixtureRepoYamlXcape
			})

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

			It("can remove a package from the cache", func() {
				cmd := exec.Command(aurgoBinary, "sync")
				err := cmd.Run(
					exec.Stdout(GinkgoWriter),
					exec.Stderr(GinkgoWriter),
					exec.Setenv("AURGOPATH", aurgoPath),
				)
				Expect(err).ToNot(HaveOccurred())

				xcapePath := filepath.Join(aurgoPath, "src", "xcape")
				Expect(xcapePath).To(BeADirectory())

				err = ioutil.WriteFile(repoYml, []byte(fixtureRepoYamlNoPackages), 0644)
				Expect(err).ToNot(HaveOccurred())

				cmd = exec.Command(aurgoBinary, "sync")
				err = cmd.Run(
					exec.Stdout(GinkgoWriter),
					exec.Stderr(GinkgoWriter),
					exec.Setenv("AURGOPATH", aurgoPath),
				)
				Expect(err).ToNot(HaveOccurred())

				Expect(xcapePath).ToNot(BeADirectory())
			})
		})

		Context("with a package with dependencies", func() {
			BeforeEach(func() {
				fixture = fixtureRepoYamlYaourt
			})

			It("downloads dependencies", func() {
				cmd := exec.Command(aurgoBinary, "sync")
				err := cmd.Run(
					exec.Stdout(GinkgoWriter),
					exec.Stderr(GinkgoWriter),
					exec.Setenv("AURGOPATH", aurgoPath),
				)
				Expect(err).ToNot(HaveOccurred())

				yaourtPkgbuildPath := filepath.Join(aurgoPath, "src", "yaourt", "PKGBUILD")
				Expect(yaourtPkgbuildPath).To(BeARegularFile())
				packageQueryPkgbuildPath := filepath.Join(aurgoPath, "src", "package-query", "PKGBUILD")
				Expect(packageQueryPkgbuildPath).To(BeARegularFile())
			})
		})

		Context("when dependencies are only met by the Provides field from pacman", func() {
			BeforeEach(func() {
				fixture = fixtureRepoYamlNtkGit
			})

			It("downloads dependencies", func() {
				cmd := exec.Command(aurgoBinary, "sync")
				err := cmd.Run(
					exec.Stdout(GinkgoWriter),
					exec.Stderr(GinkgoWriter),
					exec.Setenv("AURGOPATH", aurgoPath),
				)
				Expect(err).ToNot(HaveOccurred())

				ntkGitPkgbuildPath := filepath.Join(aurgoPath, "src", "ntk-git", "PKGBUILD")
				Expect(ntkGitPkgbuildPath).To(BeARegularFile())
			})
		})

		Context("when a package does not exist", func() {
			BeforeEach(func() {
				fixture = fixtureRepoYamlNotFound
			})

			It("returns an error", func() {
				cmd := exec.Command(aurgoBinary, "sync")
				err := cmd.Run(
					exec.Stdout(GinkgoWriter),
					exec.Stderr(GinkgoWriter),
					exec.Setenv("AURGOPATH", aurgoPath),
				)
				Expect(err).To(HaveOccurred())

				notFoundPackagePath := filepath.Join(
					aurgoPath, "src", "totally-not-a-package-i-hope",
				)
				Expect(notFoundPackagePath).ToNot(BeADirectory())
			})
		})
	})

	It("can view the version of a package in the AUR", func() {
		cmd := exec.Command(aurgoBinary, "info", "xcape")
		output, err := cmd.CombinedOutput()
		Expect(err).ToNot(HaveOccurred())

		Expect(string(output)).To(Equal("aur/xcape 1.2-1\n"))
	})
})

var fixtureRepoYamlXcape = `---
packages:
- xcape
`

var fixtureRepoYamlYaourt = `---
packages:
- yaourt
`

var fixtureRepoYamlNtkGit = `---
packages:
- ntk-git
`

var fixtureRepoYamlNoPackages = `---
packages: []
`

var fixtureRepoYamlNotFound = `---
packages:
- totally-not-a-package-i-hope
`
