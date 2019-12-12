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

package jug_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	vdbStorage "github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("jug storage keys loader", func() {
	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		storageKeysLoader storage.KeysLoader
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		storageKeysLoader = jug.NewKeysLoader(storageRepository)
	})

	It("returns value metadata for static keys", func() {
		mappings, err := storageKeysLoader.LoadMappings()

		Expect(err).NotTo(HaveOccurred())
		Expect(mappings[jug.VatKey]).To(Equal(jug.VatMetadata))
		Expect(mappings[jug.VowKey]).To(Equal(jug.VowMetadata))
		Expect(mappings[jug.BaseKey]).To(Equal(jug.BaseMetadata))
	})

	Describe("when getting ilks fails", func() {
		It("returns error", func() {
			storageRepository.GetIlksError = fakes.FakeError

			_, err := storageKeysLoader.LoadMappings()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})
	})

	Describe("when getting ilks succeeds", func() {
		It("returns value metadata for tax", func() {
			storageRepository.Ilks = []string{test_helpers.FakeIlk}
			ilkTaxKey := common.BytesToHash(crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + jug.IlkMappingIndex)))
			expectedMetadata := vdbStorage.ValueMetadata{
				Name: jug.IlkDuty,
				Keys: map[vdbStorage.Key]string{constants.Ilk: test_helpers.FakeIlk},
				Type: vdbStorage.Uint256,
			}

			mappings, err := storageKeysLoader.LoadMappings()

			Expect(err).NotTo(HaveOccurred())
			Expect(mappings[ilkTaxKey]).To(Equal(expectedMetadata))
		})

		It("returns value metadata for rho", func() {
			storageRepository.Ilks = []string{test_helpers.FakeIlk}
			ilkTaxKeyBytes := crypto.Keccak256(common.FromHex(test_helpers.FakeIlk + jug.IlkMappingIndex))
			ilkTaxAsInt := big.NewInt(0).SetBytes(ilkTaxKeyBytes)
			incrementedIlkTax := big.NewInt(0).Add(ilkTaxAsInt, big.NewInt(1))
			ilkRhoKey := common.BytesToHash(incrementedIlkTax.Bytes())
			expectedMetadata := vdbStorage.ValueMetadata{
				Name: jug.IlkRho,
				Keys: map[vdbStorage.Key]string{constants.Ilk: test_helpers.FakeIlk},
				Type: vdbStorage.Uint256,
			}

			mappings, err := storageKeysLoader.LoadMappings()

			Expect(err).NotTo(HaveOccurred())
			Expect(mappings[ilkRhoKey]).To(Equal(expectedMetadata))
		})
	})
})
