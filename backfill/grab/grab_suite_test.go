package grab_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGrab(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BackFill Grab Suite")
}
