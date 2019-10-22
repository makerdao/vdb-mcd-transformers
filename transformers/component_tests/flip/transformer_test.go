// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package flip

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flip"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strconv"
)

var _ = Describe("Executing the flip transformer", func() {
	var (
		db                *postgres.DB
		repository        = flip.FlipStorageRepository{}
		transformer       storage.Transformer
		contractAddress   = "0x43c331c0389a92af62ee726d5ae0c8a424320c31"
		storageKeysLookup = storage.NewKeysLookup(flip.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress))
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		transformer = storage.Transformer{
			HashedAddress:     utils.HexToKeccak256Hash(contractAddress),
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
		transformer.NewTransformer(db)
	})

	It("reads in a vat storage diff and persists it", func() {
		blockNumber := 11579891
		blockHash := common.HexToHash("5f2be3f6566f39dddfcfcf29784866280399ed9070af0b4fccd465509260349d")
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
		err = db.Get(&vatResult,
			`SELECT block_number, block_hash, vat AS value FROM maker.flip_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, blockNumber, blockHash.Hex(), "0x284ecB5880CdC3362D979D07D162bf1d8488975D")
	})

	It("reads in an ilk storage diff and persists it", func() {
		blockNumber := 11579891
		blockHash := common.HexToHash("5f2be3f6566f39dddfcfcf29784866280399ed9070af0b4fccd465509260349d")
		ilk := "4554482d41000000000000000000000000000000000000000000000000000000"
		diff := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHash:     blockHash,
			BlockHeight:   blockNumber,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
			StorageValue:  common.HexToHash(ilk),
		}
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var ilkResult test_helpers.VariableRes
		err = db.Get(&ilkResult,
			`SELECT block_number, block_hash, ilk_id AS value FROM maker.flip_ilk`)
		Expect(err).NotTo(HaveOccurred())
		ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
		Expect(ilkErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(ilkResult, blockNumber, blockHash.Hex(), strconv.FormatInt(ilkID, 10))
	})

	It("reads in a beg storage diff and persists it", func() {
		blockNumber := 11579891
		blockHash := common.HexToHash("5f2be3f6566f39dddfcfcf29784866280399ed9070af0b4fccd465509260349d")
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
		err = db.Get(&begResult,
			`SELECT block_number, block_hash, beg AS value FROM maker.flip_beg`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(begResult, blockNumber, blockHash.Hex(), "1050000000000000000000000000")
	})

	It("reads in a ttl storage diff and persists it", func() {
		blockNumber := 11579891
		blockHash := common.HexToHash("5f2be3f6566f39dddfcfcf29784866280399ed9070af0b4fccd465509260349d")
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
		err = db.Get(&ttlResult,
			`SELECT block_number, block_hash, ttl AS value FROM maker.flip_ttl`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(ttlResult, blockNumber, blockHash.Hex(), "10800")
	})

	It("reads in a tau storage diff and persists it", func() {
		blockNumber := 11579891
		blockHash := common.HexToHash("5f2be3f6566f39dddfcfcf29784866280399ed9070af0b4fccd465509260349d")
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
		err = db.Get(&ttlResult,
			`SELECT block_number, block_hash, tau AS value FROM maker.flip_tau`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(ttlResult, blockNumber, blockHash.Hex(), "172800")
	})

	XIt("reads in a kicks storage diff and persists it", func() {
		//TODO: update this when we get a storage diff row for Flap kicks
	})

	Describe("bids", func() {
		//TODO: update when we get real flip bid storage diffs
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
				addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
				Expect(addressErr).NotTo(HaveOccurred())

				_, writeErr := db.Exec(flip.InsertFlipKicksQuery, blockNumber, blockHash.Hex(), addressId, bidId)
				Expect(writeErr).NotTo(HaveOccurred())

				executeErr := transformer.Execute(diff)
				Expect(executeErr).NotTo(HaveOccurred())
			})

			It("reads and persists a guy diff", func() {
				var bidGuyResult test_helpers.MappingRes
				dbErr := db.Get(&bidGuyResult, `SELECT block_number, block_hash, bid_id AS key, guy AS value FROM maker.flip_bid_guy`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidGuyResult, blockNumber, blockHash.Hex(), strconv.Itoa(bidId), "0x284ecB5880CdC3362D979D07D162bf1d8488975D")
			})

			It("reads and persists a tic diff", func() {
				var bidTicResult test_helpers.MappingRes
				dbErr := db.Get(&bidTicResult, `SELECT block_number, block_hash, bid_id AS key, tic AS value FROM maker.flip_bid_tic`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidTicResult, blockNumber, blockHash.Hex(), strconv.Itoa(bidId), "10800")
			})

			It("reads and persists an end diff", func() {
				var bidEndResult test_helpers.MappingRes
				dbErr := db.Get(&bidEndResult, `SELECT block_number, block_hash, bid_id AS key, "end" AS value FROM maker.flip_bid_end`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidEndResult, blockNumber, blockHash.Hex(), strconv.Itoa(bidId), "172800")
			})
		})
	})
})
