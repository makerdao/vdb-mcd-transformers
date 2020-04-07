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

package spot

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/spot"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
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
		contractAddress   = test_data.SpotAddress()
		keccakOfAddress   = types.HexToKeccak256Hash(contractAddress)
		storageKeysLookup = storage.NewKeysLookup(spot.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress))
		repository        = spot.SpotStorageRepository{ContractAddress: contractAddress}
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
		headerID, insertHeaderErr := headerRepository.CreateOrUpdateHeader(header)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		header.Id = headerID
	})

	It("reads in a Spot Vat storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
		value := common.HexToHash("00000000000000000000000057aa8b02f5d3e28371fedcf672c8668869f9aac7")
		spotVatDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		executeErr := transformer.Execute(spotVatDiff)

		Expect(executeErr).NotTo(HaveOccurred())
		var vatResult test_helpers.VariableRes
		getErr := db.Get(&vatResult, `SELECT diff_id, header_id, vat AS value FROM maker.spot_vat`)
		Expect(getErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, spotVatDiff.ID, header.Id, "0x57aA8B02F5D3E28371FEdCf672C8668869f9AAC7")
	})

	It("reads in a Spot Par storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		value := common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000")
		spotParDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		executeErr := transformer.Execute(spotParDiff)

		Expect(executeErr).NotTo(HaveOccurred())
		var parResult test_helpers.VariableRes
		getErr := db.Get(&parResult, `SELECT diff_id, header_id, par AS value FROM maker.spot_par`)
		Expect(getErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(parResult, spotParDiff.ID, header.Id, "1000000000000000000000000000")
	})

	It("reads in a Spot Live storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		spotLiveDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		executeErr := transformer.Execute(spotLiveDiff)

		Expect(executeErr).NotTo(HaveOccurred())
		var liveResult test_helpers.VariableRes
		getErr := db.Get(&liveResult, `SELECT diff_id, header_id, live AS value FROM maker.spot_live`)
		Expect(getErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, spotLiveDiff.ID, header.Id, "1")
	})

	Describe("wards", func() {
		It("reads in a wards storage diff row and persists it", func() {
			denyLog := test_data.CreateTestLog(header.Id, db)
			denyModel := test_data.DenyModel()

			spotAddressID, spotAddressErr := shared.GetOrCreateAddress(test_data.SpotAddress(), db)
			Expect(spotAddressErr).NotTo(HaveOccurred())

			userAddress := "0x2be4b34a34c3cda056dc9e514a30040b6358bf89"
			userAddressID, userAddressErr := shared.GetOrCreateAddress(userAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())

			msgSenderAddress := "0x" + fakes.RandomString(40)
			msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
			Expect(msgSenderAddressErr).NotTo(HaveOccurred())

			denyModel.ColumnValues[event.HeaderFK] = header.Id
			denyModel.ColumnValues[event.LogFK] = denyLog.ID
			denyModel.ColumnValues[event.AddressFK] = spotAddressID
			denyModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
			denyModel.ColumnValues[constants.UsrColumn] = userAddressID
			insertErr := event.PersistModels([]event.InsertionModel{denyModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			key := common.HexToHash("acbda0c7abc278c8fb8df441982ecd46bd66bed192fdc761196288a48630eb70")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
			wardsDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			transformErr := transformer.Execute(wardsDiff)
			Expect(transformErr).NotTo(HaveOccurred())

			var wardsResult test_helpers.MappingResWithAddress
			err := db.Get(&wardsResult, `SELECT diff_id, header_id, address_id, usr AS key, wards.wards AS value FROM maker.wards`)
			Expect(err).NotTo(HaveOccurred())
			Expect(wardsResult.AddressID).To(Equal(strconv.FormatInt(spotAddressID, 10)))
			test_helpers.AssertMapping(wardsResult.MappingRes, wardsDiff.ID, header.Id, strconv.FormatInt(userAddressID, 10), "1")
		})
	})

	Describe("ilk", func() {
		var ilkID int64

		BeforeEach(func() {
			ilk := "0x434f4c352d410000000000000000000000000000000000000000000000000000"
			var ilkErr error
			ilkID, ilkErr = shared.GetOrCreateIlk(ilk, db)
			Expect(ilkErr).NotTo(HaveOccurred())
		})

		It("reads in a Spot Ilk Pip storage diff row and persists it", func() {
			key := common.HexToHash("1730ac98111482efebd8acadb14d7fa301298e0d95bf3c34c3378ef524670bc6")
			value := common.HexToHash("000000000000000000000000a53e6efb4cbed841eace02220498860905e94998")
			spotIlkPipDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			executeErr := transformer.Execute(spotIlkPipDiff)

			Expect(executeErr).NotTo(HaveOccurred())
			var ilkPipResult test_helpers.MappingRes
			getErr := db.Get(&ilkPipResult, `SELECT diff_id, header_id, ilk_id AS key, pip AS value FROM maker.spot_ilk_pip`)
			Expect(getErr).NotTo(HaveOccurred())
			test_helpers.AssertMapping(ilkPipResult, spotIlkPipDiff.ID, header.Id, strconv.FormatInt(ilkID, 10), "0xA53e6EFB4cBeD841Eace02220498860905E94998")
		})

		It("reads in a Spot Ilk Mat storage diff row and persists it", func() {
			key := common.HexToHash("1730ac98111482efebd8acadb14d7fa301298e0d95bf3c34c3378ef524670bc7")
			value := common.HexToHash("000000000000000000000000000000000000000006765c793fa10079d0000000")
			spotIlkMatDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			executeErr := transformer.Execute(spotIlkMatDiff)

			Expect(executeErr).NotTo(HaveOccurred())
			var ilkRhoResult test_helpers.MappingRes
			getErr := db.Get(&ilkRhoResult, `SELECT diff_id, header_id, ilk_id AS key, mat AS value FROM maker.spot_ilk_mat`)
			Expect(getErr).NotTo(HaveOccurred())
			test_helpers.AssertMapping(ilkRhoResult, spotIlkMatDiff.ID, header.Id, strconv.FormatInt(ilkID, 10), "2000000000000000000000000000")
		})
	})
})
