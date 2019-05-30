package vat_fork_test

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVatFork(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VatFork Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})
