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

package trigger_test

import (
	"database/sql"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip"
	. "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Updating bid_event table", func() {
	var (
		header core.Header
		logID  int64
		db     = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		header = CreateHeader(int64(rand.Int31()), rand.Int(), db)
		logID = test_data.CreateTestLog(header.Id, db).ID
	})

	Specify("inserting a flip_kick event triggers a bid_event insertion", func() {
		flipAddress := test_data.EthFlipAddress()
		addressID, addressErr := shared.GetOrCreateAddress(flipAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		flipKickModel := test_data.FlipKickModel()
		flipKickModel.ColumnValues[event.HeaderFK] = header.Id
		flipKickModel.ColumnValues[event.AddressFK] = addressID
		flipKickModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEvent(flipKickModel, "kick", flipAddress, header.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{flipKickModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, block_height FROM api.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a flop_kick event triggers a bid_event insertion", func() {
		flopAddress := test_data.FlopAddress()
		addressID, addressErr := shared.GetOrCreateAddress(flopAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(header.Id, db).ID
		flopKickModel := test_data.FlopKickModel()
		flopKickModel.ColumnValues[event.HeaderFK] = header.Id
		flopKickModel.ColumnValues[event.AddressFK] = addressID
		flopKickModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEvent(flopKickModel, "kick", flopAddress, header.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{flopKickModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, block_height FROM api.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a flap_kick event triggers a bid_event insertion", func() {
		flapAddress := test_data.FlapAddress()
		addressID, addressErr := shared.GetOrCreateAddress(flapAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(header.Id, db).ID
		flapKickModel := test_data.FlapKickModel()
		flapKickModel.ColumnValues[event.HeaderFK] = header.Id
		flapKickModel.ColumnValues[event.AddressFK] = addressID
		flapKickModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEvent(flapKickModel, "kick", flapAddress, header.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{flapKickModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, block_height FROM api.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a tend event triggers a bid_event insertion", func() {
		address := test_data.EthFlipAddress()
		addressID, addressErr := shared.GetOrCreateAddress(address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(header.Id, db).ID
		tendModel := test_data.TendModel()
		tendModel.ColumnValues[event.HeaderFK] = header.Id
		tendModel.ColumnValues[event.AddressFK] = addressID
		tendModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEvent(tendModel, "tend", address, header.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{tendModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, block_height FROM api.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a dent event triggers a bid_event insertion", func() {
		address := test_data.EthFlipAddress()
		addressID, addressErr := shared.GetOrCreateAddress(address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(header.Id, db).ID
		dentModel := test_data.DentModel()
		dentModel.ColumnValues[event.HeaderFK] = header.Id
		dentModel.ColumnValues[event.AddressFK] = addressID
		dentModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEvent(dentModel, "dent", address, header.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{dentModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, block_height FROM api.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a tick event triggers a bid_event insertion", func() {
		address := test_data.EthFlipAddress()
		addressID, addressErr := shared.GetOrCreateAddress(address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(header.Id, db).ID
		tickModel := test_data.TickModel()
		tickModel.ColumnValues[event.HeaderFK] = header.Id
		tickModel.ColumnValues[event.AddressFK] = addressID
		tickModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEventNullStrings(tickModel, "tick", address, header.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{tickModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, block_height FROM api.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a deal event triggers a bid_event insertion", func() {
		address := test_data.EthFlipAddress()
		addressID, addressErr := shared.GetOrCreateAddress(address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(header.Id, db).ID
		dealModel := test_data.DealModel()
		dealModel.ColumnValues[event.HeaderFK] = header.Id
		dealModel.ColumnValues[event.AddressFK] = addressID
		dealModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEventNullStrings(dealModel, "deal", address, header.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{dealModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, block_height FROM api.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a yank event triggers a bid_event insertion", func() {
		address := test_data.EthFlipAddress()
		addressID, addressErr := shared.GetOrCreateAddress(address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(header.Id, db).ID
		yankModel := test_data.YankModel()
		yankModel.ColumnValues[event.HeaderFK] = header.Id
		yankModel.ColumnValues[event.AddressFK] = addressID
		yankModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEventNullStrings(yankModel, "yank", address, header.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{yankModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, block_height FROM api.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Describe("inserting events after diffs", func() {
		var (
			flipAddress   = test_data.EthFlipAddress()
			flipRepo      flip.FlipStorageRepository
			flipKickModel event.InsertionModel
			diffID        int64
		)

		BeforeEach(func() {
			flipAddressID, addressErr := shared.GetOrCreateAddress(flipAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
			flipRepo = flip.FlipStorageRepository{ContractAddress: flipAddress}
			flipRepo.SetDB(db)
			diffID = CreateFakeDiffRecord(db)
			flipKickModel = test_data.FlipKickModel()
			flipKickModel.ColumnValues[event.HeaderFK] = header.Id
			flipKickModel.ColumnValues[event.AddressFK] = flipAddressID
			flipKickModel.ColumnValues[event.LogFK] = logID
		})

		It("gets the relevant ilk for event", func() {
			expectedEvent := expectedBidEvent(flipKickModel, "kick", flipAddress, header.BlockNumber)
			expectedEvent.IlkIdentifier = test_helpers.FakeIlk.Identifier

			flipIlkErr := flipRepo.Create(diffID, header.Id, flip.IlkMetadata, test_helpers.FakeIlk.Hex)
			Expect(flipIlkErr).NotTo(HaveOccurred())

			insertErr := event.PersistModels([]event.InsertionModel{flipKickModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			var bidEvents []bidEvent
			queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, ilk_identifier, block_height FROM api.bid_event`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(bidEvents).To(ConsistOf(expectedEvent))
		})

		It("gets the relevant urn for event", func() {
			expectedEvent := expectedBidEvent(flipKickModel, "kick", flipAddress, header.BlockNumber)
			expectedEvent.UrnIdentifier = common.HexToAddress("0x" + test_data.RandomString(40)).Hex()

			bidUsrMetadata := types.GetValueMetadata(storage.BidUsr,
				map[types.Key]string{constants.BidId: expectedEvent.BidID}, types.Address)
			usrErr := flipRepo.Create(diffID, header.Id, bidUsrMetadata, expectedEvent.UrnIdentifier)
			Expect(usrErr).NotTo(HaveOccurred())

			insertErr := event.PersistModels([]event.InsertionModel{flipKickModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			var bidEvents []bidEvent
			queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, urn_identifier, block_height FROM api.bid_event`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(bidEvents).To(ConsistOf(expectedEvent))
		})
	})

	Describe("when diffs are inserted after events", func() {
		var (
			flipAddress   string
			flipRepo      flip.FlipStorageRepository
			diffID        int64
			flipKickModel event.InsertionModel
		)

		BeforeEach(func() {
			flipAddress = test_data.EthFlipAddress()
			flipAddressID, addressErr := shared.GetOrCreateAddress(flipAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
			flipRepo = flip.FlipStorageRepository{ContractAddress: flipAddress}
			flipRepo.SetDB(db)
			diffID = CreateFakeDiffRecord(db)
			flipKickModel = test_data.FlipKickModel()
			flipKickModel.ColumnValues[event.HeaderFK] = header.Id
			flipKickModel.ColumnValues[event.AddressFK] = flipAddressID
			flipKickModel.ColumnValues[event.LogFK] = logID
			insertErr := event.PersistModels([]event.InsertionModel{flipKickModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())
		})

		Specify("inserting a flip_ilk diff triggers an update to the ilk_identifier of related bids", func() {
			expectedEvent := expectedBidEvent(flipKickModel, "kick", flipAddress, header.BlockNumber)
			expectedEvent.IlkIdentifier = test_helpers.FakeIlk.Identifier

			flipIlkErr := flipRepo.Create(diffID, header.Id, flip.IlkMetadata, test_helpers.FakeIlk.Hex)
			Expect(flipIlkErr).NotTo(HaveOccurred())

			var bidEvents []bidEvent
			queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, ilk_identifier, block_height FROM api.bid_event`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(bidEvents).To(ConsistOf(expectedEvent))
		})

		Specify("inserting a flip_bid_usr diff triggers an update to the urn_identifier of related bids", func() {
			usr := common.HexToAddress("0x" + test_data.RandomString(40)).Hex()
			expectedEvent := expectedBidEvent(flipKickModel, "kick", flipAddress, header.BlockNumber)
			bidUsrMetadata := types.GetValueMetadata(storage.BidUsr,
				map[types.Key]string{constants.BidId: expectedEvent.BidID}, types.Address)
			expectedEvent.UrnIdentifier = usr

			usrErr := flipRepo.Create(diffID, header.Id, bidUsrMetadata, usr)
			Expect(usrErr).NotTo(HaveOccurred())

			var bidEvents []bidEvent
			queryErr := db.Select(&bidEvents, `SELECT bid_id, contract_address, act, lot, bid_amount, urn_identifier, block_height FROM api.bid_event`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(bidEvents).To(ConsistOf(expectedEvent))
		})
	})
})

func expectedBidEvent(eventModel event.InsertionModel, bidAct, contractAddress string, blockHeight int64) bidEvent {
	return bidEvent{
		BidID:           eventModel.ColumnValues[constants.BidIDColumn].(string),
		ContractAddress: contractAddress,
		Act:             bidAct,
		Lot:             test_helpers.GetValidNullString(eventModel.ColumnValues[constants.LotColumn].(string)),
		Bid:             test_helpers.GetValidNullString(eventModel.ColumnValues[constants.BidColumn].(string)),
		BlockHeight:     strconv.FormatInt(blockHeight, 10),
	}
}

func expectedBidEventNullStrings(eventModel event.InsertionModel, bidAct, contractAddress string, blockHeight int64) bidEvent {
	return bidEvent{
		BidID:           eventModel.ColumnValues[constants.BidIDColumn].(string),
		ContractAddress: contractAddress,
		Act:             bidAct,
		Lot:             sql.NullString{Valid: false},
		Bid:             sql.NullString{Valid: false},
		BlockHeight:     strconv.FormatInt(blockHeight, 10),
	}
}

type bidEvent struct {
	BidID           string `db:"bid_id"`
	ContractAddress string `db:"contract_address"`
	Act             string
	Lot             sql.NullString
	Bid             sql.NullString `db:"bid_amount"`
	IlkIdentifier   string         `db:"ilk_identifier"`
	UrnIdentifier   string         `db:"urn_identifier"`
	BlockHeight     string         `db:"block_height"`
}
