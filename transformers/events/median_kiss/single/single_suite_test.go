package single_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMedianKissSingle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MedianKissSingle Event Transformer Suite")
}
