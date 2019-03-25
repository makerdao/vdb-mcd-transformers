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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("jug storage mappings", func() {
	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := jug.JugMappings{StorageRepository: storageRepository}

			Expect(mappings.Lookup(jug.VatKey)).To(Equal(jug.VatMetadata))
			Expect(mappings.Lookup(jug.VowKey)).To(Equal(jug.VowMetadata))
			Expect(mappings.Lookup(jug.RepoKey)).To(Equal(jug.RepoMetadata))
		})

		It("returns error if key does not exist", func() {
			mappings := jug.JugMappings{StorageRepository: &test_helpers.MockMakerStorageRepository{}}

			_, err := mappings.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from the repository if key not found", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := jug.JugMappings{StorageRepository: storageRepository}

			mappings.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetIlksCalled).To(BeTrue())
		})

		It("returns value metadata for tax when ilk in the DB", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			fakeIlk := "fakeIlk"
			storageRepository.Ilks = []string{fakeIlk}
			mappings := jug.JugMappings{StorageRepository: storageRepository}
			ilkTaxKey := common.BytesToHash(crypto.Keccak256(common.FromHex("0x" + fakeIlk + jug.IlkMappingIndex)))
			expectedMetadata := utils.StorageValueMetadata{
				Name: jug.IlkTax,
				Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
				Type: utils.Uint256,
			}

			Expect(mappings.Lookup(ilkTaxKey)).To(Equal(expectedMetadata))
		})

		It("returns value metadata for rho when ilk in the DB", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			fakeIlk := "fakeIlk"
			storageRepository.Ilks = []string{fakeIlk}
			mappings := jug.JugMappings{StorageRepository: storageRepository}
			ilkTaxKeyBytes := crypto.Keccak256(common.FromHex("0x" + fakeIlk + jug.IlkMappingIndex))
			ilkTaxAsInt := big.NewInt(0).SetBytes(ilkTaxKeyBytes)
			incrementedIlkTax := big.NewInt(0).Add(ilkTaxAsInt, big.NewInt(1))
			ilkRhoKey := common.BytesToHash(incrementedIlkTax.Bytes())
			expectedMetadata := utils.StorageValueMetadata{
				Name: jug.IlkRho,
				Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
				Type: utils.Uint48,
			}

			Expect(mappings.Lookup(ilkRhoKey)).To(Equal(expectedMetadata))
		})

		It("returns error if key not found", func() {
			storageRepository := &test_helpers.MockMakerStorageRepository{}
			mappings := jug.JugMappings{StorageRepository: storageRepository}

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})
})
