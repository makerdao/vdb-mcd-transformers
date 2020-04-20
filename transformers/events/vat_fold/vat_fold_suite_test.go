package vat_fold_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatFold(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatFold Event Transformer Suite")
}
