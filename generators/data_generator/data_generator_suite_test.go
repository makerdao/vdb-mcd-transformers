package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDataGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Data Generator suite")
}
