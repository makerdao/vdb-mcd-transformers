package main

import (
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDataGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data generator suite")
}

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})
