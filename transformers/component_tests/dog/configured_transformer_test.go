package dog

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/dog"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
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

	It("reads in a Dog Vat storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		value := common.HexToHash("000000000000000000000000acdd1ee0f74954ed8f0ac581b081b7b86bd6aad9")
		dogVatDiff := test_helpers.CreateDiffRecord(db, header, contractAddress, key, value)

		contractAddressID, contractAddressErr := repository.GetOrCreateAddress(db, contractAddress.Hex())
		Expect(contractAddressErr).NotTo(HaveOccurred())

		err := transformer.Execute(dogVatDiff)
		Expect(err).NotTo(HaveOccurred())

		diffAddressID, diffAddressErr := repository.GetOrCreateAddress(db, "0xaCdd1ee0F74954Ed8F0aC581b081B7b86bD6aad9")
		Expect(diffAddressErr).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableResWithAddress
		err = db.Get(&vatResult, `SELECT diff_id, header_id, address_id, vat AS value FROM maker.dog_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariableWithAddress(vatResult, dogVatDiff.ID, header.Id, contractAddressID, strconv.FormatInt(diffAddressID, 10))
	})

	It("reads in a Dog Vow storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
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
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004")
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
})
