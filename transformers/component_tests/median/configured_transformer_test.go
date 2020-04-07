package median_test

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
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

var _ = Describe("Executing the flip transformer", func() {
	var (
		db                = test_config.NewTestDB(test_config.NewTestNode())
		contractAddress   = test_data.EthMedianAddress()
		keccakAddress     = types.HexToKeccak256Hash(contractAddress)
		repository        = median.MedianStorageRepository{ContractAddress: contractAddress}
		storageKeysLookup = storage.NewKeysLookup(median.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress))
		header            = fakes.FakeHeader
		transformer       storage.Transformer
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
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

			medianAddressID, medianAddressErr := shared.GetOrCreateAddress(test_data.EthMedianAddress(), db)
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
			Expect(wardsResult.AddressID).To(Equal(strconv.FormatInt(medianAddressID, 10)))
			test_helpers.AssertMapping(wardsResult.MappingRes, wardsDiff.ID, header.Id, strconv.FormatInt(userAddressID, 10), "1")
		})
	})
})
