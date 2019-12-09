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
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All flip bid events query", func() {
	var (
		db                     *postgres.DB
		headerRepo             repositories.HeaderRepository
		contractAddress        = fakes.FakeAddress.Hex()
		anotherContractAddress = common.HexToAddress("0xabcdef123456789").Hex()
		addressId              int64
		bidId                  int
		blockOne, timestampOne int
		headerOne              core.Header
		flipKickEvent          event.InsertionModel
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		bidId = rand.Int()

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)

		flipKickLog := test_data.CreateTestLog(headerOne.Id, db)

		var addressErr error
		addressId, addressErr = shared.GetOrCreateAddress(contractAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		flipKickEvent = test_data.FlipKickModel()
		flipKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
		flipKickEvent.ColumnValues[event.LogFK] = flipKickLog.ID
		flipKickEvent.ColumnValues[event.AddressFK] = addressId
		flipKickEvent.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidId)
		flipKickErr := event.PersistModels([]event.InsertionModel{flipKickEvent}, db)
		Expect(flipKickErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("all_flip_bid_events", func() {
		It("returns all flip bid events when they are all in the same block", func() {
			tendLot := rand.Int()
			tendBidAmount := rand.Int()
			dentLot := rand.Int()
			dentBidAmount := rand.Int()

			flipTendLog := test_data.CreateTestLog(headerOne.Id, db)
			flipTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				BidId:           bidId,
				Lot:             tendLot,
				BidAmount:       tendBidAmount,
				TendHeaderId:    headerOne.Id,
				TendLogId:       flipTendLog.ID,
			})
			Expect(flipTendErr).NotTo(HaveOccurred())

			tickLog := test_data.CreateTestLog(headerOne.Id, db)
			tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				DB:              db,
				BidId:           bidId,
				ContractAddress: contractAddress,
				TickHeaderId:    headerOne.Id,
				TickLogId:       tickLog.ID,
			})
			Expect(tickErr).NotTo(HaveOccurred())

			flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, bidId)
			test_helpers.CreateFlip(db, headerOne, flipStorageValues,
				test_helpers.GetFlipMetadatas(strconv.Itoa(bidId)), contractAddress)

			flipDentLog := test_data.CreateTestLog(headerOne.Id, db)
			flipDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				DB:              db,
				BidId:           bidId,
				ContractAddress: contractAddress,
				Lot:             dentLot,
				BidAmount:       dentBidAmount,
				DentHeaderId:    headerOne.Id,
				DentLogId:       flipDentLog.ID,
			})
			Expect(flipDentErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:     strconv.Itoa(bidId),
					BidAmount: flipKickEvent.ColumnValues["bid"].(string),
					Lot:       flipKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(tendBidAmount), Lot: strconv.Itoa(tendLot), Act: "tend"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flipStorageValues[storage.BidBid].(string), Lot: flipStorageValues[storage.BidLot].(string), Act: "tick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(dentBidAmount), Lot: strconv.Itoa(dentLot), Act: "dent"}))
		})

		It("returns flip bid events across all blocks", func() {
			tendLot := rand.Int()
			tendBidAmount := rand.Int()
			dentLot := rand.Int()
			dentBidAmount := rand.Int()

			flipTendLog := test_data.CreateTestLog(headerOne.Id, db)
			flipTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				BidId:           bidId,
				Lot:             tendLot,
				BidAmount:       tendBidAmount,
				TendHeaderId:    headerOne.Id,
				TendLogId:       flipTendLog.ID,
			})
			Expect(flipTendErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

			tickLog := test_data.CreateTestLog(headerTwo.Id, db)
			tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				DB:              db,
				BidId:           bidId,
				ContractAddress: contractAddress,
				TickHeaderId:    headerTwo.Id,
				TickLogId:       tickLog.ID,
			})
			Expect(tickErr).NotTo(HaveOccurred())

			flipStorageValuesBlockTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, bidId)
			test_helpers.CreateFlip(db, headerTwo, flipStorageValuesBlockTwo,
				test_helpers.GetFlipMetadatas(strconv.Itoa(bidId)), contractAddress)

			headerThree := createHeader(blockOne+2, timestampOne+2, headerRepo)

			flipDentLog := test_data.CreateTestLog(headerThree.Id, db)
			flipDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				DB:              db,
				BidId:           bidId,
				ContractAddress: contractAddress,
				Lot:             dentLot,
				BidAmount:       dentBidAmount,
				DentHeaderId:    headerThree.Id,
				DentLogId:       flipDentLog.ID,
			})
			Expect(flipDentErr).NotTo(HaveOccurred())

			flipDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				DB:              db,
				BidId:           bidId,
				ContractAddress: contractAddress,
				DealHeaderId:    headerThree.Id,
			})
			Expect(flipDealErr).NotTo(HaveOccurred())

			flipStorageValuesBlockThree := test_helpers.GetFlipStorageValues(3, test_helpers.FakeIlk.Hex, bidId)
			test_helpers.CreateFlip(db, headerThree, flipStorageValuesBlockThree,
				test_helpers.GetFlipMetadatas(strconv.Itoa(bidId)), contractAddress)

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flipKickEvent.ColumnValues["bid"].(string), Lot: flipKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(tendBidAmount), Lot: strconv.Itoa(tendLot), Act: "tend"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flipStorageValuesBlockTwo[storage.BidBid].(string), Lot: flipStorageValuesBlockTwo[storage.BidLot].(string), Act: "tick"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(dentBidAmount), Lot: strconv.Itoa(dentLot), Act: "dent"},
				test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flipStorageValuesBlockThree[storage.BidBid].(string), Lot: flipStorageValuesBlockThree[storage.BidLot].(string), Act: "deal"}))
		})

		Describe("result pagination", func() {
			var updatedFlipValues map[string]interface{}

			BeforeEach(func() {
				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

				logID := test_data.CreateTestLog(headerTwo.Id, db).ID

				tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
					DB:              db,
					BidId:           bidId,
					ContractAddress: contractAddress,
					TickHeaderId:    headerTwo.Id,
					TickLogId:       logID,
				})
				Expect(tickErr).NotTo(HaveOccurred())

				updatedFlipValues = test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, bidId)
				test_helpers.CreateFlip(db, headerTwo, updatedFlipValues,
					test_helpers.GetFlipMetadatas(strconv.Itoa(bidId)), contractAddress)
			})

			It("limits result to latest blocks if max_results argument is provided", func() {
				maxResults := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events($1)`, maxResults)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualBidEvents).To(ConsistOf(
					test_helpers.BidEvent{
						BidId:     strconv.Itoa(bidId),
						BidAmount: updatedFlipValues[storage.BidBid].(string),
						Lot:       updatedFlipValues[storage.BidLot].(string),
						Act:       "tick",
					},
				))
			})

			It("offsets results if offset is provided", func() {
				maxResults := 1
				resultOffset := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events($1, $2)`, maxResults, resultOffset)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualBidEvents).To(ConsistOf(
					test_helpers.BidEvent{
						BidId:     strconv.Itoa(bidId),
						BidAmount: flipKickEvent.ColumnValues["bid"].(string),
						Lot:       flipKickEvent.ColumnValues["lot"].(string),
						Act:       "kick",
					},
				))
			})
		})

		It("returns bid events from flippers that have different bid ids", func() {
			differentBidId := rand.Int()
			differentLot := rand.Int()

			flipKickLogTwo := test_data.CreateTestLog(headerOne.Id, db)

			flipKickEventTwo := test_data.FlipKickModel()
			flipKickEventTwo.ColumnValues[event.HeaderFK] = headerOne.Id
			flipKickEventTwo.ColumnValues[event.LogFK] = flipKickLogTwo.ID
			flipKickEventTwo.ColumnValues[event.AddressFK] = addressId
			flipKickEventTwo.ColumnValues[constants.BidIDColumn] = strconv.Itoa(differentBidId)
			flipKickEventTwo.ColumnValues[constants.LotColumn] = strconv.Itoa(differentLot)
			flipKickErr := event.PersistModels([]event.InsertionModel{flipKickEventTwo}, db)
			Expect(flipKickErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:     strconv.Itoa(bidId),
					BidAmount: flipKickEvent.ColumnValues["bid"].(string),
					Lot:       flipKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{
					BidId:     flipKickEventTwo.ColumnValues["bid_id"].(string),
					BidAmount: flipKickEventTwo.ColumnValues["bid"].(string),
					Lot:       flipKickEventTwo.ColumnValues["lot"].(string),
					Act:       "kick"},
			))
		})

		It("returns bid events from different kinds of flips (flips with different contract addresses", func() {
			anotherFlipContractAddress := "DifferentFlipAddress"
			differentLot := rand.Int()
			differentBidAmount := rand.Int()

			anotherAddressId, addressErr := shared.GetOrCreateAddress(anotherFlipContractAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())

			flipKickLog := test_data.CreateTestLog(headerOne.Id, db)
			flipKickEventTwo := test_data.FlipKickModel()
			flipKickEventTwo.ColumnValues[event.HeaderFK] = headerOne.Id
			flipKickEventTwo.ColumnValues[event.LogFK] = flipKickLog.ID
			flipKickEventTwo.ColumnValues[event.AddressFK] = anotherAddressId
			flipKickEventTwo.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidId)
			flipKickEventTwo.ColumnValues[constants.LotColumn] = strconv.Itoa(differentLot)
			flipKickEventTwo.ColumnValues[constants.BidColumn] = strconv.Itoa(differentBidAmount)
			flipKickErr := event.PersistModels([]event.InsertionModel{flipKickEventTwo}, db)
			Expect(flipKickErr).NotTo(HaveOccurred())

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidEvents).To(ConsistOf(
				test_helpers.BidEvent{
					BidId:     strconv.Itoa(bidId),
					BidAmount: flipKickEvent.ColumnValues["bid"].(string),
					Lot:       flipKickEvent.ColumnValues["lot"].(string),
					Act:       "kick"},
				test_helpers.BidEvent{
					BidId:     flipKickEventTwo.ColumnValues["bid_id"].(string),
					BidAmount: flipKickEventTwo.ColumnValues["bid"].(string),
					Lot:       flipKickEventTwo.ColumnValues["lot"].(string),
					Act:       "kick"},
			))
		})

		Describe("tend", func() {
			It("returns tend events from multiple blocks", func() {
				lotOne := rand.Int()
				lotTwo := rand.Int()
				bidAmountOne := rand.Int()
				bidAmountTwo := rand.Int()

				flipTendHeaderOneLog := test_data.CreateTestLog(headerOne.Id, db)
				flipTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
					DB:              db,
					ContractAddress: contractAddress,
					BidId:           bidId,
					Lot:             lotOne,
					BidAmount:       bidAmountOne,
					TendHeaderId:    headerOne.Id,
					TendLogId:       flipTendHeaderOneLog.ID,
				})
				Expect(flipTendErr).NotTo(HaveOccurred())

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

				flipTendHeaderTwoLog := test_data.CreateTestLog(headerTwo.Id, db)
				flipTendHeaderTwoErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
					DB:              db,
					ContractAddress: contractAddress,
					BidId:           bidId,
					Lot:             lotTwo,
					BidAmount:       bidAmountTwo,
					TendHeaderId:    headerTwo.Id,
					TendLogId:       flipTendHeaderTwoLog.ID,
				})
				Expect(flipTendHeaderTwoErr).NotTo(HaveOccurred())

				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events()`)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualBidEvents).To(ConsistOf(
					test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flipKickEvent.ColumnValues["bid"].(string), Lot: flipKickEvent.ColumnValues["lot"].(string), Act: "kick"},
					test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(bidAmountOne), Lot: strconv.Itoa(lotOne), Act: "tend"},
					test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(bidAmountTwo), Lot: strconv.Itoa(lotTwo), Act: "tend"},
				))
			})

			It("ignores tend events that are not from a flip", func() {
				flapTendLog := test_data.CreateTestLog(headerOne.Id, db)
				flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
					DB:              db,
					ContractAddress: anotherContractAddress,
					BidId:           bidId,
					Lot:             rand.Int(),
					BidAmount:       rand.Int(),
					TendHeaderId:    headerOne.Id,
					TendLogId:       flapTendLog.ID,
				})
				Expect(flapTendErr).NotTo(HaveOccurred())

				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events()`)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualBidEvents).To(ConsistOf(
					test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flipKickEvent.ColumnValues["bid"].(string), Lot: flipKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				))
			})
		})

		Describe("dent", func() {
			It("returns dent events from multiple blocks", func() {
				lotOne := rand.Int()
				lotTwo := rand.Int()
				bidAmountOne := rand.Int()
				bidAmountTwo := rand.Int()

				flipDentHeaderOneLog := test_data.CreateTestLog(headerOne.Id, db)
				flipDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
					DB:              db,
					BidId:           bidId,
					ContractAddress: contractAddress,
					Lot:             lotOne,
					BidAmount:       bidAmountOne,
					DentHeaderId:    headerOne.Id,
					DentLogId:       flipDentHeaderOneLog.ID,
				})
				Expect(flipDentErr).NotTo(HaveOccurred())

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

				flipDentHeaderTwoLog := test_data.CreateTestLog(headerTwo.Id, db)
				flipDentHeaderTwoErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
					DB:              db,
					BidId:           bidId,
					ContractAddress: contractAddress,
					Lot:             lotTwo,
					BidAmount:       bidAmountTwo,
					DentHeaderId:    headerTwo.Id,
					DentLogId:       flipDentHeaderTwoLog.ID,
				})
				Expect(flipDentHeaderTwoErr).NotTo(HaveOccurred())

				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events()`)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualBidEvents).To(ConsistOf(
					test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flipKickEvent.ColumnValues["bid"].(string), Lot: flipKickEvent.ColumnValues["lot"].(string), Act: "kick"},
					test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(bidAmountOne), Lot: strconv.Itoa(lotOne), Act: "dent"},
					test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: strconv.Itoa(bidAmountTwo), Lot: strconv.Itoa(lotTwo), Act: "dent"},
				))
			})

			It("ignores dent events that are not from flip", func() {
				flapDentLog := test_data.CreateTestLog(headerOne.Id, db)
				flapDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
					DB:              db,
					BidId:           bidId,
					ContractAddress: anotherContractAddress,
					Lot:             rand.Int(),
					BidAmount:       rand.Int(),
					DentHeaderId:    headerOne.Id,
					DentLogId:       flapDentLog.ID,
				})
				Expect(flapDentErr).NotTo(HaveOccurred())

				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events()`)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualBidEvents).To(ConsistOf(
					test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flipKickEvent.ColumnValues["bid"].(string), Lot: flipKickEvent.ColumnValues["lot"].(string), Act: "kick"},
				))
			})
		})

		Describe("yank", func() {
			It("includes yank and gets values from the block where the yank occurred", func() {
				tendLot := rand.Int()
				tendBidAmount := rand.Int()

				flipTendLog := test_data.CreateTestLog(headerOne.Id, db)
				flipTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
					DB:              db,
					ContractAddress: contractAddress,
					BidId:           bidId,
					Lot:             tendLot,
					BidAmount:       tendBidAmount,
					TendHeaderId:    headerOne.Id,
					TendLogId:       flipTendLog.ID,
				})
				Expect(flipTendErr).NotTo(HaveOccurred())

				flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, bidId)
				test_helpers.CreateFlip(db, headerOne, flipStorageValues,
					test_helpers.GetFlipMetadatas(strconv.Itoa(bidId)), contractAddress)

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

				flipYankLog := test_data.CreateTestLog(headerTwo.Id, db)
				flipYankErr := test_helpers.CreateYank(test_helpers.YankCreationInput{
					DB:              db,
					BidId:           bidId,
					ContractAddress: contractAddress,
					YankHeaderId:    headerTwo.Id,
					YankLogId:       flipYankLog.ID,
				})
				Expect(flipYankErr).NotTo(HaveOccurred())

				updatedFlipStorageValues := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, bidId)
				test_helpers.CreateFlip(db, headerTwo, updatedFlipStorageValues,
					test_helpers.GetFlipMetadatas(strconv.Itoa(bidId)), contractAddress)

				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events()`)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(actualBidEvents).To(ConsistOf(
					test_helpers.BidEvent{
						BidId:     strconv.Itoa(bidId),
						BidAmount: flipKickEvent.ColumnValues["bid"].(string),
						Lot:       flipKickEvent.ColumnValues["lot"].(string),
						Act:       "kick",
					},
					test_helpers.BidEvent{
						BidId:     strconv.Itoa(bidId),
						BidAmount: strconv.Itoa(tendBidAmount),
						Lot:       strconv.Itoa(tendLot),
						Act:       "tend",
					},
					test_helpers.BidEvent{
						BidId:     strconv.Itoa(bidId),
						BidAmount: updatedFlipStorageValues[storage.BidBid].(string),
						Lot:       updatedFlipStorageValues[storage.BidLot].(string),
						Act:       "yank",
					},
				))
			})

			Describe("tick", func() {
				It("includes tick events", func() {
					flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, bidId)
					test_helpers.CreateFlip(db, headerOne, flipStorageValues,
						test_helpers.GetFlipMetadatas(strconv.Itoa(bidId)), contractAddress)
					tickLog := test_data.CreateTestLog(headerOne.Id, db)
					tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
						DB:              db,
						BidId:           bidId,
						ContractAddress: contractAddress,
						TickHeaderId:    headerOne.Id,
						TickLogId:       tickLog.ID,
					})
					Expect(tickErr).NotTo(HaveOccurred())

					var actualBidEvents []test_helpers.BidEvent
					queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events()`)
					Expect(queryErr).NotTo(HaveOccurred())

					Expect(actualBidEvents).To(ConsistOf(
						test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flipKickEvent.ColumnValues["bid"].(string), Lot: flipKickEvent.ColumnValues["lot"].(string), Act: "kick"},
						test_helpers.BidEvent{
							BidId:     strconv.Itoa(bidId),
							BidAmount: flipStorageValues[storage.BidBid].(string),
							Lot:       flipStorageValues[storage.BidLot].(string),
							Act:       "tick",
						},
					))
				})

				It("ignores tick events that aren't from flips", func() {
					tickLog := test_data.CreateTestLog(headerOne.Id, db)
					tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
						DB:              db,
						BidId:           bidId,
						ContractAddress: "flop",
						TickHeaderId:    headerOne.Id,
						TickLogId:       tickLog.ID,
					})
					Expect(tickErr).NotTo(HaveOccurred())

					var actualBidEvents []test_helpers.BidEvent
					queryErr := db.Select(&actualBidEvents, `SELECT bid_id, bid_amount, lot, act FROM api.all_flip_bid_events()`)
					Expect(queryErr).NotTo(HaveOccurred())

					// just the kick event because the tick is for a flop
					Expect(actualBidEvents).To(ConsistOf(
						test_helpers.BidEvent{BidId: strconv.Itoa(bidId), BidAmount: flipKickEvent.ColumnValues["bid"].(string), Lot: flipKickEvent.ColumnValues["lot"].(string), Act: "kick"},
					))
				})
			})
		})
	})
})
