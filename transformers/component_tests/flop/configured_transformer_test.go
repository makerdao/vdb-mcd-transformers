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

package flop

import (
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop"
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

var _ = Describe("Executing the flop transformer", func() {
	var (
		db                     = test_config.NewTestDB(test_config.NewTestNode())
		flopperContractAddress = test_data.FlopAddress()
		keccakOfAddress        = types.HexToKeccak256Hash(flopperContractAddress)
		repository             = flop.StorageRepository{ContractAddress: flopperContractAddress}
		storageKeysLookup      = storage.NewKeysLookup(flop.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, flopperContractAddress))
		header                 = fakes.FakeHeader
		transformer            storage.Transformer
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		transformer = storage.Transformer{
			Address:           common.HexToAddress(flopperContractAddress),
			StorageKeysLookup: storageKeysLookup,
			Repository:        &repository,
		}
		transformer.NewTransformer(db)
		headerRepository := repositories.NewHeaderRepository(db)
		var insertHeaderErr error
		header.Id, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
	})

	It("reads in a vat storage diff and persists it", func() {
		vat := "0x1CC5ABe5C0464F3af2a10df0c711236a8446BF75"
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")
		value := common.HexToHash("0000000000000000000000001cc5abe5c0464f3af2a10df0c711236a8446bf75")
		diff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var vatResult test_helpers.VariableRes
		err = db.Get(&vatResult, `SELECT diff_id, header_id, vat AS value from maker.flop_vat`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vatResult, diff.ID, header.Id, vat)
	})

	It("reads in a gem storage diff and persists it", func() {
		gem := "0xAaF64BFCC32d0F15873a02163e7E500671a4ffcD"
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")
		value := common.HexToHash("000000000000000000000000aaf64bfcc32d0f15873a02163e7e500671a4ffcd")
		diff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var gemResult test_helpers.VariableRes
		err = db.Get(&gemResult, `SELECT diff_id, header_id, gem AS value from maker.flop_gem`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(gemResult, diff.ID, header.Id, gem)
	})

	It("reads in a beg storage diff and persists it", func() {
		beg := "1050000000000000000000000000"
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004")
		value := common.HexToHash("000000000000000000000000000000000000000003648a260e3486a65a000000")
		diff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var begResult test_helpers.VariableRes
		err = db.Get(&begResult, `SELECT diff_id, header_id, beg AS value from maker.flop_beg`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(begResult, diff.ID, header.Id, beg)
	})

	It("reads in a pad storage diff and persists it", func() {
		pad := "1500000000000000000000000000"
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005")
		value := common.HexToHash("000000000000000000000000000000000000000004d8c55aefb8c05b5c000000")
		diff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var ttlResult test_helpers.VariableRes
		err = db.Get(&ttlResult, `SELECT diff_id, header_id, pad AS value from maker.flop_pad`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(ttlResult, diff.ID, header.Id, pad)
	})

	It("reads in a ttl storage diff and persists it", func() {
		ttl := "10800"
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000006")
		value := common.HexToHash("000000000000000000000000000000000000000000000002a300000000002a30")
		diff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var ttlResult test_helpers.VariableRes
		err = db.Get(&ttlResult, `SELECT diff_id, header_id, ttl AS value from maker.flop_ttl`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(ttlResult, diff.ID, header.Id, ttl)
	})

	It("reads in a tau storage diff and persists it", func() {
		ttl := "172800"
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000006")
		value := common.HexToHash("000000000000000000000000000000000000000000000002a300000000002a30")
		diff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var tauResult test_helpers.VariableRes
		err = db.Get(&tauResult, `SELECT diff_id, header_id, tau AS value from maker.flop_tau`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(tauResult, diff.ID, header.Id, ttl)
	})

	It("reads in a kicks storage diff and persists it", func() {
		//TODO: update this when we get a storage diff row for Flop kicks
	})

	It("reads in a live storage diff and persists it", func() {
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000008")
		value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
		diff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)
		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var liveResult test_helpers.VariableRes
		err = db.Get(&liveResult, `SELECT diff_id, header_id, live AS value from maker.flop_live`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(liveResult, diff.ID, header.Id, "1")
	})

	It("reads in a vow storage diff and persists it", func() {
		vow := "0x1CC5ABe5C0464F3af2a10df0c711236a8446BF75"
		key := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000009")
		value := common.HexToHash("0000000000000000000000001cc5abe5c0464f3af2a10df0c711236a8446bf75")
		diff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

		err := transformer.Execute(diff)
		Expect(err).NotTo(HaveOccurred())

		var vowResult test_helpers.VariableRes
		err = db.Get(&vowResult, `SELECT diff_id, header_id, vow AS value from maker.flop_vow`)
		Expect(err).NotTo(HaveOccurred())
		test_helpers.AssertVariable(vowResult, diff.ID, header.Id, vow)
	})

	Describe("wards", func() {
		It("reads in a wards storage diff row and persists it", func() {
			denyLog := test_data.CreateTestLog(header.Id, db)
			denyModel := test_data.DenyModel()

			flopAddressID, flopAddressErr := shared.GetOrCreateAddress(test_data.FlopAddress(), db)
			Expect(flopAddressErr).NotTo(HaveOccurred())

			userAddress := "0xffb0382ca7cfdc4fc4d5cc8913af1393d7ee1ef1"
			userAddressID, userAddressErr := shared.GetOrCreateAddress(userAddress, db)
			Expect(userAddressErr).NotTo(HaveOccurred())

			msgSenderAddress := "0x" + fakes.RandomString(40)
			msgSenderAddressID, msgSenderAddressErr := shared.GetOrCreateAddress(msgSenderAddress, db)
			Expect(msgSenderAddressErr).NotTo(HaveOccurred())

			denyModel.ColumnValues[event.HeaderFK] = header.Id
			denyModel.ColumnValues[event.LogFK] = denyLog.ID
			denyModel.ColumnValues[event.AddressFK] = flopAddressID
			denyModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID
			denyModel.ColumnValues[constants.UsrColumn] = userAddressID
			insertErr := event.PersistModels([]event.InsertionModel{denyModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			key := common.HexToHash("4f3fc9e802fdeddd3e9ba88447e1731d7cfb3279d1b86a2328ef7efe1d42ac84")
			value := common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")
			wardsDiff := test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

			transformErr := transformer.Execute(wardsDiff)
			Expect(transformErr).NotTo(HaveOccurred())

			var wardsResult test_helpers.WardsMappingRes
			err := db.Get(&wardsResult, `SELECT diff_id, header_id, address_id, usr AS key, wards.wards AS value FROM maker.wards`)
			Expect(err).NotTo(HaveOccurred())
			Expect(wardsResult.AddressID).To(Equal(strconv.FormatInt(flopAddressID, 10)))
			test_helpers.AssertMapping(wardsResult.MappingRes, wardsDiff.ID, header.Id, strconv.FormatInt(userAddressID, 10), "1")
		})
	})

	Describe("bids", func() {
		//TODO: update when we get real flop bid storage diffs
		Describe("guy + tic + end packed slot", func() {
			var (
				bidId int
				diff  types.PersistedDiff
			)

			BeforeEach(func() {
				bidId = 1
				key := common.HexToHash("cc69885fda6bcc1a4ace058b4a62bf5e179ea78fd58a1ccd71c22cc9b6887931")
				value := common.HexToHash("00000002a300000000002a30284ecb5880cdc3362d979d07d162bf1d8488975d")
				diff = test_helpers.CreateDiffRecord(db, header, keccakOfAddress, key, value)

				addressId, addressErr := shared.GetOrCreateAddress(flopperContractAddress, db)
				Expect(addressErr).NotTo(HaveOccurred())

				_, writeErr := db.Exec(flop.InsertFlopKicksQuery, diff.ID, header.Id, addressId, bidId)
				Expect(writeErr).NotTo(HaveOccurred())

				executeErr := transformer.Execute(diff)
				Expect(executeErr).NotTo(HaveOccurred())
			})

			It("reads and persists a guy diff", func() {
				var bidGuyResult test_helpers.MappingRes
				dbErr := db.Get(&bidGuyResult, `SELECT diff_id, header_id, bid_id AS key, guy AS value FROM maker.flop_bid_guy`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidGuyResult, diff.ID, header.Id, strconv.Itoa(bidId), "0x284ecB5880CdC3362D979D07D162bf1d8488975D")
			})

			It("reads and persists a tic diff", func() {
				var bidTicResult test_helpers.MappingRes
				dbErr := db.Get(&bidTicResult, `SELECT diff_id, header_id, bid_id AS key, tic AS value FROM maker.flop_bid_tic`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidTicResult, diff.ID, header.Id, strconv.Itoa(bidId), "10800")
			})

			It("reads and persists an end diff", func() {
				var bidEndResult test_helpers.MappingRes
				dbErr := db.Get(&bidEndResult, `SELECT diff_id, header_id, bid_id AS key, "end" AS value FROM maker.flop_bid_end`)
				Expect(dbErr).NotTo(HaveOccurred())
				test_helpers.AssertMapping(bidEndResult, diff.ID, header.Id, strconv.Itoa(bidId), "172800")
			})
		})
	})
})
