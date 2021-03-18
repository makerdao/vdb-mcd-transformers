package vow_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDogFileVow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dog File Vow Suite")
}
