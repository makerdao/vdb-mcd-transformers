package spot_poke_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSpotPoke(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SpotPoke Event Transformer Suite")
}
