package chroot_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestChroot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chroot Suite")
}
