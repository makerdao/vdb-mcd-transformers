package single_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestMedianKissSingle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MedianKissSingle Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
