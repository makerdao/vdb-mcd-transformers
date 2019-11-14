package flap_test

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flap Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
