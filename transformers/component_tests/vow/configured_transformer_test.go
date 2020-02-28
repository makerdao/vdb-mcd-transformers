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

package vow

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vow"
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
		vowAddress        = test_data.VowAddress()
		keccakOfAddress   = types.HexToKeccak256Hash(vowAddress)
		storageKeysLookup = storage.NewKeysLookup(vow.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, vowAddress))
		repository        = vow.VowStorageRepository{ContractAddress: vowAddress}
		transformer       = storage.Transformer{
			Address:           common.HexToAddress(vowAddress),
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

	Describe("wards", func() {
		It("reads in a wards storage diff row and persists it", func() {
			denyLog := test_data.CreateTestLog(header.Id, db)
			denyModel := test_data.DenyModel()

			vowAddressID, vowAddressErr := shared.GetOrCreateAddress(test_data.VowAddress(), db)
			Expect(vowAddressErr).NotTo(HaveOccurred())

			userAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
			userAddressID, userAddressErr := shared.GetOrCreateAddress(userAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())

			msgSenderAddress := "0x" + fakes.RandomString(40)
			msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
			Expect(msgSenderAddressErr).NotTo(HaveOccurred())

			denyModel.ColumnValues[event.HeaderFK] = header.Id
			denyModel.ColumnValues[event.LogFK] = denyLog.ID
			denyModel.ColumnValues[event.AddressFK] = vowAddressID
			denyModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
			denyModel.ColumnValues[constants.UsrColumn] = userAddressID
			insertErr := event.PersistModels([]event.InsertionModel{denyModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			key := common.HexToHash("614c9873ec2671d6eb30d7a22b531442a34fc10f8c24a6598ef401fe94d9cb7d")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
			wardsDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			transformErr := transformer.Execute(wardsDiff)
			Expect(transformErr).NotTo(HaveOccurred())

			var wardsResult test_helpers.WardsMappingRes
			err := db.Get(&wardsResult, `SELECT diff_id, header_id, address_id, usr AS key, wards.wards AS value FROM maker.wards`)
			Expect(err).NotTo(HaveOccurred())
			Expect(wardsResult.AddressID).To(Equal(strconv.FormatInt(vowAddressID, 10)))
			test_helpers.AssertMapping(wardsResult.MappingRes, wardsDiff.ID, header.Id, strconv.FormatInt(userAddressID, 10), "1")
		})
	})

	It("reads in a Vow.vat storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		value := common.HexToHash("00000000000000000000000067fd6c3575fc2dbe2cb596bd3bebc9edb5571fa1")
		vowVat := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(vowVat)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT diff_id, header_id, vat AS value FROM maker.vow_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, vowVat.ID, header.Id, "0x67fd6c3575Fc2dBE2CB596bD3bEbc9EDb5571fA1")
	})

	It("reads in a Vow.flapper storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
		value := common.HexToHash("000000000000000000000000b6e31ab6ea62be7c530c32daea96e84d92fe20b7")
		vowFlapper := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(vowFlapper)
		Expect(err).NotTo(HaveOccurred())

		var flapperResult test_helpers.VariableRes
		err = db.Get(&flapperResult, `SELECT diff_id, header_id, flapper AS value FROM maker.vow_flapper`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(flapperResult, vowFlapper.ID, header.Id, "0xB6e31ab6Ea62Be7c530C32DAEa96E84d92fe20B7")
	})

	It("reads in a Vow.flopper storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		value := common.HexToHash("000000000000000000000000275ec1950d6406e3ce6156f9f529c047ea41c8ce")
		vowFlopper := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(vowFlopper)
		Expect(err).NotTo(HaveOccurred())

		var flopperResult test_helpers.VariableRes
		err = db.Get(&flopperResult, `SELECT diff_id, header_id, flopper AS value FROM maker.vow_flopper`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(flopperResult, vowFlopper.ID, header.Id, "0x275eC1950D6406e3cE6156f9F529c047Ea41c8cE")
	})

	It("reads in a Vow.dump storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000008")
		value := common.HexToHash("000000000000000000000000000000000000000000000000002386f26fc10000")
		vowDump := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(vowDump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT diff_id, header_id, dump AS value FROM maker.vow_dump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, vowDump.ID, header.Id, "10000000000000000")
	})

	It("reads in a Vow.sump storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000009")
		value := common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000")
		vowSump := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(vowSump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT diff_id, header_id, sump AS value FROM maker.vow_sump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, vowSump.ID, header.Id, "100000000000000000000000000000000000000000000")
	})

	It("reads in a Vow.bump storage diff row and persists it", func() {
		key := common.HexToHash("000000000000000000000000000000000000000000000000000000000000000a")
		value := common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000")
		vowBump := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(vowBump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT diff_id, header_id, bump AS value FROM maker.vow_bump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, vowBump.ID, header.Id, "100000000000000000000000000000000000000000000")
	})

	It("reads in a Vow.hump storage diff row and persists it", func() {
		// TODO: Update with a real storage diff
		key := common.HexToHash("000000000000000000000000000000000000000000000000000000000000000b")
		value := common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000")
		vowHump := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(vowHump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT diff_id, header_id, hump AS value FROM maker.vow_hump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, vowHump.ID, header.Id, "100000000000000000000000000000000000000000000")
	})
})
