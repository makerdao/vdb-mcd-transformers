package log_matching_enabled_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogMatchingEnabled(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogMatchingEnabled Event Transformer Suite")
}
