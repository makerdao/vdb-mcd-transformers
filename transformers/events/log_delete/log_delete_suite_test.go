package log_delete_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogDelete(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogDelete Event Transformer Suite")
}
