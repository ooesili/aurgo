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

			repoTar, err := filepath.Abs("fixtures/repo.tar.gz")
			Expect(err).ToNot(HaveOccurred())

			untarCmd := exec.Command("tar", "xzf", repoTar)
			err = untarCmd.Run(
				exec.Dir(tempDir),
				exec.Stdout(GinkgoWriter),
				exec.Stderr(GinkgoWriter),
			)
			Expect(err).ToNot(HaveOccurred())

			sourceRepo = filepath.Join(tempDir, "repo")
			destRepo = filepath.Join(tempDir, "dest")
			Expect(os.Mkdir(destRepo, 0755)).To(Succeed())

			git = New()
		})

		AfterEach(func() {
			os.RemoveAll(tempDir)
		})

		It("can clone a repo", func() {
			err := git.Clone(sourceRepo, destRepo)
			Expect(err).ToNot(HaveOccurred())

			helloPath := filepath.Join(destRepo, "hello")
			Expect(helloPath).To(BeARegularFile())
		})

		It("can be run against an existing repo", func() {
			err := git.Clone(sourceRepo, destRepo)
			Expect(err).ToNot(HaveOccurred())

			err = git.Clone(sourceRepo, destRepo)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when the repo cannot be cloned", func() {
			It("returns an error", func() {
				Expect(os.RemoveAll(sourceRepo)).To(Succeed())
				err := git.Clone(sourceRepo, destRepo)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(HavePrefix("git clone failed: "))
			})
		})
	})
})
