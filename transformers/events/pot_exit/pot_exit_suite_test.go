package pot_exit_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPotExit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PotExit Suite")
}
