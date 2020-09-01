package v1_1_0_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestV110(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cat v1_1_0 Suite")
}
