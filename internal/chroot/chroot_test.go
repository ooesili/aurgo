package chroot_test

import (
	"errors"
	"os"
	"path/filepath"

	. "github.com/ooesili/aurgo/internal/chroot"
	"github.com/ooesili/aurgo/test/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Chroot", func() {
	var (
		executor   *mocks.Executor
		filesystem *mocks.Filesystem
		chroot     Chroot
	)

	BeforeEach(func() {
		executor = &mocks.Executor{}
		filesystem = &mocks.Filesystem{}
		chroot = New(executor, filesystem)
	})

	Describe("Create", func() {
		var err error

		JustBeforeEach(func() {
			err = chroot.Create("/path/to/chroot")
		})

		Context("when the chroot can be created", func() {
			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("creates the leading directory", func() {
				path := filesystem.MkdirAllCall.Recieved.Path
				mode := filesystem.MkdirAllCall.Recieved.Mode

				Expect(path).To(Equal("/path/to/chroot"))
				Expect(mode).To(Equal(os.FileMode(0755)))
			})

			It("creates the chroot using mkarchroot", func() {
				expectedPath := filepath.Join("/path/to/chroot", "root")
				received := executor.ExecuteCall.Received

				Expect(received.Command).To(Equal("mkarchroot"))
				Expect(received.Args).To(Equal([]string{
					expectedPath, "base-devel",
				}))
			})
		})

		Context("when the leading directory cannot be created", func() {
			BeforeEach(func() {
				filesystem.MkdirAllCall.Returns.Err = errors.New("shoot")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})

			It("does not create the chroot with mkarchroot", func() {
				Expect(executor.ExecuteCall.Received.Command).To(BeZero())
				Expect(executor.ExecuteCall.Received.Args).To(BeZero())
			})
		})

		Context("when the chroot cannot be created", func() {
			BeforeEach(func() {
				executor.ExecuteCall.Returns.Err = errors.New("shoot")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Exists", func() {
		var (
			exists bool
			err    error
		)

		JustBeforeEach(func() {
			exists, err = chroot.Exists("/path/to/chroot")
		})

		Context("when the chroot exists", func() {
			BeforeEach(func() {
				filesystem.ExistsCall.Returns.Exists = true
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns true", func() {
				Expect(exists).To(BeTrue())
			})

			It("checks for the .arch-chroot file", func() {
				path := filesystem.ExistsCall.Received.Path
				Expect(path).To(Equal(filepath.Join("/path/to/chroot", "root", ".arch-chroot")))
			})
		})

		Context("when the chroot does not exist", func() {
			BeforeEach(func() {
				filesystem.ExistsCall.Returns.Exists = false
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns false", func() {
				Expect(exists).To(BeFalse())
			})
		})

		Context("when checking for the .arch-chroot file fails", func() {
			BeforeEach(func() {
				filesystem.ExistsCall.Returns.Err = errors.New("shoot")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
