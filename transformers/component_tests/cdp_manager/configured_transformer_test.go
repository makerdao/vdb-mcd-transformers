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

package cdp_manager

import (
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cdp_manager"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Executing the transformer", func() {
	var (
		db                = test_config.NewTestDB(test_config.NewTestNode())
		storageKeysLookup = storage.NewKeysLookup(cdp_manager.NewKeysLoader(&mcdStorage.MakerStorageRepository{}))
		repository        = cdp_manager.CdpManagerStorageRepository{}
		contractAddress   = "7a4991c6bd1053c31f1678955ce839999d9841b1"
		transformer       = storage.Transformer{
			HashedAddress:     types.HexToKeccak256Hash(contractAddress),
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
		header = fakes.FakeHeader
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		transformer.NewTransformer(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		header.Id, insertHeaderErr = headerRepository.CreateOrUpdateHeader(header)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("reads in a vat storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")
		value := common.HexToHash("00000000000000000000000004c67ea772ebb467383772cb1b64c7a9b1e02bca")
		vatDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

		err := transformer.Execute(vatDiff)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT diff_id, header_id, vat AS value FROM maker.cdp_manager_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, vatDiff.ID, header.Id, "0x04C67ea772EBb467383772Cb1b64c7a9b1e02BCa")
	})

	It("reads in a cdpi storage diff row and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		cdpiDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

		err := transformer.Execute(cdpiDiff)
		Expect(err).NotTo(HaveOccurred())

		var cdpiResult test_helpers.VariableRes
		err = db.Get(&cdpiResult, `SELECT diff_id, header_id, cdpi AS value FROM maker.cdp_manager_cdpi`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(cdpiResult, cdpiDiff.ID, header.Id, "3")
	})

	Describe("cdpi key mappings", func() {
		cdpi := 2

		It("reads in an urns storage diff row and persists it", func() {
			key := common.HexToHash("679795a0195a1b76cdebb7c51d74e058aee92919b8c3389af86ef24535e8a28c")
			value := common.HexToHash("00000000000000000000000031f92649bf2d780be06bab1c5f591d0f1cc4b0d2")
			urnsDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			_, insertErr := db.Exec(cdp_manager.InsertCdpiQuery, urnsDiff.ID, header.Id, cdpi)
			Expect(insertErr).NotTo(HaveOccurred())

			err := transformer.Execute(urnsDiff)
			Expect(err).NotTo(HaveOccurred())

			var urnsResult test_helpers.MappingRes
			err = db.Get(&urnsResult, `SELECT diff_id, header_id, cdpi AS key, urn AS value FROM maker.cdp_manager_urns`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(urnsResult, urnsDiff.ID, header.Id, strconv.Itoa(cdpi), "0x31f92649BF2d780BE06BAB1C5F591d0f1Cc4b0D2")
		})

		It("reads in a list prev storage diff row and persists it", func() {
			key := common.HexToHash("c3a24b0501bd2c13a7e57f2db4369ec4c223447539fc0724a9d55ac4a06ebd4d")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
			listPrevDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			_, insertErr := db.Exec(cdp_manager.InsertCdpiQuery, listPrevDiff.ID, header.Id, cdpi)
			Expect(insertErr).NotTo(HaveOccurred())

			err := transformer.Execute(listPrevDiff)
			Expect(err).NotTo(HaveOccurred())

			var listPrevResult test_helpers.MappingRes
			err = db.Get(&listPrevResult, `SELECT diff_id, header_id, cdpi AS key, prev AS value FROM maker.cdp_manager_list_prev`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(listPrevResult, listPrevDiff.ID, header.Id, strconv.Itoa(cdpi), "1")
		})

		It("reads in a list next storage diff row and persists it", func() {
			key := common.HexToHash("c3a24b0501bd2c13a7e57f2db4369ec4c223447539fc0724a9d55ac4a06ebd4e")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
			listNextDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			_, insertErr := db.Exec(cdp_manager.InsertCdpiQuery, listNextDiff.ID, header.Id, cdpi)
			Expect(insertErr).NotTo(HaveOccurred())

			err := transformer.Execute(listNextDiff)
			Expect(err).NotTo(HaveOccurred())

			var listNextResult test_helpers.MappingRes
			err = db.Get(&listNextResult, `SELECT diff_id, header_id, cdpi AS key, next AS value FROM maker.cdp_manager_list_next`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(listNextResult, listNextDiff.ID, header.Id, strconv.Itoa(cdpi), "3")
		})

		It("reads in an owns storage diff row and persists it", func() {
			key := common.HexToHash("91da3fd0782e51c6b3986e9e672fd566868e71f3dbc2d6c2cd6fbb3e361af2a7")
			value := common.HexToHash("00000000000000000000000016fb96a5fa0427af0c8f7cf1eb4870231c8154b6")
			ownsDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			_, insertErr := db.Exec(cdp_manager.InsertCdpiQuery, ownsDiff.ID, header.Id, cdpi)
			Expect(insertErr).NotTo(HaveOccurred())

			err := transformer.Execute(ownsDiff)
			Expect(err).NotTo(HaveOccurred())

			var ownsResult test_helpers.MappingRes
			err = db.Get(&ownsResult, `SELECT diff_id, header_id, cdpi AS key, owner AS value FROM maker.cdp_manager_owns`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(ownsResult, ownsDiff.ID, header.Id, strconv.Itoa(cdpi), "0x16Fb96a5fa0427Af0C8F7cF1eB4870231c8154B6")
		})

		It("reads in an ilks storage diff row and persists it", func() {
			key := common.HexToHash("89832631fb3c3307a103ba2c84ab569c64d6182a18893dcd163f0f1c2090733a")
			value := common.HexToHash("4554482d41000000000000000000000000000000000000000000000000000000")
			ilksDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			_, insertErr := db.Exec(cdp_manager.InsertCdpiQuery, ilksDiff.ID, header.Id, cdpi)
			Expect(insertErr).NotTo(HaveOccurred())

			ilk := "0x4554482d41000000000000000000000000000000000000000000000000000000"
			transformErr := transformer.Execute(ilksDiff)
			Expect(transformErr).NotTo(HaveOccurred())

			var ilksResult test_helpers.MappingRes
			readErr := db.Get(&ilksResult, `SELECT diff_id, header_id, cdpi AS key, ilk_id AS value FROM maker.cdp_manager_ilks`)
			Expect(readErr).NotTo(HaveOccurred())
			ilkId, ilkErr := shared.GetOrCreateIlk(ilk, db)
			Expect(ilkErr).NotTo(HaveOccurred())
			test_helpers.AssertMapping(ilksResult, ilksDiff.ID, header.Id, strconv.Itoa(cdpi), strconv.FormatInt(ilkId, 10))
		})
	})

	Describe("owner key mappings", func() {
		owner := "0x16Fb96a5fa0427Af0C8F7cF1eB4870231c8154B6"

		It("reads in a first storage diff row and persists it", func() {
			key := common.HexToHash("361ac87b78b4b96bd716b22773b802e3ec15c69f4ba42c6d6f8cb594e4914397")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
			firstDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			_, insertErr := db.Exec(cdp_manager.InsertOwnsQuery, firstDiff.ID, header.Id, rand.Int(), owner)
			Expect(insertErr).NotTo(HaveOccurred())

			err := transformer.Execute(firstDiff)
			Expect(err).NotTo(HaveOccurred())

			var firstResult test_helpers.MappingRes
			err = db.Get(&firstResult, `SELECT diff_id, header_id, owner AS key, first AS value FROM maker.cdp_manager_first`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(firstResult, firstDiff.ID, header.Id, owner, "1")
		})

		It("reads in a last storage diff row and persists it", func() {
			key := common.HexToHash("4f62af9d63bc3c7d5e96c3d1083b2438d0fa9b6244cdfc09d00d09b1afbd7438")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
			lastDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			_, insertErr := db.Exec(cdp_manager.InsertOwnsQuery, lastDiff.ID, header.Id, rand.Int(), owner)
			Expect(insertErr).NotTo(HaveOccurred())

			err := transformer.Execute(lastDiff)
			Expect(err).NotTo(HaveOccurred())

			var lastResult test_helpers.MappingRes
			err = db.Get(&lastResult, `SELECT diff_id, header_id, owner AS key, last AS value FROM maker.cdp_manager_last`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(lastResult, lastDiff.ID, header.Id, owner, "2")
		})

		It("reads in a count storage diff row and persists it", func() {
			key := common.HexToHash("0b29a919802754cc12fe9af109d06c6ac93a8cac604ffa44ff5474c8c41bd5c0")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
			countDiff := test_helpers.CreateDiffRecord(db, header, transformer.HashedAddress, key, value)

			_, insertErr := db.Exec(cdp_manager.InsertOwnsQuery, countDiff.ID, header.Id, rand.Int(), owner)
			Expect(insertErr).NotTo(HaveOccurred())

			err := transformer.Execute(countDiff)
			Expect(err).NotTo(HaveOccurred())

			var lastResult test_helpers.MappingRes
			err = db.Get(&lastResult, `SELECT diff_id, header_id, owner AS key, count AS value FROM maker.cdp_manager_count`)
			Expect(err).NotTo(HaveOccurred())
			test_helpers.AssertMapping(lastResult, countDiff.ID, header.Id, owner, "2")
		})
	})
})
