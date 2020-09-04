package cat_claw_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCatClaw(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CatClaw Suite")
}
