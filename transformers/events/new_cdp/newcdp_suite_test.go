package new_cdp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"testing"
)

func TestNewCdp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NewCdp Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
