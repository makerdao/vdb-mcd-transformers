package chop_dunk_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChopDunk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CatFileChopDunk Event Transformer Suite")
}
