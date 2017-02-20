package pacman_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/ooesili/aurgo/internal/pacman"
)

var _ = Describe("Executor", func() {
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
		executor := NewOsExecutor()
		stdout, err = executor.Execute(command, args...)
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
