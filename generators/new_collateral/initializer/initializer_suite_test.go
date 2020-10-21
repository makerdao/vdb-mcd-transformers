package initializer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestInitializer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Initializer Suite")
}
