package medianizer

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
)

func TestMedianizer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Medianzier Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
