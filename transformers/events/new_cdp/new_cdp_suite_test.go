package new_cdp_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNewCdp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NewCdp Event Transformer Suite")
}
