package frob_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFrob(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BackFill Frob Suite")
}
