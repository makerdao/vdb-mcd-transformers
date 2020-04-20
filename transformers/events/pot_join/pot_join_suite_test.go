package pot_join_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPotJoin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PotJoin Event Transformer Suite")
}
