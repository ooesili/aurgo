package cache_test

import (
	"errors"
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
		fs      *mocks.Filesystem
		err     error
	)

	BeforeEach(func() {
		config = &mocks.Config{}
		git = &mocks.Git{}
		srcinfo = &mocks.SrcInfo{}
		fs = &mocks.Filesystem{}
	})

	JustBeforeEach(func() {
		cache = New(config, git, srcinfo, fs)
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
			srcinfoPath  string
			pkgs         []string
			srcinfoBytes []byte
			err          error
		)

		BeforeEach(func() {
			srcinfoBytes = []byte("cool bytes")
			fs.ReadFileCall.Returns.Bytes = srcinfoBytes

			pkgPath := "/path/to/coolpkg"
			config.SourcePathCall.Returns.Path = pkgPath
			srcinfoPath = filepath.Join(pkgPath, ".SRCINFO")
		})

		JustBeforeEach(func() {
			pkgs, err = cache.GetDeps("coolpkg")
		})

		Context("when all dependencies succeed", func() {
			BeforeEach(func() {
				srcinfo.ParseCall.Returns.Package = Package{
					Depends:      []string{"leftpad", "libdope"},
					Makedepends:  []string{"cmake", "maven"},
					Checkdepends: []string{"checktool", "testlib"},
				}
			})

			It("suceeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("looks up the source path for the package", func() {
				Expect(config.SourcePathCall.Received.Package).To(Equal("coolpkg"))
			})

			It("reads the correct SRCINFO file", func() {
				Expect(fs.ReadFileCall.Received.Path).To(Equal(srcinfoPath))
			})

			It("calls SrcInfo.Parse with the SRCINFO file contents", func() {
				Expect(srcinfo.ParseCall.Received.Input).To(Equal(srcinfoBytes))
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

		Context("when reading the SRCINFO file fails", func() {
			BeforeEach(func() {
				fs.ReadFileCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when parsing the sourcinfo file fails", func() {
			BeforeEach(func() {
				fs.ReadFileCall.Returns.Bytes = srcinfoBytes
				srcinfo.ParseCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("List", func() {
		var (
			pkgs []string
			err  error
		)

		BeforeEach(func() {
			config.SourceBaseCall.SourceBase = "/path/to/sourcebase"
		})

		JustBeforeEach(func() {
			pkgs, err = cache.List()
		})

		Context("when listing the files succeeds", func() {
			BeforeEach(func() {
				fs.ListFilesCall.Returns.Filenames = []string{"dopepkg", "dopelib"}
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("lists files in the source base directory", func() {
				dirname := fs.ListFilesCall.Received.Dirname
				Expect(dirname).To(Equal("/path/to/sourcebase"))
			})

			It("returns all packages in the cache", func() {
				Expect(pkgs).To(Equal([]string{"dopelib", "dopepkg"}))
			})
		})

		Context("when there is an error listing the files", func() {
			BeforeEach(func() {
				fs.ListFilesCall.Returns.Err = errors.New("darn")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Remove", func() {
		var err error

		BeforeEach(func() {
			config.SourcePathCall.Returns.Path = "/path/to/pkg"
		})

		JustBeforeEach(func() {
			err = cache.Remove("dopepkg")
		})

		Context("when removing the package's source path succeeds", func() {
			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("removes the correct directory", func() {
				Expect(fs.RemoveAllCall.Received.Path).To(Equal("/path/to/pkg"))
			})
		})

		Context("when removign the package's source path fails", func() {
			BeforeEach(func() {
				fs.RemoveAllCall.Returns.Err = errors.New("shoot")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
