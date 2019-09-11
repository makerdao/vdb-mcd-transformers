package queries

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/tend"
	"github.com/vulcanize/mcd_transformers/transformers/events/yank"
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("Flap bid events query", func() {
	var (
		db                     *postgres.DB
		flapKickRepo           flap_kick.FlapKickRepository
		tendRepo               tend.TendRepository
		dealRepo               deal.DealRepository
		yankRepo               yank.YankRepository
		headerRepo             repositories.HeaderRepository
		contractAddress        = fakes.FakeAddress.Hex()
		anotherContractAddress = common.HexToAddress("0xabcdef123456789").Hex()
		blockOne               int64
		headerOne              core.Header
		headerOneId            int64
		headerOneErr           error
		fakeBidId              int
		flapKickEvent          flap_kick.FlapKickModel
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flapKickRepo = flap_kick.FlapKickRepository{}
		flapKickRepo.SetDB(db)
		tendRepo = tend.TendRepository{}
		tendRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		yankRepo = yank.YankRepository{}
		yankRepo.SetDB(db)
		fakeBidId = rand.Int()

		blockOne = 1
		headerOne = fakes.GetFakeHeader(blockOne)
		headerOneId, headerOneErr = headerRepo.CreateOrUpdateHeader(headerOne)
		Expect(headerOneErr).NotTo(HaveOccurred())

		flapKickEvent = test_data.FlapKickModel
		flapKickEvent.ContractAddress = contractAddress
		flapKickEvent.BidId = strconv.Itoa(fakeBidId)
		flapKickErr := flapKickRepo.Create(headerOneId, []interface{}{flapKickEvent})
		Expect(flapKickErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("all_flap_bid_events", func() {
		It("returns all flap bid events (same block)", func() {
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()

			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				TendRepo:        tendRepo,
				TendHeaderId:    headerOneId,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			flapDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderId:    headerOneId,
			})
			Expect(flapDealErr).NotTo(HaveOccurred())

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flapKickEvent.BidId, BidAmount: flapKickEvent.Bid, Lot: flapKickEvent.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(fakeBidAmount), Lot: strconv.Itoa(fakeLot), Act: "tend"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: flapStorageValues[storage.BidBid].(string), Lot: flapStorageValues[storage.BidLot].(string), Act: "deal"},
			))
		})

		It("returns all flap bid events across all blocks", func() {
			fakeBidIdTwo := fakeBidId + 1

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flapKickEventTwo := test_data.FlapKickModel
			flapKickEventTwo.Bid = strconv.Itoa(rand.Int())
			flapKickEventTwo.Lot = strconv.Itoa(rand.Int())
			flapKickEventTwo.ContractAddress = contractAddress
			flapKickEventTwo.BidId = strconv.Itoa(fakeBidIdTwo)
			flapKickErr := flapKickRepo.Create(headerTwoId, []interface{}{flapKickEventTwo})
			Expect(flapKickErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flapKickEvent.BidId, BidAmount: flapKickEvent.Bid, Lot: flapKickEvent.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: flapKickEventTwo.BidId, BidAmount: flapKickEventTwo.Bid, Lot: flapKickEventTwo.Lot, Act: "kick"},
			))
		})

		It("returns bid events for multiple bid ids", func() {
			bidIdOne := fakeBidId
			bidIdTwo := rand.Int()
			lotOne := rand.Int()
			bidAmountOne := rand.Int()

			flapKickEventOne := flapKickEvent
			flapKickEventTwo := test_data.FlapKickModel
			flapKickEventTwo.ContractAddress = contractAddress
			flapKickEventTwo.BidId = strconv.Itoa(bidIdTwo)
			flapKickEventTwo.TransactionIndex = flapKickEventOne.TransactionIndex + 1
			flapKickEventTwo.LogIndex = flapKickEventOne.LogIndex + 1
			flapKickErr := flapKickRepo.Create(headerOneId, []interface{}{flapKickEventTwo})
			Expect(flapKickErr).NotTo(HaveOccurred())

			flapTendOneErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           bidIdOne,
				ContractAddress: contractAddress,
				Lot:             lotOne,
				BidAmount:       bidAmountOne,
				TendRepo:        tendRepo,
				TendHeaderId:    headerOneId,
			})
			Expect(flapTendOneErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flapKickEventOne.BidId, BidAmount: flapKickEventOne.Bid, Lot: flapKickEventOne.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: flapKickEventTwo.BidId, BidAmount: flapKickEventTwo.Bid, Lot: flapKickEventTwo.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidIdOne), BidAmount: strconv.Itoa(bidAmountOne), Lot: strconv.Itoa(lotOne), Act: "tend"},
			))
		})

		It("limits results to latest blocks if max_results argument is provided", func() {
			bidId := rand.Int()
			lotAmount := rand.Int()
			bidAmount := rand.Int()

			flapKick := test_data.FlapKickModel
			flapKick.ContractAddress = contractAddress
			flapKick.BidId = strconv.Itoa(bidId)
			flapKickErr := flapKickRepo.Create(headerOneId, []interface{}{flapKick})
			Expect(flapKickErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           bidId,
				ContractAddress: contractAddress,
				Lot:             lotAmount,
				BidAmount:       bidAmount,
				TendRepo:        tendRepo,
				TendHeaderId:    headerTwoId,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			maxResults := 1
			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events($1)`,
				maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(bidAmount), Lot: strconv.Itoa(lotAmount), Act: "tend"},
			))
		})

		It("ignores bid events from flops", func() {
			flopKickRepo := flop_kick.FlopKickRepository{}
			flopKickRepo.SetDB(db)
			flopKickEvent := test_data.FlopKickModel
			flopKickEvent.ContractAddress = "flop"
			flopKickEvent.BidId = strconv.Itoa(fakeBidId)
			flopKickErr := flopKickRepo.Create(headerOneId, []interface{}{flopKickEvent})
			Expect(flopKickErr).NotTo(HaveOccurred())
			flopKickBidEvent := test_helpers.BidEvent{BidId: flopKickEvent.BidId, BidAmount: flopKickEvent.Bid, Lot: flopKickEvent.Lot, Act: "kick", ContractAddress: flopKickEvent.ContractAddress}

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act, contract_address FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flapKickEvent.BidId, BidAmount: flapKickEvent.Bid, Lot: flapKickEvent.Lot, Act: "kick", ContractAddress: flapKickEvent.ContractAddress},
			))
			Expect(actualBidEvents).NotTo(ContainElement(flopKickBidEvent))
		})
	})

	Describe("tend", func() {
		It("returns flap tend bid events from multiple blocks", func() {
			lot := rand.Int()
			bidAmount := rand.Int()
			updatedLot := lot + 100
			updatedBidAmount := bidAmount + 100
			flapKickBlockOne := flapKickEvent

			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             lot,
				BidAmount:       bidAmount,
				TendRepo:        tendRepo,
				TendHeaderId:    headerOneId,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flapTendErr = test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             updatedLot,
				BidAmount:       updatedBidAmount,
				TendRepo:        tendRepo,
				TendHeaderId:    headerTwoId,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			headerThree := fakes.GetFakeHeaderWithTimestamp(int64(333333333), 3)
			headerThreeId, headerThreeErr := headerRepo.CreateOrUpdateHeader(headerThree)
			Expect(headerThreeErr).NotTo(HaveOccurred())

			// create irrelevant flop tend
			flopTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           fakeBidId,
				ContractAddress: anotherContractAddress,
				Lot:             lot,
				BidAmount:       bidAmount,
				TendRepo:        tendRepo,
				TendHeaderId:    headerThreeId,
			})
			Expect(flopTendErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(bidAmount), Lot: strconv.Itoa(lot), Act: "tend"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(updatedBidAmount), Lot: strconv.Itoa(updatedLot), Act: "tend"},
				test_helpers.BidEvent{BidId: flapKickBlockOne.BidId, BidAmount: flapKickBlockOne.Bid, Lot: flapKickBlockOne.Lot, Act: "kick"},
			))
		})
	})

	Describe("Deal", func() {
		It("returns bid events with lot and bid amount values from the block where the deal occurred", func() {
			blockTwo := blockOne + 1
			blockThree := blockTwo + 1

			flapKickBlockOne := flapKickEvent

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := fakes.GetFakeHeader(int64(blockTwo))
			_, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			updatedFlapStorageValues := test_helpers.GetFlapStorageValues(2, fakeBidId)
			test_helpers.CreateFlap(db, headerTwo, updatedFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerThree := fakes.GetFakeHeader(int64(blockThree))
			headerThreeId, headerThreeErr := headerRepo.CreateOrUpdateHeader(headerThree)
			Expect(headerThreeErr).NotTo(HaveOccurred())

			flapDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderId:    headerThreeId,
			})
			Expect(flapDealErr).NotTo(HaveOccurred())

			dealBlockFlapStorageValues := test_helpers.GetFlapStorageValues(0, fakeBidId)
			test_helpers.CreateFlap(db, headerThree, dealBlockFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: dealBlockFlapStorageValues[storage.BidBid].(string), Lot: dealBlockFlapStorageValues[storage.BidLot].(string), Act: "deal"},
				test_helpers.BidEvent{BidId: flapKickBlockOne.BidId, BidAmount: flapKickBlockOne.Bid, Lot: flapKickBlockOne.Lot, Act: "kick"}))
		})
	})

	Describe("Yank event", func() {
		It("includes yank in all flap bid events", func() {
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()

			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				TendRepo:        tendRepo,
				TendHeaderId:    headerOneId,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flapYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				YankRepo:        yankRepo,
				YankHeaderId:    headerTwoId,
			})
			Expect(flapYankErr).NotTo(HaveOccurred())

			updatedFlapStorageValues := test_helpers.GetFlapStorageValues(2, fakeBidId)
			test_helpers.CreateFlap(db, headerTwo, updatedFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flapKickEvent.BidId, BidAmount: flapKickEvent.Bid, Lot: flapKickEvent.Lot, Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(fakeBidAmount), Lot: strconv.Itoa(fakeLot), Act: "tend"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: updatedFlapStorageValues[storage.BidBid].(string), Lot: updatedFlapStorageValues[storage.BidLot].(string), Act: "yank"},
			))
		})

		It("ignores flop yank events", func() {
			fakeBidId := rand.Int()

			flopStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			// irrelevant flop yank
			flopYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidId:           fakeBidId,
				ContractAddress: anotherContractAddress,
				YankRepo:        yankRepo,
				YankHeaderId:    headerTwoId,
			})
			Expect(flopYankErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: flapKickEvent.BidId, BidAmount: flapKickEvent.Bid, Lot: flapKickEvent.Lot, Act: "kick"},
			))
		})
	})

})
