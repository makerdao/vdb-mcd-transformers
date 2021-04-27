package clip_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClip(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Clip Storage Transformer Suite")
}
