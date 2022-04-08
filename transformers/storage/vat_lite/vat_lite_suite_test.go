package vat_lite

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatLite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vat Lite Storage Transformer Suite")
}
