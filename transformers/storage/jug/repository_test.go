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
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"math/rand"
	"strconv"
)

var _ = Describe("Jug storage repository", func() {
	var (
		db              *postgres.DB
		repo            jug.JugStorageRepository
		fakeAddress     = "0x12345"
		fakeBlockNumber = 123
		fakeBlockHash   = "expected_block_hash"
		fakeUint256     = "12345"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = jug.JugStorageRepository{}
		repo.SetDB(db)
	})

	Describe("Ilk", func() {
		Describe("Rho", func() {
			It("writes a row", func() {
				ilkRhoMetadata := utils.GetStorageValueMetadata(jug.IlkRho, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRhoMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, rho AS VALUE FROM maker.jug_ilk_rho`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())
				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkRhoMetadata := utils.GetStorageValueMetadata(jug.IlkRho, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRhoMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRhoMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.jug_ilk_rho`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkRhoMetadata := utils.GetStorageValueMetadata(jug.IlkRho, nil, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkRhoMetadata, fakeUint256)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			Describe("updating current_ilk_state trigger table", func() {
				It("inserts a row for new ilk identifier", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					ilkRhoMetadata := utils.GetStorageValueMetadata(
						jug.IlkRho, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}

					err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRhoMetadata, fakeUint256)
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, rho, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Rho).To(Equal(fakeUint256))
					Expect(ilkState.Created).To(Equal(expectedTime))
					Expect(ilkState.Updated).To(Equal(expectedTime))
				})

				It("updates time created if new diff is from earlier block", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTimeUpdated := sql.NullString{String: formattedTimestamp, Valid: true}

					// set up old ilk state in later block
					ilkRhoMetadata := utils.GetStorageValueMetadata(jug.IlkRho, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, rho, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					// set up earlier header
					earlierBlockNumber := fakeBlockNumber - 1
					earlierTimestamp := rawTimestamp - 1
					CreateHeader(earlierTimestamp, earlierBlockNumber, db)
					formattedEarlierTimestamp := FormatTimestamp(earlierTimestamp)
					expectedTimeCreated := sql.NullString{String: formattedEarlierTimestamp, Valid: true}

					// trigger new ilk state from earlier block
					err := repo.Create(earlierBlockNumber, fakeBlockHash, ilkRhoMetadata, strconv.Itoa(rand.Int()))
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, rho, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Rho).To(Equal(fakeUint256))
					Expect(ilkState.Created).To(Equal(expectedTimeCreated))
					Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
				})

				It("updates rho and time updated if new diff is from later block", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTimeCreated := sql.NullString{String: formattedTimestamp, Valid: true}

					// set up old ilk state in earlier block
					ilkRhoMetadata := utils.GetStorageValueMetadata(jug.IlkRho, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, rho, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					// set up later header
					laterBlockNumber := fakeBlockNumber + 1
					laterTimestamp := rawTimestamp + 1
					CreateHeader(laterTimestamp, laterBlockNumber, db)
					formattedLaterTimestamp := FormatTimestamp(laterTimestamp)
					expectedTimeUpdated := sql.NullString{String: formattedLaterTimestamp, Valid: true}

					// trigger new ilk state from later block
					newRho := strconv.Itoa(rand.Int())
					err := repo.Create(laterBlockNumber, fakeBlockHash, ilkRhoMetadata, newRho)
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, rho, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Rho).To(Equal(newRho))
					Expect(ilkState.Created).To(Equal(expectedTimeCreated))
					Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
				})

				It("otherwise leaves row as is", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTime := sql.NullString{String: formattedTimestamp, Valid: true}

					ilkRhoMetadata := utils.GetStorageValueMetadata(jug.IlkRho, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, rho, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRhoMetadata, strconv.Itoa(rand.Int()))
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, rho, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Rho).To(Equal(fakeUint256))
					Expect(ilkState.Created).To(Equal(expectedTime))
					Expect(ilkState.Updated).To(Equal(expectedTime))
				})
			})
		})

		Describe("Duty", func() {
			It("writes a row", func() {
				ilkDutyMetadata := utils.GetStorageValueMetadata(jug.IlkDuty, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDutyMetadata, fakeUint256)

				Expect(err).NotTo(HaveOccurred())
				var result MappingRes
				err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS KEY, duty AS VALUE FROM maker.jug_ilk_duty`)
				Expect(err).NotTo(HaveOccurred())
				ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
				Expect(err).NotTo(HaveOccurred())

				AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
			})

			It("does not duplicate row", func() {
				ilkDutyMetadata := utils.GetStorageValueMetadata(jug.IlkDuty, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDutyMetadata, fakeUint256)
				Expect(insertOneErr).NotTo(HaveOccurred())

				insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDutyMetadata, fakeUint256)

				Expect(insertTwoErr).NotTo(HaveOccurred())
				var count int
				getCountErr := db.Get(&count, `SELECT count(*) FROM maker.jug_ilk_duty`)
				Expect(getCountErr).NotTo(HaveOccurred())
				Expect(count).To(Equal(1))
			})

			It("returns an error if metadata missing ilk", func() {
				malformedIlkDutyMetadata := utils.GetStorageValueMetadata(jug.IlkDuty, nil, utils.Uint256)

				err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkDutyMetadata, fakeUint256)
				Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
			})

			Describe("updating current_ilk_state trigger table", func() {
				It("inserts a row for new ilk identifier", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					ilkDutyMetadata := utils.GetStorageValueMetadata(
						jug.IlkDuty, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}

					err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDutyMetadata, fakeUint256)
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, duty, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Duty).To(Equal(fakeUint256))
					Expect(ilkState.Created).To(Equal(expectedTime))
					Expect(ilkState.Updated).To(Equal(expectedTime))
				})

				It("updates time created if new diff is from earlier block", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTimeUpdated := sql.NullString{String: formattedTimestamp, Valid: true}

					// set up old ilk state in later block
					ilkDutyMetadata := utils.GetStorageValueMetadata(jug.IlkDuty, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, duty, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					// set up earlier header
					earlierBlockNumber := fakeBlockNumber - 1
					earlierTimestamp := rawTimestamp - 1
					CreateHeader(earlierTimestamp, earlierBlockNumber, db)
					formattedEarlierTimestamp := FormatTimestamp(earlierTimestamp)
					expectedTimeCreated := sql.NullString{String: formattedEarlierTimestamp, Valid: true}

					// trigger new ilk state from earlier block
					err := repo.Create(earlierBlockNumber, fakeBlockHash, ilkDutyMetadata, strconv.Itoa(rand.Int()))
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, duty, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Duty).To(Equal(fakeUint256))
					Expect(ilkState.Created).To(Equal(expectedTimeCreated))
					Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
				})

				It("updates duty and time updated if new diff is from later block", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTimeCreated := sql.NullString{String: formattedTimestamp, Valid: true}

					// set up old ilk state in earlier block
					ilkDutyMetadata := utils.GetStorageValueMetadata(jug.IlkDuty, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, duty, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					// set up later header
					laterBlockNumber := fakeBlockNumber + 1
					laterTimestamp := rawTimestamp + 1
					CreateHeader(laterTimestamp, laterBlockNumber, db)
					formattedLaterTimestamp := FormatTimestamp(laterTimestamp)
					expectedTimeUpdated := sql.NullString{String: formattedLaterTimestamp, Valid: true}

					// trigger new ilk state from later block
					newDuty := strconv.Itoa(rand.Int())
					err := repo.Create(laterBlockNumber, fakeBlockHash, ilkDutyMetadata, newDuty)
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, duty, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Duty).To(Equal(newDuty))
					Expect(ilkState.Created).To(Equal(expectedTimeCreated))
					Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
				})

				It("otherwise leaves row as is", func() {
					rawTimestamp := int64(rand.Int31())
					CreateHeader(rawTimestamp, fakeBlockNumber, db)
					formattedTimestamp := FormatTimestamp(rawTimestamp)
					expectedTime := sql.NullString{String: formattedTimestamp, Valid: true}

					ilkDutyMetadata := utils.GetStorageValueMetadata(jug.IlkDuty, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
					_, insertErr := db.Exec(
						`INSERT INTO api.current_ilk_state (ilk_identifier, duty, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
						test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
					Expect(insertErr).NotTo(HaveOccurred())

					err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDutyMetadata, strconv.Itoa(rand.Int()))
					Expect(err).NotTo(HaveOccurred())

					var ilkState test_helpers.IlkState
					queryErr := db.Get(&ilkState, `SELECT ilk_identifier, duty, created, updated FROM api.current_ilk_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
					Expect(ilkState.Duty).To(Equal(fakeUint256))
					Expect(ilkState.Created).To(Equal(expectedTime))
					Expect(ilkState.Updated).To(Equal(expectedTime))
				})
			})
		})
	})

	It("persists a jug vat", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, jug.VatMetadata, fakeAddress)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vat AS value FROM maker.jug_vat`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeAddress)
	})

	It("does not duplicate jug vat", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, jug.VatMetadata, fakeAddress)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, jug.VatMetadata, fakeAddress)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.jug_vat`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a jug vow", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, jug.VowMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())
		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vow AS value FROM maker.jug_vow`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate jug vow", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, jug.VowMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, jug.VowMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.jug_vow`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists a jug base", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, jug.BaseMetadata, fakeUint256)
		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, base AS value FROM maker.jug_base`)
		Expect(err).NotTo(HaveOccurred())

		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate jug base", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, jug.BaseMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, jug.BaseMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.jug_base`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})
