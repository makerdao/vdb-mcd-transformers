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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db                = test_config.NewTestDB(test_config.NewTestNode())
		storageKeysLookup = storage.NewKeysLookup(vat.NewKeysLoader(&mcdStorage.MakerStorageRepository{}))
		repository        = vat.StorageRepository{}
		contractAddress   = test_data.VatAddress()
		keccakOfAddress   = types.HexToKeccak256Hash(contractAddress)
		transformer       = storage.Transformer{
			Address:           common.HexToAddress(contractAddress),
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
		header = fakes.FakeHeader
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		header.Id, insertHeaderErr = headerRepository.CreateOrUpdateHeader(header)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("reads in a Vat wards storage diff row and persists it", func() {
		vatDenyLog := test_data.CreateTestLog(header.Id, db)
		vatDenyModel := test_data.VatDenyModel()

		vatAddressID, vatAddressErr := shared.GetOrCreateAddress(test_data.VatAddress(), db)
		Expect(vatAddressErr).NotTo(HaveOccurred())

		userAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		userAddressID, userAddressErr := shared.GetOrCreateAddress(userAddress, db)
		Expect(userAddressErr).NotTo(HaveOccurred())

		msgSenderAddress := "0x" + fakes.RandomString(40)
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())

		vatDenyModel.ColumnValues[event.HeaderFK] = header.Id
		vatDenyModel.ColumnValues[event.LogFK] = vatDenyLog.ID
		vatDenyModel.ColumnValues[event.AddressFK] = vatAddressID
		vatDenyModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
		vatDenyModel.ColumnValues[constants.UsrColumn] = userAddressID
		insertErr := event.PersistModels([]event.InsertionModel{vatDenyModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		key := common.HexToHash("614c9873ec2671d6eb30d7a22b531442a34fc10f8c24a6598ef401fe94d9cb7d")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		wardsDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		transformErr := transformer.Execute(wardsDiff)
		Expect(transformErr).NotTo(HaveOccurred())

		var wardsResult test_helpers.MappingResWithAddress
		err := db.Get(&wardsResult, `SELECT diff_id, header_id, address_id, usr AS key, wards.wards AS value FROM maker.wards`)
		Expect(err).NotTo(HaveOccurred())
		Expect(wardsResult.AddressID).To(Equal(strconv.FormatInt(vatAddressID, 10)))
		test_helpers.AssertMapping(wardsResult.MappingRes, wardsDiff.ID, header.Id, strconv.FormatInt(userAddressID, 10), "1")
	})

	It("reads in a Vat debt storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000007")
		value := common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000")
		vatDebtDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(vatDebtDiff)
		Expect(err).NotTo(HaveOccurred())

		var debtResult test_helpers.VariableRes
		err = db.Get(&debtResult, `SELECT diff_id, header_id, debt AS value FROM maker.vat_debt`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(debtResult, vatDebtDiff.ID, header.Id, "100000000000000000000000000000000000000000000")
	})

	It("reads in a Vat Line storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000009")
		value := common.HexToHash("0000000000000000000002ac3a4edbbfb8014e3ba83411e915e8000000000000")
		vatLineDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(vatLineDiff)
		Expect(err).NotTo(HaveOccurred())

		var lineResult test_helpers.VariableRes
		err = db.Get(&lineResult, `SELECT diff_id, header_id, line AS value FROM maker.vat_line`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(lineResult, vatLineDiff.ID, header.Id, "1000000000000000000000000000000000000000000000000000")
	})

	It("reads in a Vat live storage diff row and persists it", func() {
		key := common.HexToHash("000000000000000000000000000000000000000000000000000000000000000a")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		vatLiveDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(vatLiveDiff)
		Expect(err).NotTo(HaveOccurred())

		var liveResult test_helpers.VariableRes
		err = db.Get(&liveResult, `SELECT diff_id, header_id, live AS value FROM maker.vat_live`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, vatLiveDiff.ID, header.Id, "1")
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
			key := common.HexToHash("5cd43a2b0a7e767504a508ed07c6f6d26130368a2a5ce573193b4c24eba603bb")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000de0b6b3a7640000")
			ilkArtDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			err := transformer.Execute(ilkArtDiff)
			Expect(err).NotTo(HaveOccurred())

			var artResult test_helpers.MappingRes
			err = db.Get(&artResult, `SELECT diff_id, header_id, ilk_id AS key, art AS value FROM maker.vat_ilk_art`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(artResult, ilkArtDiff.ID, header.Id, strconv.FormatInt(ilkId, 10), "1000000000000000000")
		})

		It("reads in a Vat ilk rate storage diff row and persists it", func() {
			key := common.HexToHash("5cd43a2b0a7e767504a508ed07c6f6d26130368a2a5ce573193b4c24eba603bc")
			value := common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000")
			ilkRateDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			err := transformer.Execute(ilkRateDiff)
			Expect(err).NotTo(HaveOccurred())

			var rateResult test_helpers.MappingRes
			err = db.Get(&rateResult, `SELECT diff_id, header_id, ilk_id AS key, rate AS value FROM maker.vat_ilk_rate`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(rateResult, ilkRateDiff.ID, header.Id, strconv.FormatInt(ilkId, 10), "1000000000000000000000000000")
		})

		It("reads in a Vat ilk spot storage diff row and persists it", func() {
			key := common.HexToHash("5cd43a2b0a7e767504a508ed07c6f6d26130368a2a5ce573193b4c24eba603bd")
			value := common.HexToHash("0000000000000000000000000000000000000001215a061b4dc8dbb48e000000")
			ilkSpotDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			err := transformer.Execute(ilkSpotDiff)
			Expect(err).NotTo(HaveOccurred())

			var spotResult test_helpers.MappingRes
			err = db.Get(&spotResult, `SELECT diff_id, header_id, ilk_id AS key, spot AS value FROM maker.vat_ilk_spot`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(spotResult, ilkSpotDiff.ID, header.Id, strconv.FormatInt(ilkId, 10), "89550000000000000000000000000")
		})

		It("reads in a Vat ilk line storage diff row and persists it", func() {
			key := common.HexToHash("5cd43a2b0a7e767504a508ed07c6f6d26130368a2a5ce573193b4c24eba603be")
			value := common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000")
			ilkLineDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			err := transformer.Execute(ilkLineDiff)
			Expect(err).NotTo(HaveOccurred())

			var lineResult test_helpers.MappingRes
			err = db.Get(&lineResult, `SELECT diff_id, header_id, ilk_id AS key, line AS value FROM maker.vat_ilk_line`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(lineResult, ilkLineDiff.ID, header.Id, strconv.FormatInt(ilkId, 10), "100000000000000000000000000000000000000000000")
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
			key := common.HexToHash("f61b39a22cef8e61a5dc6836ca1a1d267a584ca41782d5b2832fb973dc4731e7")
			value := common.HexToHash("000000000000000000000000000000000000000000000002b5e3af16b1880000")
			urnInkDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			err := transformer.Execute(urnInkDiff)
			Expect(err).NotTo(HaveOccurred())

			var inkResult test_helpers.MappingRes
			err = db.Get(&inkResult, `SELECT diff_id, header_id, urn_id AS key, ink AS value FROM maker.vat_urn_ink`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(inkResult, urnInkDiff.ID, header.Id, strconv.FormatInt(urnID, 10), "50000000000000000000")
		})

		It("reads in a Vat urn art storage diff row and persists it", func() {
			key := common.HexToHash("f61b39a22cef8e61a5dc6836ca1a1d267a584ca41782d5b2832fb973dc4731e8")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000de0b6b3a7640000")
			urnArtDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			err := transformer.Execute(urnArtDiff)
			Expect(err).NotTo(HaveOccurred())

			var artResult test_helpers.MappingRes
			err = db.Get(&artResult, `SELECT diff_id, header_id, urn_id AS key, art AS value FROM maker.vat_urn_art`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(artResult, urnArtDiff.ID, header.Id, strconv.FormatInt(urnID, 10), "1000000000000000000")
		})
	})

	Describe("gem", func() {
		var (
			guy = "0x118D6a283f9044Ce17b95226822e5c73F50e0B90"
			ilk = "0x434f4c312d410000000000000000000000000000000000000000000000000000"
		)

		BeforeEach(func() {
			vatFrobLog := test_data.CreateTestLog(header.Id, db)
			vatFrob := test_data.VatFrobModelWithPositiveDart()
			urnID, urnErr := shared.GetOrCreateUrn(guy, ilk, db)
			Expect(urnErr).NotTo(HaveOccurred())
			vatFrob.ColumnValues[constants.UrnColumn] = urnID
			vatFrob.ColumnValues[constants.VColumn] = guy
			vatFrob.ColumnValues[event.HeaderFK] = header.Id
			vatFrob.ColumnValues[event.LogFK] = vatFrobLog.ID
			insertErr := event.PersistModels([]event.InsertionModel{vatFrob}, db)
			Expect(insertErr).NotTo(HaveOccurred())
		})

		It("reads in a Vat gem storage diff row and persists it", func() {
			key := common.HexToHash("f4cd303dafe86407ce2225d1c13f4a50accf4db0b5187d7cf268acff9841cfa4")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
			vatGemDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			err := transformer.Execute(vatGemDiff)
			Expect(err).NotTo(HaveOccurred())

			var gemResult test_helpers.DoubleMappingRes
			err = db.Get(&gemResult, `SELECT diff_id, header_id, ilk_id AS key_one, guy AS key_two, gem AS value FROM maker.vat_gem`)
			Expect(err).NotTo(HaveOccurred())
			ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
			Expect(ilkErr).NotTo(HaveOccurred())
			test_helpers.AssertDoubleMapping(gemResult, vatGemDiff.ID, header.Id, strconv.FormatInt(ilkID, 10), guy, "0")
		})
	})

	Describe("dai", func() {
		var guy = "0x118D6a283f9044Ce17b95226822e5c73F50e0B90"

		BeforeEach(func() {
			vatFrobLog := test_data.CreateTestLog(header.Id, db)
			vatFrob := test_data.VatFrobModelWithPositiveDart()
			ilk := "0x434f4c312d410000000000000000000000000000000000000000000000000000"
			urnID, urnErr := shared.GetOrCreateUrn(guy, ilk, db)
			Expect(urnErr).NotTo(HaveOccurred())
			vatFrob.ColumnValues[constants.UrnColumn] = urnID
			vatFrob.ColumnValues[constants.WColumn] = guy
			vatFrob.ColumnValues[event.HeaderFK] = header.Id
			vatFrob.ColumnValues[event.LogFK] = vatFrobLog.ID
			insertErr := event.PersistModels([]event.InsertionModel{vatFrob}, db)
			Expect(insertErr).NotTo(HaveOccurred())
		})

		It("reads in a Vat dai storage diff row and persists it", func() {
			key := common.HexToHash("de69de809d681089810c52d4e65d3489177732d304a5eda082bcf61b015d3918")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
			vatDaiDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			err := transformer.Execute(vatDaiDiff)
			Expect(err).NotTo(HaveOccurred())

			var daiResult test_helpers.MappingRes
			err = db.Get(&daiResult, `SELECT diff_id, header_id, guy AS key, dai AS value FROM maker.vat_dai`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(daiResult, vatDaiDiff.ID, header.Id, guy, "0")
		})
	})
})
