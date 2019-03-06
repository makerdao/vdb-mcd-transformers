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

package drip_test

import (
	utils2 "github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/drip"
	. "github.com/vulcanize/mcd_transformers/transformers/storage_diffs/maker/test_helpers"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var _ = Describe("Drip storage repository", func() {
	var (
		db              *postgres.DB
		repo            drip.DripStorageRepository
		fakeAddress     = "0x12345"
		fakeBlockNumber = 123
		fakeBlockHash   = "expected_block_hash"
		fakeIlk         = "fake_ilk"
		fakeUint256     = "12345"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = drip.DripStorageRepository{}
		repo.SetDB(db)
	})

	Describe("Ilk", func() {
		Describe("Rho", func() {
			It("writes a row", func() {
				ilkRhoMetadata := utils.GetStorageValueMetadata(drip.IlkRho, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRhoMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT block_number, block_hash, ilk AS key, rho AS VALUE FROM maker.drip_ilk_rho`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := utils2.GetOrCreateIlk(fakeIlk, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkRhoMetadata := utils.GetStorageValueMetadata(drip.IlkRho, nil, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkRhoMetadata, fakeUint256)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})
		})

		Describe("Tax", func() {
			It("writes a row", func() {
				ilkTaxMetadata := utils.GetStorageValueMetadata(drip.IlkTax, map[utils.Key]string{constants.Ilk: fakeIlk}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkTaxMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT block_number, block_hash, ilk AS KEY, tax AS VALUE FROM maker.drip_ilk_tax`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := utils2.GetOrCreateIlk(fakeIlk, db)
				Expect(err).NotTo(HaveOccurred())

				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.Itoa(ilkID), fakeUint256)
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkTaxMetadata := utils.GetStorageValueMetadata(drip.IlkTax, nil, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkTaxMetadata, fakeUint256)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})
		})
	})

	It("persists a drip vat", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, drip.VatMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vat AS value FROM maker.drip_vat`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
	})

	It("persists a drip vow", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, drip.VowMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vow AS value FROM maker.drip_vow`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("persists a drip repo", func() {
		expectedRepo := "12345"

		err := repo.Create(fakeBlockNumber, fakeBlockHash, drip.RepoMetadata, expectedRepo)
		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, repo AS value FROM maker.drip_repo`)
		Expect(err).NotTo(HaveOccurred())

		AssertVariable(result, fakeBlockNumber, fakeBlockHash, expectedRepo)
	})
})
