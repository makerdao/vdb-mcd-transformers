package log_take_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogTake(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogTake Event Transformer Suite")
}
