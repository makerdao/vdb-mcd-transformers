package log_value_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogValue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogValue Event Transformer Suite")
}
