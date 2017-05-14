package sys_test

import (
	"bytes"

	. "github.com/ooesili/aurgo/internal/sys"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Executor", func() {
	var (
		defaultOut *bytes.Buffer
		defaultErr *bytes.Buffer
		executor   Executor
	)

	BeforeEach(func() {
		defaultOut = &bytes.Buffer{}
		defaultErr = &bytes.Buffer{}
		executor = NewExecutor(defaultOut, defaultErr)
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
				Expect(defaultOut.String()).To(Equal("out\n"))
			})

			It("writes to the given stderr", func() {
				Expect(defaultErr.String()).To(Equal("err\n"))
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

	Describe("ExecuteCapture", func() {
		var (
			command string
			args    []string
			stdout  *bytes.Buffer
			err     error
		)

		BeforeEach(func() {
			command = ""
			args = nil
		})

		JustBeforeEach(func() {
			stdout, err = executor.ExecuteCapture(command, args...)
		})

		Context("when the command executes correctly", func() {
			BeforeEach(func() {
				command = "echo"
				args = []string{"hello world"}
			})

			It("succeeds", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns the captured stdout in a buffer", func() {
				Expect(stdout.String()).To(Equal("hello world\n"))
			})
		})

		Context("when the command fails to execute", func() {
			BeforeEach(func() {
				command = "false"
			})

			It("returns an error", func() {
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
