package cat_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cat Suite")
}
