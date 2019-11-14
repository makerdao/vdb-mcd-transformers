package new_cdp_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

func TestNewCdp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NewCdp Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
