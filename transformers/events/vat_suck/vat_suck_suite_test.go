package vat_suck_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatSuck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatSuck Event Transformer Suite")
}
