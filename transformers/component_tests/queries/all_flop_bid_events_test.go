package queries

import (
	"github.com/vulcanize/mcd_transformers/transformers/events/tick"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/dent"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/yank"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Flop bid events query", func() {
	var (
		db              *postgres.DB
		flopKickRepo    flop_kick.FlopKickRepository
		dentRepo        dent.DentRepository
		dealRepo        deal.DealRepository
		yankRepo        yank.YankRepository
		tickRepo        tick.TickRepository
		headerRepo      repositories.HeaderRepository
		blockOne        int64
		headerOne       core.Header
		headerOneId     int64
		headerOneErr    error
		contractAddress string
		fakeBidId       int
		flopKickEvent   flop_kick.Model
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flopKickRepo = flop_kick.FlopKickRepository{}
		flopKickRepo.SetDB(db)
		dentRepo = dent.DentRepository{}
		dentRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		yankRepo = yank.YankRepository{}
		yankRepo.SetDB(db)
		tickRepo = tick.TickRepository{}
		tickRepo.SetDB(db)

		fakeBidId = rand.Int()
		contractAddress = "0x763ztv6x68exwqrgtl325e7hrcvavid4e3fcb4g"

		blockOne = 1
		headerOne = fakes.GetFakeHeader(blockOne)
		headerOneId, headerOneErr = headerRepo.CreateOrUpdateHeader(headerOne)
		Expect(headerOneErr).NotTo(HaveOccurred())
		flopKickLog := test_data.CreateTestLog(headerOneId, db)

		flopKickEvent = test_data.FlopKickModel
		flopKickEvent.ContractAddress = contractAddress
		flopKickEvent.BidId = strconv.Itoa(fakeBidId)
		flopKickEvent.HeaderID = headerOneId
		flopKickEvent.LogID = flopKickLog.ID
		flopKickErr := flopKickRepo.Create([]interface{}{flopKickEvent})
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

			flopDentLog := test_data.CreateTestLog(headerOneId, db)
			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				ContractAddress: contractAddress,
				BidId:           fakeBidId,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				DentRepo:        dentRepo,
				DentHeaderId:    headerOneId,
				DentLogID:       flopDentLog.ID,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoID, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flopDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderID:    headerTwoID,
			})
			Expect(flopDealErr).NotTo(HaveOccurred())

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerTwo, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.BidId, BidAmount: flopKickEvent.Bid, Lot: flopKickEvent.Lot, Act: "kick"},
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

			flopKickEventTwoLog := test_data.CreateTestLog(headerOneId, db)

			flopKickEventTwo := test_data.FlopKickModel
			flopKickEventTwo.ContractAddress = contractAddress
			flopKickEventTwo.BidId = strconv.Itoa(bidIdTwo)
			flopKickEventTwo.HeaderID = headerOneId
			flopKickEventTwo.LogID = flopKickEventTwoLog.ID
			flopKickErr := flopKickRepo.Create([]interface{}{flopKickEventTwo})
			Expect(flopKickErr).NotTo(HaveOccurred())

			flopDentLog := test_data.CreateTestLog(headerOneId, db)
			flopDentOneErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             lotOne,
				BidAmount:       bidAmountOne,
				DentRepo:        dentRepo,
				DentHeaderId:    headerOneId,
				DentLogID:       flopDentLog.ID,
			})
			Expect(flopDentOneErr).NotTo(HaveOccurred())

			flopDentTwoLog := test_data.CreateTestLog(headerOneId, db)
			flopDentTwoErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           bidIdTwo,
				ContractAddress: contractAddress,
				Lot:             lotTwo,
				BidAmount:       bidAmountTwo,
				DentRepo:        dentRepo,
				DentHeaderId:    headerOneId,
				DentLogID:       flopDentTwoLog.ID,
			})
			Expect(flopDentTwoErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.BidId, BidAmount: flopKickEvent.Bid, Lot: flopKickEvent.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: flopKickEventTwo.BidId, BidAmount: flopKickEventTwo.Bid, Lot: flopKickEventTwo.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(bidAmountOne), Lot: strconv.Itoa(lotOne), Act: "dent"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidIdTwo), BidAmount: strconv.Itoa(bidAmountTwo), Lot: strconv.Itoa(lotTwo), Act: "dent"},
			))
		})

		It("ignores bid events from flaps", func() {
			insertedLog := test_data.CreateTestLog(headerOneId, db)

			flapKickRepo := flap_kick.FlapKickRepository{}
			flapKickRepo.SetDB(db)
			flapKickEvent := test_data.FlapKickModel
			flapKickEvent.BidId = strconv.Itoa(fakeBidId)
			flapKickEvent.HeaderID = headerOneId
			flapKickEvent.LogID = insertedLog.ID
			flapKickErr := flapKickRepo.Create([]interface{}{flapKickEvent})
			Expect(flapKickErr).NotTo(HaveOccurred())

			flapKickBidEvent := test_helpers.BidEvent{BidId: flapKickEvent.BidId, BidAmount: flapKickEvent.Bid, Lot: flapKickEvent.Lot, Act: "kick"}

			var flopBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&flopBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(flopBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.BidId, BidAmount: flopKickEvent.Bid, Lot: flopKickEvent.Lot, Act: "kick"},
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

			flopDentHeaderOneLog := test_data.CreateTestLog(headerOneId, db)
			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             lot,
				BidAmount:       bidAmount,
				DentRepo:        dentRepo,
				DentHeaderId:    headerOneId,
				DentLogID:       flopDentHeaderOneLog.ID,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			flopDentHeaderTwoLog := test_data.CreateTestLog(headerTwoId, db)

			flopDentHeaderTwoErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             updatedLot,
				BidAmount:       updatedBidAmount,
				DentRepo:        dentRepo,
				DentHeaderId:    headerTwoId,
				DentLogID:       flopDentHeaderTwoLog.ID,
			})
			Expect(flopDentHeaderTwoErr).NotTo(HaveOccurred())

			headerThree := fakes.GetFakeHeaderWithTimestamp(int64(333333333), 3)
			headerThreeId, headerThreeErr := headerRepo.CreateOrUpdateHeader(headerThree)
			Expect(headerThreeErr).NotTo(HaveOccurred())
			flapDentLog := test_data.CreateTestLog(headerThreeId, db)

			// create irrelevant flap dent
			flapDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           fakeBidId,
				ContractAddress: "flap contract address",
				Lot:             lot,
				BidAmount:       bidAmount,
				DentRepo:        dentRepo,
				DentHeaderId:    headerThreeId,
				DentLogID:       flapDentLog.ID,
			})
			Expect(flapDentErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.BidId, BidAmount: flopKickEvent.Bid, Lot: flopKickEvent.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(bidAmount), Lot: strconv.Itoa(lot), Act: "dent"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(updatedBidAmount), Lot: strconv.Itoa(updatedLot), Act: "dent"},
			))
		})
	})

	Describe("Deal", func() {
		It("returns bid events with lot and bid amount values from the block where the deal occurred", func() {
			fakeBidId := rand.Int()
			blockOne := rand.Int()
			blockTwo := blockOne + 1
			blockThree := blockTwo + 1

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := fakes.GetFakeHeader(int64(blockTwo))
			_, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			updatedFlopStorageValues := test_helpers.GetFlopStorageValues(2, fakeBidId)
			test_helpers.CreateFlop(db, headerTwo, updatedFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerThree := fakes.GetFakeHeader(int64(blockThree))
			headerThreeId, headerThreeErr := headerRepo.CreateOrUpdateHeader(headerThree)
			Expect(headerThreeErr).NotTo(HaveOccurred())

			flopDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderID:    headerThreeId,
			})
			Expect(flopDealErr).NotTo(HaveOccurred())

			dealBlockFlopStorageValues := test_helpers.GetFlopStorageValues(0, fakeBidId)
			test_helpers.CreateFlop(db, headerThree, dealBlockFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: dealBlockFlopStorageValues[storage.BidBid].(string), Lot: dealBlockFlopStorageValues[storage.BidLot].(string), Act: "deal"},
				test_helpers.BidEvent{BidId: flopKickEvent.BidId, BidAmount: flopKickEvent.Bid, Lot: flopKickEvent.Lot, Act: "kick"}))
		})
	})

	Describe("Yank event", func() {
		It("includes yank in all flop bid events", func() {
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()

			dentLog := test_data.CreateTestLog(headerOneId, db)
			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				DentRepo:        dentRepo,
				DentHeaderId:    headerOneId,
				DentLogID:       dentLog.ID,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			insertedLogTwo := test_data.CreateTestLog(headerOneId, db)

			flopYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				YankRepo:        yankRepo,
				YankHeaderId:    headerTwoId,
				YankLogID:       insertedLogTwo.ID,
			})
			Expect(flopYankErr).NotTo(HaveOccurred())

			updatedFlopStorageValues := test_helpers.GetFlopStorageValues(2, fakeBidId)
			test_helpers.CreateFlop(db, headerTwo, updatedFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.BidId, BidAmount: flopKickEvent.Bid, Lot: flopKickEvent.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(fakeBidAmount), Lot: strconv.Itoa(fakeLot), Act: "dent"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: updatedFlopStorageValues[storage.BidBid].(string), Lot: updatedFlopStorageValues[storage.BidLot].(string), Act: "yank"},
			))
		})

		It("ignores flap yank events", func() {
			fakeBidId := rand.Int()

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			insertedLogTwo := test_data.CreateTestLog(headerTwoId, db)

			// irrelevant flap yank
			flapYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidId:           fakeBidId,
				ContractAddress: "flap contract address",
				YankRepo:        yankRepo,
				YankHeaderId:    headerTwoId,
				YankLogID:       insertedLogTwo.ID,
			})
			Expect(flapYankErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.BidId, BidAmount: flopKickEvent.Bid, Lot: flopKickEvent.Lot, Act: "kick"},
			))
		})
	})

	Describe("tick event", func() {
		It("ignores tick events from non flop contracts", func() {
			fakeBidId := rand.Int()
			tickLog := test_data.CreateTestLog(headerOneId, db)

			// irrelevant tick event
			tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				BidId:           fakeBidId,
				ContractAddress: "flip",
				TickRepo:        tickRepo,
				TickHeaderId:    headerOneId,
				TickLogId:       tickLog.ID,
			})
			Expect(tickErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.BidId, BidAmount: flopKickEvent.Bid, Lot: flopKickEvent.Lot, Act: "kick"},
			))
		})

		It("includes flop tick bid events", func() {
			fakeBidId := rand.Int()
			tickLog := test_data.CreateTestLog(headerOneId, db)

			tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				TickRepo:        tickRepo,
				TickHeaderId:    headerOneId,
				TickLogId:       tickLog.ID,
			})
			Expect(tickErr).NotTo(HaveOccurred())
			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEvent.BidId, BidAmount: flopKickEvent.Bid, Lot: flopKickEvent.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: flopStorageValues[storage.BidBid].(string), Lot: flopStorageValues[storage.BidLot].(string), Act: "tick"},
			))
		})
	})

	Describe("result pagination", func() {
		var (
			updatedBidAmount, updatedLot int
			flopKickBlockOne             flop_kick.Model
		)

		BeforeEach(func() {
			lot := rand.Int()
			bidAmount := rand.Int()
			updatedLot = lot + 100
			updatedBidAmount = bidAmount + 100

			logID := test_data.CreateTestLog(headerOneId, db).ID
			flopKickBlockOne = test_data.FlopKickModel
			flopKickBlockOne.BidId = strconv.Itoa(fakeBidId)
			flopKickBlockOne.ContractAddress = contractAddress
			flopKickBlockOne.HeaderID = headerOneId
			flopKickBlockOne.LogID = logID
			flopKickErr := flopKickRepo.Create([]interface{}{flopKickBlockOne})
			Expect(flopKickErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			logTwoID := test_data.CreateTestLog(headerTwoId, db).ID

			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             updatedLot,
				BidAmount:       updatedBidAmount,
				DentRepo:        dentRepo,
				DentHeaderId:    headerTwoId,
				DentLogID:       logTwoID,
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
					BidId:     flopKickBlockOne.BidId,
					BidAmount: flopKickBlockOne.Bid,
					Lot:       flopKickBlockOne.Lot,
					Act:       "kick",
				},
			))
		})
	})
})
