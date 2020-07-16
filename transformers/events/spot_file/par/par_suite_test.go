package par_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPar(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SpotFilePar Event Transformer Suite")
}
