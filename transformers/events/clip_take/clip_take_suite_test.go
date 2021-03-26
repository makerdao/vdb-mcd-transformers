package clip_take_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClipTake(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Clip Take Suite")
}
