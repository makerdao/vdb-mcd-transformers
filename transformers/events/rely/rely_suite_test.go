package rely_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRely(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rely Suite")
}
