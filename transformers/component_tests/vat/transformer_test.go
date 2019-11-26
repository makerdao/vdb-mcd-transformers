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

package vat

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_frob"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db                *postgres.DB
		storageKeysLookup = storage.NewKeysLookup(vat.NewKeysLoader(&mcdStorage.MakerStorageRepository{}))
		repository        = vat.VatStorageRepository{}
		contractAddress   = "48f749bd988caafacd7b951abbecc1aa31488690"
		transformer       = storage.Transformer{
			HashedAddress:     utils.HexToKeccak256Hash(contractAddress),
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
		headerID int64
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("reads in a Vat debt storage diff row and persists it", func() {
		vatDebtRow := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000007"),
			StorageValue:  common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(vatDebtRow)
		Expect(err).NotTo(HaveOccurred())

		var debtResult test_helpers.VariableRes
		err = db.Get(&debtResult, `SELECT header_id, debt AS value FROM maker.vat_debt`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(debtResult, headerID, "100000000000000000000000000000000000000000000")
	})

	It("reads in a Vat Line storage diff row and persists it", func() {
		vatLineRow := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000009"),
			StorageValue:  common.HexToHash("0000000000000000000002ac3a4edbbfb8014e3ba83411e915e8000000000000"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(vatLineRow)
		Expect(err).NotTo(HaveOccurred())

		var lineResult test_helpers.VariableRes
		err = db.Get(&lineResult, `SELECT header_id, line AS value FROM maker.vat_line`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(lineResult, headerID, "1000000000000000000000000000000000000000000000000000")
	})

	It("reads in a Vat live storage diff row and persists it", func() {
		vatLiveRow := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("000000000000000000000000000000000000000000000000000000000000000a"),
			StorageValue:  common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(vatLiveRow)
		Expect(err).NotTo(HaveOccurred())

		var liveResult test_helpers.VariableRes
		err = db.Get(&liveResult, `SELECT header_id, live AS value FROM maker.vat_live`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, headerID, "1")
	})

	Describe("ilk", func() {
		var (
			ilkId  int64
			ilkErr error
		)

		BeforeEach(func() {
			ilk := "0x4554482d41000000000000000000000000000000000000000000000000000000"
			ilkId, ilkErr = shared.GetOrCreateIlk(ilk, db)
			Expect(ilkErr).NotTo(HaveOccurred())
		})

		It("reads in a Vat ilk Art storage diff row and persists it", func() {
			ilkArtRow := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				StorageKey:    common.HexToHash("5cd43a2b0a7e767504a508ed07c6f6d26130368a2a5ce573193b4c24eba603bb"),
				StorageValue:  common.HexToHash("0000000000000000000000000000000000000000000000000de0b6b3a7640000"),
				HeaderID:      headerID,
			}
			err := transformer.Execute(ilkArtRow)
			Expect(err).NotTo(HaveOccurred())

			var artResult test_helpers.MappingRes
			err = db.Get(&artResult, `SELECT header_id, ilk_id AS key, art AS value FROM maker.vat_ilk_art`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(artResult, headerID, strconv.FormatInt(ilkId, 10), "1000000000000000000")
		})

		It("reads in a Vat ilk rate storage diff row and persists it", func() {
			ilkRateRow := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				StorageKey:    common.HexToHash("5cd43a2b0a7e767504a508ed07c6f6d26130368a2a5ce573193b4c24eba603bc"),
				StorageValue:  common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000"),
				HeaderID:      headerID,
			}
			err := transformer.Execute(ilkRateRow)
			Expect(err).NotTo(HaveOccurred())

			var rateResult test_helpers.MappingRes
			err = db.Get(&rateResult, `SELECT header_id, ilk_id AS key, rate AS value FROM maker.vat_ilk_rate`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(rateResult, headerID, strconv.FormatInt(ilkId, 10), "1000000000000000000000000000")
		})

		It("reads in a Vat ilk spot storage diff row and persists it", func() {
			ilkSpotRow := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				StorageKey:    common.HexToHash("5cd43a2b0a7e767504a508ed07c6f6d26130368a2a5ce573193b4c24eba603bd"),
				StorageValue:  common.HexToHash("0000000000000000000000000000000000000001215a061b4dc8dbb48e000000"),
				HeaderID:      headerID,
			}
			err := transformer.Execute(ilkSpotRow)
			Expect(err).NotTo(HaveOccurred())

			var spotResult test_helpers.MappingRes
			err = db.Get(&spotResult, `SELECT header_id, ilk_id AS key, spot AS value FROM maker.vat_ilk_spot`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(spotResult, headerID, strconv.FormatInt(ilkId, 10), "89550000000000000000000000000")
		})

		It("reads in a Vat ilk spot storage diff with a hashed storage key", func() {
			anotherIlk := "0x474e542d41000000000000000000000000000000000000000000000000000000"
			anotherIlkID, err := shared.GetOrCreateIlk(anotherIlk, db)
			Expect(err).NotTo(HaveOccurred())

			ilkSpotRow := utils.StorageDiff{
				HashedAddress: utils.HexToKeccak256Hash("0x26A5C505c5B8558834483d1322B5305F61b0160D"),
				StorageKey:    common.HexToHash("2165edb4e1c37b99b60fa510d84f939dd35d5cd1d1c8f299d6456ea09df65a76"),
				StorageValue:  common.HexToHash("00000000000000000000000000000000000000008b1bb2b1a88f91522d765555"),
				HeaderID:      headerID,
			}
			err = transformer.Execute(ilkSpotRow)
			Expect(err).NotTo(HaveOccurred())

			var spotResult test_helpers.MappingRes
			err = db.Get(&spotResult, `SELECT header_id, ilk_id AS key, spot AS value FROM maker.vat_ilk_spot`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(spotResult, headerID, strconv.FormatInt(anotherIlkID, 10), "43051901220750297886077900117")
		})

		It("reads in a Vat ilk line storage diff row and persists it", func() {
			ilkLineRow := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				StorageKey:    common.HexToHash("5cd43a2b0a7e767504a508ed07c6f6d26130368a2a5ce573193b4c24eba603be"),
				StorageValue:  common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
				HeaderID:      headerID,
			}
			err := transformer.Execute(ilkLineRow)
			Expect(err).NotTo(HaveOccurred())

			var lineResult test_helpers.MappingRes
			err = db.Get(&lineResult, `SELECT header_id, ilk_id AS key, line AS value FROM maker.vat_ilk_line`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(lineResult, headerID, strconv.FormatInt(ilkId, 10), "100000000000000000000000000000000000000000000")
		})
	})

	Describe("urn", func() {
		var (
			urnID  int64
			urnErr error
		)

		BeforeEach(func() {
			ilk := "0x434f4c312d410000000000000000000000000000000000000000000000000000"
			urn := "0x118D6a283f9044Ce17b95226822e5c73F50e0B90"
			urnID, urnErr = shared.GetOrCreateUrn(urn, ilk, db)
			Expect(urnErr).NotTo(HaveOccurred())
		})

		It("reads in a Vat urn ink storage diff row and persists it", func() {
			urnInkRow := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				StorageKey:    common.HexToHash("f61b39a22cef8e61a5dc6836ca1a1d267a584ca41782d5b2832fb973dc4731e7"),
				StorageValue:  common.HexToHash("000000000000000000000000000000000000000000000002b5e3af16b1880000"),
				HeaderID:      headerID,
			}
			err := transformer.Execute(urnInkRow)
			Expect(err).NotTo(HaveOccurred())

			var inkResult test_helpers.MappingRes
			err = db.Get(&inkResult, `SELECT header_id, urn_id AS key, ink AS value FROM maker.vat_urn_ink`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(inkResult, headerID, strconv.FormatInt(urnID, 10), "50000000000000000000")
		})

		It("reads in a Vat urn art storage diff row and persists it", func() {
			urnInkRow := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				StorageKey:    common.HexToHash("f61b39a22cef8e61a5dc6836ca1a1d267a584ca41782d5b2832fb973dc4731e8"),
				StorageValue:  common.HexToHash("0000000000000000000000000000000000000000000000000de0b6b3a7640000"),
				HeaderID:      headerID,
			}
			err := transformer.Execute(urnInkRow)
			Expect(err).NotTo(HaveOccurred())

			var artResult test_helpers.MappingRes
			err = db.Get(&artResult, `SELECT header_id, urn_id AS key, art AS value FROM maker.vat_urn_art`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(artResult, headerID, strconv.FormatInt(urnID, 10), "1000000000000000000")
		})
	})

	Describe("gem", func() {
		var (
			guy = "0x118D6a283f9044Ce17b95226822e5c73F50e0B90"
			ilk = "0x434f4c312d410000000000000000000000000000000000000000000000000000"
		)

		BeforeEach(func() {
			vatFrobLog := test_data.CreateTestLog(headerID, db)
			vatFrobRepository := vat_frob.VatFrobRepository{}
			vatFrobRepository.SetDB(db)
			vatFrob := test_data.VatFrobModelWithPositiveDart()
			vatFrob.ForeignKeyValues[constants.IlkFK] = ilk
			vatFrob.ColumnValues["v"] = guy
			vatFrob.ColumnValues[constants.HeaderFK] = headerID
			vatFrob.ColumnValues[constants.LogFK] = vatFrobLog.ID
			insertErr := vatFrobRepository.Create([]shared.InsertionModel{vatFrob})
			Expect(insertErr).NotTo(HaveOccurred())
		})

		It("reads in a Vat gem storage diff row and persists it", func() {
			urnInkRow := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				StorageKey:    common.HexToHash("f4cd303dafe86407ce2225d1c13f4a50accf4db0b5187d7cf268acff9841cfa4"),
				StorageValue:  common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000"),
				HeaderID:      headerID,
			}
			err := transformer.Execute(urnInkRow)
			Expect(err).NotTo(HaveOccurred())

			var gemResult test_helpers.DoubleMappingRes
			err = db.Get(&gemResult, `SELECT header_id, ilk_id AS key_one, guy AS key_two, gem AS value FROM maker.vat_gem`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
			Expect(ilkErr).NotTo(HaveOccurred())
			test_helpers.AssertDoubleMapping(gemResult, headerID, strconv.FormatInt(ilkID, 10), guy, "0")
		})
	})

	Describe("dai", func() {
		var guy = "0x118D6a283f9044Ce17b95226822e5c73F50e0B90"

		BeforeEach(func() {
			vatFrobLog := test_data.CreateTestLog(headerID, db)
			vatFrobRepository := vat_frob.VatFrobRepository{}
			vatFrobRepository.SetDB(db)
			vatFrob := test_data.VatFrobModelWithPositiveDart()
			vatFrob.ColumnValues["w"] = guy
			vatFrob.ColumnValues[constants.HeaderFK] = headerID
			vatFrob.ColumnValues[constants.LogFK] = vatFrobLog.ID
			insertErr := vatFrobRepository.Create([]shared.InsertionModel{vatFrob})
			Expect(insertErr).NotTo(HaveOccurred())
		})

		It("reads in a Vat dai storage diff row and persists it", func() {
			urnInkRow := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				StorageKey:    common.HexToHash("de69de809d681089810c52d4e65d3489177732d304a5eda082bcf61b015d3918"),
				StorageValue:  common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000"),
				HeaderID:      headerID,
			}
			err := transformer.Execute(urnInkRow)
			Expect(err).NotTo(HaveOccurred())

			var daiResult test_helpers.MappingRes
			err = db.Get(&daiResult, `SELECT header_id, guy AS key, dai AS value FROM maker.vat_dai`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(daiResult, headerID, guy, "0")
		})
	})
})
