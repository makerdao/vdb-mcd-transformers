package vow_heal_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVowHeal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VowHeal Event Transformer Suite")
}
