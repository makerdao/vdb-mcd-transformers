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

package pit_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/pit"
	"github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/test_helpers"
)

var _ = Describe("Pit storage mappings", func() {
	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := pit.PitMappings{StorageRepository: storageRepository}

			Expect(mappings.Lookup(pit.DripKey)).To(Equal(pit.DripMetadata))
			Expect(mappings.Lookup(pit.LineKey)).To(Equal(pit.LineMetadata))
			Expect(mappings.Lookup(pit.LiveKey)).To(Equal(pit.LiveMetadata))
			Expect(mappings.Lookup(pit.VatKey)).To(Equal(pit.VatMetadata))
		})

		It("returns error if key does not exist", func() {
			mappings := pit.PitMappings{StorageRepository: &test_helpers.MockMakerStorageRepository{}}

			_, err := mappings.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := pit.PitMappings{StorageRepository: storageRepository}

			mappings.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetIlksCalled).To(BeTrue())
		})

		It("returns value metadata for spot when ilk in the DB", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			fakeIlk := "fakeIlk"
			storageRepository.Ilks = []string{fakeIlk}
			mappings := pit.PitMappings{StorageRepository: storageRepository}
			ilkSpotKey := common.BytesToHash(crypto.Keccak256(common.FromHex("0x" + fakeIlk + pit.IlkSpotIndex)))
			expectedMetadata := utils.StorageValueMetadata{
				Name: pit.IlkSpot,
				Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
				Type: utils.Uint256,
			}

			Expect(mappings.Lookup(ilkSpotKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for line when ilk in the DB", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			fakeIlk := "fakeIlk"
			storageRepository.Ilks = []string{fakeIlk}
			mappings := pit.PitMappings{StorageRepository: storageRepository}
			ilkSpotKeyBytes := crypto.Keccak256(common.FromHex("0x" + fakeIlk + pit.IlkSpotIndex))
			ilkSpotAsInt := big.NewInt(0).SetBytes(ilkSpotKeyBytes)
			incrementedIlkSpot := big.NewInt(0).Add(ilkSpotAsInt, big.NewInt(1))
			ilkLineKey := common.BytesToHash(incrementedIlkSpot.Bytes())
			expectedMetadata := utils.StorageValueMetadata{
				Name: pit.IlkLine,
				Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
				Type: utils.Uint256,
			}

			Expect(mappings.Lookup(ilkLineKey)).To(Equal(expectedMetadata))
		})

		It("returns error if key not found", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := pit.PitMappings{StorageRepository: storageRepository}

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})
})
