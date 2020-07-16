package pot_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pot Storage Transformer Suite")
}
