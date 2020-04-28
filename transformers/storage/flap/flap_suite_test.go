package flap_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlap(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flap Storage Transformer Suite")
}
