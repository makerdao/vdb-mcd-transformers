package log_buy_enabled

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogBuyEnabled(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogBuyEnabled Event Transformer Suite")
}
