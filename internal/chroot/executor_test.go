package chroot_test

import (
	"bytes"

	. "github.com/ooesili/aurgo/internal/chroot"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OSExecutor", func() {
	var (
		stdout   *bytes.Buffer
		stderr   *bytes.Buffer
		executor OSExecutor
	)

	BeforeEach(func() {
		stdout = &bytes.Buffer{}
		stderr = &bytes.Buffer{}
		executor = NewOSExecutor(stdout, stderr)
	})

	Describe("Execute", func() {
		var err error

		Context("when the command succeds", func() {
			BeforeEach(func() {
				err = executor.Execute("sh", "-c", "echo out; echo err >&2")
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("writes to the given stdout", func() {
				Expect(stdout.String()).To(Equal("out\n"))
			})

			It("writes to the given stderr", func() {
				Expect(stderr.String()).To(Equal("err\n"))
			})
		})

		Context("when the command fails", func() {
			BeforeEach(func() {
				err = executor.Execute("exit 1")
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
