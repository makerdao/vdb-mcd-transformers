package clip_kick_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClipKick(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Clip Kick Suite")
}
