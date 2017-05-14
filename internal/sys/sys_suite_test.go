package sys_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSys(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sys Suite")
}
