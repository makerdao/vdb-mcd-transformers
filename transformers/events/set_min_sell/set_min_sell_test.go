package set_min_sell_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogMinSell(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SetMinSell Event Transformer Suite")
}
