package log_sorted_offer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogSortedOffer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogSortedOffer Event Transformer Suite")
}
