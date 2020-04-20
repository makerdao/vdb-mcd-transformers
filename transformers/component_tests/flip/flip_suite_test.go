package flip_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlip(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flip Component Test Suite")
}
