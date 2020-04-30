package log_median_price_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogMedianPrice(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogMedianPrice Suite")
}
