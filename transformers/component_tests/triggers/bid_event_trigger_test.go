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
		timestampOne int32
		blockOne     int
		logID        int64
		headerOne    core.Header
		db           = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		timestampOne = rand.Int31()
		blockOne = rand.Int()
		headerOne = CreateHeader(int64(timestampOne), blockOne, db)
		logID = test_data.CreateTestLog(headerOne.Id, db).ID
	})

	Specify("inserting a flip_kick event triggers a bid_event insertion", func() {
		flipAddress := test_data.FlipEthAddress()
		addressID, addressErr := shared.GetOrCreateAddress(flipAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		flipKickModel := test_data.FlipKickModel()
		flipKickModel.ColumnValues[event.HeaderFK] = headerOne.Id
		flipKickModel.ColumnValues[event.AddressFK] = addressID
		flipKickModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEvent(flipKickModel, "kick", flipAddress, headerOne.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{flipKickModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, block_height FROM maker.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a flop_kick event triggers a bid_event insertion", func() {
		flopAddress := test_data.FlopAddress()
		addressID, addressErr := shared.GetOrCreateAddress(flopAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(headerOne.Id, db).ID
		flopKickModel := test_data.FlopKickModel()
		flopKickModel.ColumnValues[event.HeaderFK] = headerOne.Id
		flopKickModel.ColumnValues[event.AddressFK] = addressID
		flopKickModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEvent(flopKickModel, "kick", flopAddress, headerOne.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{flopKickModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, block_height FROM maker.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a flap_kick event triggers a bid_event insertion", func() {
		flapAddress := test_data.FlapAddress()
		addressID, addressErr := shared.GetOrCreateAddress(flapAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(headerOne.Id, db).ID
		flapKickModel := test_data.FlapKickModel()
		flapKickModel.ColumnValues[event.HeaderFK] = headerOne.Id
		flapKickModel.ColumnValues[event.AddressFK] = addressID
		flapKickModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEvent(flapKickModel, "kick", flapAddress, headerOne.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{flapKickModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, block_height FROM maker.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a tend event triggers a bid_event insertion", func() {
		address := test_data.FlipEthAddress()
		addressID, addressErr := shared.GetOrCreateAddress(address, db)
		Expect(addressErr).NotTo(HaveOccurred())

		tendLog := test_data.CreateTestLogFromEventLog(headerOne.Id, test_data.TendEventLog.Log, db)
		tendModel := test_data.TendModel()
		tendModel.ColumnValues[event.HeaderFK] = headerOne.Id
		tendModel.ColumnValues[event.AddressFK] = addressID
		tendModel.ColumnValues[event.LogFK] = tendLog.ID
		test_data.AssignMessageSenderID(tendLog, tendModel, db)
		expectedEvent := expectedBidEvent(tendModel, "tend", address, headerOne.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{tendModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, block_height FROM maker.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a dent event triggers a bid_event insertion", func() {
		address := test_data.FlipEthAddress()
		addressID, addressErr := shared.GetOrCreateAddress(address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(headerOne.Id, db).ID
		dentModel := test_data.DentModel()
		dentModel.ColumnValues[event.HeaderFK] = headerOne.Id
		dentModel.ColumnValues[event.AddressFK] = addressID
		dentModel.ColumnValues[event.LogFK] = logID
		test_data.AssignMessageSenderID(test_data.DentEventLog, dentModel, db)
		expectedEvent := expectedBidEvent(dentModel, "dent", address, headerOne.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{dentModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, block_height FROM maker.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a tick event triggers a bid_event insertion", func() {
		address := test_data.FlipEthAddress()
		addressID, addressErr := shared.GetOrCreateAddress(address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(headerOne.Id, db).ID
		tickModel := test_data.TickModel()
		tickModel.ColumnValues[event.HeaderFK] = headerOne.Id
		tickModel.ColumnValues[event.AddressFK] = addressID
		tickModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEventNullStrings(tickModel, "tick", address, headerOne.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{tickModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, block_height FROM maker.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a deal event triggers a bid_event insertion", func() {
		address := test_data.FlipEthAddress()
		addressID, addressErr := shared.GetOrCreateAddress(address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(headerOne.Id, db).ID
		dealModel := test_data.DealModel()
		dealModel.ColumnValues[event.HeaderFK] = headerOne.Id
		dealModel.ColumnValues[event.AddressFK] = addressID
		dealModel.ColumnValues[event.LogFK] = logID
		test_data.AssignMessageSenderID(test_data.DealEventLog, dealModel, db)
		expectedEvent := expectedBidEventNullStrings(dealModel, "deal", address, headerOne.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{dealModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, block_height FROM maker.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Specify("inserting a yank event triggers a bid_event insertion", func() {
		address := test_data.FlipEthAddress()
		addressID, addressErr := shared.GetOrCreateAddress(address, db)
		Expect(addressErr).NotTo(HaveOccurred())
		logID := test_data.CreateTestLog(headerOne.Id, db).ID
		yankModel := test_data.YankModel()
		yankModel.ColumnValues[event.HeaderFK] = headerOne.Id
		yankModel.ColumnValues[event.AddressFK] = addressID
		yankModel.ColumnValues[event.LogFK] = logID
		expectedEvent := expectedBidEventNullStrings(yankModel, "yank", address, headerOne.BlockNumber)

		insertErr := event.PersistModels([]event.InsertionModel{yankModel}, db)
		Expect(insertErr).NotTo(HaveOccurred())

		var bidEvents []bidEvent
		queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, block_height FROM maker.bid_event`)
		Expect(queryErr).NotTo(HaveOccurred())
		Expect(bidEvents).To(ConsistOf(expectedEvent))
	})

	Describe("inserting events after flip-specific diffs", func() {
		var (
			flipAddress   = test_data.FlipEthAddress()
			flipRepo      flip.StorageRepository
			flipKickModel event.InsertionModel
			diffID        int64
		)

		BeforeEach(func() {
			flipAddressID, addressErr := shared.GetOrCreateAddress(flipAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
			flipRepo = flip.StorageRepository{ContractAddress: flipAddress}
			flipRepo.SetDB(db)
			diffID = CreateFakeDiffRecord(db)
			flipKickModel = test_data.FlipKickModel()
			flipKickModel.ColumnValues[event.HeaderFK] = headerOne.Id
			flipKickModel.ColumnValues[event.AddressFK] = flipAddressID
			flipKickModel.ColumnValues[event.LogFK] = logID
		})

		It("gets the relevant ilk for event", func() {
			expectedEvent := expectedBidEvent(flipKickModel, "kick", flipAddress, headerOne.BlockNumber)
			expectedEvent.IlkIdentifier = test_helpers.FakeIlk.Identifier

			flipIlkErr := flipRepo.Create(diffID, headerOne.Id, flip.IlkMetadata, test_helpers.FakeIlk.Hex)
			Expect(flipIlkErr).NotTo(HaveOccurred())

			insertErr := event.PersistModels([]event.InsertionModel{flipKickModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			var bidEvents []bidEvent
			queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, ilk_identifier, block_height FROM maker.bid_event`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(bidEvents).To(ConsistOf(expectedEvent))
		})

		It("gets the relevant urn for event", func() {
			expectedEvent := expectedBidEvent(flipKickModel, "kick", flipAddress, headerOne.BlockNumber)
			expectedEvent.UrnIdentifier = common.HexToAddress("0x" + test_data.RandomString(40)).Hex()

			bidUsrMetadata := types.GetValueMetadata(storage.BidUsr,
				map[types.Key]string{constants.BidId: expectedEvent.BidID}, types.Address)
			usrErr := flipRepo.Create(diffID, headerOne.Id, bidUsrMetadata, expectedEvent.UrnIdentifier)
			Expect(usrErr).NotTo(HaveOccurred())

			insertErr := event.PersistModels([]event.InsertionModel{flipKickModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())

			var bidEvents []bidEvent
			queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, urn_identifier, block_height FROM maker.bid_event`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(bidEvents).To(ConsistOf(expectedEvent))
		})
	})

	Describe("when flip-specific diffs are inserted after events", func() {
		var (
			flipAddress   string
			flipRepo      flip.StorageRepository
			diffID        int64
			flipKickModel event.InsertionModel
		)

		BeforeEach(func() {
			flipAddress = test_data.FlipEthAddress()
			flipAddressID, addressErr := shared.GetOrCreateAddress(flipAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
			flipRepo = flip.StorageRepository{ContractAddress: flipAddress}
			flipRepo.SetDB(db)
			diffID = CreateFakeDiffRecord(db)
			flipKickModel = test_data.FlipKickModel()
			flipKickModel.ColumnValues[event.HeaderFK] = headerOne.Id
			flipKickModel.ColumnValues[event.AddressFK] = flipAddressID
			flipKickModel.ColumnValues[event.LogFK] = logID
			insertErr := event.PersistModels([]event.InsertionModel{flipKickModel}, db)
			Expect(insertErr).NotTo(HaveOccurred())
		})

		Specify("inserting a flip_ilk diff triggers an update to the ilk_identifier of related bids", func() {
			var initialIlks []sql.NullString
			initialQueryErr := db.Select(&initialIlks, `SELECT ilk_identifier FROM maker.bid_event`)
			Expect(initialQueryErr).NotTo(HaveOccurred())
			Expect(initialIlks).To(ConsistOf(sql.NullString{Valid: false}))
			expectedEvent := expectedBidEvent(flipKickModel, "kick", flipAddress, headerOne.BlockNumber)
			expectedEvent.IlkIdentifier = test_helpers.FakeIlk.Identifier

			flipIlkErr := flipRepo.Create(diffID, headerOne.Id, flip.IlkMetadata, test_helpers.FakeIlk.Hex)
			Expect(flipIlkErr).NotTo(HaveOccurred())

			var bidEvents []bidEvent
			queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, ilk_identifier, block_height FROM maker.bid_event`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(bidEvents).To(ConsistOf(expectedEvent))
		})

		Specify("inserting a flip_bid_usr diff triggers an update to the urn_identifier of related bids", func() {
			var initialUrns []sql.NullString
			initialQueryErr := db.Select(&initialUrns, `SELECT urn_identifier FROM maker.bid_event`)
			Expect(initialQueryErr).NotTo(HaveOccurred())
			Expect(initialUrns).To(ConsistOf(sql.NullString{Valid: false}))
			usr := common.HexToAddress("0x" + test_data.RandomString(40)).Hex()
			expectedEvent := expectedBidEvent(flipKickModel, "kick", flipAddress, headerOne.BlockNumber)
			bidUsrMetadata := types.GetValueMetadata(storage.BidUsr,
				map[types.Key]string{constants.BidId: expectedEvent.BidID}, types.Address)
			expectedEvent.UrnIdentifier = usr

			usrErr := flipRepo.Create(diffID, headerOne.Id, bidUsrMetadata, usr)
			Expect(usrErr).NotTo(HaveOccurred())

			var bidEvents []bidEvent
			queryErr := db.Select(&bidEvents, `SELECT log_id, bid_id, contract_address, act, lot, bid_amount, urn_identifier, block_height FROM maker.bid_event`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(bidEvents).To(ConsistOf(expectedEvent))
		})
	})

	Describe("deleting records", func() {
		var (
			headerTwo          core.Header
			flipAddress        string
			usrOne, usrTwo     string
			flipRepo           flip.StorageRepository
			bidOneID, bidTwoID int
			diffID,
			logTwoID,
			logThreeID int64
			flipKickModelOne,
			tendModel,
			flipKickModelTwo event.InsertionModel
		)

		BeforeEach(func() {
			headerTwo = CreateHeader(int64(timestampOne+1), blockOne+1, db)
			logTwoID = test_data.CreateTestLog(headerTwo.Id, db).ID
			logThreeID = test_data.CreateTestLog(headerTwo.Id, db).ID

			flipAddress = test_data.FlipEthAddress()
			ethFlipAddressID, ethFlipAddressErr := shared.GetOrCreateAddress(flipAddress, db)
			Expect(ethFlipAddressErr).NotTo(HaveOccurred())

			flipRepo = flip.StorageRepository{ContractAddress: flipAddress}
			flipRepo.SetDB(db)
			diffID = CreateFakeDiffRecord(db)

			bidOneID = rand.Int()
			bidTwoID = bidOneID + 1
			usrOne = common.HexToAddress("0x" + test_data.RandomString(40)).Hex()
			usrTwo = common.HexToAddress("0x" + test_data.RandomString(40)).Hex()

			flipIlkErr := flipRepo.Create(diffID, headerOne.Id, flip.IlkMetadata, test_helpers.FakeIlk.Hex)
			Expect(flipIlkErr).NotTo(HaveOccurred())

			bidUsrMetadataOne := types.GetValueMetadata(storage.BidUsr,
				map[types.Key]string{constants.BidId: strconv.Itoa(bidOneID)}, types.Address)
			usrOneErr := flipRepo.Create(diffID, headerOne.Id, bidUsrMetadataOne, usrOne)
			Expect(usrOneErr).NotTo(HaveOccurred())

			bidUsrMetadataTwo := types.GetValueMetadata(storage.BidUsr,
				map[types.Key]string{constants.BidId: strconv.Itoa(bidTwoID)}, types.Address)
			usrTwoErr := flipRepo.Create(diffID, headerOne.Id, bidUsrMetadataTwo, usrTwo)
			Expect(usrTwoErr).NotTo(HaveOccurred())

			flipKickModelOne = test_data.FlipKickModel()
			flipKickModelOne.ColumnValues[event.HeaderFK] = headerOne.Id
			flipKickModelOne.ColumnValues[event.AddressFK] = ethFlipAddressID
			flipKickModelOne.ColumnValues[event.LogFK] = logID
			flipKickModelOne.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidOneID)
			insertKickErrOne := event.PersistModels([]event.InsertionModel{flipKickModelOne}, db)
			Expect(insertKickErrOne).NotTo(HaveOccurred())

			tendLog := test_data.CreateTestLogFromEventLog(headerOne.Id, test_data.TendEventLog.Log, db)
			tendModel = test_data.TendModel()
			tendModel.ColumnValues[event.HeaderFK] = headerOne.Id
			tendModel.ColumnValues[event.AddressFK] = ethFlipAddressID
			tendModel.ColumnValues[event.LogFK] = logTwoID
			tendModel.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidOneID)
			test_data.AssignMessageSenderID(tendLog, tendModel, db)
			insertTendErr := event.PersistModels([]event.InsertionModel{tendModel}, db)
			Expect(insertTendErr).NotTo(HaveOccurred())

			flipKickModelTwo = test_data.FlipKickModel()
			flipKickModelTwo.ColumnValues[event.HeaderFK] = headerTwo.Id
			flipKickModelTwo.ColumnValues[event.AddressFK] = ethFlipAddressID
			flipKickModelTwo.ColumnValues[event.LogFK] = logThreeID
			flipKickModelTwo.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidTwoID)
			insertKickErrTwo := event.PersistModels([]event.InsertionModel{flipKickModelTwo}, db)
			Expect(insertKickErrTwo).NotTo(HaveOccurred())
		})

		Specify("deleting a flip_ilk sets corresponding events' ilk_identifier to null", func() {
			var initialIlks []string
			initialIlkErr := db.Select(&initialIlks, `SELECT ilk_identifier FROM maker.bid_event`)
			Expect(initialIlkErr).NotTo(HaveOccurred())
			ilkIdentifier := test_helpers.FakeIlk.Identifier
			Expect(initialIlks).To(ConsistOf(ilkIdentifier, ilkIdentifier, ilkIdentifier))

			flipAddressID, flipAddressErr := shared.GetOrCreateAddress(flipAddress, db)
			Expect(flipAddressErr).NotTo(HaveOccurred())

			_, deleteIlkErr := db.Exec(`DELETE FROM maker.flip_ilk WHERE address_id = $1`, flipAddressID)
			Expect(deleteIlkErr).NotTo(HaveOccurred())

			var eventIlks []sql.NullString
			ilkErr := db.Select(&eventIlks, `SELECT ilk_identifier FROM maker.bid_event`)
			Expect(ilkErr).NotTo(HaveOccurred())
			nullString := sql.NullString{Valid: false}
			Expect(eventIlks).To(ConsistOf(nullString, nullString, nullString))
		})

		Specify("deleting a flip_bid_usr sets corresponding events' urn_identifier to null", func() {
			var initialUrns []string
			initialUrnErr := db.Select(&initialUrns, `SELECT urn_identifier FROM maker.bid_event`)
			Expect(initialUrnErr).NotTo(HaveOccurred())
			Expect(initialUrns).To(ConsistOf(usrOne, usrOne, usrTwo))

			_, deleteUrnErr := db.Exec(`DELETE FROM maker.flip_bid_usr WHERE bid_id = $1`, bidTwoID)
			Expect(deleteUrnErr).NotTo(HaveOccurred())

			var eventUrns []sql.NullString
			urnErr := db.Select(&eventUrns, `SELECT urn_identifier FROM maker.bid_event`)
			Expect(urnErr).NotTo(HaveOccurred())
			expectedUsrOne := test_helpers.GetValidNullString(usrOne)
			Expect(eventUrns).To(ConsistOf(expectedUsrOne, expectedUsrOne, sql.NullString{Valid: false}))
		})
	})
})

func expectedBidEvent(eventModel event.InsertionModel, bidAct, contractAddress string, blockHeight int64) bidEvent {
	return bidEvent{
		LogID:           strconv.FormatInt(eventModel.ColumnValues[event.LogFK].(int64), 10),
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
		LogID:           strconv.FormatInt(eventModel.ColumnValues[event.LogFK].(int64), 10),
		BidID:           eventModel.ColumnValues[constants.BidIDColumn].(string),
		ContractAddress: contractAddress,
		Act:             bidAct,
		Lot:             sql.NullString{Valid: false},
		Bid:             sql.NullString{Valid: false},
		BlockHeight:     strconv.FormatInt(blockHeight, 10),
	}
}

type bidEvent struct {
	LogID           string `db:"log_id"`
	BidID           string `db:"bid_id"`
	ContractAddress string `db:"contract_address"`
	Act             string
	Lot             sql.NullString
	Bid             sql.NullString `db:"bid_amount"`
	IlkIdentifier   string         `db:"ilk_identifier"`
	UrnIdentifier   string         `db:"urn_identifier"`
	BlockHeight     string         `db:"block_height"`
}
