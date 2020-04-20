package deal_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDeal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deal Event Transformer Suite")
}
