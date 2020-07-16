package tend_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTend(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tend Event Transformer Suite")
}
