package aurgo_test

import (
	"errors"

	. "github.com/ooesili/aurgo/internal/aurgo"
	"github.com/ooesili/aurgo/test/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BuildManager", func() {
	var (
		buildEnv     *mocks.BuildEnv
		buildManager BuildManager
		config       *mocks.Config
	)

	BeforeEach(func() {
		buildEnv = &mocks.BuildEnv{}
		config = &mocks.Config{}
		config.ChrootPathCall.ChrootPath = "/path/to/chroot"

		buildManager = NewBuildManager(buildEnv, config)
	})

	Describe("Provision", func() {
		var err error

		JustBeforeEach(func() {
			err = buildManager.Provision()
		})

		Context("when the build environment does not exist", func() {
			BeforeEach(func() {
				buildEnv.ExistsCall.Returns.Exists = false
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("checks to see if the environment exists", func() {
				location := buildEnv.ExistsCall.Received.Location
				Expect(location).To(Equal("/path/to/chroot"))
			})

			It("creates the build environment", func() {
				location := buildEnv.CreateCall.Received.Location
				Expect(location).To(Equal("/path/to/chroot"))
			})
		})

		Context("when the build environment does exist", func() {
			BeforeEach(func() {
				buildEnv.ExistsCall.Returns.Exists = true
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("checks to see if the environment exists", func() {
				location := buildEnv.ExistsCall.Received.Location
				Expect(location).To(Equal("/path/to/chroot"))
			})

			It("does not create the build environment", func() {
				location := buildEnv.CreateCall.Received.Location
				Expect(location).To(BeEmpty())
			})
		})

		Context("when checking if the build environment exists fails", func() {
			BeforeEach(func() {
				buildEnv.ExistsCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when creating the build environment fails", func() {
			BeforeEach(func() {
				buildEnv.ExistsCall.Returns.Exists = false
				buildEnv.CreateCall.Returns.Err = errors.New("dang")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
