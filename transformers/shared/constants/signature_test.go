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

package constants

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Signature constants", func() {
	It("generates bite signature", func() {
		Expect(BiteSignature()).To(Equal("0xa716da86bc1fb6d43d1493373f34d7a418b619681cd7b90f7ea667ba1489be28"))
	})

	It("generates cat file chop lump signature", func() {
		Expect(CatFileChopLumpSignature()).To(Equal("0x000000000000000000000000000000000000000000000000000000001a0b287e"))
	})

	It("generates cat file flip signature", func() {
		Expect(CatFileFlipSignature()).To(Equal("0x00000000000000000000000000000000000000000000000000000000ebecb39d"))
	})

	It("generates cat file vow signature", func() {
		Expect(CatFileVowSignature()).To(Equal("0x00000000000000000000000000000000000000000000000000000000d4e8be83"))
	})

	It("generates deal signature", func() {
		Expect(DealSignature()).To(Equal("0x00000000000000000000000000000000000000000000000000000000c959c42b"))
	})

	It("generates dent signature", func() {
		Expect(DentSignature()).To(Equal("0x000000000000000000000000000000000000000000000000000000005ff3a382"))
	})

	It("generates flap kick signature", func() {
		Expect(FlapKickSignature()).To(Equal("0xefa52d9342a199cb30efd2692463f2c2bef63cd7186b50382d4fb94ad207880e"))
	})

	It("generates flip kick signature", func() {
		Expect(FlipKickSignature()).To(Equal("0xbac86238bdba81d21995024470425ecb370078fa62b7271b90cf28cbd1e3e87e"))
	})

	It("generates flop kick signature", func() {
		Expect(FlopKickSignature()).To(Equal("0xefa52d9342a199cb30efd2692463f2c2bef63cd7186b50382d4fb94ad207880e"))
	})

	It("generates jug drip signature", func() {
		Expect(JugDripSignature()).To(Equal("0x0000000000000000000000000000000000000000000000000000000044e2a5a8"))
	})

	It("generates jug file base signature", func() {
		Expect(JugFileBaseSignature()).To(Equal("0x0000000000000000000000000000000000000000000000000000000029ae8114"))
	})

	It("generates jug file ilk signature", func() {
		Expect(JugFileIlkSignature()).To(Equal("0x000000000000000000000000000000000000000000000000000000001a0b287e"))
	})

	It("generates jug file vow signature", func() {
		Expect(JugFileVowSignature()).To(Equal("0x00000000000000000000000000000000000000000000000000000000d4e8be83"))
	})

	It("generates pip log value signature", func() {
		Expect(PipLogValueSignature()).To(Equal("0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527"))
	})

	It("generates tend signature", func() {
		Expect(TendSignature()).To(Equal("0x000000000000000000000000000000000000000000000000000000004b43ed12"))
	})

	It("generates vat file debt ceiling signature", func() {
		Expect(VatFileDebtCeilingSignature()).To(Equal("0x0000000000000000000000000000000000000000000000000000000029ae8114"))
	})

	It("generates vat file ilk signature", func() {
		Expect(VatFileIlkSignature()).To(Equal("0x000000000000000000000000000000000000000000000000000000001a0b287e"))
	})

	It("generates vat flux signature", func() {
		Expect(VatFluxSignature()).To(Equal("0x000000000000000000000000000000000000000000000000000000006111be2e"))
	})

	It("generates vat fold signature", func() {
		Expect(VatFoldSignature()).To(Equal("0x00000000000000000000000000000000000000000000000000000000b65337df"))
	})

	It("generates vat frob signature", func() {
		Expect(VatFrobSignature()).To(Equal("0x0000000000000000000000000000000000000000000000000000000076088703"))
	})

	It("generates vat grab signature", func() {
		Expect(VatGrabSignature()).To(Equal("0x000000000000000000000000000000000000000000000000000000007bab3f40"))
	})

	It("generates vat heal signature", func() {
		Expect(VatHealSignature()).To(Equal("0x00000000000000000000000000000000000000000000000000000000f37ac61c"))
	})

	It("generates vat init signature", func() {
		Expect(VatInitSignature()).To(Equal("0x000000000000000000000000000000000000000000000000000000003b663195"))
	})

	It("generates vat move signature", func() {
		Expect(VatMoveSignature()).To(Equal("0x00000000000000000000000000000000000000000000000000000000bb35783b"))
	})

	It("generates vat slip signature", func() {
		Expect(VatSlipSignature()).To(Equal("0x000000000000000000000000000000000000000000000000000000007cdd3fde"))
	})

	It("generates vat suck signature", func() {
		Expect(VatSuckSignature()).To(Equal("0x00000000000000000000000000000000000000000000000000000000f24e23eb"))
	})

	It("generates vow fess signature", func() {
		Expect(VowFessSignature()).To(Equal("0x00000000000000000000000000000000000000000000000000000000697efb78"))
	})

	It("generates vow file signature", func() {
		Expect(VowFileSignature()).To(Equal("0x0000000000000000000000000000000000000000000000000000000029ae8114"))
	})

	It("generates vow flog signature", func() {
		Expect(VowFlogSignature()).To(Equal("0x0000000000000000000000000000000000000000000000000000000035aee16f"))
	})
})
