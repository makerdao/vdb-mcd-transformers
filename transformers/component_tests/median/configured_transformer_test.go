package median_test

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/median"
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

var _ = Describe("Executing the median transformer", func() {
	var (
		db                = test_config.NewTestDB(test_config.NewTestNode())
		contractAddress   = test_data.MedianEthAddress()
		keccakAddress     = types.HexToKeccak256Hash(contractAddress)
		storageKeysLookup = storage.NewKeysLookup(median.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress))
		header            = fakes.FakeHeader
		transformer       storage.Transformer
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		var repository = median.MedianStorageRepository{ContractAddress: contractAddress}
		transformer = storage.Transformer{
			Address:           common.HexToAddress(contractAddress),
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
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

			medianAddressID, medianAddressErr := shared.GetOrCreateAddress(test_data.MedianEthAddress(), db)
			Expect(medianAddressErr).NotTo(HaveOccurred())

			userAddress := "0xffb0382ca7cfdc4fc4d5cc8913af1393d7ee1ef1"
			userAddressID, userAddressErr := shared.GetOrCreateAddress(userAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())

			msgSenderAddress := "0x" + fakes.RandomString(40)
			msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
			Expect(msgSenderAddressErr).NotTo(HaveOccurred())

			denyModel.ColumnValues[event.HeaderFK] = header.Id
			denyModel.ColumnValues[event.LogFK] = denyLog.ID
			denyModel.ColumnValues[event.AddressFK] = medianAddressID
			denyModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
			denyModel.ColumnValues[constants.UsrColumn] = userAddressID
			insertErr := event.PersistModels([]event.InsertionModel{denyModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			key := common.HexToHash("4f3fc9e802fdeddd3e9ba88447e1731d7cfb3279d1b86a2328ef7efe1d42ac84")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
			wardsDiff := test_helpers.CreateDiffRecord(db, header, keccakAddress, key, value)

			transformErr := transformer.Execute(wardsDiff)
			Expect(transformErr).NotTo(HaveOccurred())

			var wardsResult test_helpers.MappingResWithAddress
			err := db.Get(&wardsResult, `SELECT diff_id, header_id, address_id, usr AS key, wards.wards AS value FROM maker.wards`)
			Expect(err).NotTo(HaveOccurred())
			Expect(wardsResult.AddressID).To(Equal(medianAddressID))
			test_helpers.AssertMapping(wardsResult.MappingRes, wardsDiff.ID, header.Id, strconv.FormatInt(userAddressID, 10), "1")
		})
	})

	Describe("bud", func() {
		It("reads in a bud storage diff row and persists it", func() {
			kissLog := test_data.CreateTestLog(header.Id, db)
			kissModel := test_data.MedianKissSingleModel()

			medianAddressID, medianAddressErr := shared.GetOrCreateAddress(test_data.MedianEthAddress(), db)
			Expect(medianAddressErr).NotTo(HaveOccurred())

			aAddress := "0xffb0382ca7cfdc4fc4d5cc8913af1393d7ee1ef1"
			aAddressID, aAddressErr := shared.GetOrCreateAddress(aAddress, db)
			Expect(aAddressErr).NotTo(HaveOccurred())

			msgSenderAddress := "0x" + fakes.RandomString(40)
			msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
			Expect(msgSenderAddressErr).NotTo(HaveOccurred())

			kissModel.ColumnValues[event.HeaderFK] = header.Id
			kissModel.ColumnValues[event.LogFK] = kissLog.ID
			kissModel.ColumnValues[event.AddressFK] = medianAddressID
			kissModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
			kissModel.ColumnValues[constants.AColumn] = aAddressID
			insertErr := event.PersistModels([]event.InsertionModel{kissModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			key := common.HexToHash("6e8bbf796f21b82c83c834b2cacf88452e5bba3a2fb53ad9e5b2c0e6c54820fd")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
			wardsDiff := test_helpers.CreateDiffRecord(db, header, keccakAddress, key, value)

			transformErr := transformer.Execute(wardsDiff)
			Expect(transformErr).NotTo(HaveOccurred())

			var budResult test_helpers.MappingResWithAddress
			err := db.Get(&budResult, `SELECT diff_id, header_id, address_id, a AS key, bud AS value FROM maker.median_bud`)
			Expect(err).NotTo(HaveOccurred())
			Expect(budResult.AddressID).To(Equal(medianAddressID))
			test_helpers.AssertMapping(budResult.MappingRes, wardsDiff.ID, header.Id, strconv.FormatInt(aAddressID, 10), "1")
		})
	})

	Describe("orcl", func() {
		It("reads in an orcl storage diff row and persists it", func() {
			liftLog := test_data.CreateTestLog(header.Id, db)
			liftModel := test_data.MedianLiftModelWithOneAccount()

			medianAddressID, medianAddressErr := shared.GetOrCreateAddress(test_data.MedianEthAddress(), db)
			Expect(medianAddressErr).NotTo(HaveOccurred())

			aAddress := "0xaC8519b3495d8A3E3E44c041521cF7aC3f8F63B3"
			aAddressID, aAddressErr := shared.GetOrCreateAddress(aAddress, db)
			Expect(aAddressErr).NotTo(HaveOccurred())

			msgSenderAddress := "0x" + fakes.RandomString(40)
			msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
			Expect(msgSenderAddressErr).NotTo(HaveOccurred())

			liftModel.ColumnValues[event.HeaderFK] = header.Id
			liftModel.ColumnValues[event.LogFK] = liftLog.ID
			liftModel.ColumnValues[event.AddressFK] = medianAddressID
			liftModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
			liftModel.ColumnValues[constants.AColumn] = pq.Array([]string{aAddress})
			insertErr := event.PersistModels([]event.InsertionModel{liftModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			key := common.HexToHash("6e8810a330507229748898345becb3182d8632868d2bd2df00dfbd0f623252f9")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
			liftDiff := test_helpers.CreateDiffRecord(db, header, keccakAddress, key, value)

			transformErr := transformer.Execute(liftDiff)
			Expect(transformErr).NotTo(HaveOccurred())

			var orclResult test_helpers.MappingResWithAddress
			err := db.Get(&orclResult, `SELECT diff_id, header_id, address_id, a AS key, orcl AS value from maker.median_orcl`)
			Expect(err).NotTo(HaveOccurred())
			Expect(orclResult.AddressID).To(Equal(medianAddressID))
			test_helpers.AssertMapping(orclResult.MappingRes, liftDiff.ID, header.Id, strconv.FormatInt(aAddressID, 10), "1")
		})
	})

	Describe("slot", func() {
		It("reads a slot diff row and persists it", func() {
			liftLog := test_data.CreateTestLog(header.Id, db)
			liftModel := test_data.MedianLiftModelWithOneAccount()

			medianAddressID, medianAddressErr := shared.GetOrCreateAddress(test_data.MedianEthAddress(), db)
			Expect(medianAddressErr).NotTo(HaveOccurred())

			aAddress := "0xaC8519b3495d8A3E3E44c041521cF7aC3f8F63B3"
			_, aAddressErr := shared.GetOrCreateAddress(aAddress, db)
			Expect(aAddressErr).NotTo(HaveOccurred())

			msgSenderAddress := "0x" + fakes.RandomString(40)
			msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
			Expect(msgSenderAddressErr).NotTo(HaveOccurred())

			liftModel.ColumnValues[event.HeaderFK] = header.Id
			liftModel.ColumnValues[event.LogFK] = liftLog.ID
			liftModel.ColumnValues[event.AddressFK] = medianAddressID
			liftModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
			liftModel.ColumnValues[constants.AColumn] = pq.Array([]string{aAddress})
			insertErr := event.PersistModels([]event.InsertionModel{liftModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			//key is keccak hash aAddress passed through solidity bitshift uint8(uint256(a[i]) >> 152) + the index
			key := common.HexToHash("2944b9af8d962e2b5d171cd2b530c03b245945580d9e2a1c9efc472e2e5ec88b")
			value := common.HexToHash("000000000000000000000000ac8519b3495d8a3e3e44c041521cf7ac3f8f63b3")
			liftDiff := test_helpers.CreateDiffRecord(db, header, keccakAddress, key, value)

			transformErr := transformer.Execute(liftDiff)
			Expect(transformErr).NotTo(HaveOccurred())

			var slotResult test_helpers.MappingResWithAddress
			err := db.Get(&slotResult, `SELECT diff_id, header_id, address_id, slot_id AS key, slot AS value from maker.median_slot`)
			Expect(err).NotTo(HaveOccurred())
			slotValueAddressID, slotErr := shared.GetOrCreateAddress(value.String(), db)
			Expect(slotErr).NotTo(HaveOccurred())
			Expect(slotResult.AddressID).To(Equal(medianAddressID))
			solidityKey := "172"
			test_helpers.AssertMapping(slotResult.MappingRes, liftDiff.ID, header.Id, solidityKey, strconv.FormatInt(slotValueAddressID, 10))
		})
	})
})
