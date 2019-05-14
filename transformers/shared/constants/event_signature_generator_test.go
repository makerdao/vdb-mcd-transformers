// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package constants_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Event signature generator", func() {
	Describe("generating non-anonymous event signatures", func() {
		It("generates bite event signature", func() {
			expected := test_data.KovanBiteSignature
			actual := constants.GetEventTopicZero("Bite(bytes32,address,uint256,uint256,uint256,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates the flap kick event signature", func() {
			expected := test_data.KovanFlapKickSignature
			actual := constants.GetEventTopicZero("Kick(uint256,uint256,uint256,address,uint48)")

			Expect(expected).To(Equal(actual))
		})

		It("generates flip kick event signature", func() {
			expected := test_data.KovanFlipKickSignature
			actual := constants.GetEventTopicZero("Kick(uint256,uint256,uint256,address,uint48,bytes32,uint256)")

			Expect(expected).To(Equal(actual))
		})

		It("generates pip log value signature", func() {
			expected := test_data.KovanPipLogValueSignature
			actual := constants.GetEventTopicZero("LogValue(bytes32)")

			Expect(expected).To(Equal(actual))
		})
	})

	Describe("generating LogNote event signatures", func() {
		It("generates cat file chop lump event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithZeroPadding("file(bytes32,bytes32,uint256)")

			Expect(test_data.KovanCatFileChopLumpSignature).To(Equal(actual))
		})

		It("generates cat file flip event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithZeroPadding("file(bytes32,bytes32,address)")

			Expect(test_data.KovanCatFileFlipSignature).To(Equal(actual))
		})

		It("generates cat file vow event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithZeroPadding("file(bytes32,address)")

			Expect(test_data.KovanCatFileVowSignature).To(Equal(actual))
		})

		It("generates flip tend event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithZeroPadding("tend(uint256,uint256,uint256)")

			Expect(test_data.KovanTendSignature).To(Equal(actual))
		})

		It("generates the jug file drip signature", func() {
			actual := constants.GetLogNoteTopicZeroWithZeroPadding("drip(bytes32)")

			Expect(test_data.KovanJugDripSignature).To(Equal(actual))
		})

		It("generates the jug file base signature", func() {
			actual := constants.GetLogNoteTopicZeroWithZeroPadding("file(bytes32,uint256)")

			Expect(test_data.KovanJugFileBaseSignature).To(Equal(actual))
		})

		It("generates the jug file ilk signature", func() {
			actual := constants.GetLogNoteTopicZeroWithZeroPadding("file(bytes32,bytes32,uint256)")

			Expect(test_data.KovanJugFileIlkSignature).To(Equal(actual))
		})

		It("generates the jug file vow signature", func() {
			actual := constants.GetLogNoteTopicZeroWithZeroPadding("file(bytes32,address)")

			Expect(test_data.KovanJugFileVowSignature).To(Equal(actual))
		})

		It("generates vat file ilk event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithLeadingZeros("file(bytes32,bytes32,uint256)")

			Expect(test_data.KovanVatFileIlkSignature).To(Equal(actual))
		})

		It("generates the vat file debt ceiling event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithLeadingZeros("file(bytes32,uint256)")

			Expect(test_data.KovanVatFileDebtCeilingSignature).To(Equal(actual))
		})

		It("generates the vat flux event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithLeadingZeros("flux(bytes32,address,address,uint256)")

			Expect(test_data.KovanVatFluxSignature).To(Equal(actual))
		})

		It("generates the vat fold event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithLeadingZeros("fold(bytes32,address,int256)")

			Expect(test_data.KovanVatFoldSignature).To(Equal(actual))
		})

		It("generates the vat frob event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithLeadingZeros("frob(bytes32,address,address,address,int256,int256)")

			Expect(test_data.KovanVatFrobSignature).To(Equal(actual))
		})

		It("generates the vat grab event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithLeadingZeros("grab(bytes32,address,address,address,int256,int256)")

			Expect(test_data.KovanVatGrabSignature).To(Equal(actual))
		})

		It("generates the vat heal event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithLeadingZeros("heal(address,address,int256)")

			Expect(test_data.KovanVatHealSignature).To(Equal(actual))
		})

		It("generates the vat init event signature", func() {
			actual := constants.GetLogNoteTopicZeroWithLeadingZeros("init(bytes32)")

			Expect(test_data.KovanVatInitSignature).To(Equal(actual))
		})

		It("generates the vat move signature", func() {
			actual := constants.GetLogNoteTopicZeroWithLeadingZeros("move(address,address,uint256)")

			Expect(test_data.KovanVatMoveSignature).To(Equal(actual))
		})

		It("generates the vat slip signature", func() {
			actual := constants.GetLogNoteTopicZeroWithLeadingZeros("slip(bytes32,address,int256)")

			Expect(test_data.KovanVatSlipSignature).To(Equal(actual))
		})
	})

	Describe("getting the solidity method/event signature from the abi", func() {
		Describe("it handles methods", func() {
			Describe("from the cat contract", func() {
				It("gets the file method signature", func() {
					expected := "file(bytes32,bytes32,address)"
					actual := constants.GetSolidityFunctionSignature(constants.CatABI(), "file")

					Expect(expected).To(Equal(actual))
				})
			})

			Describe("from the jug contract", func() {
				It("gets the drip method signature", func() {
					expected := "drip(bytes32)"
					actual := constants.GetSolidityFunctionSignature(constants.JugABI(), "drip")

					Expect(expected).To(Equal(actual))
				})
			})

			Describe("from the flipper contract", func() {
				It("gets the deal method signature", func() {
					expected := "deal(uint256)"
					actual := constants.GetSolidityFunctionSignature(constants.FlipperABI(), "deal")

					Expect(expected).To(Equal(actual))
				})

				It("gets the dent method signature", func() {
					expected := "dent(uint256,uint256,uint256)"
					actual := constants.GetSolidityFunctionSignature(constants.FlipperABI(), "dent")

					Expect(expected).To(Equal(actual))
				})

				It("gets the tend method signature", func() {
					expected := "tend(uint256,uint256,uint256)"
					actual := constants.GetSolidityFunctionSignature(constants.FlipperABI(), "tend")

					Expect(expected).To(Equal(actual))
				})
			})

			Describe("from the jug contract", func() {
				It("gets the file (vow) method signature", func() {
					expected := "file(bytes32,address)"
					actual := constants.GetSolidityFunctionSignature(constants.JugABI(), "file")

					Expect(expected).To(Equal(actual))
				})
			})

			Describe("from the vat contract", func() {
				It("gets the flux method signature", func() {
					expected := "flux(bytes32,address,address,uint256)"
					actual := constants.GetSolidityFunctionSignature(constants.VatABI(), "flux")

					Expect(expected).To(Equal(actual))
				})

				It("gets the fold method signature", func() {
					expected := "fold(bytes32,address,int256)"
					actual := constants.GetSolidityFunctionSignature(constants.VatABI(), "fold")

					Expect(expected).To(Equal(actual))
				})

				It("gets the frob method signature", func() {
					expected := "frob(bytes32,address,address,address,int256,int256)"
					actual := constants.GetSolidityFunctionSignature(constants.VatABI(), "frob")

					Expect(expected).To(Equal(actual))
				})

				It("gets the grab method signature", func() {
					expected := "grab(bytes32,address,address,address,int256,int256)"
					actual := constants.GetSolidityFunctionSignature(constants.VatABI(), "grab")

					Expect(expected).To(Equal(actual))
				})

				It("gets the heal method signature", func() {
					expected := "heal(address,address,int256)"
					actual := constants.GetSolidityFunctionSignature(constants.VatABI(), "heal")

					Expect(expected).To(Equal(actual))
				})

				It("gets the init method signature", func() {
					expected := "init(bytes32)"
					actual := constants.GetSolidityFunctionSignature(constants.VatABI(), "init")

					Expect(expected).To(Equal(actual))
				})

				It("gets the move method signature", func() {
					expected := "move(address,address,uint256)"
					actual := constants.GetSolidityFunctionSignature(constants.VatABI(), "move")

					Expect(expected).To(Equal(actual))
				})

				It("gets the slip method signature", func() {
					expected := "slip(bytes32,address,int256)"
					actual := constants.GetSolidityFunctionSignature(constants.VatABI(), "slip")

					Expect(expected).To(Equal(actual))
				})
			})

			Describe("from the vow contract", func() {
				It("gets the fess method signature", func() {
					expected := "fess(uint256)"
					actual := constants.GetSolidityFunctionSignature(constants.VowABI(), "fess")

					Expect(expected).To(Equal(actual))
				})

				It("gets the flog method signature", func() {
					expected := "flog(uint48)"
					actual := constants.GetSolidityFunctionSignature(constants.VowABI(), "flog")

					Expect(expected).To(Equal(actual))
				})
			})
		})

		Describe("it handles events", func() {
			It("gets the Bite event signature", func() {
				expected := "Bite(bytes32,address,uint256,uint256,uint256,uint256)"
				actual := constants.GetSolidityFunctionSignature(constants.CatABI(), "Bite")

				Expect(expected).To(Equal(actual))
			})

			It("gets the flap Kick event signature", func() {
				expected := "Kick(uint256,uint256,uint256,address,uint48)"
				actual := constants.GetSolidityFunctionSignature(constants.FlapperABI(), "Kick")

				Expect(expected).To(Equal(actual))
			})

			It("gets the flip Kick event signature", func() {
				expected := "Kick(uint256,uint256,uint256,address,uint48,bytes32,uint256)"
				actual := constants.GetSolidityFunctionSignature(constants.FlipperABI(), "Kick")

				Expect(expected).To(Equal(actual))
			})

			It("gets the flop Kick event signature", func() {
				expected := "Kick(uint256,uint256,uint256,address,uint48)"
				actual := constants.GetSolidityFunctionSignature(constants.FlopperABI(), "Kick")

				Expect(expected).To(Equal(actual))
			})

			It("gets the log value method signature", func() {
				expected := "LogValue(bytes32)"
				actual := constants.GetSolidityFunctionSignature(constants.PipABI(), "LogValue")

				Expect(expected).To(Equal(actual))
			})
		})
	})
})
