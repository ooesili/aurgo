package pacman_test

import (
	"bytes"
	"errors"

	. "github.com/ooesili/aurgo/internal/pacman"
	"github.com/ooesili/aurgo/test/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pacman", func() {
	var (
		executor *mocks.Executor
		pacman   Pacman
		err      error
	)

	BeforeEach(func() {
		executor = &mocks.Executor{}
	})

	JustBeforeEach(func() {
		pacman, err = New(executor)
	})

	Context("when the pacman command executes successfully", func() {
		BeforeEach(func() {
			stdout := bytes.NewBufferString(fixtureRealOutput)
			executor.ExecuteCall.Returns.Stdout = stdout
		})

		It("succeeds", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("executes the correct command", func() {
			recieved := executor.ExecuteCall.Recieved
			Expect(recieved.Command).To(Equal("pacman"))
			Expect(recieved.Args).To(Equal([]string{"-Si"}))
		})

		It("can return the correct list provided packages", func() {
			pkgs := pacman.ListAvailable()
			Expect(pkgs).To(ConsistOf(
				"cronie",
				"cron",
				"curl",
				"libcurl.so",
				"grub",
				"grub-common",
				"grub-bios",
				"grub-emu",
				"grub-efi-x86_64",
				"rust",
			))
		})
	})

	Context("when pacman fails to execute", func() {
		BeforeEach(func() {
			executor.ExecuteCall.Returns.Err = errors.New("dang")
		})

		It("returns an error", func() {
			Expect(err).To(HaveOccurred())
		})
	})
})
