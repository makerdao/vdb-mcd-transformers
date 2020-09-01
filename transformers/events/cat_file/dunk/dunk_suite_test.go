package dunk_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDunk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CatFileDunk Event Transformer Suite")
}
