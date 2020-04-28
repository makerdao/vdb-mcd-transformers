package osm_change_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOsmChange(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OsmChange Event Transformer Suite")
}
