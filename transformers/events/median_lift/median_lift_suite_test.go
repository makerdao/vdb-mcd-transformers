package median_lift

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMedianLift(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MedianLift Event Transformer Suite")
}
