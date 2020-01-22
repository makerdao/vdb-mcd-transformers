package bid_creation_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestTriggers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bid Creation Trigger Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
