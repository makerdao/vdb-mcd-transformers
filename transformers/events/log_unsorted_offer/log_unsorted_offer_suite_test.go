package log_unsorted_offer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogUnsortedOffer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogUnsortedOffer Event Transformer Suite")
}
