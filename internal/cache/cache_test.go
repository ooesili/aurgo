package cache_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/ooesili/aurgo/internal/cache"
	"github.com/ooesili/aurgo/test/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	var (
		cache   Cache
		config  *mocks.Config
		git     *mocks.Git
		srcinfo *mocks.SrcInfo
		err     error
	)

	BeforeEach(func() {
		config = &mocks.Config{}
		git = &mocks.Git{}
		srcinfo = &mocks.SrcInfo{}
	})

	JustBeforeEach(func() {
		cache = New(config, git, srcinfo)
	})

	Describe("Sync", func() {
		JustBeforeEach(func() {
			err = cache.Sync("coolpkg")
		})

		Context("when all dependencies succeed", func() {
			BeforeEach(func() {
				config.SourcePathCall.Returns.Path = "/path/to/coolpkg"
				config.AurRepoURLCall.Returns.URL = "https://example.com/coolpkg.git"
			})

			It("looks up the source path for the package", func() {
				Expect(config.SourcePathCall.Received.Package).To(Equal("coolpkg"))
			})

			It("looks up the aur repo URL for the package", func() {
				Expect(config.AurRepoURLCall.Received.Package).To(Equal("coolpkg"))
			})

			It("cones the git repo from the URL", func() {
				Expect(git.CloneCall.Received.Path).To(Equal("/path/to/coolpkg"))
				Expect(git.CloneCall.Received.URL).To(Equal("https://example.com/coolpkg.git"))
			})
		})

		Context("when looking up the source path for the package fails", func() {
			BeforeEach(func() {
				config.SourcePathCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when cloning the package fails", func() {
			BeforeEach(func() {
				config.SourcePathCall.Returns.Path = "/path/to/coolpkg"
				config.AurRepoURLCall.Returns.URL = "https://example.com/coolpkg.git"
				git.CloneCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetDeps", func() {
		var (
			tempDir      string
			pkgPath      string
			srcinfoPath  string
			expectedPkg  Package
			pkgs         []string
			srcinfoBytes []byte
			err          error
		)

		BeforeEach(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "aurgo-srcinfo-test-")
			Expect(err).ToNot(HaveOccurred())

			srcinfoBytes = []byte("cool bytes")
			expectedPkg = Package{
				Depends:      []string{"leftpad", "libdope"},
				Makedepends:  []string{"cmake", "maven"},
				Checkdepends: []string{"checktool", "testlib"},
			}

			pkgPath = filepath.Join(tempDir, "coolpkg")
			Expect(os.Mkdir(pkgPath, 0755)).To(Succeed())

			srcinfoPath = filepath.Join(pkgPath, ".SRCINFO")
			err = ioutil.WriteFile(srcinfoPath, srcinfoBytes, 0644)
			Expect(err).ToNot(HaveOccurred())
		})

		JustBeforeEach(func() {
			pkgs, err = cache.GetDeps("coolpkg")
		})

		AfterEach(func() {
			os.RemoveAll(tempDir)
		})

		Context("when all dependencies succeed", func() {
			BeforeEach(func() {
				config.SourcePathCall.Returns.Path = pkgPath
				srcinfo.ParseCall.Returns.Package = expectedPkg
			})

			It("suceeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("looks up the source path for the package", func() {
				Expect(config.SourcePathCall.Received.Package).To(Equal("coolpkg"))
			})

			It("calls SrcInfo.Parse with the SRCINFO file contents", func() {
				Expect(srcinfo.ParseCall.Recieved.Input).To(Equal(srcinfoBytes))
			})

			It("returns all dependencies from the package", func() {
				Expect(pkgs).To(ConsistOf(
					"libdope",
					"leftpad",
					"cmake",
					"maven",
					"checktool",
					"testlib",
				))
			})
		})

		Context("when looking up the source path for the package fails", func() {
			BeforeEach(func() {
				config.SourcePathCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when the SRCINFO file does not exist", func() {
			BeforeEach(func() {
				config.SourcePathCall.Returns.Path = pkgPath
				Expect(os.Remove(srcinfoPath)).To(Succeed())
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when parsing the sourcinfo file fails", func() {
			BeforeEach(func() {
				config.SourcePathCall.Returns.Path = pkgPath
				srcinfo.ParseCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
