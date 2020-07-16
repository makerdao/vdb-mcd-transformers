package jug_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJug(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jug Component Test Suite")
}
