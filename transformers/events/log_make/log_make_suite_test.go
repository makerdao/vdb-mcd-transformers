package log_make_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogMake(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogMake Event Transformer Suite")
}
