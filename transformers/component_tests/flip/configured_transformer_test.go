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
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip"
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
		contractAddress   = test_data.FlipEthV100Address()
		keccakAddress     = types.HexToKeccak256Hash(contractAddress)
		repository        = flip.StorageRepository{ContractAddress: contractAddress}
		storageKeysLookup = storage.NewKeysLookup(flip.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress))
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

	It("reads in a vat storage diff and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
		value := common.HexToHash("000000000000000000000000284ecb5880cdc3362d979d07d162bf1d8488975d")
		diff := test_helpers.CreateDiffRecord(db, header, keccakAddress, key, value)

		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult,
			`SELECT diff_id, header_id, vat AS value FROM maker.flip_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, diff.ID, header.Id, "0x284ecB5880CdC3362D979D07D162bf1d8488975D")
	})

	It("reads in an ilk storage diff and persists it", func() {
		ilk := "4554482d41000000000000000000000000000000000000000000000000000000"
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		value := common.HexToHash(ilk)
		diff := test_helpers.CreateDiffRecord(db, header, keccakAddress, key, value)

		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var ilkResult test_helpers.VariableRes
		err = db.Get(&ilkResult,
			`SELECT diff_id, header_id, ilk_id AS value FROM maker.flip_ilk`)
		Expect(err).NotTo(HaveOccurred())
		ilkID, ilkErr := shared.GetOrCreateIlk(ilk, db)
		Expect(ilkErr).NotTo(HaveOccurred())
		test_helpers.AssertVariable(ilkResult, diff.ID, header.Id, strconv.FormatInt(ilkID, 10))
	})

	It("reads in a beg storage diff and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004")
		value := common.HexToHash("000000000000000000000000000000000000000003648a260e3486a65a000000")
		diff := test_helpers.CreateDiffRecord(db, header, keccakAddress, key, value)

		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var begResult test_helpers.VariableRes
		err = db.Get(&begResult,
			`SELECT diff_id, header_id, beg AS value FROM maker.flip_beg`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(begResult, diff.ID, header.Id, "1050000000000000000000000000")
	})

	It("reads in a ttl storage diff and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005")
		value := common.HexToHash("000000000000000000000000000000000000000000000002a300000000002a30")
		diff := test_helpers.CreateDiffRecord(db, header, keccakAddress, key, value)

		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var ttlResult test_helpers.VariableRes
		err = db.Get(&ttlResult,
			`SELECT diff_id, header_id, ttl AS value FROM maker.flip_ttl`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(ttlResult, diff.ID, header.Id, "10800")
	})

	It("reads in a tau storage diff and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005")
		value := common.HexToHash("000000000000000000000000000000000000000000000002a300000000002a30")
		diff := test_helpers.CreateDiffRecord(db, header, keccakAddress, key, value)

		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var ttlResult test_helpers.VariableRes
		err = db.Get(&ttlResult,
			`SELECT diff_id, header_id, tau AS value FROM maker.flip_tau`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(ttlResult, diff.ID, header.Id, "172800")
	})

	XIt("reads in a kicks storage diff and persists it", func() {
		//TODO: update this when we get a storage diff row for Flap kicks
	})

	Describe("wards", func() {
		It("reads in a wards storage diff row and persists it", func() {
			denyLog := test_data.CreateTestLog(header.Id, db)
			denyModel := test_data.DenyModel()

			flipAddressID, flipAddressErr := shared.GetOrCreateAddress(test_data.FlipEthV100Address(), db)
			Expect(flipAddressErr).NotTo(HaveOccurred())

			userAddress := "0xffb0382ca7cfdc4fc4d5cc8913af1393d7ee1ef1"
			userAddressID, userAddressErr := shared.GetOrCreateAddress(userAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())

			msgSenderAddress := "0x" + fakes.RandomString(40)
			msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
			Expect(msgSenderAddressErr).NotTo(HaveOccurred())

			denyModel.ColumnValues[event.HeaderFK] = header.Id
			denyModel.ColumnValues[event.LogFK] = denyLog.ID
			denyModel.ColumnValues[event.AddressFK] = flipAddressID
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
			Expect(wardsResult.AddressID).To(Equal(flipAddressID))
			test_helpers.AssertMapping(wardsResult.MappingRes, wardsDiff.ID, header.Id, strconv.FormatInt(userAddressID, 10), "1")
		})
	})

	Describe("bids", func() {
		Describe("guy + tic + end packed slot", func() {
			var (
				bidId int
				diff  types.PersistedDiff
			)

			BeforeEach(func() {
				bidId = 1
				key := common.HexToHash("cc69885fda6bcc1a4ace058b4a62bf5e179ea78fd58a1ccd71c22cc9b6887931")
				value := common.HexToHash("00000002a300000000002a30284ecb5880cdc3362d979d07d162bf1d8488975d")
				diff = test_helpers.CreateDiffRecord(db, header, keccakAddress, key, value)

				addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
				Expect(addressErr).NotTo(HaveOccurred())

				_, writeErr := db.Exec(flip.InsertFlipKicksQuery, diff.ID, header.Id, addressId, bidId)
				Expect(writeErr).NotTo(HaveOccurred())

				executeErr := transformer.Execute(diff)
				Expect(executeErr).NotTo(HaveOccurred())
			})

			It("reads and persists a guy diff", func() {
				var bidGuyResult test_helpers.MappingRes
				dbErr := db.Get(&bidGuyResult, `SELECT diff_id, header_id, bid_id AS key, guy AS value FROM maker.flip_bid_guy`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidGuyResult, diff.ID, header.Id, strconv.Itoa(bidId), "0x284ecB5880CdC3362D979D07D162bf1d8488975D")
			})

			It("reads and persists a tic diff", func() {
				var bidTicResult test_helpers.MappingRes
				dbErr := db.Get(&bidTicResult, `SELECT diff_id, header_id, bid_id AS key, tic AS value FROM maker.flip_bid_tic`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidTicResult, diff.ID, header.Id, strconv.Itoa(bidId), "10800")
			})

			It("reads and persists an end diff", func() {
				var bidEndResult test_helpers.MappingRes
				dbErr := db.Get(&bidEndResult, `SELECT diff_id, header_id, bid_id AS key, "end" AS value FROM maker.flip_bid_end`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidEndResult, diff.ID, header.Id, strconv.Itoa(bidId), "172800")
			})
		})
	})
})
