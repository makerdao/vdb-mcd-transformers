package flip_tick_test

import (
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

func TestFlipTick(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flip tick Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})
