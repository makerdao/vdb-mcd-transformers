package bid_creation_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTriggers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bid Creation Trigger Component Test Suite")
}
