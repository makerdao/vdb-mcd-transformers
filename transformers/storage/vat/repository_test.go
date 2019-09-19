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
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	. "github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
)

var _ = Describe("Vat storage repository", func() {
	var (
		db              *postgres.DB
		repo            vat.VatStorageRepository
		fakeBlockNumber = 123
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
