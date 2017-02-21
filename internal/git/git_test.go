package git_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ooesili/aurgo/internal/git"
	"github.com/pkg/exec"
)

var _ = Describe("Git", func() {
	Describe("Clone", func() {
		var (
			tempDir    string
			sourceRepo string
			destRepo   string
			git        Git
		)

		BeforeEach(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "aurgo-git-test-")
			Expect(err).ToNot(HaveOccurred())

			Expect(untarFixture("repo.tar.gz", tempDir)).To(Succeed())

			sourceRepo = filepath.Join(tempDir, "repo")
			destRepo = filepath.Join(tempDir, "dest")
			Expect(os.Mkdir(destRepo, 0755)).To(Succeed())

			git = New(GinkgoWriter, GinkgoWriter)
		})

		AfterEach(func() {
			os.RemoveAll(tempDir)
		})

		Describe("cloning a new repo", func() {
			It("can clone a repo", func() {
				err := git.Clone(sourceRepo, destRepo)
				Expect(err).ToNot(HaveOccurred())

				helloPath := filepath.Join(destRepo, "hello")
				Expect(helloPath).To(BeARegularFile())
			})

			Context("when the source repo can not be found", func() {
				It("returns an error", func() {
					Expect(os.RemoveAll(sourceRepo)).To(Succeed())
					err := git.Clone(sourceRepo, destRepo)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(HavePrefix("git clone failed: "))
				})
			})
		})

		Describe("updating an existing repo", func() {
			It("pulls changes when the repo already exists", func() {
				err := git.Clone(sourceRepo, destRepo)
				Expect(err).ToNot(HaveOccurred())

				Expect(os.RemoveAll(sourceRepo)).To(Succeed())
				Expect(untarFixture("repo-updated.tar.gz", tempDir)).To(Succeed())

				err = git.Clone(sourceRepo, destRepo)
				Expect(err).ToNot(HaveOccurred())

				hiPath := filepath.Join(destRepo, "hi")
				Expect(hiPath).To(BeARegularFile())
			})

			Context("when the source repo cannot be found", func() {
				It("pulls changes when the repo already exists", func() {
					err := git.Clone(sourceRepo, destRepo)
					Expect(err).ToNot(HaveOccurred())

					Expect(os.RemoveAll(sourceRepo)).To(Succeed())

					err = git.Clone(sourceRepo, destRepo)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(HavePrefix("git pull failed: "))
				})
			})
		})
	})
})

func untarFixture(tarName, destDir string) error {
	repoTar, err := filepath.Abs(filepath.Join("fixtures", tarName))
	if err != nil {
		return err
	}

	untarCmd := exec.Command("tar", "xzf", repoTar)
	return untarCmd.Run(
		exec.Dir(destDir),
		exec.Stdout(GinkgoWriter),
		exec.Stderr(GinkgoWriter),
	)
}
