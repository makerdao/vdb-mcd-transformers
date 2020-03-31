package median_lift

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestMedianLift(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MedianLift Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
