package transformer_exporter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTransformerExporter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TransformerExporter Suite")
}
