package pot_drip_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPotDrip(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PotDrip Event Transformer Suite")
}
