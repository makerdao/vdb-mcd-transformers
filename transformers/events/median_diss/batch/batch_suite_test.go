package batch_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMedianDissBatch(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MedianDissBatch Event Transformer Suite")
}
