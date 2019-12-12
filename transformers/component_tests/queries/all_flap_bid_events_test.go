package queries

import (
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flap bid events query", func() {
	var (
		headerRepo             repositories.HeaderRepository
		contractAddress        = fakes.FakeAddress.Hex()
		anotherContractAddress = common.HexToAddress("0xabcdef123456789").Hex()
		blockOne, timestampOne int
		headerOne              core.Header
		fakeBidId              int
		flapKickEvent          event.InsertionModel
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		fakeBidId = rand.Int()

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
		flapKickLog := test_data.CreateTestLog(headerOne.Id, db)
		addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		flapKickEvent = test_data.FlapKickModel()
		flapKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
		flapKickEvent.ColumnValues[event.LogFK] = flapKickLog.ID
		flapKickEvent.ColumnValues[event.AddressFK] = addressId
		flapKickEvent.ColumnValues[constants.BidIDColumn] = strconv.Itoa(fakeBidId)
		flapKickErr := event.PersistModels([]event.InsertionModel{flapKickEvent}, db)
		Expect(flapKickErr).NotTo(HaveOccurred())
	})

	Describe("all_flap_bid_events", func() {
		It("returns all flap bid events (same block)", func() {
			fakeLot := rand.Int()
			fakeBidAmount := rand.Int()
			tendLog := test_data.CreateTestLog(headerOne.Id, db)

			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				BidId:           fakeBidId,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				TendHeaderId:    headerOne.Id,
				TendLogId:       tendLog.ID,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			tickLog := test_data.CreateTestLog(headerOne.Id, db)
			flapTickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				TickHeaderId:    headerOne.Id,
				TickLogId:       tickLog.ID,
			})
			Expect(flapTickErr).NotTo(HaveOccurred())

			flapDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealHeaderId:    headerOne.Id,
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

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

			flapKickEventTwoLog := test_data.CreateTestLog(headerTwo.Id, db)
			addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
			flapKickEventTwo := test_data.FlapKickModel()
			flapKickEventTwo.ColumnValues[event.HeaderFK] = headerTwo.Id
			flapKickEventTwo.ColumnValues[event.LogFK] = flapKickEventTwoLog.ID
			flapKickEventTwo.ColumnValues[event.AddressFK] = addressId
			flapKickEventTwo.ColumnValues[constants.BidIDColumn] = strconv.Itoa(fakeBidIdTwo)
			flapKickEventTwo.ColumnValues[constants.LotColumn] = strconv.Itoa(rand.Int())
			flapKickEventTwo.ColumnValues[constants.BidColumn] = strconv.Itoa(rand.Int())
			flapKickErr := event.PersistModels([]event.InsertionModel{flapKickEventTwo}, db)
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

			flapKickEventTwoLog := test_data.CreateTestLog(headerOne.Id, db)
			addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())
			flapKickEventTwo := test_data.FlapKickModel()
			flapKickEventTwo.ColumnValues[event.HeaderFK] = headerOne.Id
			flapKickEventTwo.ColumnValues[event.LogFK] = flapKickEventTwoLog.ID
			flapKickEventTwo.ColumnValues[event.AddressFK] = addressId
			flapKickEventTwo.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidIdTwo)
			flapKickErr := event.PersistModels([]event.InsertionModel{flapKickEventTwo}, db)
			Expect(flapKickErr).NotTo(HaveOccurred())

			flapTendOneLog := test_data.CreateTestLog(headerOne.Id, db)
			flapTendOneErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				BidId:           bidIdOne,
				Lot:             lotOne,
				BidAmount:       bidAmountOne,
				TendHeaderId:    headerOne.Id,
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

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
				tendLogId := test_data.CreateTestLog(headerTwo.Id, db).ID

				flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
					DB:              db,
					ContractAddress: contractAddress,
					BidId:           fakeBidId,
					Lot:             lotAmount,
					BidAmount:       bidAmount,
					TendHeaderId:    headerTwo.Id,
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
			flopKickLog := test_data.CreateTestLog(headerOne.Id, db)
			flopAddress := test_data.RandomString(40)
			addressId, addressErr := shared.GetOrCreateAddress(flopAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())

			flopKickEvent := test_data.FlopKickModel()
			flopKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
			flopKickEvent.ColumnValues[event.LogFK] = flopKickLog.ID
			flopKickEvent.ColumnValues[event.AddressFK] = addressId
			flopKickEvent.ColumnValues[constants.BidIDColumn] = strconv.Itoa(fakeBidId)
			flopKickErr := event.PersistModels([]event.InsertionModel{flopKickEvent}, db)
			Expect(flopKickErr).NotTo(HaveOccurred())
			flopKickBidEvent := test_helpers.BidEvent{
				BidId:           flopKickEvent.ColumnValues[constants.BidIDColumn].(string),
				BidAmount:       flopKickEvent.ColumnValues[constants.BidColumn].(string),
				Lot:             flopKickEvent.ColumnValues[constants.LotColumn].(string),
				Act:             "kick",
				ContractAddress: flopAddress,
			}

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act, contract_address FROM api.all_flap_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:           flapKickEvent.ColumnValues["bid_id"].(string),
					BidAmount:       flapKickEvent.ColumnValues["bid"].(string),
					Lot:             flapKickEvent.ColumnValues["lot"].(string),
					Act:             "kick",
					ContractAddress: contractAddress,
				},
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
			flapTendOneLog := test_data.CreateTestLog(headerOne.Id, db)

			flapTendOneErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				BidId:           fakeBidId,
				Lot:             lot,
				BidAmount:       bidAmount,
				TendHeaderId:    headerOne.Id,
				TendLogId:       flapTendOneLog.ID,
			})
			Expect(flapTendOneErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			flapTendTwoLog := test_data.CreateTestLog(headerTwo.Id, db)

			flapTendTwoErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				BidId:           fakeBidId,
				Lot:             updatedLot,
				BidAmount:       updatedBidAmount,
				TendHeaderId:    headerTwo.Id,
				TendLogId:       flapTendTwoLog.ID,
			})
			Expect(flapTendTwoErr).NotTo(HaveOccurred())

			headerThree := createHeader(blockOne+2, timestampOne+2, headerRepo)
			tendLog := test_data.CreateTestLog(headerThree.Id, db)

			// create irrelevant flop tend
			flopTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				DB:              db,
				ContractAddress: anotherContractAddress,
				BidId:           fakeBidId,
				Lot:             lot,
				BidAmount:       bidAmount,
				TendHeaderId:    headerThree.Id,
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
			tickLog := test_data.CreateTestLog(headerOne.Id, db)

			tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				TickHeaderId:    headerOne.Id,
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

			headerTwo := createHeader(blockTwo, timestampOne+1, headerRepo)
			headerThree := createHeader(blockThree, timestampOne+2, headerRepo)

			flapKickBlockOne := flapKickEvent

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			updatedFlapStorageValues := test_helpers.GetFlapStorageValues(2, fakeBidId)
			test_helpers.CreateFlap(db, headerTwo, updatedFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flapDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealHeaderId:    headerThree.Id,
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

			tendLog := test_data.CreateTestLog(headerOne.Id, db)
			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				BidId:           fakeBidId,
				Lot:             fakeLot,
				BidAmount:       fakeBidAmount,
				TendHeaderId:    headerOne.Id,
				TendLogId:       tendLog.ID,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			flapYankLog := test_data.CreateTestLog(headerTwo.Id, db)

			flapYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				YankHeaderId:    headerTwo.Id,
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

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			yankLog := test_data.CreateTestLog(headerTwo.Id, db)

			// irrelevant flop yank
			flopYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: anotherContractAddress,
				YankHeaderId:    headerTwo.Id,
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
