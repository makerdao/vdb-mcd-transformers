package val_poke_test

import (
	"testing"

	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

func TestValPoke(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ValPoke Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})
