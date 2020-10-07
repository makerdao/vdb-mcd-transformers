package cat_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlipFileCat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flip File Cat Suite")
}
