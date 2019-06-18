package spot_poke_test

import (
	"io/ioutil"
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSpotPoke(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SpotPoke Suite")
}

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
})
