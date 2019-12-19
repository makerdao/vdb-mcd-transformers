package cat_rely_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCatRely(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cat Rely Suite")
}
