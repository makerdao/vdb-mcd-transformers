package cat_test

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cat Suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
