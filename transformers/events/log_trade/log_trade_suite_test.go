package log_trade_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogTrade(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogTrade Event Transformer Suite")
}
