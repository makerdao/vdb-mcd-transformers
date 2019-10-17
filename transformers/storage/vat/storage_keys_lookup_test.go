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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/big"
)

var _ = Describe("Vat storage mappings", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLookup vat.StorageKeysLookup
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLookup = vat.StorageKeysLookup{StorageRepository: storageRepository}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(storageKeysLookup.Lookup(vat.DebtKey)).To(Equal(vat.DebtMetadata))
			Expect(storageKeysLookup.Lookup(vat.ViceKey)).To(Equal(vat.ViceMetadata))
			Expect(storageKeysLookup.Lookup(vat.LineKey)).To(Equal(vat.LineMetadata))
			Expect(storageKeysLookup.Lookup(vat.LiveKey)).To(Equal(vat.LiveMetadata))
		})

		It("returns value metadata if keccak of key exists", func() {
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vat.DebtKey[:]))).To(Equal(vat.DebtMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vat.ViceKey[:]))).To(Equal(vat.ViceMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vat.LineKey[:]))).To(Equal(vat.LineMetadata))
			Expect(storageKeysLookup.Lookup(crypto.Keccak256Hash(vat.LiveKey[:]))).To(Equal(vat.LiveMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := storageKeysLookup.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetDaiKeysCalled).To(BeTrue())
			Expect(storageRepository.GetGemKeysCalled).To(BeTrue())
			Expect(storageRepository.GetIlksCalled).To(BeTrue())
			Expect(storageRepository.GetVatSinKeysCalled).To(BeTrue())
			Expect(storageRepository.GetUrnsCalled).To(BeTrue())
		})

		It("returns error if dai keys lookup fails", func() {
			storageRepository.GetDaiKeysError = fakes.FakeError

			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if gem keys lookup fails", func() {
			storageRepository.GetGemKeysError = fakes.FakeError

			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if ilks lookup fails", func() {
			storageRepository.GetIlksError = fakes.FakeError

			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if sin keys lookup fails", func() {
			storageRepository.GetVatSinKeysError = fakes.FakeError

			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if urns lookup fails", func() {
			storageRepository.GetUrnsError = fakes.FakeError

			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if lookups return addresses not of length 42", func() {
			storageRepository.DaiKeys = []string{"0xshortAddress"}

			_, err := storageKeysLookup.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
		})

		Describe("ilk", func() {
			var ilkArtKey common.Hash
			var ilkArtAsInt *big.Int

			BeforeEach(func() {
				storageRepository.Ilks = []string{test_helpers.FakeIlk}
				ilkArtKey = common.BytesToHash(crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + vat.IlksMappingIndex)))
				ilkArtAsInt = big.NewInt(0).SetBytes(ilkArtKey.Bytes())
			})

			It("returns value metadata for ilk Art", func() {
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.IlkArt,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}

				Expect(storageKeysLookup.Lookup(ilkArtKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk rate", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(1))
				ilkRateKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.IlkRate,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}

				Expect(storageKeysLookup.Lookup(ilkRateKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk spot", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(2))
				ilkSpotKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.IlkSpot,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}

				Expect(storageKeysLookup.Lookup(ilkSpotKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk line", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(3))
				ilkLineKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.IlkLine,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}

				Expect(storageKeysLookup.Lookup(ilkLineKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk dust", func() {
				incrementedIlkArt := big.NewInt(0).Add(ilkArtAsInt, big.NewInt(4))
				ilkDustKey := common.BytesToHash(incrementedIlkArt.Bytes())
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.IlkDust,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk},
					Type: utils.Uint256,
				}

				Expect(storageKeysLookup.Lookup(ilkDustKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("urn", func() {
			It("returns value metadata for urn ink", func() {
				storageRepository.Urns = []storage.Urn{{Ilk: test_helpers.FakeIlk, Identifier: test_helpers.FakeAddress}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + vat.UrnsMappingIndex))
				paddedUrnGuy := common.FromHex("0x000000000000000000000000" + test_helpers.FakeAddress[2:])
				encodedSecondaryMapIndex := crypto.Keccak256(paddedUrnGuy, encodedPrimaryMapIndex)
				urnInkKey := common.BytesToHash(encodedSecondaryMapIndex)
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.UrnInk,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk, constants.Guy: test_helpers.FakeAddress},
					Type: utils.Uint256,
				}

				Expect(storageKeysLookup.Lookup(urnInkKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for urn art", func() {
				storageRepository.Urns = []storage.Urn{{Ilk: test_helpers.FakeIlk, Identifier: test_helpers.FakeAddress}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + vat.UrnsMappingIndex))
				paddedUrnGuy := common.FromHex("0x000000000000000000000000" + test_helpers.FakeAddress[2:])
				urnInkKey := crypto.Keccak256(paddedUrnGuy, encodedPrimaryMapIndex)
				urnInkAsInt := big.NewInt(0).SetBytes(urnInkKey)
				incrementedUrnInk := big.NewInt(0).Add(urnInkAsInt, big.NewInt(1))
				urnArtKey := common.BytesToHash(incrementedUrnInk.Bytes())
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.UrnArt,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk, constants.Guy: test_helpers.FakeAddress},
					Type: utils.Uint256,
				}

				Expect(storageKeysLookup.Lookup(urnArtKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("gem", func() {
			It("returns value metadata for gem", func() {
				storageRepository.GemKeys = []storage.Urn{{Ilk: test_helpers.FakeIlk, Identifier: test_helpers.FakeAddress}}
				encodedPrimaryMapIndex := crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + vat.GemsMappingIndex))
				paddedGemAddress := common.FromHex("0x000000000000000000000000" + test_helpers.FakeAddress[2:])
				encodedSecondaryMapIndex := crypto.Keccak256(paddedGemAddress, encodedPrimaryMapIndex)
				gemKey := common.BytesToHash(encodedSecondaryMapIndex)
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.Gem,
					Keys: map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk, constants.Guy: test_helpers.FakeAddress},
					Type: utils.Uint256,
				}

				Expect(storageKeysLookup.Lookup(gemKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("dai", func() {
			It("returns value metadata for dai", func() {
				storageRepository.DaiKeys = []string{test_helpers.FakeAddress}
				paddedDaiAddress := "0x000000000000000000000000" + test_helpers.FakeAddress[2:]
				daiKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedDaiAddress + vat.DaiMappingIndex)))
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.Dai,
					Keys: map[utils.Key]string{constants.Guy: test_helpers.FakeAddress},
					Type: utils.Uint256,
				}

				Expect(storageKeysLookup.Lookup(daiKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("sin", func() {
			It("returns value metadata for sin", func() {
				storageRepository.SinKeys = []string{test_helpers.FakeAddress}
				paddedSinAddress := "0x000000000000000000000000" + test_helpers.FakeAddress[2:]
				sinKey := common.BytesToHash(crypto.Keccak256(common.FromHex(paddedSinAddress + vat.SinMappingIndex)))
				expectedMetadata := utils.StorageValueMetadata{
					Name: vat.Sin,
					Keys: map[utils.Key]string{constants.Guy: test_helpers.FakeAddress},
					Type: utils.Uint256,
				}

				Expect(storageKeysLookup.Lookup(sinKey)).To(Equal(expectedMetadata))
			})
		})
	})
})
