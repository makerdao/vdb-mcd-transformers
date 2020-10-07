package box_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBox(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Box Suite")
}
