package aur_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAur(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aur Suite")
}
