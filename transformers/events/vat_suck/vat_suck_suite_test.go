package vat_suck_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

func TestVatSuck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vat Suck Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})
