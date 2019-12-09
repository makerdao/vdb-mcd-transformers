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
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vat storage repository", func() {
	var (
		db           = test_config.NewTestDB(test_config.NewTestNode())
		repo         vat.VatStorageRepository
		fakeGuy      = "fake_urn"
		fakeUint256  = "12345"
		fakeHeaderID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = vat.VatStorageRepository{}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	Describe("dai", func() {
		It("writes a row", func() {
			daiMetadata := utils.GetStorageValueMetadata(vat.Dai, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeHeaderID, daiMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT header_id, guy AS key, dai AS value FROM maker.vat_dai`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeHeaderID, fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			daiMetadata := utils.GetStorageValueMetadata(vat.Dai, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeHeaderID, daiMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeHeaderID, daiMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_dai`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing guy", func() {
			malformedDaiMetadata := utils.GetStorageValueMetadata(vat.Dai, nil, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedDaiMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("gem", func() {
		It("writes row", func() {
			gemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeHeaderID, gemMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `SELECT header_id, ilk_id AS key_one, guy AS key_two, gem AS value FROM maker.vat_gem`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			gemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeHeaderID, gemMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeHeaderID, gemMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_gem`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedGemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedGemMetadata := utils.GetStorageValueMetadata(vat.Gem, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("ilk Art", func() {
		It("writes row", func() {
			ilkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeHeaderID, ilkArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT header_id, ilk_id AS key, art AS value FROM maker.vat_ilk_art`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeHeaderID, ilkArtMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeHeaderID, ilkArtMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_art`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkArtMetadata := utils.GetStorageValueMetadata(vat.IlkArt, nil, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedIlkArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      utils.GetStorageValueMetadata(vat.IlkArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
			PropertyName:  "Art",
			PropertyValue: strconv.Itoa(rand.Int()),
			TableName:     "maker.vat_ilk_art",
		})
	})

	Describe("ilk dust", func() {
		It("writes row", func() {
			ilkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeHeaderID, ilkDustMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT header_id, ilk_id AS key, dust AS value FROM maker.vat_ilk_dust`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeHeaderID, ilkDustMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeHeaderID, ilkDustMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_dust`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkDustMetadata := utils.GetStorageValueMetadata(vat.IlkDust, nil, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedIlkDustMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      utils.GetStorageValueMetadata(vat.IlkDust, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
			PropertyName:  "Dust",
			PropertyValue: strconv.Itoa(rand.Int()),
			TableName:     "maker.vat_ilk_dust",
		})
	})

	Describe("ilk line", func() {
		It("writes row", func() {
			ilkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeHeaderID, ilkLineMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT header_id, ilk_id AS key, line AS value FROM maker.vat_ilk_line`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeHeaderID, ilkLineMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeHeaderID, ilkLineMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_line`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkLineMetadata := utils.GetStorageValueMetadata(vat.IlkLine, nil, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedIlkLineMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      utils.GetStorageValueMetadata(vat.IlkLine, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
			PropertyName:  "Line",
			PropertyValue: strconv.Itoa(rand.Int()),
			TableName:     "maker.vat_ilk_line",
		})
	})

	Describe("ilk rate", func() {
		It("writes row", func() {
			ilkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeHeaderID, ilkRateMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT header_id, ilk_id AS key, rate AS value FROM maker.vat_ilk_rate`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeHeaderID, ilkRateMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeHeaderID, ilkRateMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_rate`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkRateMetadata := utils.GetStorageValueMetadata(vat.IlkRate, nil, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedIlkRateMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      utils.GetStorageValueMetadata(vat.IlkRate, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
			PropertyName:  "Rate",
			PropertyValue: strconv.Itoa(rand.Int()),
			TableName:     "maker.vat_ilk_rate",
		})
	})

	Describe("ilk spot", func() {
		It("writes row", func() {
			ilkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeHeaderID, ilkSpotMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT header_id, ilk_id AS key, spot AS value FROM maker.vat_ilk_spot`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)
			insertOneErr := repo.Create(fakeHeaderID, ilkSpotMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeHeaderID, ilkSpotMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_ilk_spot`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkSpotMetadata := utils.GetStorageValueMetadata(vat.IlkSpot, nil, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedIlkSpotMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      utils.GetStorageValueMetadata(vat.IlkSpot, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256),
			PropertyName:  "Spot",
			PropertyValue: strconv.Itoa(rand.Int()),
			TableName:     "maker.vat_ilk_spot",
		})
	})

	Describe("sin", func() {
		It("writes a row", func() {
			sinMetadata := utils.GetStorageValueMetadata(vat.Sin, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeHeaderID, sinMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			err = db.Get(&result, `SELECT header_id, guy AS key, sin AS value FROM maker.vat_sin`)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, fakeHeaderID, fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			sinMetadata := utils.GetStorageValueMetadata(vat.Sin, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeHeaderID, sinMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeHeaderID, sinMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_sin`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing guy", func() {
			malformedSinMetadata := utils.GetStorageValueMetadata(vat.Sin, nil, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedSinMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("urn art", func() {
		It("writes row", func() {
			urnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeHeaderID, urnArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `
				SELECT header_id, ilks.id AS key_one, urns.identifier AS key_two, art AS value
				FROM maker.vat_urn_art
				INNER JOIN maker.urns ON maker.urns.id = maker.vat_urn_art.urn_id
				INNER JOIN maker.ilks on maker.urns.ilk_id = maker.ilks.id
			`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			urnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeHeaderID, urnArtMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeHeaderID, urnArtMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_urn_art`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})

		Describe("updating historical_ilk_state trigger table", func() {
			var (
				blockOne,
				blockTwo,
				blockThree int
				rawTimestampOne, rawTimestampTwo, rawTimestampThree int64
				headerOne,
				headerTwo core.Header
				hashOne        = common.BytesToHash([]byte{1, 2, 3, 4, 5})
				hashTwo        = common.BytesToHash([]byte{5, 4, 3, 2, 1})
				hashThree      = common.BytesToHash([]byte{9, 8, 7, 6, 5})
				getStateQuery  = `SELECT urn_identifier, ilk_identifier, block_height, ink, art, created, updated FROM api.historical_urn_state ORDER BY block_height`
				getArtQuery    = `SELECT art FROM api.historical_urn_state ORDER BY block_height`
				insertArtQuery = `INSERT INTO api.historical_urn_state (urn_identifier, ilk_identifier, block_height, art, updated) VALUES ($1, $2, $3, $4, NOW())`
				deleteRowQuery = `DELETE FROM maker.vat_urn_art WHERE header_id = $1`
				urnArtMetadata = utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
			)

			BeforeEach(func() {
				blockOne = rand.Int()
				blockTwo = blockOne + 1
				blockThree = blockTwo + 1
				rawTimestampOne = int64(rand.Int31())
				rawTimestampTwo = rawTimestampOne + 1
				rawTimestampThree = rawTimestampTwo + 1
				headerOne = CreateHeaderWithHash(hashOne.String(), rawTimestampOne, blockOne, db)
				headerTwo = CreateHeaderWithHash(hashTwo.String(), rawTimestampTwo, blockTwo, db)
				CreateHeaderWithHash(hashThree.String(), rawTimestampThree, blockThree, db)
			})

			It("inserts time of first ink diff into created", func() {
				urnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
				setupErrOne := repo.Create(headerOne.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(setupErrOne).NotTo(HaveOccurred())
				expectedTimeCreated := test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))

				err := repo.Create(headerTwo.Id, urnArtMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var timeCreated sql.NullString
				queryErr := db.Get(&timeCreated,
					`SELECT created FROM api.historical_urn_state WHERE block_height = $1`, headerTwo.BlockNumber)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(timeCreated).To(Equal(expectedTimeCreated))
			})

			It("inserts a row for new urn-block", func() {
				initialUrnValues := test_helpers.GetUrnSetupData()
				newArt := rand.Int()
				test_helpers.CreateUrn(initialUrnValues, headerOne.Id, test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy), repo)

				err := repo.Create(headerTwo.Id, urnArtMetadata, strconv.Itoa(newArt))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getStateQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].BlockHeight).To(Equal(int(headerTwo.BlockNumber)))
				Expect(urnStates[1].Ink).To(Equal(strconv.Itoa(initialUrnValues[vat.UrnInk])))
				Expect(urnStates[1].Art).To(Equal(strconv.Itoa(newArt)))
				Expect(urnStates[1].Created).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))))
				Expect(urnStates[1].Updated).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampTwo))))
			})

			It("updates art if urn-block combination already exists in table", func() {
				initialUrnValues := test_helpers.GetUrnSetupData()
				newArt := rand.Int()
				test_helpers.CreateUrn(initialUrnValues, headerOne.Id, test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy), repo)

				err := repo.Create(headerOne.Id, urnArtMetadata, strconv.Itoa(newArt))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getStateQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(1))
				Expect(urnStates[0].BlockHeight).To(Equal(int(headerOne.BlockNumber)))
				Expect(urnStates[0].Ink).To(Equal(strconv.Itoa(initialUrnValues[vat.UrnInk])))
				Expect(urnStates[0].Art).To(Equal(strconv.Itoa(newArt)))
				Expect(urnStates[0].Created).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))))
				Expect(urnStates[0].Updated).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))))
			})

			It("updates all consecutive arts if record is the latest art diff", func() {
				initialArt := rand.Int()
				newArt := initialArt + 1
				_, setupErr := db.Exec(insertArtQuery,
					fakeGuy, test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, initialArt)
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(headerOne.Id, urnArtMetadata, strconv.Itoa(newArt))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getArtQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].Art).To(Equal(strconv.Itoa(newArt)))
			})

			It("ignores rows from blocks after the next art diff", func() {
				initialArt := strconv.Itoa(rand.Int())
				setupErr := repo.Create(headerTwo.Id, urnArtMetadata, initialArt)
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(headerOne.Id, urnArtMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getArtQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].Art).To(Equal(initialArt))
			})

			It("ignores rows from different urn", func() {
				initialArt := rand.Int()
				_, setupErr := db.Exec(insertArtQuery,
					fakeGuy, test_helpers.AnotherFakeIlk.Identifier, headerTwo.BlockNumber, initialArt)
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(headerOne.Id, urnArtMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getArtQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].Art).To(Equal(strconv.Itoa(initialArt)))
			})

			It("ignores rows from earlier blocks", func() {
				initialArt := rand.Int()
				_, setupErr := db.Exec(insertArtQuery,
					fakeGuy, test_helpers.FakeIlk.Identifier, headerOne.BlockNumber, initialArt)
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(headerTwo.Id, urnArtMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getArtQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[0].Art).To(Equal(strconv.Itoa(initialArt)))
			})

			Describe("when diff is deleted", func() {
				It("updates art to previous value until block number of next diff", func() {
					initialArt := rand.Int()
					setupErrOne := repo.Create(headerOne.Id, urnArtMetadata, strconv.Itoa(initialArt))
					Expect(setupErrOne).NotTo(HaveOccurred())

					subsequentArt := initialArt + 1
					setupErrTwo := repo.Create(headerTwo.Id, urnArtMetadata, strconv.Itoa(subsequentArt))
					Expect(setupErrTwo).NotTo(HaveOccurred())
					_, setupErrThree := db.Exec(insertArtQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockThree, subsequentArt)
					Expect(setupErrThree).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteRowQuery, headerTwo.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var urnStates []test_helpers.UrnState
					queryErr := db.Select(&urnStates, getArtQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(urnStates[1].Art).To(Equal(strconv.Itoa(initialArt)))
				})

				It("sets field in subsequent rows to null if no previous diff exists", func() {
					initialArt := strconv.Itoa(rand.Int())
					setupErrOne := repo.Create(headerOne.Id, urnArtMetadata, initialArt)
					Expect(setupErrOne).NotTo(HaveOccurred())
					_, setupErrTwo := db.Exec(insertArtQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockTwo, initialArt)
					Expect(setupErrTwo).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteRowQuery, headerOne.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var actualArts []sql.NullString
					queryErr := db.Select(&actualArts, getArtQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(actualArts)).To(Equal(1))
					Expect(actualArts[0].Valid).To(BeFalse())
				})

				It("deletes ilk state associated with diff if identical to previous state", func() {
					initialArt := rand.Int()
					setupErrOne := repo.Create(headerOne.Id, urnArtMetadata, strconv.Itoa(initialArt))
					Expect(setupErrOne).NotTo(HaveOccurred())

					subsequentArt := initialArt + 1
					setupErrTwo := repo.Create(headerTwo.Id, urnArtMetadata, strconv.Itoa(subsequentArt))
					Expect(setupErrTwo).NotTo(HaveOccurred())
					_, setupErrThree := db.Exec(insertArtQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockThree, subsequentArt)
					Expect(setupErrThree).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteRowQuery, headerTwo.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var urnStates []test_helpers.UrnState
					queryErr := db.Select(&urnStates, getArtQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(urnStates)).To(Equal(2))
				})

				It("deletes ilk state associated with diff if it's the earliest state in the table", func() {
					initialArt := rand.Int()
					setupErrOne := repo.Create(headerOne.Id, urnArtMetadata, strconv.Itoa(initialArt))
					Expect(setupErrOne).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteRowQuery, headerOne.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var urnStates []test_helpers.UrnState
					queryErr := db.Select(&urnStates, getArtQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(urnStates)).To(Equal(0))
				})
			})
		})
	})

	Describe("urn ink", func() {
		It("writes row", func() {
			urnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeHeaderID, urnInkMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			err = db.Get(&result, `
				SELECT header_id, ilks.id AS key_one, urns.identifier AS key_two, ink AS value
				FROM maker.vat_urn_ink
				INNER JOIN maker.urns ON maker.urns.id = maker.vat_urn_ink.urn_id
				INNER JOIN maker.ilks on maker.urns.ilk_id = maker.ilks.id
			`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			urnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
			insertOneErr := repo.Create(fakeHeaderID, urnInkMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(fakeHeaderID, urnInkMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_urn_ink`)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Guy: fakeGuy}, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnInkMetadata := utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, utils.Uint256)

			err := repo.Create(fakeHeaderID, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrMetadataMalformed{MissingData: constants.Guy}))
		})

		Describe("updating historical_ilk_state trigger table", func() {
			var (
				blockOne,
				blockTwo,
				blockThree int
				rawTimestampOne, rawTimestampTwo, rawTimestampThree int64
				headerOne,
				headerTwo core.Header
				hashOne        = common.BytesToHash([]byte{1, 2, 3, 4, 5})
				hashTwo        = common.BytesToHash([]byte{5, 4, 3, 2, 1})
				hashThree      = common.BytesToHash([]byte{9, 8, 7, 6, 5})
				getStateQuery  = `SELECT urn_identifier, ilk_identifier, block_height, ink, art, created, updated FROM api.historical_urn_state ORDER BY block_height`
				getInkQuery    = `SELECT ink FROM api.historical_urn_state ORDER BY block_height`
				insertInkQuery = `INSERT INTO api.historical_urn_state (urn_identifier, ilk_identifier, block_height, ink, updated) VALUES ($1, $2, $3, $4, NOW())`
				deleteRowQuery = `DELETE FROM maker.vat_urn_ink WHERE header_id = $1`
				urnInkMetadata = utils.GetStorageValueMetadata(vat.UrnInk, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
			)

			BeforeEach(func() {
				blockOne = rand.Int()
				blockTwo = blockOne + 1
				blockThree = blockTwo + 1
				rawTimestampOne = int64(rand.Int31())
				rawTimestampTwo = rawTimestampOne + 1
				headerOne = CreateHeaderWithHash(hashOne.String(), rawTimestampOne, blockOne, db)
				headerTwo = CreateHeaderWithHash(hashTwo.String(), rawTimestampTwo, blockTwo, db)
				CreateHeaderWithHash(hashThree.String(), rawTimestampThree, blockThree, db)
			})

			It("inserts time of first ink diff into created", func() {
				setupErrOne := repo.Create(headerOne.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(setupErrOne).NotTo(HaveOccurred())
				setupErrTwo := repo.Create(headerTwo.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(setupErrTwo).NotTo(HaveOccurred())
				expectedTimeCreated := test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))

				var timeCreatedValues []sql.NullString
				queryErr := db.Select(&timeCreatedValues, `SELECT created FROM api.historical_urn_state`)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(len(timeCreatedValues)).To(Equal(2))
				Expect(timeCreatedValues[0]).To(Equal(expectedTimeCreated))
				Expect(timeCreatedValues[1]).To(Equal(expectedTimeCreated))
			})

			It("updates time created for all the urn's states when new ink is added", func() {
				urnArtMetadata := utils.GetStorageValueMetadata(vat.UrnArt, map[utils.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, utils.Uint256)
				setupErr := repo.Create(headerOne.Id, urnArtMetadata, strconv.Itoa(rand.Int()))
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(headerTwo.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())
				expectedTimeCreated := test_helpers.GetValidNullString(FormatTimestamp(rawTimestampTwo))

				var timeCreatedValues []sql.NullString
				queryErr := db.Select(&timeCreatedValues, `SELECT created FROM api.historical_urn_state`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(timeCreatedValues)).To(Equal(2))
				Expect(timeCreatedValues[0]).To(Equal(expectedTimeCreated))
				Expect(timeCreatedValues[1]).To(Equal(expectedTimeCreated))
			})

			It("inserts a row for new urn-block", func() {
				initialUrnValues := test_helpers.GetUrnSetupData()
				newInk := rand.Int()
				test_helpers.CreateUrn(initialUrnValues, headerOne.Id, test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy), repo)

				err := repo.Create(headerTwo.Id, urnInkMetadata, strconv.Itoa(newInk))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getStateQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].BlockHeight).To(Equal(int(headerTwo.BlockNumber)))
				Expect(urnStates[1].Ink).To(Equal(strconv.Itoa(newInk)))
				Expect(urnStates[1].Art).To(Equal(strconv.Itoa(initialUrnValues[vat.UrnArt])))
				Expect(urnStates[1].Created).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))))
				Expect(urnStates[1].Updated).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampTwo))))
			})

			It("updates ink if urn-block combination already exists in table", func() {
				initialUrnValues := test_helpers.GetUrnSetupData()
				newInk := rand.Int()
				test_helpers.CreateUrn(initialUrnValues, headerOne.Id, test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy), repo)

				err := repo.Create(headerOne.Id, urnInkMetadata, strconv.Itoa(newInk))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getStateQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(1))
				Expect(urnStates[0].BlockHeight).To(Equal(int(headerOne.BlockNumber)))
				Expect(urnStates[0].Ink).To(Equal(strconv.Itoa(newInk)))
				Expect(urnStates[0].Art).To(Equal(strconv.Itoa(initialUrnValues[vat.UrnArt])))
				Expect(urnStates[0].Created).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))))
				Expect(urnStates[0].Updated).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))))
			})

			It("updates all consecutive inks if record is the latest ink diff", func() {
				initialInk := rand.Int()
				newInk := initialInk + 1
				_, setupErr := db.Exec(insertInkQuery,
					fakeGuy, test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, initialInk)
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(headerOne.Id, urnInkMetadata, strconv.Itoa(newInk))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getInkQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].Ink).To(Equal(strconv.Itoa(newInk)))
			})

			It("ignores rows from blocks after the next ink diff", func() {
				initialInk := strconv.Itoa(rand.Int())
				setupErr := repo.Create(headerTwo.Id, urnInkMetadata, initialInk)
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(headerOne.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getInkQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].Ink).To(Equal(initialInk))
			})

			It("ignores rows from different urn", func() {
				initialInk := rand.Int()
				_, setupErr := db.Exec(insertInkQuery,
					fakeGuy, test_helpers.AnotherFakeIlk.Identifier, headerTwo.BlockNumber, initialInk)
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(headerOne.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getInkQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].Ink).To(Equal(strconv.Itoa(initialInk)))
			})

			It("ignores rows from earlier blocks", func() {
				initialInk := rand.Int()
				_, setupErr := db.Exec(insertInkQuery,
					fakeGuy, test_helpers.FakeIlk.Identifier, headerOne.BlockNumber, initialInk)
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(headerTwo.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getInkQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[0].Ink).To(Equal(strconv.Itoa(initialInk)))
			})

			Describe("when diff is deleted", func() {
				It("updates ink to previous value until block number of next diff", func() {
					initialInk := rand.Int()
					setupErrOne := repo.Create(headerOne.Id, urnInkMetadata, strconv.Itoa(initialInk))
					Expect(setupErrOne).NotTo(HaveOccurred())

					subsequentInk := initialInk + 1
					setupErrTwo := repo.Create(headerTwo.Id, urnInkMetadata, strconv.Itoa(subsequentInk))
					Expect(setupErrTwo).NotTo(HaveOccurred())
					_, setupErrThree := db.Exec(insertInkQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockThree, subsequentInk)
					Expect(setupErrThree).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteRowQuery, headerTwo.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var urnStates []test_helpers.UrnState
					queryErr := db.Select(&urnStates, getInkQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(urnStates[1].Ink).To(Equal(strconv.Itoa(initialInk)))
				})

				It("sets field in subsequent rows to null if no previous diff exists", func() {
					initialInk := strconv.Itoa(rand.Int())
					setupErrOne := repo.Create(headerOne.Id, urnInkMetadata, initialInk)
					Expect(setupErrOne).NotTo(HaveOccurred())
					_, setupErrTwo := db.Exec(insertInkQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockTwo, initialInk)
					Expect(setupErrTwo).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteRowQuery, headerOne.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var actualInks []sql.NullString
					queryErr := db.Select(&actualInks, getInkQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(actualInks)).To(Equal(1))
					Expect(actualInks[0].Valid).To(BeFalse())
				})

				It("deletes ilk state associated with diff if identical to previous state", func() {
					initialInk := rand.Int()
					setupErrOne := repo.Create(headerOne.Id, urnInkMetadata, strconv.Itoa(initialInk))
					Expect(setupErrOne).NotTo(HaveOccurred())

					subsequentInk := initialInk + 1
					setupErrTwo := repo.Create(headerTwo.Id, urnInkMetadata, strconv.Itoa(subsequentInk))
					Expect(setupErrTwo).NotTo(HaveOccurred())
					_, setupErrThree := db.Exec(insertInkQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockThree, subsequentInk)
					Expect(setupErrThree).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteRowQuery, headerTwo.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var urnStates []test_helpers.UrnState
					queryErr := db.Select(&urnStates, getInkQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(urnStates)).To(Equal(2))
				})

				It("deletes ilk state associated with diff if it's the earliest state in the table", func() {
					initialInk := rand.Int()
					setupErrOne := repo.Create(headerOne.Id, urnInkMetadata, strconv.Itoa(initialInk))
					Expect(setupErrOne).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteRowQuery, headerOne.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var urnStates []test_helpers.UrnState
					queryErr := db.Select(&urnStates, getInkQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(urnStates)).To(Equal(0))
				})

				It("updates time created for all urn's states", func() {
					initialInk := rand.Int()
					setupErrOne := repo.Create(headerOne.Id, urnInkMetadata, strconv.Itoa(initialInk))
					Expect(setupErrOne).NotTo(HaveOccurred())

					subsequentInk := initialInk + 1
					setupErrTwo := repo.Create(headerTwo.Id, urnInkMetadata, strconv.Itoa(subsequentInk))
					Expect(setupErrTwo).NotTo(HaveOccurred())
					expectedTimeCreated := test_helpers.GetValidNullString(FormatTimestamp(rawTimestampTwo))
					_, setupErrThree := db.Exec(insertInkQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockThree, subsequentInk)
					Expect(setupErrThree).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteRowQuery, headerOne.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var actualCreatedValues []sql.NullString
					queryErr := db.Select(&actualCreatedValues, `SELECT created FROM api.historical_urn_state`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(actualCreatedValues)).To(Equal(2))
					Expect(actualCreatedValues[0]).To(Equal(expectedTimeCreated))
					Expect(actualCreatedValues[1]).To(Equal(expectedTimeCreated))
				})
			})
		})
	})

	It("persists vat debt", func() {
		err := repo.Create(fakeHeaderID, vat.DebtMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT header_id, debt AS value FROM maker.vat_debt`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vat debt", func() {
		insertOneErr := repo.Create(fakeHeaderID, vat.DebtMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeHeaderID, vat.DebtMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_debt`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat vice", func() {
		err := repo.Create(fakeHeaderID, vat.ViceMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT header_id, vice AS value FROM maker.vat_vice`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vat vice", func() {
		insertOneErr := repo.Create(fakeHeaderID, vat.ViceMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeHeaderID, vat.ViceMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_vice`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat Line", func() {
		err := repo.Create(fakeHeaderID, vat.LineMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT header_id, line AS value FROM maker.vat_line`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vat Line", func() {
		insertOneErr := repo.Create(fakeHeaderID, vat.LineMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeHeaderID, vat.LineMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_line`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat live", func() {
		err := repo.Create(fakeHeaderID, vat.LiveMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		err = db.Get(&result, `SELECT header_id, live AS value FROM maker.vat_live`)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vat live", func() {
		insertOneErr := repo.Create(fakeHeaderID, vat.LiveMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(fakeHeaderID, vat.LiveMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		getCountErr := db.Get(&count, `SELECT count(*) FROM maker.vat_live`)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})
