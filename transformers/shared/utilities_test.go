package shared_test

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

var _ = Describe("Shared utilities", func() {
	Describe("getting log note data bytes at index", func() {
		Describe("extracting DSNote third argument from data", func() {
			It("accounts for topic zero's signature padding being appended to the end of the data", func() {
				actual := shared.GetDSNoteThirdArgument(test_data.EthCatFileChopLog.Data)

				dataInt := shared.ConvertUint256HexToBigInt(hexutil.Encode(actual))
				Expect(test_data.CatFileChopModel.Data).To(Equal(dataInt.String()))
			})
		})

		Describe("extracting Vat Note data", func() {
			It("returns error if index less than four", func() {
				_, err := shared.GetVatNoteDataBytesAtIndex(3, []byte{})

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(shared.ErrInvalidIndex(3)))
			})

			It("returns error if index greater than six", func() {
				_, err := shared.GetVatNoteDataBytesAtIndex(7, []byte{})

				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(shared.ErrInvalidIndex(7)))
			})

			It("extracts fourth argument of four arguments", func() {
				wadBytes, err := shared.GetVatNoteDataBytesAtIndex(4, test_data.EthVatFluxLog.Data)

				Expect(err).NotTo(HaveOccurred())
				wadInt := shared.ConvertUint256HexToBigInt(hexutil.Encode(wadBytes))
				Expect(wadInt.String()).To(Equal(test_data.VatFluxModel.Wad))
			})

			It("extracts fourth of five arguments", func() {
				dinkBytes, err := shared.GetVatNoteDataBytesAtIndex(4, test_data.EthVatForkLogWithNegativeDart.Data)

				Expect(err).NotTo(HaveOccurred())
				dinkInt := shared.ConvertInt256HexToBigInt(hexutil.Encode(dinkBytes))
				Expect(dinkInt.String()).To(Equal(test_data.VatForkModelWithNegativeDart.Dink))
			})

			It("extracts fifth of five arguments", func() {
				dartBytes, err := shared.GetVatNoteDataBytesAtIndex(5, test_data.EthVatForkLogWithNegativeDart.Data)

				Expect(err).NotTo(HaveOccurred())
				dartInt := shared.ConvertInt256HexToBigInt(hexutil.Encode(dartBytes))
				Expect(dartInt.String()).To(Equal(test_data.VatForkModelWithNegativeDart.Dart))
			})

			It("extracts the fourth of six arguments", func() {
				wBytes, err := shared.GetVatNoteDataBytesAtIndex(4, test_data.EthVatGrabLogWithPositiveDink.Data)

				Expect(err).NotTo(HaveOccurred())
				wAddress := common.BytesToAddress(wBytes)
				Expect(wAddress.String()).To(Equal(test_data.VatGrabModelWithPositiveDink.W))
			})

			It("extracts the fifth of six arguments", func() {
				dinkBytes, err := shared.GetVatNoteDataBytesAtIndex(5, test_data.EthVatGrabLogWithPositiveDink.Data)

				Expect(err).NotTo(HaveOccurred())
				dinkInt := shared.ConvertInt256HexToBigInt(hexutil.Encode(dinkBytes))
				Expect(dinkInt.String()).To(Equal(test_data.VatGrabModelWithPositiveDink.Dink))
			})

			It("extracts the sixth of six arguments", func() {
				dartBytes, err := shared.GetVatNoteDataBytesAtIndex(6, test_data.EthVatGrabLogWithPositiveDink.Data)

				Expect(err).NotTo(HaveOccurred())
				dartInt := shared.ConvertInt256HexToBigInt(hexutil.Encode(dartBytes))
				Expect(dartInt.String()).To(Equal(test_data.VatGrabModelWithPositiveDink.Dart))
			})
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

	Describe("decoding ilk name", func() {
		It("handles hex ilk", func() {
			actualIlkName, err := shared.DecodeIlkName(test_helpers.FakeIlk.Hex)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualIlkName).To(Equal(test_helpers.FakeIlk.Name))
		})

		It("handles hex ilk with leading 0x", func() {
			hexIlk := "0x" + test_helpers.FakeIlk.Hex
			actualIlkName, err := shared.DecodeIlkName(hexIlk)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualIlkName).To(Equal(test_helpers.FakeIlk.Name))
		})

		It("discards zero bytes", func() {
			hexIlk := "0x000000"
			actualIlkName, err := shared.DecodeIlkName(hexIlk)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualIlkName).To(Equal(""))
		})
	})
})
