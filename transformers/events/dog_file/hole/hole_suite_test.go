package hole_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDogFileHole(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DogFileHole Suite")
}
