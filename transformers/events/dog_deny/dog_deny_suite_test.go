package dog_deny_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDogDeny(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DogDeny Suite")
}
