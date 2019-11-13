package dent_test

import (
	"testing"

	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

func TestDent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dent Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})
