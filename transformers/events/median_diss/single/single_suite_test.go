package single_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMedianDissSingle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MedianDissSingle Event Transformer Suite")
}
