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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/spot"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/big"
)

var _ = Describe("spot storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = spot.NewKeysLoader(storageRepository)
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[spot.VatKey]).To(Equal(spot.VatMetadata))
		Expect(mappings[spot.ParKey]).To(Equal(spot.ParMetadata))
	})

	Describe("when getting ilks fails", func() {
		It("returns error", func() {
			storageRepository.GetIlksError = fakes.FakeError

			_, err := storageKeysLoader.LoadMappings()

			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("when getting ilks succeeds", func() {
		It("returns value metadata for pip", func() {
			fakeIlk := "fakeIlk"
			storageRepository.Ilks = []string{fakeIlk}
			ilkPipKey := common.BytesToHash(crypto.Keccak256(common.FromHex("0x" + fakeIlk + spot.IlkMappingIndex)))
			expectedMetadata := utils.StorageValueMetadata{
				Name: spot.IlkPip,
				Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
				Type: utils.Address,
			}

			mappings, err := storageKeysLoader.LoadMappings()

			Expect(err).NotTo(HaveOccurred())
			Expect(mappings[ilkPipKey]).To(Equal(expectedMetadata))
		})

		It("returns value metadata for mat", func() {
			fakeIlk := "fakeIlk"
			storageRepository.Ilks = []string{fakeIlk}
			ilkPipKeyBytes := crypto.Keccak256(common.FromHex("0x" + fakeIlk + spot.IlkMappingIndex))
			ilkPipAsInt := big.NewInt(0).SetBytes(ilkPipKeyBytes)
			incrementedIlkPip := big.NewInt(0).Add(ilkPipAsInt, big.NewInt(1))
			ilkMatKey := common.BytesToHash(incrementedIlkPip.Bytes())
			expectedMetadata := utils.StorageValueMetadata{
				Name: spot.IlkMat,
				Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
				Type: utils.Uint256,
			}

			mappings, err := storageKeysLoader.LoadMappings()

			Expect(err).NotTo(HaveOccurred())
			Expect(mappings[ilkMatKey]).To(Equal(expectedMetadata))
		})
	})
})
