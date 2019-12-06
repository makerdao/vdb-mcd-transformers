package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flop_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flop bid events query", func() {
	var (
		db                     *postgres.DB
		flopKickRepo           flop_kick.FlopKickRepository
		headerRepo             repositories.HeaderRepository
		blockOne, timestampOne int
		headerOne              core.Header
		contractAddress        string
		fakeBidId              int
		flopKickEvent          shared.InsertionModel
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flopKickRepo = flop_kick.FlopKickRepository{}
		flopKickRepo.SetDB(db)

		fakeBidId = rand.Int()
		contractAddress = "0x763ztv6x68exwqrgtl325e7hrcvavid4e3fcb4g"

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)

		flopKickLog := test_data.CreateTestLog(headerOne.Id, db)

		flopKickEvent = test_data.FlopKickModel()
		flopKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
		flopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
		flopKickEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
		flopKickEvent.ColumnValues[constants.LogFK] = flopKickLog.ID
		flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEvent})
		Expect(flopKickErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("all_flop_bid_events", func() {
		It("returns all flop bid events", func() {
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()

			flopDentLog := test_data.CreateTestLog(headerOne.Id, db)
			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				BidId:           fakeBidId,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				DentHeaderId:    headerOne.Id,
				DentLogId:       flopDentLog.ID,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

			flopDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealHeaderId:    headerTwo.Id,
			})
			Expect(flopDealErr).NotTo(HaveOccurred())

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerTwo, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(fakeBidAmount), Lot: strconv.Itoa(fakeLot), Act: "dent"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: flopStorageValues[storage.BidBid].(string), Lot: flopStorageValues[storage.BidLot].(string), Act: "deal"},
			))
		})

		It("returns bid events from floppers that have different bid ids", func() {
			lotOne := rand.Int()
			bidAmountOne := rand.Int()

			bidIdTwo := rand.Int()
			lotTwo := rand.Int()
			bidAmountTwo := rand.Int()

			flopKickEventTwoLog := test_data.CreateTestLog(headerOne.Id, db)

			flopKickEventTwo := test_data.FlopKickModel()
			flopKickEventTwo.ForeignKeyValues[constants.AddressFK] = contractAddress
			flopKickEventTwo.ColumnValues["bid_id"] = strconv.Itoa(bidIdTwo)
			flopKickEventTwo.ColumnValues[constants.HeaderFK] = headerOne.Id
			flopKickEventTwo.ColumnValues[constants.LogFK] = flopKickEventTwoLog.ID
			flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEventTwo})

			Expect(flopKickErr).NotTo(HaveOccurred())

			flopDentLog := test_data.CreateTestLog(headerOne.Id, db)
			flopDentOneErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             lotOne,
				BidAmount:       bidAmountOne,
				DentHeaderId:    headerOne.Id,
				DentLogId:       flopDentLog.ID,
			})
			Expect(flopDentOneErr).NotTo(HaveOccurred())

			flopDentTwoLog := test_data.CreateTestLog(headerOne.Id, db)
			flopDentTwoErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				DB:              db,
				BidId:           bidIdTwo,
				ContractAddress: contractAddress,
				Lot:             lotTwo,
				BidAmount:       bidAmountTwo,
				DentHeaderId:    headerOne.Id,
				DentLogId:       flopDentTwoLog.ID,
			})
			Expect(flopDentTwoErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidId: flopKickEventTwo.ColumnValues["bid_id"].(string), BidAmount: flopKickEventTwo.ColumnValues["bid"].(string), Lot: flopKickEventTwo.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(bidAmountOne), Lot: strconv.Itoa(lotOne), Act: "dent"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidIdTwo), BidAmount: strconv.Itoa(bidAmountTwo), Lot: strconv.Itoa(lotTwo), Act: "dent"},
			))
		})

		It("ignores bid events from flaps", func() {
			flapKickLog := test_data.CreateTestLog(headerOne.Id, db)

			addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
			flapKickEvent := test_data.FlapKickModel()
			flapKickEvent.ColumnValues[event.AddressFK] = addressId
			flapKickEvent.ColumnValues[constants.BidIdColumn] = strconv.Itoa(fakeBidId)
			flapKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
			flapKickEvent.ColumnValues[event.LogFK] = flapKickLog.ID
			flapKickErr := event.PersistModels([]event.InsertionModel{flapKickEvent}, db)
			Expect(flapKickErr).NotTo(HaveOccurred())

			flapKickBidEvent := test_helpers.BidEvent{
				BidId:     flapKickEvent.ColumnValues["bid_id"].(string),
				BidAmount: flapKickEvent.ColumnValues["bid"].(string),
				Lot:       flapKickEvent.ColumnValues["lot"].(string),
				Act:       "kick"}

			var flopBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&flopBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(flopBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
			))
			Expect(flopBidEvents).NotTo(ContainElement(flapKickBidEvent))
		})
	})

	Describe("dent", func() {
		It("returns flop dent bid events from multiple blocks", func() {
			lot := rand.Int()
			bidAmount := rand.Int()
			updatedLot := lot + 100
			updatedBidAmount := bidAmount + 100

			flopDentHeaderOneLog := test_data.CreateTestLog(headerOne.Id, db)
			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             lot,
				BidAmount:       bidAmount,
				DentHeaderId:    headerOne.Id,
				DentLogId:       flopDentHeaderOneLog.ID,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			flopDentHeaderTwoLog := test_data.CreateTestLog(headerTwo.Id, db)

			flopDentHeaderTwoErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             updatedLot,
				BidAmount:       updatedBidAmount,
				DentHeaderId:    headerTwo.Id,
				DentLogId:       flopDentHeaderTwoLog.ID,
			})
			Expect(flopDentHeaderTwoErr).NotTo(HaveOccurred())

			headerThree := createHeader(blockOne+2, timestampOne+2, headerRepo)
			flapDentLog := test_data.CreateTestLog(headerThree.Id, db)

			// create irrelevant flap dent
			flapDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: "flap contract address",
				Lot:             lot,
				BidAmount:       bidAmount,
				DentHeaderId:    headerThree.Id,
				DentLogId:       flapDentLog.ID,
			})
			Expect(flapDentErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(bidAmount), Lot: strconv.Itoa(lot), Act: "dent"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(updatedBidAmount), Lot: strconv.Itoa(updatedLot), Act: "dent"},
			))
		})
	})

	Describe("Deal", func() {
		It("returns bid events with lot and bid amount values from the block where the deal occurred", func() {
			fakeBidId := rand.Int()

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

			updatedFlopStorageValues := test_helpers.GetFlopStorageValues(2, fakeBidId)
			test_helpers.CreateFlop(db, headerTwo, updatedFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerThree := createHeader(blockOne+2, timestampOne+2, headerRepo)

			flopDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealHeaderId:    headerThree.Id,
			})
			Expect(flopDealErr).NotTo(HaveOccurred())

			dealBlockFlopStorageValues := test_helpers.GetFlopStorageValues(0, fakeBidId)
			test_helpers.CreateFlop(db, headerThree, dealBlockFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: dealBlockFlopStorageValues[storage.BidBid].(string), Lot: dealBlockFlopStorageValues[storage.BidLot].(string), Act: "deal"},
				test_helpers.BidEvent{BidId: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"}))
		})
	})

	Describe("Yank event", func() {
		It("includes yank in all flop bid events", func() {
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()

			dentLog := test_data.CreateTestLog(headerOne.Id, db)
			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				DentHeaderId:    headerOne.Id,
				DentLogId:       dentLog.ID,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			flopYankLog := test_data.CreateTestLog(headerOne.Id, db)

			flopYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				YankHeaderId:    headerTwo.Id,
				YankLogId:       flopYankLog.ID,
			})
			Expect(flopYankErr).NotTo(HaveOccurred())

			updatedFlopStorageValues := test_helpers.GetFlopStorageValues(2, fakeBidId)
			test_helpers.CreateFlop(db, headerTwo, updatedFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(fakeBidAmount), Lot: strconv.Itoa(fakeLot), Act: "dent"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: updatedFlopStorageValues[storage.BidBid].(string), Lot: updatedFlopStorageValues[storage.BidLot].(string), Act: "yank"},
			))
		})

		It("ignores flap yank events", func() {
			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			flapYankLog := test_data.CreateTestLog(headerTwo.Id, db)

			// irrelevant flap yank
			flapYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: "flap contract address",
				YankHeaderId:    headerTwo.Id,
				YankLogId:       flapYankLog.ID,
			})
			Expect(flapYankErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
			))
		})
	})

	Describe("tick event", func() {
		It("ignores tick events from non flop contracts", func() {
			fakeBidId := rand.Int()
			tickLog := test_data.CreateTestLog(headerOne.Id, db)

			// irrelevant tick event
			tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: "flip",
				TickHeaderId:    headerOne.Id,
				TickLogId:       tickLog.ID,
			})
			Expect(tickErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
			))
		})

		It("includes flop tick bid events", func() {
			fakeBidId := rand.Int()
			tickLog := test_data.CreateTestLog(headerOne.Id, db)

			tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				TickHeaderId:    headerOne.Id,
				TickLogId:       tickLog.ID,
			})
			Expect(tickErr).NotTo(HaveOccurred())
			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.ColumnValues["bid_id"].(string), BidAmount: flopKickEvent.ColumnValues["bid"].(string), Lot: flopKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: flopStorageValues[storage.BidBid].(string), Lot: flopStorageValues[storage.BidLot].(string), Act: "tick"},
			))
		})
	})

	Describe("result pagination", func() {
		var (
			updatedBidAmount, updatedLot int
			flopKickBlockOne             shared.InsertionModel
		)

		BeforeEach(func() {
			lot := rand.Int()
			bidAmount := rand.Int()
			updatedLot = lot + 100
			updatedBidAmount = bidAmount + 100

			logID := test_data.CreateTestLog(headerOne.Id, db).ID
			flopKickBlockOne = test_data.FlopKickModel()
			flopKickBlockOne.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
			flopKickBlockOne.ForeignKeyValues[constants.AddressFK] = contractAddress
			flopKickBlockOne.ColumnValues[constants.HeaderFK] = headerOne.Id
			flopKickBlockOne.ColumnValues[constants.LogFK] = logID
			flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickBlockOne})
			Expect(flopKickErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			logTwoID := test_data.CreateTestLog(headerTwo.Id, db).ID

			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             updatedLot,
				BidAmount:       updatedBidAmount,
				DentHeaderId:    headerTwo.Id,
				DentLogId:       logTwoID,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())
		})

		It("limits result to latest blocks if max_results argument is provided", func() {
			maxResults := 1
			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events($1)`,
				maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:     strconv.Itoa(fakeBidId),
					BidAmount: strconv.Itoa(updatedBidAmount),
					Lot:       strconv.Itoa(updatedLot),
					Act:       "dent",
				},
			))
		})

		It("offsets results if offset is provided", func() {
			maxResults := 1
			resultOffset := 1
			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events($1, $2)`,
				maxResults, resultOffset)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:     flopKickBlockOne.ColumnValues["bid_id"].(string),
					BidAmount: flopKickBlockOne.ColumnValues["bid"].(string),
					Lot:       flopKickBlockOne.ColumnValues["lot"].(string),
					Act:       "kick",
				},
			))
		})
	})
})
