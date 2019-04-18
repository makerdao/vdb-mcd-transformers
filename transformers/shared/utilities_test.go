package shared_test

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

var _ = Describe("Shared utilities", func() {
	Describe("getting log note data bytes at index", func() {
		It("accounts for topic zero's signature padding being appended to the end of the data", func() {
			logData := hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c441fd3ef74554480000000000000000000000000000000000000000000000000000000000da15dce70ab462e66779f23ee14f21d993789ee3000000000000000000000000da15dce70ab462e66779f23ee14f21d993789ee3000000000000000000000000da15dce70ab462e66779f23ee14f21d993789ee30000000000000000000000000000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
			addressBytes := common.HexToAddress("0xda15dce70ab462e66779f23ee14f21d993789ee3").Bytes()
			// common.address.Bytes() returns [20]byte{}, need [32]byte{}
			expected := append(addressBytes, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}...)
			actual := shared.GetLogNoteDataBytesAtIndex(-3, logData)

			Expect(expected[:]).To(Equal(actual))
		})
	})

	Describe("getting hex without prefix", func() {
		It("returns bytes as hex without 0x prefix", func() {
			raw := common.HexToHash("0x4554480000000000000000000000000000000000000000000000000000000000").Bytes()
			result := shared.GetHexWithoutPrefix(raw)
			Expect(result).To(Equal("4554480000000000000000000000000000000000000000000000000000000000"))
		})
	})

	Describe("converting int256 hex to big int", func() {
		It("correctly converts positive number", func() {
			result := shared.ConvertInt256HexToBigInt("0x00000000000000000000000000000000000000000000000007a1fe1602770000")

			Expect(result.String()).To(Equal("550000000000000000"))
		})

		It("correctly converts negative number", func() {
			result := shared.ConvertInt256HexToBigInt("0xffffffffffffffffffffffffffffffffffffffffffffffffff4e5d43d13b0000")

			Expect(result.String()).To(Equal("-50000000000000000"))
		})

		It("correctly converts another negative number", func() {
			result := shared.ConvertInt256HexToBigInt("0xfffffffffffffffffffffffffffffffffffffffffffffffffe9cba87a2760000")

			Expect(result.String()).To(Equal("-100000000000000000"))
		})
	})
})
