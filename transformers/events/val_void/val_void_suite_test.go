package val_void_test

import (
	"testing"

	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

func TestValVoid(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ValVoid Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})
