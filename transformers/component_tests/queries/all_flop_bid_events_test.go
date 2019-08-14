package queries

import (
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
	"math/rand"
	"strconv"
	"time"
)

var _ = Describe("Flop bid events query", func() {
	var (
		db           *postgres.DB
		flopKickRepo flop_kick.FlopKickRepository
		dentRepo        dent.DentRepository
		dealRepo        deal.DealRepository
		yankRepo        yank.YankRepository
		headerRepo      repositories.HeaderRepository
		contractAddress = "flop contract address"
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
		rand.Seed(time.Now().UnixNano())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("all_flop_bid_events", func() {
		It("returns all flop bid events", func() {
			fakeBidId := rand.Int()
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()

			header := fakes.GetFakeHeader(1)
			headerId, headerErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			flopKickEvent := test_data.FlopKickModel
			flopKickEvent.ContractAddress = contractAddress
			flopKickEvent.BidId = strconv.Itoa(fakeBidId)
			flopKickErr := flopKickRepo.Create(headerId, []interface{}{flopKickEvent})
			Expect(flopKickErr).NotTo(HaveOccurred())

			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				ContractAddress: contractAddress,
				BidId:           fakeBidId,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				DentRepo:        dentRepo,
				DentHeaderId:    headerId,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flopDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderId:    headerTwoId,
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

		It("returns bid events from multiple floppers", func() {
			bidIdOne := rand.Int()
			lotOne := rand.Int()
			bidAmountOne := rand.Int()

			bidIdTwo := rand.Int()
			lotTwo := rand.Int()
			bidAmountTwo := rand.Int()

			header := fakes.GetFakeHeader(1)
			headerId, headerErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			flopKickEventOne := test_data.FlopKickModel
			flopKickEventOne.ContractAddress = contractAddress
			flopKickEventOne.BidId = strconv.Itoa(bidIdOne)
			flopKickErr := flopKickRepo.Create(headerId, []interface{}{flopKickEventOne})
			Expect(flopKickErr).NotTo(HaveOccurred())

			flopKickEventTwo := test_data.FlopKickModel
			flopKickEventTwo.ContractAddress = contractAddress
			flopKickEventTwo.BidId = strconv.Itoa(bidIdTwo)
			flopKickEventTwo.TransactionIndex = 11
			flopKickEventTwo.LogIndex = 12
			flopKickErr = flopKickRepo.Create(headerId, []interface{}{flopKickEventTwo})
			Expect(flopKickErr).NotTo(HaveOccurred())

			flopDentOneErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           bidIdOne,
				ContractAddress: contractAddress,
				Lot:             lotOne,
				BidAmount:       bidAmountOne,
				DentRepo:        dentRepo,
				DentHeaderId:    headerId,
			})
			Expect(flopDentOneErr).NotTo(HaveOccurred())

			flopDentTwoErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           bidIdTwo,
				ContractAddress: contractAddress,
				Lot:             lotTwo,
				BidAmount:       bidAmountTwo,
				TxIndex:         21,
				LogIndex:        22,
				DentRepo:        dentRepo,
				DentHeaderId:    headerId,
			})
			Expect(flopDentTwoErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flopKickEventOne.BidId, BidAmount: flopKickEventOne.Bid, Lot: flopKickEventOne.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: flopKickEventTwo.BidId, BidAmount: flopKickEventTwo.Bid, Lot: flopKickEventTwo.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidIdOne), BidAmount: strconv.Itoa(bidAmountOne), Lot: strconv.Itoa(lotOne), Act: "dent"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidIdTwo), BidAmount: strconv.Itoa(bidAmountTwo), Lot: strconv.Itoa(lotTwo), Act: "dent"},
			))
		})

		It("ignores bid events from flaps", func() {
			fakeBidId := rand.Int()
			header := fakes.GetFakeHeader(1)
			headerId, headerErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			flapKickRepo := flap_kick.FlapKickRepository{}
			flapKickRepo.SetDB(db)
			flapKickEvent := test_data.FlapKickModel
			flapKickEvent.BidId = strconv.Itoa(fakeBidId)
			flapKickErr := flapKickRepo.Create(headerId, []interface{}{flapKickEvent})
			Expect(flapKickErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(len(actualBidEvents)).To(Equal(0))
		})
	})

	Describe("dent", func() {
		It("returns flop dent bid events from multiple blocks", func() {
			fakeBidId := rand.Int()
			lot := rand.Int()
			bidAmount := rand.Int()
			updatedLot := lot + 100
			updatedBidAmount := bidAmount + 100

			headerOne := fakes.GetFakeHeaderWithTimestamp(int64(111111111), 1)
			headerOneId, headerOneErr := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(headerOneErr).NotTo(HaveOccurred())

			flopKickBlockOne := test_data.FlopKickModel
			flopKickBlockOne.BidId = strconv.Itoa(fakeBidId)
			flopKickBlockOne.ContractAddress = contractAddress
			flopKickErr := flopKickRepo.Create(headerOneId, []interface{}{flopKickBlockOne})
			Expect(flopKickErr).NotTo(HaveOccurred())

			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             lot,
				BidAmount:       bidAmount,
				DentRepo:        dentRepo,
				DentHeaderId:    headerOneId,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flopDentErr = test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             updatedLot,
				BidAmount:       updatedBidAmount,
				DentRepo:        dentRepo,
				DentHeaderId:    headerTwoId,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			headerThree := fakes.GetFakeHeaderWithTimestamp(int64(333333333), 3)
			headerThreeId, headerThreeErr := headerRepo.CreateOrUpdateHeader(headerThree)
			Expect(headerThreeErr).NotTo(HaveOccurred())

			// create irrelevant flap tend
			flapDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           fakeBidId,
				ContractAddress: "flap contract address",
				Lot:             lot,
				BidAmount:       bidAmount,
				DentRepo:        dentRepo,
				DentHeaderId:    headerThreeId,
			})
			Expect(flapDentErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(bidAmount), Lot: strconv.Itoa(lot), Act: "dent"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(updatedBidAmount), Lot: strconv.Itoa(updatedLot), Act: "dent"},
				test_helpers.BidEvent{BidId: flopKickBlockOne.BidId, BidAmount: flopKickBlockOne.Bid, Lot: flopKickBlockOne.Lot, Act: "kick"},
			))
		})
	})

	Describe("Deal", func() {
		It("returns bid events with lot and bid amount values from the block where the deal occurred", func() {
			fakeBidId := rand.Int()
			blockOne := rand.Int()
			blockTwo := blockOne + 1
			blockThree := blockTwo + 1

			headerOne := fakes.GetFakeHeader(int64(blockOne))
			headerOneId, headerOneErr := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(headerOneErr).NotTo(HaveOccurred())

			flopKickBlockOne := test_data.FlopKickModel
			flopKickBlockOne.BidId = strconv.Itoa(fakeBidId)
			flopKickBlockOne.ContractAddress = contractAddress
			flopKickErr := flopKickRepo.Create(headerOneId, []interface{}{flopKickBlockOne})
			Expect(flopKickErr).NotTo(HaveOccurred())

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
				DealHeaderId:    headerThreeId,
			})
			Expect(flopDealErr).NotTo(HaveOccurred())

			dealBlockFlopStorageValues := test_helpers.GetFlopStorageValues(0, fakeBidId)
			test_helpers.CreateFlop(db, headerThree, dealBlockFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flop_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: dealBlockFlopStorageValues[storage.BidBid].(string), Lot: dealBlockFlopStorageValues[storage.BidLot].(string), Act: "deal"},
				test_helpers.BidEvent{BidId: flopKickBlockOne.BidId, BidAmount: flopKickBlockOne.Bid, Lot: flopKickBlockOne.Lot, Act: "kick"}))
		})
	})

	Describe("Yank event", func() {
		It("includes yank in all flop bid events", func() {
			fakeBidId := rand.Int()
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()

			header := fakes.GetFakeHeader(1)
			headerId, headerErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			flopKickEvent := test_data.FlopKickModel
			flopKickEvent.ContractAddress = contractAddress
			flopKickEvent.BidId = strconv.Itoa(fakeBidId)
			flopKickErr := flopKickRepo.Create(headerId, []interface{}{flopKickEvent})
			Expect(flopKickErr).NotTo(HaveOccurred())

			flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				DentRepo:        dentRepo,
				DentHeaderId:    headerId,
			})
			Expect(flopDentErr).NotTo(HaveOccurred())

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, header, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flopYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				YankRepo:        yankRepo,
				YankHeaderId:    headerTwoId,
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

			header := fakes.GetFakeHeader(1)
			headerId, headerErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			flopKickEvent := test_data.FlopKickModel
			flopKickEvent.ContractAddress = contractAddress
			flopKickEvent.BidId = strconv.Itoa(fakeBidId)
			flopKickErr := flopKickRepo.Create(headerId, []interface{}{flopKickEvent})
			Expect(flopKickErr).NotTo(HaveOccurred())

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, header, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			// irrelevant flap yank
			flapYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidId:           fakeBidId,
				ContractAddress: "flap contract address",
				YankRepo:        yankRepo,
				YankHeaderId:    headerTwoId,
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
})
