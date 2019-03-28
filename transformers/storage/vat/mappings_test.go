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

package vat_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
)

var _ = Describe("Vat storage mappings", func() {
	var (
		fakeIlk           = "fakeIlk"
		fakeGuy           = "fakeGuy"
		storageRepository *test_helpers.MockMakerStorageRepository
		mappings          vat.VatMappings
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		mappings = vat.VatMappings{StorageRepository: storageRepository}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(mappings.Lookup(vat.DebtKey)).To(Equal(vat.DebtMetadata))
			Expect(mappings.Lookup(vat.ViceKey)).To(Equal(vat.ViceMetadata))
			Expect(mappings.Lookup(vat.LineKey)).To(Equal(vat.LineMetadata))
			Expect(mappings.Lookup(vat.LiveKey)).To(Equal(vat.LiveMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := mappings.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			mappings.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetDaiKeysCalled).To(BeTrue())
			Expect(storageRepository.GetGemKeysCalled).To(BeTrue())
			Expect(storageRepository.GetIlksCalled).To(BeTrue())
			Expect(storageRepository.GetSinKeysCalled).To(BeTrue())
			Expect(storageRepository.GetUrnsCalled).To(BeTrue())
		})

		It("returns error if dai keys lookup fails", func() {
			storageRepository.GetDaiKeysError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if gem keys lookup fails", func() {
			storageRepository.GetGemKeysError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if ilks lookup fails", func() {
			storageRepository.GetIlksError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if sin keys lookup fails", func() {
			storageRepository.GetSinKeysError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if urns lookup fails", func() {
			storageRepository.GetUrnsError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		Describe("ilk", func() {
			var ilkArtKey common.Hash
			var ilkArtAsInt *big.Int

			BeforeEach(func() {
				storageRepository.Ilks = []string{fakeIlk}
				ilkArtKey = common.BytesToHash(crypto.Keccak256(common.FromHex("0x" + fakeIlk + vat.IlksMappingIndex)))
				ilkArtAsInt = big.NewInt(0).SetBytes(ilkArtKey.Bytes())
			})

			It("returns value metadata for ilk Art", func() {
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.IlkArt,
					Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
					Type: utils.Uint256,
				}

				Expect(mappings.Lookup(ilkArtKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk rate", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(1))
				ilkRateKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.IlkRate,
					Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
					Type: utils.Uint256,
				}

				Expect(mappings.Lookup(ilkRateKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk spot", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(2))
				ilkSpotKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.IlkSpot,
					Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
					Type: utils.Uint256,
				}

				Expect(mappings.Lookup(ilkSpotKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk line", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(3))
				ilkLineKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.IlkLine,
					Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
					Type: utils.Uint256,
				}

				Expect(mappings.Lookup(ilkLineKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk dust", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(4))
				ilkDustKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.IlkDust,
					Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
					Type: utils.Uint256,
				}

				Expect(mappings.Lookup(ilkDustKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("urn", func() {
			It("returns value metadata for urn ink", func() {
				storageRepository.Urns = []storage.Urn{{Ilk: fakeIlk, Guy: fakeGuy}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex("0x" + fakeIlk + vat.UrnsMappingIndex))
				encodedSecondaryMapIndex := crypto.Keccak256(common.FromHex(fakeGuy), encodedPrimaryMapIndex)
				urnInkKey := common.BytesToHash(encodedSecondaryMapIndex)
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.UrnInk,
					Keys: map[utils.Key]string{constants.Ilk: fakeIlk, constants.Guy: fakeGuy},
					Type: utils.Uint256,
				}

				Expect(mappings.Lookup(urnInkKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for urn art", func() {
				storageRepository.Urns = []storage.Urn{{Ilk: fakeIlk, Guy: fakeGuy}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex("0x" + fakeIlk + vat.UrnsMappingIndex))
				urnInkAsInt := big.NewInt(0).SetBytes(crypto.Keccak256(common.FromHex(fakeGuy), encodedPrimaryMapIndex))
				incrementedUrnInk := big.NewInt(0).Add(urnInkAsInt, big.NewInt(1))
				urnArtKey := common.BytesToHash(incrementedUrnInk.Bytes())
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.UrnArt,
					Keys: map[utils.Key]string{constants.Ilk: fakeIlk, constants.Guy: fakeGuy},
					Type: utils.Uint256,
				}

				Expect(mappings.Lookup(urnArtKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("gem", func() {
			It("returns value metadata for gem", func() {
				storageRepository.GemKeys = []storage.Urn{{Ilk: fakeIlk, Guy: fakeGuy}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex("0x" + fakeIlk + vat.GemsMappingIndex))
				encodedSecondaryMapIndex := crypto.Keccak256(common.FromHex(fakeGuy), encodedPrimaryMapIndex)
				gemKey := common.BytesToHash(encodedSecondaryMapIndex)
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.Gem,
					Keys: map[utils.Key]string{constants.Ilk: fakeIlk, constants.Guy: fakeGuy},
					Type: utils.Uint256,
				}

				Expect(mappings.Lookup(gemKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("dai", func() {
			It("returns value metadata for dai", func() {
				storageRepository.DaiKeys = []string{fakeGuy}
				daiKey := common.BytesToHash(crypto.Keccak256(common.FromHex("0x" + fakeGuy + vat.DaiMappingIndex)))
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.Dai,
					Keys: map[utils.Key]string{constants.Guy: fakeGuy},
					Type: utils.Uint256,
				}

				Expect(mappings.Lookup(daiKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("when sin key exists in the db", func() {
			It("returns value metadata for sin", func() {
				storageRepository.SinKeys = []string{fakeGuy}
				sinKey := common.BytesToHash(crypto.Keccak256(common.FromHex("0x" + fakeGuy + vat.SinMappingIndex)))
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.Sin,
					Keys: map[utils.Key]string{constants.Guy: fakeGuy},
					Type: utils.Uint256,
				}

				Expect(mappings.Lookup(sinKey)).To(Equal(expectedMetadata))
			})
		})
	})
})
