package auction_address_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVowFileAuctionAddress(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VowFileAuctionAddress Suite")
}
