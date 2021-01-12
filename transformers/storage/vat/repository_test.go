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
	"fmt"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	mcdShared "github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data/shared_behaviors"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vat storage repository", func() {
	var (
		db                   = test_config.NewTestDB(test_config.NewTestNode())
		repo                 vat.StorageRepository
		fakeGuy              = "fake_urn"
		fakeUint256          = "12345"
		diffID, fakeHeaderID int64
		deleteHeaderQuery    = `DELETE from public.headers WHERE id =$1`
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		repo = vat.StorageRepository{}
		repo.SetDB(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		fakeHeaderID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		diffID = CreateFakeDiffRecord(db)
	})

	Describe("Wards mapping", func() {
		It("writes a row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)

			setupErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(setupErr).NotTo(HaveOccurred())

			var result MappingResWithAddress
			query := fmt.Sprintf(`SELECT diff_id, header_id, address_id, usr AS key, wards AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			err := db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, test_data.VatAddress())
			Expect(contractAddressErr).NotTo(HaveOccurred())
			userAddressID, userAddressErr := repository.GetOrCreateAddress(db, fakeUserAddress)
			Expect(userAddressErr).NotTo(HaveOccurred())
			AssertMappingWithAddress(result, diffID, fakeHeaderID, contractAddressID, strconv.FormatInt(userAddressID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			fakeUserAddress := "0x" + fakes.RandomString(40)
			wardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{constants.User: fakeUserAddress}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, wardsMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.WardsTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns an error if metadata missing user", func() {
			malformedWardsMetadata := types.GetValueMetadata(wards.Wards, map[types.Key]string{}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedWardsMetadata, fakeUint256)
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.User}))
		})
	})

	Describe("dai", func() {
		It("writes a row", func() {
			daiMetadata := types.GetValueMetadata(vat.Dai, map[types.Key]string{constants.Guy: fakeGuy}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, daiMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, guy AS KEY, dai AS VALUE FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatDaiTable))
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, diffID, fakeHeaderID, fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			daiMetadata := types.GetValueMetadata(vat.Dai, map[types.Key]string{constants.Guy: fakeGuy}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, daiMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, daiMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatDaiTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing guy", func() {
			malformedDaiMetadata := types.GetValueMetadata(vat.Dai, nil, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedDaiMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("gem", func() {
		It("writes row", func() {
			gemMetadata := types.GetValueMetadata(vat.Gem, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, gemMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS key_one, guy AS key_two, gem AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatGemTable))
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := mcdShared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			gemMetadata := types.GetValueMetadata(vat.Gem, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, gemMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, gemMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatGemTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedGemMetadata := types.GetValueMetadata(vat.Gem, map[types.Key]string{constants.Guy: fakeGuy}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedGemMetadata := types.GetValueMetadata(vat.Gem, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedGemMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("ilk Art", func() {
		It("writes row", func() {
			ilkArtMetadata := types.GetValueMetadata(vat.IlkArt, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, ilkArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS key, art AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatIlkArtTable))
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := mcdShared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkArtMetadata := types.GetValueMetadata(vat.IlkArt, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, ilkArtMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkArtMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatIlkArtTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkArtMetadata := types.GetValueMetadata(vat.IlkArt, nil, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedIlkArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      types.GetValueMetadata(vat.IlkArt, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256),
			PropertyName:  "Art",
			PropertyValue: strconv.Itoa(rand.Int()),
			Schema:        constants.MakerSchema,
			TableName:     constants.VatIlkArtTable,
		})
	})

	Describe("ilk dust", func() {
		It("writes row", func() {
			ilkDustMetadata := types.GetValueMetadata(vat.IlkDust, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, ilkDustMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS key, dust AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatIlkDustTable))
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := mcdShared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkDustMetadata := types.GetValueMetadata(vat.IlkDust, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, ilkDustMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkDustMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatIlkDustTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkDustMetadata := types.GetValueMetadata(vat.IlkDust, nil, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedIlkDustMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      types.GetValueMetadata(vat.IlkDust, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256),
			PropertyName:  "Dust",
			PropertyValue: strconv.Itoa(rand.Int()),
			Schema:        constants.MakerSchema,
			TableName:     constants.VatIlkDustTable,
		})
	})

	Describe("ilk line", func() {
		It("writes row", func() {
			ilkLineMetadata := types.GetValueMetadata(vat.IlkLine, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, ilkLineMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS key, line AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatIlkLineTable))
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := mcdShared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkLineMetadata := types.GetValueMetadata(vat.IlkLine, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, ilkLineMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkLineMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatIlkLineTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkLineMetadata := types.GetValueMetadata(vat.IlkLine, nil, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedIlkLineMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      types.GetValueMetadata(vat.IlkLine, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256),
			PropertyName:  "Line",
			PropertyValue: strconv.Itoa(rand.Int()),
			Schema:        constants.MakerSchema,
			TableName:     constants.VatIlkLineTable,
		})
	})

	Describe("ilk rate", func() {
		It("writes row", func() {
			ilkRateMetadata := types.GetValueMetadata(vat.IlkRate, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, ilkRateMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS key, rate AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatIlkRateTable))
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := mcdShared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkRateMetadata := types.GetValueMetadata(vat.IlkRate, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, ilkRateMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkRateMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatIlkRateTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkRateMetadata := types.GetValueMetadata(vat.IlkRate, nil, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedIlkRateMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      types.GetValueMetadata(vat.IlkRate, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256),
			PropertyName:  "Rate",
			PropertyValue: strconv.Itoa(rand.Int()),
			Schema:        constants.MakerSchema,
			TableName:     constants.VatIlkRateTable,
		})
	})

	Describe("ilk spot", func() {
		It("writes row", func() {
			ilkSpotMetadata := types.GetValueMetadata(vat.IlkSpot, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, ilkSpotMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, ilk_id AS key, spot AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatIlkSpotTable))
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := mcdShared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeUint256)
		})

		It("does not duplicate row", func() {
			ilkSpotMetadata := types.GetValueMetadata(vat.IlkSpot, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, ilkSpotMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, ilkSpotMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatIlkSpotTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedIlkSpotMetadata := types.GetValueMetadata(vat.IlkSpot, nil, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedIlkSpotMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		shared_behaviors.SharedIlkTriggerTests(shared_behaviors.IlkTriggerTestInput{
			Repository:    &repo,
			Metadata:      types.GetValueMetadata(vat.IlkSpot, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256),
			PropertyName:  "Spot",
			PropertyValue: strconv.Itoa(rand.Int()),
			Schema:        constants.MakerSchema,
			TableName:     constants.VatIlkSpotTable,
		})
	})

	Describe("sin", func() {
		It("writes a row", func() {
			sinMetadata := types.GetValueMetadata(vat.Sin, map[types.Key]string{constants.Guy: fakeGuy}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, sinMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result MappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, guy AS KEY, sin AS VALUE FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatSinTable))
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			AssertMapping(result, diffID, fakeHeaderID, fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			sinMetadata := types.GetValueMetadata(vat.Sin, map[types.Key]string{constants.Guy: fakeGuy}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, sinMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, sinMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatSinTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing guy", func() {
			malformedSinMetadata := types.GetValueMetadata(vat.Sin, nil, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedSinMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Guy}))
		})
	})

	Describe("urn art", func() {
		It("writes row", func() {
			urnArtMetadata := types.GetValueMetadata(vat.UrnArt, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, urnArtMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, ilks.id AS key_one, urns.identifier AS key_two, art AS value
				FROM %s
				INNER JOIN maker.urns ON maker.urns.id = maker.vat_urn_art.urn_id
				INNER JOIN maker.ilks on maker.urns.ilk_id = maker.ilks.id`,
				shared.GetFullTableName(constants.MakerSchema, constants.VatUrnArtTable))
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := mcdShared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			urnArtMetadata := types.GetValueMetadata(vat.UrnArt, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, urnArtMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, urnArtMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatUrnArtTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnArtMetadata := types.GetValueMetadata(vat.UrnArt, map[types.Key]string{constants.Guy: fakeGuy}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnArtMetadata := types.GetValueMetadata(vat.UrnArt, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedUrnArtMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Guy}))
		})

		Describe("updating urn_snapshot trigger table", func() {
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
				getStateQuery  = `SELECT urn_identifier, ilk_identifier, block_height, ink, art, created, updated FROM api.urn_snapshot ORDER BY block_height`
				getArtQuery    = `SELECT art FROM api.urn_snapshot ORDER BY block_height`
				insertArtQuery = `INSERT INTO api.urn_snapshot (urn_identifier, ilk_identifier, block_height, art, updated) VALUES ($1, $2, $3, $4, NOW())`
				urnArtMetadata = types.GetValueMetadata(vat.UrnArt, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, types.Uint256)
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
				urnInkMetadata := types.GetValueMetadata(vat.UrnInk, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, types.Uint256)
				setupErrOne := repo.Create(diffID, headerOne.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(setupErrOne).NotTo(HaveOccurred())
				expectedTimeCreated := test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))

				err := repo.Create(diffID, headerTwo.Id, urnArtMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var timeCreated sql.NullString
				queryErr := db.Get(&timeCreated,
					`SELECT created FROM api.urn_snapshot WHERE block_height = $1`, headerTwo.BlockNumber)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(timeCreated).To(Equal(expectedTimeCreated))
			})

			It("inserts a row for new urn-block", func() {
				initialUrnValues := test_helpers.GetUrnSetupData()
				newArt := rand.Int()
				test_helpers.CreateUrn(db, initialUrnValues, headerOne, test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy), repo)

				diffIdTwo := CreateFakeDiffRecordWithHeader(db, headerTwo)
				err := repo.Create(diffIdTwo, headerTwo.Id, urnArtMetadata, strconv.Itoa(newArt))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getStateQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].BlockHeight).To(Equal(int(headerTwo.BlockNumber)))
				Expect(urnStates[1].Ink).To(Equal(strconv.Itoa(initialUrnValues[vat.UrnInk].(int))))
				Expect(urnStates[1].Art).To(Equal(strconv.Itoa(newArt)))
				Expect(urnStates[1].Created).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))))
				Expect(urnStates[1].Updated).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampTwo))))
			})

			It("updates art if urn-block combination already exists in table", func() {
				initialUrnValues := test_helpers.GetUrnSetupData()
				newArt := rand.Int()
				test_helpers.CreateUrn(db, initialUrnValues, headerOne, test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy), repo)

				err := repo.Create(diffID, headerOne.Id, urnArtMetadata, strconv.Itoa(newArt))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getStateQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(1))
				Expect(urnStates[0].BlockHeight).To(Equal(int(headerOne.BlockNumber)))
				Expect(urnStates[0].Ink).To(Equal(strconv.Itoa(initialUrnValues[vat.UrnInk].(int))))
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

				diffIdOne := CreateFakeDiffRecordWithHeader(db, headerTwo)
				err := repo.Create(diffIdOne, headerOne.Id, urnArtMetadata, strconv.Itoa(newArt))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getArtQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].Art).To(Equal(strconv.Itoa(newArt)))
			})

			It("ignores rows from blocks after the next art diff", func() {
				initialArt := strconv.Itoa(rand.Int())
				setupErr := repo.Create(diffID, headerTwo.Id, urnArtMetadata, initialArt)
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(diffID, headerOne.Id, urnArtMetadata, strconv.Itoa(rand.Int()))
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

				err := repo.Create(diffID, headerOne.Id, urnArtMetadata, strconv.Itoa(rand.Int()))
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

				diffIdTwo := CreateFakeDiffRecordWithHeader(db, headerTwo)
				err := repo.Create(diffIdTwo, headerTwo.Id, urnArtMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getArtQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[0].Art).To(Equal(strconv.Itoa(initialArt)))
			})

			Describe("when diff is deleted", func() {
				var diffIdOne, diffIdTwo int64

				BeforeEach(func() {
					diffIdOne = CreateFakeDiffRecordWithHeader(db, headerOne)
					diffIdTwo = CreateFakeDiffRecordWithHeader(db, headerTwo)
				})

				It("updates art to previous value until block number of next diff", func() {
					initialArt := rand.Int()
					setupErrOne := repo.Create(diffIdOne, headerOne.Id, urnArtMetadata, strconv.Itoa(initialArt))
					Expect(setupErrOne).NotTo(HaveOccurred())

					subsequentArt := initialArt + 1
					setupErrTwo := repo.Create(diffIdTwo, headerTwo.Id, urnArtMetadata, strconv.Itoa(subsequentArt))
					Expect(setupErrTwo).NotTo(HaveOccurred())

					_, setupErrThree := db.Exec(insertArtQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockThree, subsequentArt)
					Expect(setupErrThree).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteHeaderQuery, headerTwo.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var urnStates []test_helpers.UrnState
					queryErr := db.Select(&urnStates, getArtQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(urnStates)).To(Equal(2))
					Expect(urnStates[1].Art).To(Equal(strconv.Itoa(initialArt)))
				})

				It("sets field in subsequent rows to null if no previous diff exists", func() {
					initialArt := strconv.Itoa(rand.Int())
					setupErrOne := repo.Create(diffIdOne, headerOne.Id, urnArtMetadata, initialArt)
					Expect(setupErrOne).NotTo(HaveOccurred())

					_, setupErrTwo := db.Exec(insertArtQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockTwo, initialArt)
					Expect(setupErrTwo).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteHeaderQuery, headerOne.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var actualArts []sql.NullString
					queryErr := db.Select(&actualArts, getArtQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(actualArts)).To(Equal(1))
					Expect(actualArts[0].Valid).To(BeFalse())
				})

				It("deletes urn state associated with diff if identical to previous state", func() {
					initialArt := rand.Int()
					setupErrOne := repo.Create(diffIdOne, headerOne.Id, urnArtMetadata, strconv.Itoa(initialArt))
					Expect(setupErrOne).NotTo(HaveOccurred())

					subsequentArt := initialArt + 1
					setupErrTwo := repo.Create(diffIdTwo, headerTwo.Id, urnArtMetadata, strconv.Itoa(subsequentArt))
					Expect(setupErrTwo).NotTo(HaveOccurred())

					_, setupErrThree := db.Exec(insertArtQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockThree, subsequentArt)
					Expect(setupErrThree).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteHeaderQuery, headerTwo.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var urnStates []test_helpers.UrnState
					queryErr := db.Select(&urnStates, getArtQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(urnStates)).To(Equal(2))
				})

				It("deletes urn state associated with diff if it's the earliest state in the table", func() {
					initialArt := rand.Int()
					setupErrOne := repo.Create(diffIdOne, headerOne.Id, urnArtMetadata, strconv.Itoa(initialArt))
					Expect(setupErrOne).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteHeaderQuery, headerOne.Id)
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
			urnInkMetadata := types.GetValueMetadata(vat.UrnInk, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, urnInkMetadata, fakeUint256)

			Expect(err).NotTo(HaveOccurred())

			var result DoubleMappingRes
			query := fmt.Sprintf(`SELECT diff_id, header_id, ilks.id AS key_one, urns.identifier AS key_two, ink AS value
				FROM %s
				INNER JOIN maker.urns ON maker.urns.id = maker.vat_urn_ink.urn_id
				INNER JOIN maker.ilks on maker.urns.ilk_id = maker.ilks.id`, shared.GetFullTableName(constants.MakerSchema, constants.VatUrnInkTable))
			err = db.Get(&result, query)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := mcdShared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
			Expect(err).NotTo(HaveOccurred())
			AssertDoubleMapping(result, diffID, fakeHeaderID, strconv.FormatInt(ilkID, 10), fakeGuy, fakeUint256)
		})

		It("does not duplicate row", func() {
			urnInkMetadata := types.GetValueMetadata(vat.UrnInk, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, types.Uint256)
			insertOneErr := repo.Create(diffID, fakeHeaderID, urnInkMetadata, fakeUint256)
			Expect(insertOneErr).NotTo(HaveOccurred())

			insertTwoErr := repo.Create(diffID, fakeHeaderID, urnInkMetadata, fakeUint256)

			Expect(insertTwoErr).NotTo(HaveOccurred())
			var count int
			query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatUrnInkTable))
			getCountErr := db.Get(&count, query)
			Expect(getCountErr).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("returns error if metadata missing ilk", func() {
			malformedUrnInkMetadata := types.GetValueMetadata(vat.UrnInk, map[types.Key]string{constants.Guy: fakeGuy}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Ilk}))
		})

		It("returns error if metadata missing guy", func() {
			malformedUrnInkMetadata := types.GetValueMetadata(vat.UrnInk, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex}, types.Uint256)

			err := repo.Create(diffID, fakeHeaderID, malformedUrnInkMetadata, fakeUint256)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(types.ErrMetadataMalformed{MissingData: constants.Guy}))
		})

		Describe("updating urn_snapshot trigger table", func() {
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
				getStateQuery  = `SELECT urn_identifier, ilk_identifier, block_height, ink, art, created, updated FROM api.urn_snapshot ORDER BY block_height`
				getInkQuery    = `SELECT ink FROM api.urn_snapshot ORDER BY block_height`
				insertInkQuery = `INSERT INTO api.urn_snapshot (urn_identifier, ilk_identifier, block_height, ink, updated) VALUES ($1, $2, $3, $4, NOW())`
				urnInkMetadata = types.GetValueMetadata(vat.UrnInk, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, types.Uint256)
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
				setupErrOne := repo.Create(diffID, headerOne.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(setupErrOne).NotTo(HaveOccurred())
				setupErrTwo := repo.Create(diffID, headerTwo.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(setupErrTwo).NotTo(HaveOccurred())
				expectedTimeCreated := test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))

				var timeCreatedValues []sql.NullString
				queryErr := db.Select(&timeCreatedValues, `SELECT created FROM api.urn_snapshot`)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(len(timeCreatedValues)).To(Equal(2))
				Expect(timeCreatedValues[0]).To(Equal(expectedTimeCreated))
				Expect(timeCreatedValues[1]).To(Equal(expectedTimeCreated))
			})

			It("updates time created for all the urn's states when new ink is added", func() {
				urnArtMetadata := types.GetValueMetadata(vat.UrnArt, map[types.Key]string{constants.Ilk: test_helpers.FakeIlk.Hex, constants.Guy: fakeGuy}, types.Uint256)
				setupErr := repo.Create(diffID, headerOne.Id, urnArtMetadata, strconv.Itoa(rand.Int()))
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(diffID, headerTwo.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())
				expectedTimeCreated := test_helpers.GetValidNullString(FormatTimestamp(rawTimestampTwo))

				var timeCreatedValues []sql.NullString
				queryErr := db.Select(&timeCreatedValues, `SELECT created FROM api.urn_snapshot`)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(timeCreatedValues)).To(Equal(2))
				Expect(timeCreatedValues[0]).To(Equal(expectedTimeCreated))
				Expect(timeCreatedValues[1]).To(Equal(expectedTimeCreated))
			})

			It("inserts a row for new urn-block", func() {
				initialUrnValues := test_helpers.GetUrnSetupData()
				newInk := rand.Int()
				test_helpers.CreateUrn(db, initialUrnValues, headerOne, test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy), repo)

				diffIdTwo := CreateFakeDiffRecordWithHeader(db, headerTwo)
				err := repo.Create(diffIdTwo, headerTwo.Id, urnInkMetadata, strconv.Itoa(newInk))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getStateQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].BlockHeight).To(Equal(int(headerTwo.BlockNumber)))
				Expect(urnStates[1].Ink).To(Equal(strconv.Itoa(newInk)))
				Expect(urnStates[1].Art).To(Equal(strconv.Itoa(initialUrnValues[vat.UrnArt].(int))))
				Expect(urnStates[1].Created).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))))
				Expect(urnStates[1].Updated).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampTwo))))
			})

			It("updates ink if urn-block combination already exists in table", func() {
				initialUrnValues := test_helpers.GetUrnSetupData()
				newInk := rand.Int()
				test_helpers.CreateUrn(db, initialUrnValues, headerOne, test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy), repo)

				err := repo.Create(diffID, headerOne.Id, urnInkMetadata, strconv.Itoa(newInk))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getStateQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(1))
				Expect(urnStates[0].BlockHeight).To(Equal(int(headerOne.BlockNumber)))
				Expect(urnStates[0].Ink).To(Equal(strconv.Itoa(newInk)))
				Expect(urnStates[0].Art).To(Equal(strconv.Itoa(initialUrnValues[vat.UrnArt].(int))))
				Expect(urnStates[0].Created).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))))
				Expect(urnStates[0].Updated).To(Equal(test_helpers.GetValidNullString(FormatTimestamp(rawTimestampOne))))
			})

			It("updates all consecutive inks if record is the latest ink diff", func() {
				initialInk := rand.Int()
				newInk := initialInk + 1
				_, setupErr := db.Exec(insertInkQuery,
					fakeGuy, test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber, initialInk)
				Expect(setupErr).NotTo(HaveOccurred())

				diffIdOne := CreateFakeDiffRecordWithHeader(db, headerOne)
				err := repo.Create(diffIdOne, headerOne.Id, urnInkMetadata, strconv.Itoa(newInk))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getInkQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[1].Ink).To(Equal(strconv.Itoa(newInk)))
			})

			It("ignores rows from blocks after the next ink diff", func() {
				initialInk := strconv.Itoa(rand.Int())
				setupErr := repo.Create(diffID, headerTwo.Id, urnInkMetadata, initialInk)
				Expect(setupErr).NotTo(HaveOccurred())

				err := repo.Create(diffID, headerOne.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
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

				err := repo.Create(diffID, headerOne.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
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

				diffIdTwo := CreateFakeDiffRecordWithHeader(db, headerTwo)
				err := repo.Create(diffIdTwo, headerTwo.Id, urnInkMetadata, strconv.Itoa(rand.Int()))
				Expect(err).NotTo(HaveOccurred())

				var urnStates []test_helpers.UrnState
				queryErr := db.Select(&urnStates, getInkQuery)
				Expect(queryErr).NotTo(HaveOccurred())
				Expect(len(urnStates)).To(Equal(2))
				Expect(urnStates[0].Ink).To(Equal(strconv.Itoa(initialInk)))
			})

			Describe("when diff is deleted", func() {
				var diffIdOne, diffIdTwo int64

				BeforeEach(func() {
					diffIdOne = CreateFakeDiffRecordWithHeader(db, headerOne)
					diffIdTwo = CreateFakeDiffRecordWithHeader(db, headerTwo)
				})

				It("updates ink to previous value until block number of next diff", func() {
					initialInk := rand.Int()
					setupErrOne := repo.Create(diffIdOne, headerOne.Id, urnInkMetadata, strconv.Itoa(initialInk))
					Expect(setupErrOne).NotTo(HaveOccurred())

					subsequentInk := initialInk + 1
					setupErrTwo := repo.Create(diffIdTwo, headerTwo.Id, urnInkMetadata, strconv.Itoa(subsequentInk))
					Expect(setupErrTwo).NotTo(HaveOccurred())
					_, setupErrThree := db.Exec(insertInkQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockThree, subsequentInk)
					Expect(setupErrThree).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteHeaderQuery, headerTwo.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var urnStates []test_helpers.UrnState
					queryErr := db.Select(&urnStates, getInkQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(urnStates[1].Ink).To(Equal(strconv.Itoa(initialInk)))
				})

				It("sets field in subsequent rows to null if no previous diff exists", func() {
					initialInk := strconv.Itoa(rand.Int())
					setupErrOne := repo.Create(diffIdOne, headerOne.Id, urnInkMetadata, initialInk)
					Expect(setupErrOne).NotTo(HaveOccurred())
					_, setupErrTwo := db.Exec(insertInkQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockTwo, initialInk)
					Expect(setupErrTwo).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteHeaderQuery, headerOne.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var actualInks []sql.NullString
					queryErr := db.Select(&actualInks, getInkQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(actualInks)).To(Equal(1))
					Expect(actualInks[0].Valid).To(BeFalse())
				})

				It("deletes urn state associated with diff if identical to previous state", func() {
					initialInk := rand.Int()
					setupErrOne := repo.Create(diffIdOne, headerOne.Id, urnInkMetadata, strconv.Itoa(initialInk))
					Expect(setupErrOne).NotTo(HaveOccurred())

					subsequentInk := initialInk + 1
					setupErrTwo := repo.Create(diffIdTwo, headerTwo.Id, urnInkMetadata, strconv.Itoa(subsequentInk))
					Expect(setupErrTwo).NotTo(HaveOccurred())
					_, setupErrThree := db.Exec(insertInkQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockThree, subsequentInk)
					Expect(setupErrThree).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteHeaderQuery, headerTwo.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var urnStates []test_helpers.UrnState
					queryErr := db.Select(&urnStates, getInkQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(urnStates)).To(Equal(2))
				})

				It("deletes urn state associated with diff if it's the earliest state in the table", func() {
					initialInk := rand.Int()
					setupErrOne := repo.Create(diffIdOne, headerOne.Id, urnInkMetadata, strconv.Itoa(initialInk))
					Expect(setupErrOne).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteHeaderQuery, headerOne.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var urnStates []test_helpers.UrnState
					queryErr := db.Select(&urnStates, getInkQuery)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(urnStates)).To(Equal(0))
				})

				It("updates time created for all urn's states", func() {
					initialInk := rand.Int()
					setupErrOne := repo.Create(diffIdOne, headerOne.Id, urnInkMetadata, strconv.Itoa(initialInk))
					Expect(setupErrOne).NotTo(HaveOccurred())

					subsequentInk := initialInk + 1
					setupErrTwo := repo.Create(diffIdTwo, headerTwo.Id, urnInkMetadata, strconv.Itoa(subsequentInk))
					Expect(setupErrTwo).NotTo(HaveOccurred())
					expectedTimeCreated := test_helpers.GetValidNullString(FormatTimestamp(rawTimestampTwo))
					_, setupErrThree := db.Exec(insertInkQuery,
						fakeGuy, test_helpers.FakeIlk.Identifier, blockThree, subsequentInk)
					Expect(setupErrThree).NotTo(HaveOccurred())

					_, deleteErr := db.Exec(deleteHeaderQuery, headerOne.Id)
					Expect(deleteErr).NotTo(HaveOccurred())

					var actualCreatedValues []sql.NullString
					queryErr := db.Select(&actualCreatedValues, `SELECT created FROM api.urn_snapshot`)
					Expect(queryErr).NotTo(HaveOccurred())
					Expect(len(actualCreatedValues)).To(Equal(2))
					Expect(actualCreatedValues[0]).To(Equal(expectedTimeCreated))
					Expect(actualCreatedValues[1]).To(Equal(expectedTimeCreated))
				})
			})
		})
	})

	It("persists vat debt", func() {
		err := repo.Create(diffID, fakeHeaderID, vat.DebtMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, debt AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatDebtTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vat debt", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vat.DebtMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vat.DebtMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatDebtTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat vice", func() {
		err := repo.Create(diffID, fakeHeaderID, vat.ViceMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, vice AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatViceTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vat vice", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vat.ViceMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vat.ViceMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatViceTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat Line", func() {
		err := repo.Create(diffID, fakeHeaderID, vat.LineMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, line AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatLineTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vat Line", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vat.LineMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vat.LineMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatLineTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("persists vat live", func() {
		err := repo.Create(diffID, fakeHeaderID, vat.LiveMetadata, fakeUint256)

		Expect(err).NotTo(HaveOccurred())

		var result VariableRes
		query := fmt.Sprintf(`SELECT diff_id, header_id, live AS value FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatLiveTable))
		err = db.Get(&result, query)
		Expect(err).NotTo(HaveOccurred())
		AssertVariable(result, diffID, fakeHeaderID, fakeUint256)
	})

	It("does not duplicate vat live", func() {
		insertOneErr := repo.Create(diffID, fakeHeaderID, vat.LiveMetadata, fakeUint256)
		Expect(insertOneErr).NotTo(HaveOccurred())

		insertTwoErr := repo.Create(diffID, fakeHeaderID, vat.LiveMetadata, fakeUint256)

		Expect(insertTwoErr).NotTo(HaveOccurred())
		var count int
		query := fmt.Sprintf(`SELECT count(*) FROM %s`, shared.GetFullTableName(constants.MakerSchema, constants.VatLiveTable))
		getCountErr := db.Get(&count, query)
		Expect(getCountErr).NotTo(HaveOccurred())
		Expect(count).To(Equal(1))
	})
})
