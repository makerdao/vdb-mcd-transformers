package log_item_update_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogItemUpdate(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogItemUpdate Event Transformer Suite")
}
