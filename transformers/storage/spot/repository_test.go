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
	"database/sql"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/spot"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("Spot storage repository", func() {
	var (
		db                 *postgres.DB
		repo               spot.SpotStorageRepository
		fakeAddress        = "0x12345"
		anotherFakeAddress = "0xedcba"
		fakeBlockNumber    = 123
		fakeBlockHash      = "expected_block_hash"
		fakeUint256        = "12345"
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

			Describe("updating current_ilk_state trigger table", func() {
				It("inserts a row for new ilk identifier", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					ilkPipMetadata := utils.GetStorageValueMetadata(
						spot.IlkPip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address)
					expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}

					err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkPipMetadata, fakeAddress)
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, pip, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Pip).To(Equal(fakeAddress))
					Expect(ilkState.Created).To(Equal(expectedTime))
					Expect(ilkState.Updated).To(Equal(expectedTime))
				})

				It("updates time created if new diff is from earlier block", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTimeUpdated := sql.NullString{String: formattedTimestamp, Valid: true}

					// set up old ilk state in later block
					ilkPipMetadata := utils.GetStorageValueMetadata(spot.IlkPip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, pip, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeAddress, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					// set up earlier header
					earlierBlockNumber := fakeBlockNumber - 1
					earlierTimestamp := rawTimestamp - 1
					CreateHeader(earlierTimestamp, earlierBlockNumber, db)
					formattedEarlierTimestamp := FormatTimestamp(earlierTimestamp)
					expectedTimeCreated := sql.NullString{String: formattedEarlierTimestamp, Valid: true}

					// trigger new ilk state from earlier block
					err := repo.Create(earlierBlockNumber, fakeBlockHash, ilkPipMetadata, anotherFakeAddress)
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, pip, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Pip).To(Equal(fakeAddress))
					Expect(ilkState.Created).To(Equal(expectedTimeCreated))
					Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
				})

				It("updates pip and time updated if new diff is from later block", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTimeCreated := sql.NullString{String: formattedTimestamp, Valid: true}

					// set up old ilk state in earlier block
					ilkPipMetadata := utils.GetStorageValueMetadata(spot.IlkPip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, pip, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeAddress, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					// set up later header
					laterBlockNumber := fakeBlockNumber + 1
					laterTimestamp := rawTimestamp + 1
					CreateHeader(laterTimestamp, laterBlockNumber, db)
					formattedLaterTimestamp := FormatTimestamp(laterTimestamp)
					expectedTimeUpdated := sql.NullString{String: formattedLaterTimestamp, Valid: true}

					// trigger new ilk state from later block
					newPip := anotherFakeAddress
					err := repo.Create(laterBlockNumber, fakeBlockHash, ilkPipMetadata, newPip)
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, pip, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Pip).To(Equal(newPip))
					Expect(ilkState.Created).To(Equal(expectedTimeCreated))
					Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
				})

				It("otherwise leaves row as is", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTime := sql.NullString{String: formattedTimestamp, Valid: true}

					ilkPipMetadata := utils.GetStorageValueMetadata(spot.IlkPip, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Address)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, pip, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeAddress, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkPipMetadata, fakeAddress)
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, pip, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Pip).To(Equal(fakeAddress))
					Expect(ilkState.Created).To(Equal(expectedTime))
					Expect(ilkState.Updated).To(Equal(expectedTime))
				})
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

			Describe("updating current_ilk_state trigger table", func() {
				It("inserts a row for new ilk identifier", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					ilkMatMetadata := utils.GetStorageValueMetadata(
						spot.IlkMat, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}

					err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkMatMetadata, fakeUint256)
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, mat, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Mat).To(Equal(fakeUint256))
					Expect(ilkState.Created).To(Equal(expectedTime))
					Expect(ilkState.Updated).To(Equal(expectedTime))
				})

				It("updates time created if new diff is from earlier block", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTimeUpdated := sql.NullString{String: formattedTimestamp, Valid: true}

					// set up old ilk state in later block
					ilkMatMetadata := utils.GetStorageValueMetadata(spot.IlkMat, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, mat, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					// set up earlier header
					earlierBlockNumber := fakeBlockNumber - 1
					earlierTimestamp := rawTimestamp - 1
					CreateHeader(earlierTimestamp, earlierBlockNumber, db)
					formattedEarlierTimestamp := FormatTimestamp(earlierTimestamp)
					expectedTimeCreated := sql.NullString{String: formattedEarlierTimestamp, Valid: true}

					// trigger new ilk state from earlier block
					err := repo.Create(earlierBlockNumber, fakeBlockHash, ilkMatMetadata, strconv.Itoa(rand.Int()))
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, mat, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Mat).To(Equal(fakeUint256))
					Expect(ilkState.Created).To(Equal(expectedTimeCreated))
					Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
				})

				It("updates mat and time updated if new diff is from later block", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTimeCreated := sql.NullString{String: formattedTimestamp, Valid: true}

					// set up old ilk state in earlier block
					ilkMatMetadata := utils.GetStorageValueMetadata(spot.IlkMat, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, mat, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					// set up later header
					laterBlockNumber := fakeBlockNumber + 1
					laterTimestamp := rawTimestamp + 1
					CreateHeader(laterTimestamp, laterBlockNumber, db)
					formattedLaterTimestamp := FormatTimestamp(laterTimestamp)
					expectedTimeUpdated := sql.NullString{String: formattedLaterTimestamp, Valid: true}

					// trigger new ilk state from later block
					newMat := strconv.Itoa(rand.Int())
					err := repo.Create(laterBlockNumber, fakeBlockHash, ilkMatMetadata, newMat)
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, mat, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Mat).To(Equal(newMat))
					Expect(ilkState.Created).To(Equal(expectedTimeCreated))
					Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
				})

				It("otherwise leaves row as is", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTime := sql.NullString{String: formattedTimestamp, Valid: true}

					ilkMatMetadata := utils.GetStorageValueMetadata(spot.IlkMat, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, mat, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkMatMetadata, strconv.Itoa(rand.Int()))
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, mat, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Mat).To(Equal(fakeUint256))
					Expect(ilkState.Created).To(Equal(expectedTime))
					Expect(ilkState.Updated).To(Equal(expectedTime))
				})
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
