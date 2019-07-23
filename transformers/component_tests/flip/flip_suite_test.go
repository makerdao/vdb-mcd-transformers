package flip_test

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlip(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flip Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
