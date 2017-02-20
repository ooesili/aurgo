package srcinfo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSrcinfo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Srcinfo Suite")
}
