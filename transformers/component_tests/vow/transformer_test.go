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
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vow"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db                = test_config.NewTestDB(test_config.NewTestNode())
		storageKeysLookup = storage.NewKeysLookup(vow.NewKeysLoader(&mcdStorage.MakerStorageRepository{}))
		repository        = vow.VowStorageRepository{}
		contractAddress   = "4afcab85f27dd2e1a5ec1008b5b294e44e487f90"
		transformer       = storage.Transformer{
			HashedAddress:     utils.HexToKeccak256Hash(contractAddress),
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
		headerID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("reads in a Vow.vat storage diff row and persists it", func() {
		vowVat := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001"),
			StorageValue:  common.HexToHash("00000000000000000000000067fd6c3575fc2dbe2cb596bd3bebc9edb5571fa1"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(vowVat)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT header_id, vat AS value FROM maker.vow_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, headerID, "0x67fd6c3575Fc2dBE2CB596bD3bEbc9EDb5571fA1")
	})

	It("reads in a Vow.flapper storage diff row and persists it", func() {
		vowFlapper := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002"),
			StorageValue:  common.HexToHash("000000000000000000000000b6e31ab6ea62be7c530c32daea96e84d92fe20b7"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(vowFlapper)
		Expect(err).NotTo(HaveOccurred())

		var flapperResult test_helpers.VariableRes
		err = db.Get(&flapperResult, `SELECT header_id, flapper AS value FROM maker.vow_flapper`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(flapperResult, headerID, "0xB6e31ab6Ea62Be7c530C32DAEa96E84d92fe20B7")
	})

	It("reads in a Vow.flopper storage diff row and persists it", func() {
		vowFlopper := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003"),
			StorageValue:  common.HexToHash("000000000000000000000000275ec1950d6406e3ce6156f9f529c047ea41c8ce"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(vowFlopper)
		Expect(err).NotTo(HaveOccurred())

		var flopperResult test_helpers.VariableRes
		err = db.Get(&flopperResult, `SELECT header_id, flopper AS value FROM maker.vow_flopper`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(flopperResult, headerID, "0x275eC1950D6406e3cE6156f9F529c047Ea41c8cE")
	})

	It("reads in a Vow.dump storage diff row and persists it", func() {
		vowDump := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000008"),
			StorageValue:  common.HexToHash("000000000000000000000000000000000000000000000000002386f26fc10000"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(vowDump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT header_id, dump AS value FROM maker.vow_dump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, headerID, "10000000000000000")
	})

	It("reads in a Vow.sump storage diff row and persists it", func() {
		vowSump := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("0000000000000000000000000000000000000000000000000000000000000009"),
			StorageValue:  common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(vowSump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT header_id, sump AS value FROM maker.vow_sump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, headerID, "100000000000000000000000000000000000000000000")
	})

	It("reads in a Vow.bump storage diff row and persists it", func() {
		vowBump := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("000000000000000000000000000000000000000000000000000000000000000a"),
			StorageValue:  common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(vowBump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT header_id, bump AS value FROM maker.vow_bump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, headerID, "100000000000000000000000000000000000000000000")
	})

	It("reads in a Vow.hump storage diff row and persists it", func() {
		// TODO: Update with a real storage diff
		vowHump := utils.StorageDiff{
			HashedAddress: transformer.HashedAddress,
			StorageKey:    common.HexToHash("000000000000000000000000000000000000000000000000000000000000000b"),
			StorageValue:  common.HexToHash("00000000000000000000000000047bf19673df52e37f2410011d100000000000"),
			HeaderID:      headerID,
		}
		err := transformer.Execute(vowHump)
		Expect(err).NotTo(HaveOccurred())

		var rowResult test_helpers.VariableRes
		err = db.Get(&rowResult, `SELECT header_id, hump AS value FROM maker.vow_hump`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(rowResult, headerID, "100000000000000000000000000000000000000000000")
	})
})
