package vat_tune_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

func TestVatTune(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatTune Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})
