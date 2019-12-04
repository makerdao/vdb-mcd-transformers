package pot_cage_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPotCage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PotCage Suite")
}
