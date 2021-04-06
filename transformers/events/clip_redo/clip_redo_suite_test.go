package clip_redo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClipRedo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Clip Redo Suite")
}
