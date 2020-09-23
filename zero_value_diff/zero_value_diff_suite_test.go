package zero_value_diff_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestZeroValueDiff(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ZeroValueDiff Suite")
}
