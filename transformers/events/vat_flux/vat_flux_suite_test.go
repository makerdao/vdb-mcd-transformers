package vat_flux_test

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatFlux(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatFlux Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})
