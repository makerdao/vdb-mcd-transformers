package clip_yank_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClipYank(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Clip Yank Suite")
}
