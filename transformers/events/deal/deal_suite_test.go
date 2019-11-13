package deal_test

import (
	"testing"

	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

func TestDeal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deal Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})
