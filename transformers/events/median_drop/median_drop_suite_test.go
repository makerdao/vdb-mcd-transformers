package median_drop_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMedianDrop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MedianDrop Event Transformer Suite")
}
