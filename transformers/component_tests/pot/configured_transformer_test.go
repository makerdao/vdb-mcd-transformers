package pot

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/pot"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db                = test_config.NewTestDB(test_config.NewTestNode())
		storageKeysLookup = storage.NewKeysLookup(pot.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, test_data.PotAddress()))
		contractAddress   = common.HexToAddress(test_data.PotAddress())
		repository        = pot.StorageRepository{ContractAddress: contractAddress.Hex()}
		transformer       = storage.Transformer{
			Address:           contractAddress,
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

	It("reads in a Pot user pie storage diff row and persists it", func() {
		potJoinLog := test_data.CreateTestLog(header.Id, db)
		potJoin := test_data.PotJoinModel()
		userAddress := "0x57aA8B02F5D3E28371FEdCf672C8668869f9AAC7"
		addressID, addressErr := shared.GetOrCreateAddress(userAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		potJoin.ColumnValues[event.HeaderFK] = header.Id
		potJoin.ColumnValues[event.LogFK] = potJoinLog.ID
		potJoin.ColumnValues[constants.MsgSenderColumn] = addressID
		insertErr := event.PersistModels([]event.InsertionModel{potJoin}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		key := common.HexToHash("727c626d1f473b2fe91b66946682443a691c9d83c3d5816a53bc3365286f89f4")
		value := common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000")
		potPieDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		executeErr := transformer.Execute(potPieDiff)
		Expect(executeErr).NotTo(HaveOccurred())

		var pieResult test_helpers.MappingRes
		getErr := db.Get(&pieResult, `SELECT diff_id, header_id, "user" AS key, pie AS value FROM maker.pot_user_pie`)
		Expect(getErr).NotTo(HaveOccurred())
		test_helpers.AssertMapping(pieResult, potPieDiff.ID, header.Id, strconv.FormatInt(addressID, 10), "1000000000000000000000000000")
	})

	It("reads in a Pot Pie storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
		value := common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000")
		potPieDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		executeErr := transformer.Execute(potPieDiff)

		Expect(executeErr).NotTo(HaveOccurred())
		var pieResult test_helpers.VariableRes
		getErr := db.Get(&pieResult, `SELECT diff_id, header_id, pie AS value FROM maker.pot_pie`)
		Expect(getErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(pieResult, potPieDiff.ID, header.Id, "1000000000000000000000000000")
	})

	It("reads in a Pot dsr storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		value := common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000")
		potDsrDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		executeErr := transformer.Execute(potDsrDiff)

		Expect(executeErr).NotTo(HaveOccurred())
		var dsrResult test_helpers.VariableRes
		getErr := db.Get(&dsrResult, `SELECT diff_id, header_id, dsr AS value FROM maker.pot_dsr`)
		Expect(getErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(dsrResult, potDsrDiff.ID, header.Id, "1000000000000000000000000000")
	})

	It("reads in a Pot chi storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004")
		value := common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000")
		potChiDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		executeErr := transformer.Execute(potChiDiff)

		Expect(executeErr).NotTo(HaveOccurred())
		var chiResult test_helpers.VariableRes
		getErr := db.Get(&chiResult, `SELECT diff_id, header_id, chi AS value FROM maker.pot_chi`)
		Expect(getErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(chiResult, potChiDiff.ID, header.Id, "1000000000000000000000000000")
	})

	It("reads in a Pot vat storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005")
		value := common.HexToHash("00000000000000000000000057aa8b02f5d3e28371fedcf672c8668869f9aac7")
		potVatDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		executeErr := transformer.Execute(potVatDiff)
		Expect(executeErr).NotTo(HaveOccurred())

		addressID, addressErr := shared.GetOrCreateAddress("0x57aA8B02F5D3E28371FEdCf672C8668869f9AAC7", db)
		Expect(addressErr).NotTo(HaveOccurred())
		var vatResult test_helpers.VariableRes
		getErr := db.Get(&vatResult, `SELECT diff_id, header_id, vat AS value FROM maker.pot_vat`)
		Expect(getErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, potVatDiff.ID, header.Id, strconv.FormatInt(addressID, 10))
	})

	It("reads in a Pot vow storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000006")
		value := common.HexToHash("00000000000000000000000057aa8b02f5d3e28371fedcf672c8668869f9aac7")
		potVowDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		executeErr := transformer.Execute(potVowDiff)
		Expect(executeErr).NotTo(HaveOccurred())

		addressID, addressErr := shared.GetOrCreateAddress("0x57aA8B02F5D3E28371FEdCf672C8668869f9AAC7", db)
		Expect(addressErr).NotTo(HaveOccurred())
		var vowResult test_helpers.VariableRes
		getErr := db.Get(&vowResult, `SELECT diff_id, header_id, vow AS value FROM maker.pot_vow`)
		Expect(getErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vowResult, potVowDiff.ID, header.Id, strconv.FormatInt(addressID, 10))
	})

	It("reads in a Pot rho storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000007")
		value := common.HexToHash("0000000000000000000000000000000000000000033b2e3c9fd0803ce8000000")
		potRhoDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		executeErr := transformer.Execute(potRhoDiff)

		Expect(executeErr).NotTo(HaveOccurred())
		var rhoResult test_helpers.VariableRes
		getErr := db.Get(&rhoResult, `SELECT diff_id, header_id, rho AS value FROM maker.pot_rho`)
		Expect(getErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rhoResult, potRhoDiff.ID, header.Id, "1000000000000000000000000000")
	})

	Describe("wards", func() {
		It("reads in a wards storage diff row and persists it", func() {
			denyLog := test_data.CreateTestLog(header.Id, db)
			denyModel := test_data.DenyModel()

			potAddressID, potAddressErr := shared.GetOrCreateAddress(contractAddress.Hex(), db)
			Expect(potAddressErr).NotTo(HaveOccurred())

			userAddress := "0x39ad5d336a4c08fac74879f796e1ea0af26c1521"
			userAddressID, userAddressErr := shared.GetOrCreateAddress(userAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())

			msgSenderAddress := "0x" + fakes.RandomString(40)
			msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
			Expect(msgSenderAddressErr).NotTo(HaveOccurred())

			denyModel.ColumnValues[event.HeaderFK] = header.Id
			denyModel.ColumnValues[event.LogFK] = denyLog.ID
			denyModel.ColumnValues[event.AddressFK] = potAddressID
			denyModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
			denyModel.ColumnValues[constants.UsrColumn] = userAddressID
			insertErr := event.PersistModels([]event.InsertionModel{denyModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			key := common.HexToHash("b6d2a4300cc4010859f67ce7c804312ce9cc8f1032cdeb24e96d4b5562a4d01b")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
			wardsDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

			transformErr := transformer.Execute(wardsDiff)
			Expect(transformErr).NotTo(HaveOccurred())

			var wardsResult test_helpers.MappingResWithAddress
			err := db.Get(&wardsResult, `SELECT diff_id, header_id, address_id, usr AS key, wards.wards AS value FROM maker.wards`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMappingWithAddress(wardsResult, wardsDiff.ID, header.Id, potAddressID, strconv.FormatInt(userAddressID, 10), "1")
		})
	})

	It("reads in a Pot live storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000008")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		potLiveDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		executeErr := transformer.Execute(potLiveDiff)

		Expect(executeErr).NotTo(HaveOccurred())
		var liveResult test_helpers.VariableRes
		getErr := db.Get(&liveResult, `SELECT diff_id, header_id, live AS value FROM maker.pot_live`)
		Expect(getErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, potLiveDiff.ID, header.Id, "1")
	})
})
