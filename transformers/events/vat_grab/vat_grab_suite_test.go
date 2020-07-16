package vat_grab_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatGrab(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatGrab Event Transformer Suite")
}
