package v1_0_10_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cat v1.0.10 Storage Transformer Suite")
}
