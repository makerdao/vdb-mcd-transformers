package vow

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	storage2 "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vow"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db          *postgres.DB
		mappings    = vow.VowMappings{StorageRepository: &storage2.MakerStorageRepository{}}
		repository  = vow.VowStorageRepository{}
		contractAddress = "4afcab85f27dd2e1a5ec1008b5b294e44e487f90"
		transformer = storage.Transformer{
			HashedAddress: utils.HexToKeccak256Hash(contractAddress),
			Mappings:      &mappings,
			Repository:    &repository,
		}
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
	})

	It("reads in a Vow.vat storage diff row and persists it", func() {
		blockNumber := 10501125
		vowVat := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash("1822bb271ce246212f0d097e59b3b04e0302819da3a2bd80e85b91e8c89fc883"),
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
			StorageValue:  common.HexToHash("00000000000000000000000067fd6c3575fc2dbe2cb596bd3bebc9edb5571fa1"),
		}
		err := transformer.Execute(vowVat)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT block_number, block_hash, vat AS value FROM maker.vow_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, blockNumber, "0x1822bb271ce246212f0d097e59b3b04e0302819da3a2bd80e85b91e8c89fc883", "0x67fd6c3575Fc2dBE2CB596bD3bEbc9EDb5571fA1")
	})

	It("reads in a Vow.flapper storage diff row and persists it", func() {
		blockNumber := 10980004
		vowFlapper := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash("44c07814be2cd81491f4d815ac922cc6590184e8777a5f0e3982c3b9ea83600e"),
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002"),
			StorageValue:  common.HexToHash("000000000000000000000000b6e31ab6ea62be7c530c32daea96e84d92fe20b7"),
		}
		err := transformer.Execute(vowFlapper)
		Expect(err).NotTo(HaveOccurred())

		var flapperResult test_helpers.VariableRes
		err = db.Get(&flapperResult, `SELECT block_number, block_hash, flapper AS value FROM maker.vow_flapper`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(flapperResult, blockNumber, "0x44c07814be2cd81491f4d815ac922cc6590184e8777a5f0e3982c3b9ea83600e", "0xB6e31ab6Ea62Be7c530C32DAEa96E84d92fe20B7")
	})

	It("reads in a Vow.flopper storage diff row and persists it", func() {
		blockNumber := 10980004
		vowFlopper := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash("44c07814be2cd81491f4d815ac922cc6590184e8777a5f0e3982c3b9ea83600e"),
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
			StorageValue:  common.HexToHash("000000000000000000000000275ec1950d6406e3ce6156f9f529c047ea41c8ce"),
		}
		err := transformer.Execute(vowFlopper)
		Expect(err).NotTo(HaveOccurred())

		var flopperResult test_helpers.VariableRes
		err = db.Get(&flopperResult, `SELECT block_number, block_hash, flopper AS value FROM maker.vow_flopper`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(flopperResult, blockNumber, "0x44c07814be2cd81491f4d815ac922cc6590184e8777a5f0e3982c3b9ea83600e", "0x275eC1950D6406e3cE6156f9F529c047Ea41c8cE")
	})

	It("reads in a Vow.sump storage diff row and persists it", func() {
		blockNumber := 10869770
		vowSump := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash("fe124bd8436290c364692b928a59f02f4d458c642b40398cbae173252b54093c"),
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000008"),
			StorageValue:  common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
		}
		err := transformer.Execute(vowSump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT block_number, block_hash, sump AS value FROM maker.vow_sump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, blockNumber, "0xfe124bd8436290c364692b928a59f02f4d458c642b40398cbae173252b54093c", "100000000000000000000000000000000000000000000")
	})

	It("reads in a Vow.bump storage diff row and persists it", func() {
		blockNumber := 10869768
		vowBump := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash("a750d8cf2317bb6d65b43b96ff24a179ed8c3a237f874c0e867987180b2527a8"),
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000009"),
			StorageValue:  common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
		}
		err := transformer.Execute(vowBump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT block_number, block_hash, bump AS value FROM maker.vow_bump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, blockNumber, "0xa750d8cf2317bb6d65b43b96ff24a179ed8c3a237f874c0e867987180b2527a8", "100000000000000000000000000000000000000000000")
	})
})
