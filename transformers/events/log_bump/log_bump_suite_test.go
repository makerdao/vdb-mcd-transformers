package log_bump_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogBump(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogBump Event Transformer Suite")
}
