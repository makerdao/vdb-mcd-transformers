package log_insert_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLogInsert(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LogInsert Event Transformer Suite")
}
