package cdp_manager_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCdpManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CDP Manager Component Test Suite")
}
