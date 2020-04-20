package vow_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vow Component Test Suite")
}
