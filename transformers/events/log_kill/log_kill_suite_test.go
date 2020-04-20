package log_kill_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogKill(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogKill Event Transformer Suite")
}
