package dog_digs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDogDigs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DogDigs Suite")
}
