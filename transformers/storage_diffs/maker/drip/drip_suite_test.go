package drip_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDrip(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Drip Suite")
}
