package dog

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dog Component Test Suite")
}
