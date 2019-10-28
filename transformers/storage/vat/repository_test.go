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
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"math/rand"
	"strconv"
)

var _ = Describe("Vat storage repository", func() {
	var (
		db              *postgres.DB
		repo            vat.VatStorageRepository
		fakeBlockNumber = rand.Int()
		fakeBlockHash   = "expected_block_hash"
		fakeGuy         = "fake_urn"
		fakeUint256     = "12345"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = vat.VatStorageRepository{}
		repo.SetDB(db)
	})

	Describe("dai", func() {
		It("writes a row", func() {
			daiMetadata := utils.GetStorageValueMetadata(vat.Dai, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, daiMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, guy AS key, dai AS value FROM maker.vat_dai`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			daiMetadata := utils.GetStorageValueMetadata(vat.Dai, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, daiMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, daiMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_dai`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing guy", func() {
			malformedDaiMetadata := utils.GetStorageValueMetadata(vat.Dai, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedDaiMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("gem", func() {
		It("writes row", func() {
			gemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, gemMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key_one, guy AS key_two, gem AS value FROM maker.vat_gem`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			gemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, gemMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, gemMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_gem`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedGemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedGemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("ilk Art", func() {
		It("writes row", func() {
			ilkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, art AS value FROM maker.vat_ilk_art`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkArtMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkArtMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_art`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		Describe("updating current_ilk_state trigger table", func() {
			It("inserts a row for new ilk identifier", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				ilkArtMetadata := utils.GetStorageValueMetadata(
					vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkArtMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, art, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Art).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTime))
				Expect(ilkState.Updated).To(Equal(expectedTime))
			})

			It("updates time created if new diff is from earlier block", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTimeUpdated := sql.NullString{String: formattedTimestamp, Valid: true}

				// set up old ilk state in later block
				ilkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, art, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				// set up earlier header
				earlierBlockNumber := fakeBlockNumber - 1
				earlierTimestamp := rawTimestamp - 1
				CreateHeader(earlierTimestamp, earlierBlockNumber, db)
				formattedEarlierTimestamp := FormatTimestamp(earlierTimestamp)
				expectedTimeCreated := sql.NullString{String: formattedEarlierTimestamp, Valid: true}

				// trigger new ilk state from earlier block
				err := repo.Create(earlierBlockNumber, fakeBlockHash, ilkArtMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, art, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Art).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTimeCreated))
				Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
			})

			It("updates art and time updated if new diff is from later block", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTimeCreated := sql.NullString{String: formattedTimestamp, Valid: true}

				// set up old ilk state in earlier block
				ilkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, art, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				// set up later header
				laterBlockNumber := fakeBlockNumber + 1
				laterTimestamp := rawTimestamp + 1
				CreateHeader(laterTimestamp, laterBlockNumber, db)
				formattedLaterTimestamp := FormatTimestamp(laterTimestamp)
				expectedTimeUpdated := sql.NullString{String: formattedLaterTimestamp, Valid: true}

				// trigger new ilk state from later block
				newArt := strconv.Itoa(rand.Int())
				err := repo.Create(laterBlockNumber, fakeBlockHash, ilkArtMetadata, newArt)
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, art, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Art).To(Equal(newArt))
				Expect(ilkState.Created).To(Equal(expectedTimeCreated))
				Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
			})

			It("otherwise leaves row as is", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTime := sql.NullString{String: formattedTimestamp, Valid: true}

				ilkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, art, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkArtMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, art, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Art).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTime))
				Expect(ilkState.Updated).To(Equal(expectedTime))
			})
		})
	})

	Describe("ilk dust", func() {
		It("writes row", func() {
			ilkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDustMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, dust AS value FROM maker.vat_ilk_dust`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDustMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDustMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_dust`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkDustMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		Describe("updating current_ilk_state trigger table", func() {
			It("inserts a row for new ilk identifier", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				ilkDustMetadata := utils.GetStorageValueMetadata(
					vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDustMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, dust, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Dust).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTime))
				Expect(ilkState.Updated).To(Equal(expectedTime))
			})

			It("updates time created if new diff is from earlier block", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTimeUpdated := sql.NullString{String: formattedTimestamp, Valid: true}

				// set up old ilk state in later block
				ilkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, dust, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				// set up earlier header
				earlierBlockNumber := fakeBlockNumber - 1
				earlierTimestamp := rawTimestamp - 1
				CreateHeader(earlierTimestamp, earlierBlockNumber, db)
				formattedEarlierTimestamp := FormatTimestamp(earlierTimestamp)
				expectedTimeCreated := sql.NullString{String: formattedEarlierTimestamp, Valid: true}

				// trigger new ilk state from earlier block
				err := repo.Create(earlierBlockNumber, fakeBlockHash, ilkDustMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, dust, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Dust).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTimeCreated))
				Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
			})

			It("updates dust and time updated if new diff is from later block", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTimeCreated := sql.NullString{String: formattedTimestamp, Valid: true}

				// set up old ilk state in earlier block
				ilkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, dust, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				// set up later header
				laterBlockNumber := fakeBlockNumber + 1
				laterTimestamp := rawTimestamp + 1
				CreateHeader(laterTimestamp, laterBlockNumber, db)
				formattedLaterTimestamp := FormatTimestamp(laterTimestamp)
				expectedTimeUpdated := sql.NullString{String: formattedLaterTimestamp, Valid: true}

				// trigger new ilk state from later block
				newDust := strconv.Itoa(rand.Int())
				err := repo.Create(laterBlockNumber, fakeBlockHash, ilkDustMetadata, newDust)
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, dust, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Dust).To(Equal(newDust))
				Expect(ilkState.Created).To(Equal(expectedTimeCreated))
				Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
			})

			It("otherwise leaves row as is", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTime := sql.NullString{String: formattedTimestamp, Valid: true}

				ilkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, dust, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkDustMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, dust, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Dust).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTime))
				Expect(ilkState.Updated).To(Equal(expectedTime))
			})
		})
	})

	Describe("ilk line", func() {
		It("writes row", func() {
			ilkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLineMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, line AS value FROM maker.vat_ilk_line`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLineMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLineMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_line`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkLineMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		Describe("updating current_ilk_state trigger table", func() {
			It("inserts a row for new ilk identifier", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				ilkLineMetadata := utils.GetStorageValueMetadata(
					vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLineMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, line, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Line).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTime))
				Expect(ilkState.Updated).To(Equal(expectedTime))
			})

			It("updates time created if new diff is from earlier block", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTimeUpdated := sql.NullString{String: formattedTimestamp, Valid: true}

				// set up old ilk state in later block
				ilkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, line, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				// set up earlier header
				earlierBlockNumber := fakeBlockNumber - 1
				earlierTimestamp := rawTimestamp - 1
				CreateHeader(earlierTimestamp, earlierBlockNumber, db)
				formattedEarlierTimestamp := FormatTimestamp(earlierTimestamp)
				expectedTimeCreated := sql.NullString{String: formattedEarlierTimestamp, Valid: true}

				// trigger new ilk state from earlier block
				err := repo.Create(earlierBlockNumber, fakeBlockHash, ilkLineMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, line, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Line).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTimeCreated))
				Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
			})

			It("updates line and time updated if new diff is from later block", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTimeCreated := sql.NullString{String: formattedTimestamp, Valid: true}

				// set up old ilk state in earlier block
				ilkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, line, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				// set up later header
				laterBlockNumber := fakeBlockNumber + 1
				laterTimestamp := rawTimestamp + 1
				CreateHeader(laterTimestamp, laterBlockNumber, db)
				formattedLaterTimestamp := FormatTimestamp(laterTimestamp)
				expectedTimeUpdated := sql.NullString{String: formattedLaterTimestamp, Valid: true}

				// trigger new ilk state from later block
				newLine := strconv.Itoa(rand.Int())
				err := repo.Create(laterBlockNumber, fakeBlockHash, ilkLineMetadata, newLine)
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, line, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Line).To(Equal(newLine))
				Expect(ilkState.Created).To(Equal(expectedTimeCreated))
				Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
			})

			It("otherwise leaves row as is", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTime := sql.NullString{String: formattedTimestamp, Valid: true}

				ilkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, line, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkLineMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, line, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Line).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTime))
				Expect(ilkState.Updated).To(Equal(expectedTime))
			})
		})
	})

	Describe("ilk rate", func() {
		It("writes row", func() {
			ilkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRateMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, rate AS value FROM maker.vat_ilk_rate`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRateMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRateMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_rate`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkRateMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		Describe("updating current_ilk_state trigger table", func() {
			It("inserts a row for new ilk identifier", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				ilkRateMetadata := utils.GetStorageValueMetadata(
					vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRateMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, rate, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Rate).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTime))
				Expect(ilkState.Updated).To(Equal(expectedTime))
			})

			It("updates time created if new diff is from earlier block", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTimeUpdated := sql.NullString{String: formattedTimestamp, Valid: true}

				// set up old ilk state in later block
				ilkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, rate, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				// set up earlier header
				earlierBlockNumber := fakeBlockNumber - 1
				earlierTimestamp := rawTimestamp - 1
				CreateHeader(earlierTimestamp, earlierBlockNumber, db)
				formattedEarlierTimestamp := FormatTimestamp(earlierTimestamp)
				expectedTimeCreated := sql.NullString{String: formattedEarlierTimestamp, Valid: true}

				// trigger new ilk state from earlier block
				err := repo.Create(earlierBlockNumber, fakeBlockHash, ilkRateMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, rate, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Rate).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTimeCreated))
				Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
			})

			It("updates rate and time updated if new diff is from later block", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTimeCreated := sql.NullString{String: formattedTimestamp, Valid: true}

				// set up old ilk state in earlier block
				ilkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, rate, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				// set up later header
				laterBlockNumber := fakeBlockNumber + 1
				laterTimestamp := rawTimestamp + 1
				CreateHeader(laterTimestamp, laterBlockNumber, db)
				formattedLaterTimestamp := FormatTimestamp(laterTimestamp)
				expectedTimeUpdated := sql.NullString{String: formattedLaterTimestamp, Valid: true}

				// trigger new ilk state from later block
				newRate := strconv.Itoa(rand.Int())
				err := repo.Create(laterBlockNumber, fakeBlockHash, ilkRateMetadata, newRate)
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, rate, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Rate).To(Equal(newRate))
				Expect(ilkState.Created).To(Equal(expectedTimeCreated))
				Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
			})

			It("otherwise leaves row as is", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTime := sql.NullString{String: formattedTimestamp, Valid: true}

				ilkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, rate, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkRateMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, rate, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Rate).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTime))
				Expect(ilkState.Updated).To(Equal(expectedTime))
			})
		})
	})

	Describe("ilk spot", func() {
		It("writes row", func() {
			ilkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkSpotMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, ilk_id AS key, spot AS value FROM maker.vat_ilk_spot`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkSpotMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, ilkSpotMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_spot`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedIlkSpotMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		Describe("updating current_ilk_state trigger table", func() {
			It("inserts a row for new ilk identifier", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				ilkSpotMetadata := utils.GetStorageValueMetadata(
					vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				expectedTime := sql.NullString{String: FormatTimestamp(rawTimestamp), Valid: true}

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkSpotMetadata, fakeUint256)
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, spot, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Spot).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTime))
				Expect(ilkState.Updated).To(Equal(expectedTime))
			})

			It("updates time created if new diff is from earlier block", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTimeUpdated := sql.NullString{String: formattedTimestamp, Valid: true}

				// set up old ilk state in later block
				ilkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, spot, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				// set up earlier header
				earlierBlockNumber := fakeBlockNumber - 1
				earlierTimestamp := rawTimestamp - 1
				CreateHeader(earlierTimestamp, earlierBlockNumber, db)
				formattedEarlierTimestamp := FormatTimestamp(earlierTimestamp)
				expectedTimeCreated := sql.NullString{String: formattedEarlierTimestamp, Valid: true}

				// trigger new ilk state from earlier block
				err := repo.Create(earlierBlockNumber, fakeBlockHash, ilkSpotMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, spot, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Spot).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTimeCreated))
				Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
			})

			It("updates spot and time updated if new diff is from later block", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTimeCreated := sql.NullString{String: formattedTimestamp, Valid: true}

				// set up old ilk state in earlier block
				ilkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, spot, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				// set up later header
				laterBlockNumber := fakeBlockNumber + 1
				laterTimestamp := rawTimestamp + 1
				CreateHeader(laterTimestamp, laterBlockNumber, db)
				formattedLaterTimestamp := FormatTimestamp(laterTimestamp)
				expectedTimeUpdated := sql.NullString{String: formattedLaterTimestamp, Valid: true}

				// trigger new ilk state from later block
				newSpot := strconv.Itoa(rand.Int())
				err := repo.Create(laterBlockNumber, fakeBlockHash, ilkSpotMetadata, newSpot)
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, spot, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Spot).To(Equal(newSpot))
				Expect(ilkState.Created).To(Equal(expectedTimeCreated))
				Expect(ilkState.Updated).To(Equal(expectedTimeUpdated))
			})

			It("otherwise leaves row as is", func() {
				rawTimestamp := int64(rand.Int31())
				CreateHeader(rawTimestamp, fakeBlockNumber, db)
				formattedTimestamp := FormatTimestamp(rawTimestamp)
				expectedTime := sql.NullString{String: formattedTimestamp, Valid: true}

				ilkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
				_, insertErr := db.Exec(
					`INSERT INTO api.current_ilk_state (ilk_identifier, spot, created, updated) VALUES ($1, $2, $3::TIMESTAMP, $3::TIMESTAMP)`,
					test_helpers.FakeIlk.Identifier, fakeUint256, formattedTimestamp)
				Expect(insertErr).NotTo(HaveOccurred())

				err := repo.Create(fakeBlockNumber, fakeBlockHash, ilkSpotMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var ilkState test_helpers.IlkState
				queryErr := db.Get(&ilkState, `SELECT ilk_identifier, spot, created, updated FROM api.current_ilk_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(ilkState.IlkIdentifier).To(Equal(test_helpers.FakeIlk.Identifier))
				Expect(ilkState.Spot).To(Equal(fakeUint256))
				Expect(ilkState.Created).To(Equal(expectedTime))
				Expect(ilkState.Updated).To(Equal(expectedTime))
			})
		})
	})

	Describe("sin", func() {
		It("writes a row", func() {
			sinMetadata := utils.GetStorageValueMetadata(vat.Sin, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, sinMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT block_number, block_hash, guy AS key, sin AS value FROM maker.vat_sin`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeBlockNumber, fakeBlockHash, fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			sinMetadata := utils.GetStorageValueMetadata(vat.Sin, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, sinMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, sinMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_sin`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing guy", func() {
			malformedSinMetadata := utils.GetStorageValueMetadata(vat.Sin, nil, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedSinMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("urn art", func() {
		It("writes row", func() {
			urnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, urnArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `
				SELECT block_number, block_hash, ilks.id AS key_one, urns.identifier AS key_two, art AS value
				FROM maker.vat_urn_art
				INNER JOIN maker.urns ON maker.urns.id = maker.vat_urn_art.urn_id
				INNER JOIN maker.ilks on maker.urns.ilk_id = maker.ilks.id
			`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			urnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, urnArtMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, urnArtMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_urn_art`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("urn ink", func() {
		It("writes row", func() {
			urnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, urnInkMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `
				SELECT block_number, block_hash, ilks.id AS key_one, urns.identifier AS key_two, ink AS value
				FROM maker.vat_urn_ink
				INNER JOIN maker.urns ON maker.urns.id = maker.vat_urn_ink.urn_id
				INNER JOIN maker.ilks on maker.urns.ilk_id = maker.ilks.id
			`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeBlockNumber, fakeBlockHash, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			urnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, urnInkMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, urnInkMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_urn_ink`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeBlockNumber, fakeBlockHash, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	It("persists vat debt", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.DebtMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, debt AS value FROM maker.vat_debt`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vat debt", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.DebtMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.DebtMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_debt`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat vice", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.ViceMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, vice AS value FROM maker.vat_vice`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vat vice", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.ViceMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.ViceMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_vice`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat Line", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LineMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, line AS value FROM maker.vat_line`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vat Line", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LineMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LineMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_line`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat live", func() {
		err := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LiveMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT block_number, block_hash, live AS value FROM maker.vat_live`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeBlockNumber, fakeBlockHash, fakeUint256)
	})

	It("does not duplicate vat live", func() {
		insertOneErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LiveMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeBlockNumber, fakeBlockHash, vat.LiveMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_live`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})
