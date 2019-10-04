package flop

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	storageFactory "github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strconv"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flop"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("Executing the flop transformer", func() {
	var (
		db                     *postgres.DB
		transformer            storageFactory.Transformer
		flopperContractAddress = "0x70b1a0fa8cc13cd3ce3cd65064c226dd9bc65f49"
		repository             = flop.FlopStorageRepository{ContractAddress: flopperContractAddress}
		storageKeyLookup       = flop.StorageKeysLookup{StorageRepository: &storage.MakerStorageRepository{}, ContractAddress: flopperContractAddress}
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		transformer = storageFactory.Transformer{
			HashedAddress: utils.HexToKeccak256Hash(flopperContractAddress),
			Mappings:      &storageKeyLookup,
			Repository:    &repository,
		}
		transformer.NewTransformer(db)
	})

	It("reads in a vat storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHash:     blockHash,
			BlockHeight:   blockNumber,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002"),
			StorageValue:  common.HexToHash("000000000000000000000000284ecb5880cdc3362d979d07d162bf1d8488975d"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT block_number, block_hash, vat AS value from maker.flop_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, blockNumber, blockHash.Hex(), "0x284ecB5880CdC3362D979D07D162bf1d8488975D")
	})

	It("reads in a gem storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHash:     blockHash,
			BlockHeight:   blockNumber,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
			StorageValue:  common.HexToHash("000000000000000000000000a90843676a7f747a3c8adda142471369346369c1"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var gemResult test_helpers.VariableRes
		err = db.Get(&gemResult, `SELECT block_number, block_hash, gem AS value from maker.flop_gem`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(gemResult, blockNumber, blockHash.Hex(), "0xa90843676A7F747A3c8aDDa142471369346369c1")
	})

	It("reads in a beg storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHash:     blockHash,
			BlockHeight:   blockNumber,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004"),
			StorageValue:  common.HexToHash("000000000000000000000000000000000000000003648a260e3486a65a000000"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var begResult test_helpers.VariableRes
		err = db.Get(&begResult, `SELECT block_number, block_hash, beg AS value from maker.flop_beg`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(begResult, blockNumber, blockHash.Hex(), "1050000000000000000000000000")
	})

	It("reads in a ttl storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHash:     blockHash,
			BlockHeight:   blockNumber,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005"),
			StorageValue:  common.HexToHash("000000000000000000000000000000000000000000000002a300000000002a30"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var ttlResult test_helpers.VariableRes
		err = db.Get(&ttlResult, `SELECT block_number, block_hash, ttl AS value from maker.flop_ttl`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(ttlResult, blockNumber, blockHash.Hex(), "10800")
	})

	It("reads in a tau storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHash:     blockHash,
			BlockHeight:   blockNumber,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005"),
			StorageValue:  common.HexToHash("000000000000000000000000000000000000000000000002a300000000002a30"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var tauResult test_helpers.VariableRes
		err = db.Get(&tauResult, `SELECT block_number, block_hash, tau AS value from maker.flop_tau`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(tauResult, blockNumber, blockHash.Hex(), "172800")
	})

	XIt("reads in a kicks storage diff and persists it", func() {
		//TODO: update this when we get a storage diff row for Flop kicks
	})

	It("reads in a live storage diff and persists it", func() {
		blockNumber := 11579860
		blockHash := common.HexToHash("3d8fd744457d476c3a1a9e4cbbaa0d951e0280416988fe87528a0aaac50186a8")
		diff := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHash:     blockHash,
			BlockHeight:   blockNumber,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000007"),
			StorageValue:  common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var liveResult test_helpers.VariableRes
		err = db.Get(&liveResult, `SELECT block_number, block_hash, live AS value from maker.flop_live`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, blockNumber, blockHash.Hex(), "1")
	})

	Describe("bids", func() {
		//TODO: update when we get real flop bid storage diffs
		Describe("guy + tic + end packed slot", func() {
			bidId := 1
			blockNumber := 11579891
			blockHash := common.HexToHash("5f2be3f6566f39dddfcfcf29784866280399ed9070af0b4fccd465509260349d")
			diff := utils.StorageDiff{
				HashedAddress: transformer.HashedAddress,
				BlockHash:     blockHash,
				BlockHeight:   blockNumber,
				StorageKey:    common.HexToHash("cc69885fda6bcc1a4ace058b4a62bf5e179ea78fd58a1ccd71c22cc9b6887931"),
				StorageValue:  common.HexToHash("00000002a300000000002a30284ecb5880cdc3362d979d07d162bf1d8488975d"),
			}

			BeforeEach(func() {
				addressId, addressErr := shared.GetOrCreateAddress(flopperContractAddress, db)
				Expect(addressErr).NotTo(HaveOccurred())

				_, writeErr := db.Exec(flop.InsertFlopKicksQuery, blockNumber, blockHash.Hex(), addressId, bidId)
				Expect(writeErr).NotTo(HaveOccurred())

				executeErr := transformer.Execute(diff)
				Expect(executeErr).NotTo(HaveOccurred())
			})

			It("reads and persists a guy diff", func() {
				var bidGuyResult test_helpers.MappingRes
				dbErr := db.Get(&bidGuyResult, `SELECT block_number, block_hash, bid_id AS key, guy AS value FROM maker.flop_bid_guy`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidGuyResult, blockNumber, blockHash.Hex(), strconv.Itoa(bidId), "0x284ecB5880CdC3362D979D07D162bf1d8488975D")
			})

			It("reads and persists a tic diff", func() {
				var bidTicResult test_helpers.MappingRes
				dbErr := db.Get(&bidTicResult, `SELECT block_number, block_hash, bid_id AS key, tic AS value FROM maker.flop_bid_tic`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidTicResult, blockNumber, blockHash.Hex(), strconv.Itoa(bidId), "10800")
			})

			It("reads and persists an end diff", func() {
				var bidEndResult test_helpers.MappingRes
				dbErr := db.Get(&bidEndResult, `SELECT block_number, block_hash, bid_id AS key, "end" AS value FROM maker.flop_bid_end`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidEndResult, blockNumber, blockHash.Hex(), strconv.Itoa(bidId), "172800")
			})
		})
	})
})
