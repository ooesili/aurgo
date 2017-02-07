package config_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ooesili/aurgo/internal/config"
)

var _ = Describe("Config", func() {
	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = ioutil.TempDir("", "aurgo-config-test-")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Context("when the AURGOPATH does not exist", func() {
		It("returns an error", func() {
			aurgoPath := filepath.Join(tempDir, "not-a-dir")
			_, err := New(aurgoPath)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when the repo.yml is missing", func() {
		It("returns an error", func() {
			_, err := New(tempDir)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when the repo.yml is not valid yaml", func() {
		It("returns an error", func() {
			repoYamlPath := filepath.Join(tempDir, "repo.yml")
			err := ioutil.WriteFile(repoYamlPath, []byte("oops"), 0644)
			Expect(err).ToNot(HaveOccurred())

			_, err = New(tempDir)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when the repo.yml is valid", func() {
		var config Config

		BeforeEach(func() {
			repoYamlPath := filepath.Join(tempDir, "repo.yml")
			err := ioutil.WriteFile(repoYamlPath, []byte(fixtureRepoYaml), 0644)
			Expect(err).ToNot(HaveOccurred())

			config, err = New(tempDir)
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("Packages", func() {
			It("can list the packages from the yaml file", func() {
				packages, err := config.Packages()
				Expect(err).ToNot(HaveOccurred())
				Expect(packages).To(ConsistOf("package1", "package2"))
			})
		})

		Describe("SourcePath", func() {
			var (
				sourcePath   string
				err          error
				expectedPath string
			)

			JustBeforeEach(func() {
				expectedPath = filepath.Join(tempDir, "src", "package1")
				sourcePath, err = config.SourcePath("package1")
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("can list the source path for a package", func() {
				Expect(sourcePath).To(Equal(expectedPath))
			})

			It("creates the path to the package on the filesystem", func() {
				Expect(sourcePath).To(BeADirectory())
			})

			It("can be called a second time", func() {
				sourcePath, err := config.SourcePath("package1")
				Expect(err).ToNot(HaveOccurred())
				Expect(sourcePath).To(Equal(expectedPath))
			})

			Context("when the path can not be created", func() {
				BeforeEach(func() {
					srcPath := filepath.Join(tempDir, "src")
					Expect(ioutil.WriteFile(srcPath, nil, 0644)).To(Succeed())
				})

				It("returns an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("AurRepoURL", func() {
			It("can generate the git url for a package", func() {
				url := config.AurRepoURL("package1")
				Expect(url).To(Equal("https://aur.archlinux.org/package1.git"))
			})
		})
	})
})
