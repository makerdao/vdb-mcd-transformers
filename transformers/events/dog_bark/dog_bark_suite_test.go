package dog_bark_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDogBark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DogBark Suite")
}
