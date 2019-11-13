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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/spot"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"math/rand"
	"strconv"
)

var _ = Describe("Spot storage repository", func() {
	var (
		db              *postgres.DB
		repo            spot.SpotStorageRepository
		fakeAddress     = "0x12345"
		fakeBlockNumber = 123
		fakeBlockHash   = "expected_block_hash"
		fakeUint256     = "12345"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = spot.SpotStorageRepository{}
		repo.SetDB(db)
	})

	Describe("Ilk", func() {
		Describe("Pip", func() {
			It("writes a row", func() {
				ilkPipMetadata := utils.GetStorageValueMetadata(spot.IlkPip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkPipMetadata, fakeAddress)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, pip AS VALUE FROM maker.spot_ilk_pip`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeAddress)
			})

			It("does not duplicate row", func() {
				ilkPipMetadata := utils.GetStorageValueMetadata(spot.IlkPip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address)
				insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkPipMetadata, fakeAddress)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkPipMetadata, fakeAddress)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.spot_ilk_pip`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkPipMetadata := utils.GetStorageValueMetadata(spot.IlkPip, nil, utils.Address)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkPipMetadata, fakeAddress)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      utils.GetStorageValueMetadata(spot.IlkPip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address),
				PropertyName:  "Pip",
				PropertyValue: fakeAddress,
			})
		})

		Describe("Mat", func() {
			It("writes a row", func() {
				ilkMatMetadata := utils.GetStorageValueMetadata(spot.IlkMat, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkMatMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS KEY, mat AS VALUE FROM maker.spot_ilk_mat`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())

				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkMatMetadata := utils.GetStorageValueMetadata(spot.IlkMat, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkMatMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkMatMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.spot_ilk_mat`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkMatMetadata := utils.GetStorageValueMetadata(spot.IlkMat, nil, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkMatMetadata, fakeUint256)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
				Repository:    &repo,
				Metadata:      utils.GetStorageValueMetadata(spot.IlkMat, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
				PropertyName:  "Mat",
				PropertyValue: strconv.Itoa(rand.Int()),
			})
		})
	})

	It("persists a spot vat", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, spot.VatMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vat AS value FROM maker.spot_vat`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
	})

	It("does not duplicate spot vat", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, spot.VatMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, spot.VatMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.spot_vat`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a spot par", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, spot.ParMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, par AS value FROM maker.spot_par`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate spot par", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, spot.ParMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, spot.ParMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.spot_par`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})
