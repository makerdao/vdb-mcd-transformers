package dog

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/dog"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db              = test_config.NewTestDB(test_config.NewTestNode())
		contractAddress = common.HexToAddress(test_data.Dog130Address())
		transformer     storage.Transformer
		header          = fakes.FakeHeader
	)

	BeforeEach(func() {
		storageKeysLookup := storage.NewKeysLookup(dog.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress.Hex()))
		repository := dog.StorageRepository{ContractAddress: contractAddress.Hex()}
		transformer = storage.Transformer{
			Address:           contractAddress,
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		header.Id, insertHeaderErr = headerRepository.CreateOrUpdateHeader(header)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("reads in a Dog Vow storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
		value := common.HexToHash("00000000000000000000000021444ac712ccd21ce82af24ea1aec64cf07361d2")
		dogVowDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddress.Hex())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		err := transformer.Execute(dogVowDiff)
		Expect(err).NotTo(HaveOccurred())

		diffAddressID, diffAddressErr := repository.GetOrCreateAddress(db, "0x21444ac712ccd21ce82af24ea1aec64cf07361d2")
		Expect(diffAddressErr).NotTo(HaveOccurred())

		var vowResult test_helpers.VariableResWithAddress
		err = db.Get(&vowResult, `SELECT diff_id, header_id, address_id, vow AS value FROM maker.dog_vow`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariableWithAddress(vowResult, dogVowDiff.ID, header.Id, contractAddressID, strconv.FormatInt(diffAddressID, 10))
	})

	It("reads in a Dog Live storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		dogLiveDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddress.Hex())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		err := transformer.Execute(dogLiveDiff)
		Expect(err).NotTo(HaveOccurred())

		var liveResult test_helpers.VariableResWithAddress
		err = db.Get(&liveResult, `SELECT diff_id, header_id, address_id, live AS value FROM maker.dog_live`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariableWithAddress(liveResult, dogLiveDiff.ID, header.Id, contractAddressID, "1")
	})

	It("reads in a Dog Hole storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		dogHoleDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddress.Hex())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		err := transformer.Execute(dogHoleDiff)
		Expect(err).NotTo(HaveOccurred())

		var holeResult test_helpers.VariableResWithAddress
		err = db.Get(&holeResult, `SELECT diff_id, header_id, address_id, hole AS value FROM maker.dog_hole`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariableWithAddress(holeResult, dogHoleDiff.ID, header.Id, contractAddressID, "3")
	})

	It("reads in a Dog Dirt storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000015")
		dogHoleDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddress.Hex())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		err := transformer.Execute(dogHoleDiff)
		Expect(err).NotTo(HaveOccurred())

		var dirtResult test_helpers.VariableResWithAddress
		err = db.Get(&dirtResult, `SELECT diff_id, header_id, address_id, dirt AS value FROM maker.dog_dirt`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariableWithAddress(dirtResult, dogHoleDiff.ID, header.Id, contractAddressID, "21")
	})

	Describe("wards", func() {
		It("reads in a wards storage diff row and persists it", func() {
			denyLog := test_data.CreateTestLog(header.Id, db)
			denyModel := test_data.DenyModel()

			dogAddressID, dogAddressErr := repository.GetOrCreateAddress(db, contractAddress.Hex())
			Expect(dogAddressErr).NotTo(HaveOccurred())

			userAddress := "0x39ad5d336a4c08fac74879f796e1ea0af26c1521"
			userAddressID, userAddressErr := repository.GetOrCreateAddress(db, userAddress)
			Expect(userAddressErr).NotTo(HaveOccurred())

			msgSenderAddress := "0x" + fakes.RandomString(40)
			msgSenderAddressID, msgSenderAddressErr := repository.GetOrCreateAddress(db, msgSenderAddress)
			Expect(msgSenderAddressErr).NotTo(HaveOccurred())

			denyModel.ColumnValues[event.HeaderFK] = header.Id
			denyModel.ColumnValues[event.LogFK] = denyLog.ID
			denyModel.ColumnValues[event.AddressFK] = dogAddressID
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
			test_helpers.AssertMappingWithAddress(wardsResult, wardsDiff.ID, header.Id, dogAddressID, strconv.FormatInt(userAddressID, 10), "1")
		})
	})
})
