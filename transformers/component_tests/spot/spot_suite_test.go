package spot_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSpot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Spot Component Test Suite")
}
