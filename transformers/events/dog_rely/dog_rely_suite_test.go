package dog_rely_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDogRely(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dog Rely Suite")
}
