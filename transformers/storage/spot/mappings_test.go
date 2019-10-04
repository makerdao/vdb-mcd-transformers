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

package spot_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/spot"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("spot storage mappings", func() {
	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := spot.SpotMappings{StorageRepository: storageRepository}

			Expect(mappings.Lookup(spot.VatKey)).To(Equal(spot.VatMetadata))
			Expect(mappings.Lookup(spot.ParKey)).To(Equal(spot.ParMetadata))
		})

		It("returns value metadata if keccak of key exists", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := spot.SpotMappings{StorageRepository: storageRepository}

			Expect(mappings.Lookup(crypto.Keccak256Hash(spot.VatKey[:]))).To(Equal(spot.VatMetadata))
			Expect(mappings.Lookup(crypto.Keccak256Hash(spot.ParKey[:]))).To(Equal(spot.ParMetadata))
		})

		It("returns error if key does not exist", func() {
			mappings := spot.SpotMappings{StorageRepository: &test_helpers.MockMakerStorageRepository{}}

			_, err := mappings.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from the repository if key not found", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := spot.SpotMappings{StorageRepository: storageRepository}

			mappings.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetIlksCalled).To(BeTrue())
		})

		It("returns value metadata for pip when ilk in the DB", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			fakeIlk := "fakeIlk"
			storageRepository.Ilks = []string{fakeIlk}
			mappings := spot.SpotMappings{StorageRepository: storageRepository}
			ilkPipKey := common.BytesToHash(crypto.Keccak256(common.FromHex("0x" + fakeIlk + spot.IlkMappingIndex)))
			expectedMetadata := utils.StorageValueMetadata{
				Name: spot.IlkPip,
				Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
				Type: utils.Address,
			}

			Expect(mappings.Lookup(ilkPipKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for mat when ilk in the DB", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			fakeIlk := "fakeIlk"
			storageRepository.Ilks = []string{fakeIlk}
			mappings := spot.SpotMappings{StorageRepository: storageRepository}
			ilkPipKeyBytes := crypto.Keccak256(common.FromHex("0x" + fakeIlk + spot.IlkMappingIndex))
			ilkPipAsInt := big.NewInt(0).SetBytes(ilkPipKeyBytes)
			incrementedIlkPip := big.NewInt(0).Add(ilkPipAsInt, big.NewInt(1))
			ilkMatKey := common.BytesToHash(incrementedIlkPip.Bytes())
			expectedMetadata := utils.StorageValueMetadata{
				Name: spot.IlkMat,
				Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
				Type: utils.Uint256,
			}

			Expect(mappings.Lookup(ilkMatKey)).To(Equal(expectedMetadata))
		})

		It("returns error if key not found", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := spot.SpotMappings{StorageRepository: storageRepository}

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})
})
