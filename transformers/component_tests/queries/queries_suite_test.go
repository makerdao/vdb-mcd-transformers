package queries

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQueries(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Queries Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
