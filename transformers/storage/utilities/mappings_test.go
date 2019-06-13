package utilities_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/utilities"
)

var _ = Describe("Storage mappings utilities", func() {
	Describe("padding an address", func() {
		It("returns error if input not 20 bytes", func() {
			shortAddress := "0x23"

			_, err := utilities.PadAddress(shortAddress)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utilities.ErrInvalidAddress(shortAddress)))
		})

		It("returns address with padding to 32 bytes", func() {
			res, err := utilities.PadAddress(test_helpers.FakeAddress)

			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(Equal("0x000000000000000000000000" + test_helpers.FakeAddress[2:]))
		})
	})
})
