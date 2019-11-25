package dsr_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDsr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PotFileDsr Suite")
}
