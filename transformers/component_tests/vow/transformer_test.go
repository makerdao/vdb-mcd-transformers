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

package vow

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vow"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db                *postgres.DB
		storageKeysLookup = mcdStorage.NewKeysLookup(vow.NewKeysLoader(&mcdStorage.MakerStorageRepository{}))
		repository        = vow.VowStorageRepository{}
		contractAddress   = "4afcab85f27dd2e1a5ec1008b5b294e44e487f90"
		transformer       = storage.Transformer{
			HashedAddress: utils.HexToKeccak256Hash(contractAddress),
			Mappings:      storageKeysLookup,
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

	It("reads in a Vow.dump storage diff row and persists it", func() {
		blockNumber := 13475028
		blockHash := "ce700b74d34213a26cf368152c10897d07413102de31c9f89aaac453809a5106"
		vowDump := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash(blockHash),
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000008"),
			StorageValue:  common.HexToHash("000000000000000000000000000000000000000000000000002386f26fc10000"),
		}
		err := transformer.Execute(vowDump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT block_number, block_hash, dump AS value FROM maker.vow_dump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, blockNumber, "0x"+blockHash, "10000000000000000")
	})

	It("reads in a Vow.sump storage diff row and persists it", func() {
		blockNumber := 13475031
		vowSump := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash("1139ef7db215364b5f1bc8711e46f9e5ec4645d9d29a5baa54750af10cc035fc"),
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000009"),
			StorageValue:  common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
		}
		err := transformer.Execute(vowSump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT block_number, block_hash, sump AS value FROM maker.vow_sump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, blockNumber, "0x1139ef7db215364b5f1bc8711e46f9e5ec4645d9d29a5baa54750af10cc035fc", "100000000000000000000000000000000000000000000")
	})

	It("reads in a Vow.bump storage diff row and persists it", func() {
		blockNumber := 13475024
		vowBump := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash("b0753160aff9d2e44f677a5198f2baae36101a33df85f110af543f6878dc3f43"),
			StorageKey:    common.HexToHash("000000000000000000000000000000000000000000000000000000000000000a"),
			StorageValue:  common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
		}
		err := transformer.Execute(vowBump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT block_number, block_hash, bump AS value FROM maker.vow_bump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, blockNumber, "0xb0753160aff9d2e44f677a5198f2baae36101a33df85f110af543f6878dc3f43", "100000000000000000000000000000000000000000000")
	})

	It("reads in a Vow.hump storage diff row and persists it", func() {
		// TODO: Update with a real storage diff
		blockNumber := 10869768
		vowHump := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			BlockHeight:   blockNumber,
			BlockHash:     common.HexToHash("a750d8cf2317bb6d65b43b96ff24a179ed8c3a237f874c0e867987180b2527a8"),
			StorageKey:    common.HexToHash("000000000000000000000000000000000000000000000000000000000000000b"),
			StorageValue:  common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
		}
		err := transformer.Execute(vowHump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT block_number, block_hash, hump AS value FROM maker.vow_hump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, blockNumber, "0xa750d8cf2317bb6d65b43b96ff24a179ed8c3a237f874c0e867987180b2527a8", "100000000000000000000000000000000000000000000")
	})
})
