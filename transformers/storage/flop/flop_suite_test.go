package flop_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flop Storage Transformer Suite")
}
