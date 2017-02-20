package pacman_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPacman(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pacman Suite")
}
