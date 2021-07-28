package sale_creation_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTriggers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sale Creation Trigger Component Test Suite")
}
