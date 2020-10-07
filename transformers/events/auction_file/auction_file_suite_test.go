package auction_file_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlipFile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuctionFile Event Transformer Suite")
}
