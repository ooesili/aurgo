package aurgo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAurgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aurgo Suite")
}
