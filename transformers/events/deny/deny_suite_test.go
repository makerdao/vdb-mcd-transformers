package deny_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDeny(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deny Suite")
}
