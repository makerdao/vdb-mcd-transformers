// VulcanizeDB
// Copyright © 2019 Vulcanize

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
	It("generates auction file signature", func() {
		Expect(AuctionFileSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates bite signature", func() {
		Expect(BiteSignature()).To(Equal("0xa716da86bc1fb6d43d1493373f34d7a418b619681cd7b90f7ea667ba1489be28"))
	})

	It("generates cat file box signature", func() {
		Expect(CatFileBoxSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates cat file chop lump dunk signature", func() {
		Expect(CatFileChopLumpDunkSignature()).To(Equal("0x1a0b287e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates cat file flip signature", func() {
		Expect(CatFileFlipSignature()).To(Equal("0xebecb39d00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates cat file vow signature", func() {
		Expect(CatFileVowSignature()).To(Equal("0xd4e8be8300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates cat claw signature", func() {
		Expect(CatClawSignature()).To(Equal("0xe66d279b00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates deal signature", func() {
		Expect(DealSignature()).To(Equal("0xc959c42b00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates dent signature", func() {
		Expect(DentSignature()).To(Equal("0x5ff3a38200000000000000000000000000000000000000000000000000000000"))
	})

	It("generates deny signature", func() {
		Expect(DenySignature()).To(Equal("0x9c52a7f100000000000000000000000000000000000000000000000000000000"))
	})

	It("generates flap kick signature", func() {
		Expect(FlapKickSignature()).To(Equal("0xe6dde59cbc017becba89714a037778d234a84ce7f0a137487142a007e580d609"))
	})

	It("generates flip file cat signature", func() {
		Expect(FlipFileCatSignature()).To(Equal("0xd4e8be8300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates flip kick signature", func() {
		Expect(FlipKickSignature()).To(Equal("0xc84ce3a1172f0dec3173f04caaa6005151a4bfe40d4c9f3ea28dba5f719b2a7a"))
	})

	It("generates flop kick signature", func() {
		Expect(FlopKickSignature()).To(Equal("0x7e8881001566f9f89aedb9c5dc3d856a2b81e5235a8196413ed484be91cc0df6"))
	})

	It("generates jug drip signature", func() {
		Expect(JugDripSignature()).To(Equal("0x44e2a5a800000000000000000000000000000000000000000000000000000000"))
	})

	It("generates jug file base signature", func() {
		Expect(JugFileBaseSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates jug file ilk signature", func() {
		Expect(JugFileIlkSignature()).To(Equal("0x1a0b287e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates jug file vow signature", func() {
		Expect(JugFileVowSignature()).To(Equal("0xd4e8be8300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates jug init signature", func() {
		Expect(JugInitSignature()).To(Equal("0x3b66319500000000000000000000000000000000000000000000000000000000"))
	})

	It("generates log median price", func() {
		Expect(LogMedianPriceSignature()).To(Equal("0xb78ebc573f1f889ca9e1e0fb62c843c836f3d3a2e1f43ef62940e9b894f4ea4c"))
	})

	It("generates log value signature", func() {
		Expect(LogValueSignature()).To(Equal("0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527"))
	})

	It("generates median diss (batch) signature", func() {
		Expect(MedianDissBatchSignature()).To(Equal("0x46d4577d00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates median diss (single) signature", func() {
		Expect(MedianDissSingleSignature()).To(Equal("0x65c4ce7a00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates median drop signature", func() {
		Expect(MedianDropSignature()).To(Equal("0x8ef5eaf000000000000000000000000000000000000000000000000000000000"))
	})

	It("generates median kiss (batch) signature", func() {
		Expect(MedianKissBatchSignature()).To(Equal("0x1b25b65f00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates median kiss (single) signature", func() {
		Expect(MedianKissSingleSignature()).To(Equal("0xf29c29c400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates median lift signature", func() {
		Expect(MedianLiftSignature()).To(Equal("0x9431810600000000000000000000000000000000000000000000000000000000"))
	})

	It("generates new cdp signature", func() {
		Expect(NewCdpSignature()).To(Equal("0xd6be0bc178658a382ff4f91c8c68b542aa6b71685b8fe427966b87745c3ea7a2"))
	})

	It("generates pot cage signature", func() {
		Expect(PotCageSignature()).To(Equal("0x6924500900000000000000000000000000000000000000000000000000000000"))
	})

	It("generates osm change signature", func() {
		Expect(OsmChangeSignature()).To(Equal("0x1e77933e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates pot drip signature", func() {
		Expect(PotDripSignature()).To(Equal("0x9f678cca00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates pot exit signature", func() {
		Expect(PotExitSignature()).To(Equal("0x7f8661a100000000000000000000000000000000000000000000000000000000"))
	})

	It("generates pot file dsr signature", func() {
		Expect(PotFileDSRSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates pot file vow signature", func() {
		Expect(PotFileVowSignature()).To(Equal("0xd4e8be8300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates pot join signature", func() {
		Expect(PotJoinSignature()).To(Equal("0x049878f300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates rely signature", func() {
		Expect(RelySignature()).To(Equal("0x65fae35e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates spot file mat signature", func() {
		Expect(SpotFileMatSignature()).To(Equal("0x1a0b287e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates spot file par signature", func() {
		Expect(SpotFileParSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates spot file pip signature", func() {
		Expect(SpotFilePipSignature()).To(Equal("0xebecb39d00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates spot poke signature", func() {
		Expect(SpotPokeSignature()).To(Equal("0xdfd7467e425a8107cfd368d159957692c25085aacbcf5228ce08f10f2146486e"))
	})

	It("generates tend signature", func() {
		Expect(TendSignature()).To(Equal("0x4b43ed1200000000000000000000000000000000000000000000000000000000"))
	})

	It("generates tick signature", func() {
		Expect(TickSignature()).To(Equal("0xfc7b6aee00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat file debt ceiling signature", func() {
		Expect(VatFileDebtCeilingSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat file ilk signature", func() {
		Expect(VatFileIlkSignature()).To(Equal("0x1a0b287e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat flux signature", func() {
		Expect(VatFluxSignature()).To(Equal("0x6111be2e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat fold signature", func() {
		Expect(VatFoldSignature()).To(Equal("0xb65337df00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat fork signature", func() {
		Expect(VatForkSignature()).To(Equal("0x870c616d00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat frob signature", func() {
		Expect(VatFrobSignature()).To(Equal("0x7608870300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat grab signature", func() {
		Expect(VatGrabSignature()).To(Equal("0x7bab3f4000000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat heal signature", func() {
		Expect(VatHealSignature()).To(Equal("0xf37ac61c00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat hope signature", func() {
		Expect(VatHopeSignature()).To(Equal("0xa3b22fc400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat init signature", func() {
		Expect(VatInitSignature()).To(Equal("0x3b66319500000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat move signature", func() {
		Expect(VatMoveSignature()).To(Equal("0xbb35783b00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat nope signature", func() {
		Expect(VatNopeSignature()).To(Equal("0xdc4d20fa00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat slip signature", func() {
		Expect(VatSlipSignature()).To(Equal("0x7cdd3fde00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat suck signature", func() {
		Expect(VatSuckSignature()).To(Equal("0xf24e23eb00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vow fess signature", func() {
		Expect(VowFessSignature()).To(Equal("0x697efb7800000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vow file auction attributes signature", func() {
		Expect(VowFileAuctionAttributesSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vow file auction address signature", func() {
		Expect(VowFileAuctionAddressSignature()).To(Equal("0xd4e8be8300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vow flog signature", func() {
		Expect(VowFlogSignature()).To(Equal("0xd7ee674b00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vow heal signature", func() {
		Expect(VowHealSignature()).To(Equal("0xf37ac61c00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates yank signature", func() {
		Expect(YankSignature()).To(Equal("0x26e027f100000000000000000000000000000000000000000000000000000000"))
	})
})
