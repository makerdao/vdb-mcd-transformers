package batch_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMedianKissBatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MedianKissBatch Event Transformer Suite")
}
