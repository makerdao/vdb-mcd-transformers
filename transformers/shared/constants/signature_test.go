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

package constants_test

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Signature constants", func() {
	It("generates auction file signature", func() {
		Expect(constants.AuctionFileSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates bite signature", func() {
		Expect(constants.BiteSignature()).To(Equal("0xa716da86bc1fb6d43d1493373f34d7a418b619681cd7b90f7ea667ba1489be28"))
	})

	It("generates cat claw signature", func() {
		Expect(constants.CatClawSignature()).To(Equal("0xe66d279b00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates cat file box signature", func() {
		Expect(constants.CatFileBoxSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates cat file chop lump dunk signature", func() {
		Expect(constants.CatFileChopLumpDunkSignature()).To(Equal("0x1a0b287e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates cat file flip signature", func() {
		Expect(constants.CatFileFlipSignature()).To(Equal("0xebecb39d00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates cat file vow signature", func() {
		Expect(constants.CatFileVowSignature()).To(Equal("0xd4e8be8300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates clip kick signature", func() {
		Expect(constants.ClipKickSignature()).To(Equal("0x7c5bfdc0a5e8192f6cd4972f382cec69116862fb62e6abff8003874c58e064b8"))
	})

	It("generates clip take signature", func() {
		Expect(constants.ClipTakeSignature()).To(Equal("0x05e309fd6ce72f2ab888a20056bb4210df08daed86f21f95053deb19964d86b1"))
	})

	It("generates clip redo signature", func() {
		Expect(constants.ClipRedoSignature()).To(Equal("0x275de7ecdd375b5e8049319f8b350686131c219dd4dc450a08e9cf83b03c865f"))
	})

	It("generates clip yank signature", func() {
		Expect(constants.ClipYankSignature()).To(Equal("0x2c5d2826eb5903b8fc201cf48094b858f42f61c7eaac9aaf43ebed490138144e"))
	})

	It("generates deal signature", func() {
		Expect(constants.DealSignature()).To(Equal("0xc959c42b00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates dent signature", func() {
		Expect(constants.DentSignature()).To(Equal("0x5ff3a38200000000000000000000000000000000000000000000000000000000"))
	})

	It("generates deny signature", func() {
		Expect(constants.DenySignature()).To(Equal("0x9c52a7f100000000000000000000000000000000000000000000000000000000"))
	})

	It("generates dog bark signature", func() {
		Expect(constants.DogBarkSignature()).To(Equal("0x85258d09e1e4ef299ff3fc11e74af99563f022d21f3f940db982229dc2a3358c"))
	})

	It("generates dog digs signature", func() {
		Expect(constants.DogDigsSignature()).To(Equal("0x54f095dc7308776bf01e8580e4dd40fd959ea4bf50b069975768320ef8d77d8a"))
	})

	It("generates dog file hole signature", func() {
		Expect(constants.DogFileHoleSignature()).To(Equal("0xe986e40cc8c151830d4f61050f4fb2e4add8567caad2d5f5496f9158e91fe4c7"))
	})

	It("generates dog file ilk clip signature", func() {
		Expect(constants.DogFileIlkClipSignature()).To(Equal("0x4ff2caaa972a7c6629ea01fae9c93d73cc307d13ea4c369f9bbbb7f9b7e9461d"))
	})

	It("generates dog file vow signature", func() {
		Expect(constants.DogFileVowSignature()).To(Equal("0x8fef588b5fc1afbf5b2f06c1a435d513f208da2e6704c3d8f0e0ec91167066ba"))
	})

	It("generates dog deny signature", func() {
		Expect(constants.DogDenySignature()).To(Equal("0x184450df2e323acec0ed3b5c7531b81f9b4cdef7914dfd4c0a4317416bb5251b"))
	})

	It("generates dog rely signature", func() {
		Expect(constants.DogRelySignature()).To(Equal("0xdd0e34038ac38b2a1ce960229778ac48a8719bc900b6c4f8d0475c6e8b385a60"))
	})

	It("generates dog file ilk chop hole signature", func() {
		Expect(constants.DogFileIlkChopHoleSignature()).To(Equal("0x851aa1caf4888170ad8875449d18f0f512fd6deb2a6571ea1a41fb9f95acbcd1"))
	})

	It("generates flap kick signature", func() {
		Expect(constants.FlapKickSignature()).To(Equal("0xe6dde59cbc017becba89714a037778d234a84ce7f0a137487142a007e580d609"))
	})

	It("generates flip file cat signature", func() {
		Expect(constants.FlipFileCatSignature()).To(Equal("0xd4e8be8300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates flip kick signature", func() {
		Expect(constants.FlipKickSignature()).To(Equal("0xc84ce3a1172f0dec3173f04caaa6005151a4bfe40d4c9f3ea28dba5f719b2a7a"))
	})

	It("generates flop kick signature", func() {
		Expect(constants.FlopKickSignature()).To(Equal("0x7e8881001566f9f89aedb9c5dc3d856a2b81e5235a8196413ed484be91cc0df6"))
	})

	It("generates jug drip signature", func() {
		Expect(constants.JugDripSignature()).To(Equal("0x44e2a5a800000000000000000000000000000000000000000000000000000000"))
	})

	It("generates jug file base signature", func() {
		Expect(constants.JugFileBaseSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates jug file ilk signature", func() {
		Expect(constants.JugFileIlkSignature()).To(Equal("0x1a0b287e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates jug file vow signature", func() {
		Expect(constants.JugFileVowSignature()).To(Equal("0xd4e8be8300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates jug init signature", func() {
		Expect(constants.JugInitSignature()).To(Equal("0x3b66319500000000000000000000000000000000000000000000000000000000"))
	})

	It("generates log median price", func() {
		Expect(constants.LogMedianPriceSignature()).To(Equal("0xb78ebc573f1f889ca9e1e0fb62c843c836f3d3a2e1f43ef62940e9b894f4ea4c"))
	})

	It("generates log value signature", func() {
		Expect(constants.LogValueSignature()).To(Equal("0x296ba4ca62c6c21c95e828080cb8aec7481b71390585605300a8a76f9e95b527"))
	})

	It("generates median diss (batch) signature", func() {
		Expect(constants.MedianDissBatchSignature()).To(Equal("0x46d4577d00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates median diss (single) signature", func() {
		Expect(constants.MedianDissSingleSignature()).To(Equal("0x65c4ce7a00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates median drop signature", func() {
		Expect(constants.MedianDropSignature()).To(Equal("0x8ef5eaf000000000000000000000000000000000000000000000000000000000"))
	})

	It("generates median kiss (batch) signature", func() {
		Expect(constants.MedianKissBatchSignature()).To(Equal("0x1b25b65f00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates median kiss (single) signature", func() {
		Expect(constants.MedianKissSingleSignature()).To(Equal("0xf29c29c400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates median lift signature", func() {
		Expect(constants.MedianLiftSignature()).To(Equal("0x9431810600000000000000000000000000000000000000000000000000000000"))
	})

	It("generates new cdp signature", func() {
		Expect(constants.NewCdpSignature()).To(Equal("0xd6be0bc178658a382ff4f91c8c68b542aa6b71685b8fe427966b87745c3ea7a2"))
	})

	It("generates osm change signature", func() {
		Expect(constants.OsmChangeSignature()).To(Equal("0x1e77933e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates pot cage signature", func() {
		Expect(constants.PotCageSignature()).To(Equal("0x6924500900000000000000000000000000000000000000000000000000000000"))
	})

	It("generates pot drip signature", func() {
		Expect(constants.PotDripSignature()).To(Equal("0x9f678cca00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates pot exit signature", func() {
		Expect(constants.PotExitSignature()).To(Equal("0x7f8661a100000000000000000000000000000000000000000000000000000000"))
	})

	It("generates pot file dsr signature", func() {
		Expect(constants.PotFileDSRSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates pot file vow signature", func() {
		Expect(constants.PotFileVowSignature()).To(Equal("0xd4e8be8300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates pot join signature", func() {
		Expect(constants.PotJoinSignature()).To(Equal("0x049878f300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates rely signature", func() {
		Expect(constants.RelySignature()).To(Equal("0x65fae35e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates spot file mat signature", func() {
		Expect(constants.SpotFileMatSignature()).To(Equal("0x1a0b287e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates spot file par signature", func() {
		Expect(constants.SpotFileParSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates spot file pip signature", func() {
		Expect(constants.SpotFilePipSignature()).To(Equal("0xebecb39d00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates spot poke signature", func() {
		Expect(constants.SpotPokeSignature()).To(Equal("0xdfd7467e425a8107cfd368d159957692c25085aacbcf5228ce08f10f2146486e"))
	})

	It("generates tend signature", func() {
		Expect(constants.TendSignature()).To(Equal("0x4b43ed1200000000000000000000000000000000000000000000000000000000"))
	})

	It("generates tick signature", func() {
		Expect(constants.TickSignature()).To(Equal("0xfc7b6aee00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat file debt ceiling signature", func() {
		Expect(constants.VatFileDebtCeilingSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat file ilk signature", func() {
		Expect(constants.VatFileIlkSignature()).To(Equal("0x1a0b287e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat flux signature", func() {
		Expect(constants.VatFluxSignature()).To(Equal("0x6111be2e00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat fold signature", func() {
		Expect(constants.VatFoldSignature()).To(Equal("0xb65337df00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat fork signature", func() {
		Expect(constants.VatForkSignature()).To(Equal("0x870c616d00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat frob signature", func() {
		Expect(constants.VatFrobSignature()).To(Equal("0x7608870300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat grab signature", func() {
		Expect(constants.VatGrabSignature()).To(Equal("0x7bab3f4000000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat heal signature", func() {
		Expect(constants.VatHealSignature()).To(Equal("0xf37ac61c00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat hope signature", func() {
		Expect(constants.VatHopeSignature()).To(Equal("0xa3b22fc400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat init signature", func() {
		Expect(constants.VatInitSignature()).To(Equal("0x3b66319500000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat move signature", func() {
		Expect(constants.VatMoveSignature()).To(Equal("0xbb35783b00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat nope signature", func() {
		Expect(constants.VatNopeSignature()).To(Equal("0xdc4d20fa00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat slip signature", func() {
		Expect(constants.VatSlipSignature()).To(Equal("0x7cdd3fde00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vat suck signature", func() {
		Expect(constants.VatSuckSignature()).To(Equal("0xf24e23eb00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vow fess signature", func() {
		Expect(constants.VowFessSignature()).To(Equal("0x697efb7800000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vow file auction address signature", func() {
		Expect(constants.VowFileAuctionAddressSignature()).To(Equal("0xd4e8be8300000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vow file auction attributes signature", func() {
		Expect(constants.VowFileAuctionAttributesSignature()).To(Equal("0x29ae811400000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vow flog signature", func() {
		Expect(constants.VowFlogSignature()).To(Equal("0xd7ee674b00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates vow heal signature", func() {
		Expect(constants.VowHealSignature()).To(Equal("0xf37ac61c00000000000000000000000000000000000000000000000000000000"))
	})

	It("generates yank signature", func() {
		Expect(constants.YankSignature()).To(Equal("0x26e027f100000000000000000000000000000000000000000000000000000000"))
	})
})
