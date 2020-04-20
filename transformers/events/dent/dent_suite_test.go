package dent_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dent Event Transformer Suite")
}
