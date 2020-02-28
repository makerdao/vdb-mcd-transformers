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

package jug

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
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
		contractAddress   = test_data.JugAddress()
		keccakOfAddress   = types.HexToKeccak256Hash(contractAddress)
		storageKeysLookup = storage.NewKeysLookup(jug.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress))
		repository        = jug.JugStorageRepository{ContractAddress: contractAddress}
		transformer       = storage.Transformer{
			Address:           common.HexToAddress(contractAddress),
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
		header = fakes.FakeHeader
		ilkID  int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
		ilk := "0x4554480000000000000000000000000000000000000000000000000000000000"
		var ilkErr error
		ilkID, ilkErr = shared.GetOrCreateIlk(ilk, db)
		Expect(ilkErr).NotTo(HaveOccurred())
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		header.Id, insertHeaderErr = headerRepository.CreateOrUpdateHeader(header)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("reads in a Jug Vat storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
		value := common.HexToHash("00000000000000000000000067fd6c3575fc2dbe2cb596bd3bebc9edb5571fa1")
		jugVatDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(jugVatDiff)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT diff_id, header_id, vat AS value FROM maker.jug_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, jugVatDiff.ID, header.Id, "0x67fd6c3575Fc2dBE2CB596bD3bEbc9EDb5571fA1")
	})

	It("reads in a Jug Vow storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		value := common.HexToHash("17560834075da3db54f737db74377e799c865821000000000000000000000000")
		jugVowDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(jugVowDiff)
		Expect(err).NotTo(HaveOccurred())

		var vowResult test_helpers.VariableRes
		err = db.Get(&vowResult, `SELECT diff_id, header_id, vow AS value FROM maker.jug_vow`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vowResult, jugVowDiff.ID, header.Id, "0x17560834075da3db54f737db74377e799c865821000000000000000000000000")
	})

	It("reads in a wards storage diff row and persists it", func() {
		denyLog := test_data.CreateTestLog(header.Id, db)
		denyModel := test_data.DenyModel()

		jugAddressID, jugAddressErr := shared.GetOrCreateAddress(contractAddress, db)
		Expect(jugAddressErr).NotTo(HaveOccurred())

		userAddress := "0x13141b8a5e4a82ebc6b636849dd6a515185d6236"
		userAddressID, userAddressErr := shared.GetOrCreateAddress(userAddress, db)
		Expect(userAddressErr).NotTo(HaveOccurred())

		msgSenderAddress := "0x" + fakes.RandomString(40)
		msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())

		denyModel.ColumnValues[event.HeaderFK] = header.Id
		denyModel.ColumnValues[event.LogFK] = denyLog.ID
		denyModel.ColumnValues[event.AddressFK] = jugAddressID
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
		Expect(wardsResult.AddressID).To(Equal(strconv.FormatInt(jugAddressID, 10)))
		test_helpers.AssertMapping(wardsResult.MappingRes, wardsDiff.ID, header.Id, strconv.FormatInt(userAddressID, 10), "1")
	})

	It("reads in a Jug Ilk Duty storage diff row and persists it", func() {
		key := common.HexToHash("a27f5adbce3dcb790941ebd020e02078a61e6c9748376e52ec0929d58babf97a")
		value := common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000")
		jugIlkDutyDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(jugIlkDutyDiff)
		Expect(err).NotTo(HaveOccurred())

		var ilkDutyResult test_helpers.MappingRes
		err = db.Get(&ilkDutyResult, `SELECT diff_id, header_id, ilk_id AS key, duty AS value FROM maker.jug_ilk_duty`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(ilkDutyResult, jugIlkDutyDiff.ID, header.Id, strconv.FormatInt(ilkID, 10), "1000000000000000000000000000")
	})

	It("reads in a Jug Ilk Rho storage diff row and persists it", func() {
		key := common.HexToHash("a27f5adbce3dcb790941ebd020e02078a61e6c9748376e52ec0929d58babf97b")
		value := common.HexToHash("000000000000000000000000000000000000000000000000000000005c812808")
		jugIlkRhoDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(jugIlkRhoDiff)
		Expect(err).NotTo(HaveOccurred())

		var ilkRhoResult test_helpers.MappingRes
		err = db.Get(&ilkRhoResult, `SELECT diff_id, header_id, ilk_id AS key, rho AS value FROM maker.jug_ilk_rho`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertMapping(ilkRhoResult, jugIlkRhoDiff.ID, header.Id, strconv.FormatInt(ilkID, 10), "1551968264")
	})
})
