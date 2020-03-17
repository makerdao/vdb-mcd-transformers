package median_drop_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestMedianDrop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MedianDrop Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
