package flap

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	storage_factory "github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flap"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("Executing the flap transformer", func() {
	var (
		db               *postgres.DB
		storageKeyLookup = flap.StorageKeysLookup{StorageRepository: &storage.MakerStorageRepository{}}
		repository       = flap.FlapStorageRepository{}
		transformer      storage_factory.Transformer
	)
	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		transformer = storage_factory.Transformer{
			Address:    common.HexToAddress("0x164a942d9d7A269B2Dc8551C8dFad32e8fFd0b80"),
			Mappings:   &storageKeyLookup,
			Repository: &repository,
		}
		transformer.NewTransformer(db)
	})

	It("reads in a vat storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiffRow{
			Contract:     transformer.Address,
			BlockHash:    blockHash,
			BlockHeight:  blockNumber,
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002"),
			StorageValue: common.HexToHash("000000000000000000000000284ecb5880cdc3362d979d07d162bf1d8488975d"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult,
			`SELECT block_number, block_hash, vat AS value FROM maker.flap_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, blockNumber, blockHash.Hex(), "0x284ecB5880CdC3362D979D07D162bf1d8488975D")
	})

	It("reads in a gem storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiffRow{
			Contract:     transformer.Address,
			BlockHash:    blockHash,
			BlockHeight:  blockNumber,
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
			StorageValue: common.HexToHash("000000000000000000000000a90843676a7f747a3c8adda142471369346369c1"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var gemResult test_helpers.VariableRes
		err = db.Get(&gemResult,
			`SELECT block_number, block_hash, gem AS value FROM maker.flap_gem`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(gemResult, blockNumber, blockHash.Hex(), "0xa90843676A7F747A3c8aDDa142471369346369c1")
	})

	It("reads in a beg storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiffRow{
			Contract:     transformer.Address,
			BlockHash:    blockHash,
			BlockHeight:  blockNumber,
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004"),
			StorageValue: common.HexToHash("000000000000000000000000000000000000000003648a260e3486a65a000000"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var begResult test_helpers.VariableRes
		err = db.Get(&begResult,
			`SELECT block_number, block_hash, beg AS value FROM maker.flap_beg`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(begResult, blockNumber, blockHash.Hex(), "1050000000000000000000000000")
	})

	It("reads in a ttl storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiffRow{
			Contract:     transformer.Address,
			BlockHash:    blockHash,
			BlockHeight:  blockNumber,
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005"),
			StorageValue: common.HexToHash("000000000000000000000000000000000000000000000002a300000000002a30"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var ttlResult test_helpers.VariableRes
		err = db.Get(&ttlResult,
			`SELECT block_number, block_hash, ttl AS value FROM maker.flap_ttl`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(ttlResult, blockNumber, blockHash.Hex(), "10800")
	})

	It("reads in a tau storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiffRow{
			Contract:     transformer.Address,
			BlockHash:    blockHash,
			BlockHeight:  blockNumber,
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005"),
			StorageValue: common.HexToHash("000000000000000000000000000000000000000000000002a300000000002a30"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var ttlResult test_helpers.VariableRes
		err = db.Get(&ttlResult,
			`SELECT block_number, block_hash, tau AS value FROM maker.flap_tau`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(ttlResult, blockNumber, blockHash.Hex(), "172800")
	})

	XIt("reads in a kicks storage diff and persists it", func() {
		//TODO: update this when we get a storage diff row for Flap kicks
	})

	It("reads in a live storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiffRow{
			Contract:     transformer.Address,
			BlockHash:    blockHash,
			BlockHeight:  blockNumber,
			StorageKey:   common.HexToHash("0000000000000000000000000000000000000000000000000000000000000007"),
			StorageValue: common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var liveResult test_helpers.VariableRes
		err = db.Get(&liveResult,
			`SELECT block_number, block_hash, live AS value FROM maker.flap_live`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, blockNumber, blockHash.Hex(), "1")
	})

	Describe("bids", func() {
		//TODO: update this when we get a storage diff row for Flap bids mapping
		//storage keys for bids with bid_id 0 will likely start with 0xc13ad76448cbefd1ee83b801bcd8f33061f2577d6118395e7b44ea21c7ef62e0
	})
})
