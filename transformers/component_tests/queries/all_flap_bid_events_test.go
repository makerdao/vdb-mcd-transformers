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
	"github.com/vulcanize/mcd_transformers/transformers/events/tick"
	"github.com/vulcanize/mcd_transformers/transformers/events/yank"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
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
		tickRepo               tick.TickRepository
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
		flapKickEvent          shared.InsertionModel
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flapKickRepo = flap_kick.FlapKickRepository{}
		flapKickRepo.SetDB(db)
		tendRepo = tend.TendRepository{}
		tendRepo.SetDB(db)
		tickRepo = tick.TickRepository{}
		tickRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		yankRepo = yank.YankRepository{}
		yankRepo.SetDB(db)
		fakeBidId = rand.Int()

		blockOne = 1
		headerOne = fakes.GetFakeHeader(blockOne)
		headerOneId, headerOneErr = headerRepo.CreateOrUpdateHeader(headerOne)
		Expect(headerOneErr).NotTo(HaveOccurred())
		flapKickLog := test_data.CreateTestLog(headerOneId, db)

		flapKickEvent = test_data.FlapKickModel()
		flapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
		flapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
		flapKickEvent.ColumnValues[constants.HeaderFK] = headerOneId
		flapKickEvent.ColumnValues[constants.LogFK] = flapKickLog.ID
		flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEvent})
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
			tendLog := test_data.CreateTestLog(headerOneId, db)

			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				TendRepo:        tendRepo,
				TendHeaderId:    headerOneId,
				TendLogId:       tendLog.ID,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			tickLog := test_data.CreateTestLog(headerOneId, db)
			flapTickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				TickRepo:        tickRepo,
				TickHeaderId:    headerOneId,
				TickLogId:       tickLog.ID,
			})
			Expect(flapTickErr).NotTo(HaveOccurred())

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
				test_helpers.BidEvent{
					BidId:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(fakeBidAmount), Lot: strconv.Itoa(fakeLot), Act: "tend"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: flapStorageValues[storage.BidBid].(string), Lot: flapStorageValues[storage.BidLot].(string), Act: "tick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: flapStorageValues[storage.BidBid].(string), Lot: flapStorageValues[storage.BidLot].(string), Act: "deal"},
			))
		})

		It("returns all flap bid events across all blocks", func() {
			fakeBidIdTwo := fakeBidId + 1

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			flapKickEventTwoLog := test_data.CreateTestLog(headerTwoId, db)
			flapKickEventTwo := test_data.FlapKickModel()
			flapKickEventTwo.ColumnValues["bid"] = strconv.Itoa(rand.Int())
			flapKickEventTwo.ColumnValues["lot"] = strconv.Itoa(rand.Int())
			flapKickEventTwo.ColumnValues["bid_id"] = strconv.Itoa(fakeBidIdTwo)
			flapKickEventTwo.ColumnValues[constants.HeaderFK] = headerTwoId
			flapKickEventTwo.ColumnValues[constants.LogFK] = flapKickEventTwoLog.ID
			flapKickEventTwo.ForeignKeyValues[constants.AddressFK] = contractAddress
			flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEventTwo})
			Expect(flapKickErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{
					BidId:     flapKickEventTwo.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEventTwo.ColumnValues["bid"].(string),
					Lot:       flapKickEventTwo.ColumnValues["lot"].(string),
					Act:       "kick"},
			))
		})

		It("returns bid events for multiple bid ids", func() {
			bidIdOne := fakeBidId
			bidIdTwo := rand.Int()
			lotOne := rand.Int()
			bidAmountOne := rand.Int()

			flapKickEventTwoLog := test_data.CreateTestLog(headerOneId, db)
			flapKickEventTwo := test_data.FlapKickModel()
			flapKickEventTwo.ColumnValues["bid_id"] = strconv.Itoa(bidIdTwo)
			flapKickEventTwo.ColumnValues[constants.HeaderFK] = headerOneId
			flapKickEventTwo.ColumnValues[constants.LogFK] = flapKickEventTwoLog.ID
			flapKickEventTwo.ForeignKeyValues[constants.AddressFK] = contractAddress
			flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEventTwo})
			Expect(flapKickErr).NotTo(HaveOccurred())

			flapTendOneLog := test_data.CreateTestLog(headerOneId, db)
			flapTendOneErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           bidIdOne,
				ContractAddress: contractAddress,
				Lot:             lotOne,
				BidAmount:       bidAmountOne,
				TendRepo:        tendRepo,
				TendHeaderId:    headerOneId,
				TendLogId:       flapTendOneLog.ID,
			})
			Expect(flapTendOneErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{
					BidId:     flapKickEventTwo.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEventTwo.ColumnValues["bid"].(string),
					Lot:       flapKickEventTwo.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{
					BidId:     strconv.Itoa(bidIdOne),
					BidAmount: strconv.Itoa(bidAmountOne),
					Lot:       strconv.Itoa(lotOne),
					Act:       "tend"},
			))
		})

		Describe("result pagination", func() {
			var (
				bidAmount, lotAmount int
			)

			BeforeEach(func() {
				lotAmount = rand.Int()
				bidAmount = rand.Int()

				headerTwo := fakes.GetFakeHeader(2)
				headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
				Expect(headerTwoErr).NotTo(HaveOccurred())
				tendLogId := test_data.CreateTestLog(headerTwoId, db).ID

				flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
					Lot:             lotAmount,
					BidAmount:       bidAmount,
					TendRepo:        tendRepo,
					TendHeaderId:    headerTwoId,
					TendLogId:       tendLogId,
				})
				Expect(flapTendErr).NotTo(HaveOccurred())
			})

			It("limits results to latest blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events($1)`,
					maxResults)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualBidEvents).To(ConsistOf(
					test_helpers.BidEvent{
						BidId:     strconv.Itoa(fakeBidId),
						BidAmount: strconv.Itoa(bidAmount),
						Lot:       strconv.Itoa(lotAmount),
						Act:       "tend",
					},
				))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events($1, $2)`,
					maxResults, resultOffset)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualBidEvents).To(ConsistOf(
					test_helpers.BidEvent{
						BidId:     flapKickEvent.ColumnValues["bid_id"].(string),
						BidAmount: flapKickEvent.ColumnValues["bid"].(string),
						Lot:       flapKickEvent.ColumnValues["lot"].(string),
						Act:       "kick",
					},
				))
			})
		})

		It("ignores bid events from flops", func() {
			flopKickLog := test_data.CreateTestLog(headerOneId, db)
			flopKickRepo := flop_kick.FlopKickRepository{}
			flopKickRepo.SetDB(db)

			flopKickEvent := test_data.FlopKickModel()
			flopKickEvent.ForeignKeyValues[constants.AddressFK] = "flop"
			flopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
			flopKickEvent.ColumnValues[constants.HeaderFK] = headerOneId
			flopKickEvent.ColumnValues[constants.LogFK] = flopKickLog.ID
			flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEvent})
			Expect(flopKickErr).NotTo(HaveOccurred())
			flopKickBidEvent := test_helpers.BidEvent{
				BidId:           flopKickEvent.ColumnValues["bid_id"].(string),
				BidAmount:       flopKickEvent.ColumnValues["bid"].(string),
				Lot:             flopKickEvent.ColumnValues["lot"].(string),
				Act:             "kick",
				ContractAddress: flopKickEvent.ForeignKeyValues[constants.AddressFK]}

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act, contract_address FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:           flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount:       flapKickEvent.ColumnValues["bid"].(string),
					Lot:             flapKickEvent.ColumnValues["lot"].(string),
					Act:             "kick",
					ContractAddress: flapKickEvent.ForeignKeyValues[constants.AddressFK]},
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
			flapTendOneLog := test_data.CreateTestLog(headerOneId, db)

			flapTendOneErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             lot,
				BidAmount:       bidAmount,
				TendRepo:        tendRepo,
				TendHeaderId:    headerOneId,
				TendLogId:       flapTendOneLog.ID,
			})
			Expect(flapTendOneErr).NotTo(HaveOccurred())

			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			flapTendTwoLog := test_data.CreateTestLog(headerTwoId, db)

			flapTendTwoErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             updatedLot,
				BidAmount:       updatedBidAmount,
				TendRepo:        tendRepo,
				TendHeaderId:    headerTwoId,
				TendLogId:       flapTendTwoLog.ID,
			})
			Expect(flapTendTwoErr).NotTo(HaveOccurred())

			headerThree := fakes.GetFakeHeaderWithTimestamp(int64(333333333), 3)
			headerThreeId, headerThreeErr := headerRepo.CreateOrUpdateHeader(headerThree)
			Expect(headerThreeErr).NotTo(HaveOccurred())
			tendLog := test_data.CreateTestLog(headerThreeId, db)

			// create irrelevant flop tend
			flopTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           fakeBidId,
				ContractAddress: anotherContractAddress,
				Lot:             lot,
				BidAmount:       bidAmount,
				TendRepo:        tendRepo,
				TendHeaderId:    headerThreeId,
				TendLogId:       tendLog.ID,
			})
			Expect(flopTendErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(bidAmount), Lot: strconv.Itoa(lot), Act: "tend"},
				test_helpers.BidEvent{BidId: strconv.Itoa(fakeBidId), BidAmount: strconv.Itoa(updatedBidAmount), Lot: strconv.Itoa(updatedLot), Act: "tend"},
				test_helpers.BidEvent{
					BidId:     flapKickBlockOne.ColumnValues["bid_id"].(string),
					BidAmount: flapKickBlockOne.ColumnValues["bid"].(string),
					Lot:       flapKickBlockOne.ColumnValues["lot"].(string),
					Act:       "kick"},
			))
		})
	})

	Describe("tick event", func() {
		It("ignores tick events from non flap contracts", func() {
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
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick",
				},
			))
		})

		It("includes flap tick bid events", func() {
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
			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{
					BidId:     strconv.Itoa(fakeBidId),
					BidAmount: flapStorageValues[storage.BidBid].(string),
					Lot:       flapStorageValues[storage.BidLot].(string),
					Act:       "tick"},
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
				test_helpers.BidEvent{
					BidId:     flapKickBlockOne.ColumnValues["bid_id"].(string),
					BidAmount: flapKickBlockOne.ColumnValues["bid"].(string),
					Lot:       flapKickBlockOne.ColumnValues["lot"].(string),
					Act:       "kick"}))
		})
	})

	Describe("Yank event", func() {
		It("includes yank in all flap bid events", func() {
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()

			tendLog := test_data.CreateTestLog(headerOneId, db)
			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				TendRepo:        tendRepo,
				TendHeaderId:    headerOneId,
				TendLogId:       tendLog.ID,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := fakes.GetFakeHeader(2)
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			flapYankLog := test_data.CreateTestLog(headerTwoId, db)

			flapYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				YankRepo:        yankRepo,
				YankHeaderId:    headerTwoId,
				YankLogId:       flapYankLog.ID,
			})
			Expect(flapYankErr).NotTo(HaveOccurred())

			updatedFlapStorageValues := test_helpers.GetFlapStorageValues(2, fakeBidId)
			test_helpers.CreateFlap(db, headerTwo, updatedFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
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
			yankLog := test_data.CreateTestLog(headerTwoId, db)

			// irrelevant flop yank
			flopYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				BidId:           fakeBidId,
				ContractAddress: anotherContractAddress,
				YankRepo:        yankRepo,
				YankHeaderId:    headerTwoId,
				YankLogId:       yankLog.ID,
			})
			Expect(flopYankErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:     flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount: flapKickEvent.ColumnValues["bid"].(string),
					Lot:       flapKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
			))
		})
	})
})
