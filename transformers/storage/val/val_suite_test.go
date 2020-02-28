package val_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
)

func TestVal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Val Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
