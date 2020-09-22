package v1_0_0_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestV100(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cat v1_0_0 Suite")
}
