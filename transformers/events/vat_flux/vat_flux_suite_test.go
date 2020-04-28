package vat_flux_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatFlux(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatFlux Event Transformer Suite")
}
