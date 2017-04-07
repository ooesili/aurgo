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
				packages := config.Packages()
				Expect(packages).To(ConsistOf("package1", "package2"))
			})
		})

		Describe("SourcePath", func() {
			It("can list the source path for a package", func() {
				sourcePath := config.SourcePath("package1")

				expectedPath := filepath.Join(tempDir, "src", "package1")
				Expect(sourcePath).To(Equal(expectedPath))
			})
		})

		Describe("AurRepoURL", func() {
			It("can generate the git url for a package", func() {
				url := config.AurRepoURL("package1")
				Expect(url).To(Equal("https://aur.archlinux.org/package1.git"))
			})
		})

		Describe("SourceBase", func() {
			It("can list the source directory", func() {
				sourceBase := config.SourceBase()
				Expect(sourceBase).To(Equal(filepath.Join(tempDir, "src")))
			})
		})

		Describe("ChrootPath", func() {
			It("returns the path to the chroot directory", func() {
				chrootPath := config.ChrootPath()
				Expect(chrootPath).To(Equal(filepath.Join(tempDir, "chroot")))
			})
		})
	})
})
