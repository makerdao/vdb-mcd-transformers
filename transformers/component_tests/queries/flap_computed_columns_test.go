package queries

import (
	"math/rand"
	"strconv"

	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flap computed columns", func() {
	var (
		headerRepo             repositories.HeaderRepository
		contractAddress        = fakes.FakeAddress.Hex()
		fakeBidId              = rand.Int()
		blockOne, timestampOne int
		headerOne              core.Header
		diffID                 int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		timestampOne := int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)

		diffID = storage_helper.CreateFakeDiffRecord(db)
	})

	Describe("flap_bid_events", func() {
		It("returns the bid events for a flap", func() {
			flapKickLog := test_data.CreateTestLog(headerOne.Id, db)

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, diffID, headerOne.Id, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)
			addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())

			flapKickEvent := test_data.FlapKickModel()
			flapKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
			flapKickEvent.ColumnValues[event.LogFK] = flapKickLog.ID
			flapKickEvent.ColumnValues[event.AddressFK] = addressId
			flapKickEvent.ColumnValues[constants.BidIDColumn] = strconv.Itoa(fakeBidId)
			flapKickErr := event.PersistModels([]event.InsertionModel{flapKickEvent}, db)
			Expect(flapKickErr).NotTo(HaveOccurred())

			expectedBidEvents := test_helpers.BidEvent{
				BidId:           strconv.Itoa(fakeBidId),
				Lot:             flapKickEvent.ColumnValues["lot"].(string),
				BidAmount:       flapKickEvent.ColumnValues["bid"].(string),
				Act:             "kick",
				ContractAddress: contractAddress,
			}
			var actualBidEvents test_helpers.BidEvent
			queryErr := db.Get(&actualBidEvents,
				`SELECT bid_id, bid_amount, lot, act, contract_address FROM api.flap_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flap_state
    					FROM api.all_flaps()))`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(Equal(expectedBidEvents))
		})

		It("does not include bid events for a different flap", func() {
			flapKickLog := test_data.CreateTestLog(headerOne.Id, db)

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, diffID, headerOne.Id, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)
			addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
			Expect(addressErr).NotTo(HaveOccurred())

			flapKickEvent := test_data.FlapKickModel()
			flapKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
			flapKickEvent.ColumnValues[event.LogFK] = flapKickLog.ID
			flapKickEvent.ColumnValues[event.AddressFK] = addressId
			flapKickEvent.ColumnValues[constants.BidIDColumn] = strconv.Itoa(fakeBidId)
			flapKickErr := event.PersistModels([]event.InsertionModel{flapKickEvent}, db)
			Expect(flapKickErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			irrelevantFlipKickLog := test_data.CreateTestLog(headerTwo.Id, db)

			irrelevantBidId := fakeBidId + 9999999999999
			irrelevantFlapStorageValues := test_helpers.GetFlapStorageValues(2, irrelevantBidId)
			test_helpers.CreateFlap(db, diffID, headerTwo.Id, irrelevantFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(irrelevantBidId)), contractAddress)

			irrelevantFlapKickEvent := test_data.FlapKickModel()
			irrelevantFlapKickEvent.ColumnValues[event.HeaderFK] = headerTwo.Id
			irrelevantFlapKickEvent.ColumnValues[event.LogFK] = irrelevantFlipKickLog.ID
			irrelevantFlapKickEvent.ColumnValues[event.AddressFK] = addressId
			irrelevantFlapKickEvent.ColumnValues[constants.BidIDColumn] = strconv.Itoa(irrelevantBidId)
			flapKickErr = event.PersistModels([]event.InsertionModel{irrelevantFlapKickEvent}, db)
			Expect(flapKickErr).NotTo(HaveOccurred())

			expectedBidEvents := test_helpers.BidEvent{
				BidId:           strconv.Itoa(fakeBidId),
				Lot:             flapKickEvent.ColumnValues["lot"].(string),
				BidAmount:       flapKickEvent.ColumnValues["bid"].(string),
				Act:             "kick",
				ContractAddress: contractAddress,
			}

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents,
				`SELECT bid_id, bid_amount, lot, act, contract_address FROM api.flap_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flap_state
    					FROM api.all_flaps() WHERE bid_id = $1))`, fakeBidId)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(ConsistOf(expectedBidEvents))
		})

		Describe("result pagination", func() {
			var (
				tendBid, tendLot int
				flapKickEvent    event.InsertionModel
			)

			BeforeEach(func() {
				logId := test_data.CreateTestLog(headerOne.Id, db).ID

				flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
				test_helpers.CreateFlap(db, diffID, headerOne.Id, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)
				addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
				Expect(addressErr).NotTo(HaveOccurred())

				flapKickEvent = test_data.FlapKickModel()
				flapKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
				flapKickEvent.ColumnValues[event.LogFK] = logId
				flapKickEvent.ColumnValues[event.AddressFK] = addressId
				flapKickEvent.ColumnValues[constants.BidIDColumn] = strconv.Itoa(fakeBidId)
				flapKickErr := event.PersistModels([]event.InsertionModel{flapKickEvent}, db)
				Expect(flapKickErr).NotTo(HaveOccurred())

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
				logTwoId := test_data.CreateTestLog(headerTwo.Id, db).ID

				tendBid = rand.Int()
				tendLot = rand.Int()
				flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
					DB:              db,
					ContractAddress: contractAddress,
					BidId:           fakeBidId,
					Lot:             tendLot,
					BidAmount:       tendBid,
					TendHeaderId:    headerTwo.Id,
					TendLogId:       logTwoId,
				})
				Expect(flapTendErr).NotTo(HaveOccurred())
			})

			It("limits result to most recent block if max_results argument is provided", func() {
				expectedBidEvent := test_helpers.BidEvent{
					BidId:           strconv.Itoa(fakeBidId),
					Lot:             strconv.Itoa(tendLot),
					BidAmount:       strconv.Itoa(tendBid),
					Act:             "tend",
					ContractAddress: contractAddress,
				}

				maxResults := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents,
					`SELECT bid_id, bid_amount, lot, act, contract_address FROM api.flap_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flap_state
    					FROM api.all_flaps() WHERE bid_id = $1), $2)`, fakeBidId, maxResults)

				Expect(queryErr).NotTo(HaveOccurred())
				Expect(actualBidEvents).To(ConsistOf(expectedBidEvent))
			})

			It("offsets results if offset is provided", func() {
				expectedBidEvent := test_helpers.BidEvent{
					BidId:           strconv.Itoa(fakeBidId),
					Lot:             flapKickEvent.ColumnValues[constants.LotColumn].(string),
					BidAmount:       flapKickEvent.ColumnValues[constants.BidColumn].(string),
					Act:             "kick",
					ContractAddress: contractAddress,
				}

				maxResults := 1
				resultOffset := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents,
					`SELECT bid_id, bid_amount, lot, act, contract_address FROM api.flap_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flap_state
    					FROM api.all_flaps() WHERE bid_id = $1), $2, $3)`, fakeBidId, maxResults, resultOffset)

				Expect(queryErr).NotTo(HaveOccurred())
				Expect(actualBidEvents).To(ConsistOf(expectedBidEvent))
			})
		})
	})
})
