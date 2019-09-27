package queries

import (
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/dent"
	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Flop computed columns", func() {
	var (
		db              *postgres.DB
		flopKickRepo    flop_kick.FlopKickRepository
		headerRepo      repositories.HeaderRepository
		contractAddress = "0x763ztv6x68exwqrgtl325e7hrcvavid4e3fcb4g"

		fakeBidId      = rand.Int()
		blockOne       = rand.Int()
		timestampOne   = int(rand.Int31())
		blockOneHeader = fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flopKickRepo = flop_kick.FlopKickRepository{}
		flopKickRepo.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("flop_bid_events", func() {
		It("returns the bid events for flop", func() {
			headerId, headerErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerErr).NotTo(HaveOccurred())
			flopKickLog := test_data.CreateTestLog(headerId, db)

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, blockOneHeader, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flopKickEvent := test_data.FlopKickModel()
			flopKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			flopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
			flopKickEvent.ColumnValues[constants.HeaderFK] = headerId
			flopKickEvent.ColumnValues[constants.LogFK] = flopKickLog.ID
			flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEvent})
			Expect(flopKickErr).NotTo(HaveOccurred())

			expectedBidEvents := test_helpers.BidEvent{
				BidId:     strconv.Itoa(fakeBidId),
				Lot:       flopKickEvent.ColumnValues["lot"].(string),
				BidAmount: flopKickEvent.ColumnValues["bid"].(string),
				Act:       "kick",
			}
			var actualBidEvents test_helpers.BidEvent
			queryErr := db.Get(&actualBidEvents,
				`SELECT bid_id, bid_amount, lot, act FROM api.flop_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flop_state 
    					FROM api.all_flops()))`)
			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(Equal(expectedBidEvents))
		})

		It("does not include bid events for a different flop", func() {
			headerId, headerErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerErr).NotTo(HaveOccurred())
			flopKickLog := test_data.CreateTestLog(headerId, db)

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, blockOneHeader, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flopKickEvent := test_data.FlopKickModel()
			flopKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			flopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
			flopKickEvent.ColumnValues[constants.HeaderFK] = headerId
			flopKickEvent.ColumnValues[constants.LogFK] = flopKickLog.ID
			flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEvent})
			Expect(flopKickErr).NotTo(HaveOccurred())

			blockTwo := blockOne + 1
			timestampTwo := timestampOne + 111111
			blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			irrelevantFlopKickLog := test_data.CreateTestLog(headerId, db)

			irrelevantBidId := fakeBidId + 9999999999999
			irrelevantFlopStorageValues := test_helpers.GetFlopStorageValues(2, irrelevantBidId)
			test_helpers.CreateFlop(db, blockTwoHeader, irrelevantFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(irrelevantBidId)), contractAddress)

			irrelevantFlopKickEvent := test_data.FlopKickModel()
			irrelevantFlopKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			irrelevantFlopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(irrelevantBidId)
			irrelevantFlopKickEvent.ColumnValues[constants.HeaderFK] = headerTwoId
			irrelevantFlopKickEvent.ColumnValues[constants.LogFK] = irrelevantFlopKickLog.ID
			flopKickErr = flopKickRepo.Create([]shared.InsertionModel{flopKickEvent})

			Expect(flopKickErr).NotTo(HaveOccurred())

			expectedBidEvents := test_helpers.BidEvent{
				BidId:     strconv.Itoa(fakeBidId),
				Lot:       flopKickEvent.ColumnValues["lot"].(string),
				BidAmount: flopKickEvent.ColumnValues["bid"].(string),
				Act:       "kick",
			}

			var actualBidEvents []test_helpers.BidEvent
			queryErr := db.Select(&actualBidEvents,
				`SELECT bid_id, bid_amount, lot, act FROM api.flop_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flop_state
    					FROM api.all_flops() WHERE bid_id = $1))`, fakeBidId)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBidEvents).To(ConsistOf(expectedBidEvents))
		})

		Describe("result pagination", func() {
			var (
				dentLot, dentBid int
				flopKickEvent    shared.InsertionModel
			)

			BeforeEach(func() {
				headerId, headerErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
				Expect(headerErr).NotTo(HaveOccurred())
				logId := test_data.CreateTestLog(headerId, db).ID

				flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
				test_helpers.CreateFlop(db, blockOneHeader, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

				flopKickEvent = test_data.FlopKickModel()
				flopKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
				flopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
				flopKickEvent.ColumnValues[constants.HeaderFK] = headerId
				flopKickEvent.ColumnValues[constants.LogFK] = logId
				flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEvent})

				Expect(flopKickErr).NotTo(HaveOccurred())

				blockTwo := blockOne + 1
				timestampTwo := timestampOne + 111111
				blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
				headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
				Expect(headerTwoErr).NotTo(HaveOccurred())
				logTwoId := test_data.CreateTestLog(headerTwoId, db).ID

				dentLot = rand.Int()
				dentBid = rand.Int()
				dentRepo := dent.DentRepository{}
				dentRepo.SetDB(db)
				flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
					Lot:             dentLot,
					BidAmount:       dentBid,
					DentRepo:        dentRepo,
					DentHeaderId:    headerTwoId,
					DentLogId:       logTwoId,
				})
				Expect(flopDentErr).NotTo(HaveOccurred())
			})

			It("limits result to most recent block if max_results argument is provided", func() {
				expectedBidEvents := test_helpers.BidEvent{
					BidId:     strconv.Itoa(fakeBidId),
					Lot:       strconv.Itoa(dentLot),
					BidAmount: strconv.Itoa(dentBid),
					Act:       "dent",
				}

				maxResults := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents,
					`SELECT bid_id, bid_amount, lot, act FROM api.flop_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flop_state
    					FROM api.all_flops() WHERE bid_id = $1), $2)`,
					fakeBidId, maxResults)

				Expect(queryErr).NotTo(HaveOccurred())
				Expect(actualBidEvents).To(ConsistOf(expectedBidEvents))
			})

			It("offsets results if offset is provided", func() {
				expectedBidEvents := test_helpers.BidEvent{
					BidId:     strconv.Itoa(fakeBidId),
					Lot:       flopKickEvent.ColumnValues["lot"].(string),
					BidAmount: flopKickEvent.ColumnValues["bid"].(string),
					Act:       "kick",
				}

				maxResults := 1
				resultOffset := 1
				var actualBidEvents []test_helpers.BidEvent
				queryErr := db.Select(&actualBidEvents,
					`SELECT bid_id, bid_amount, lot, act FROM api.flop_state_bid_events(
    					(SELECT (bid_id, guy, tic, "end", lot, bid, dealt, created, updated)::api.flop_state
    					FROM api.all_flops() WHERE bid_id = $1), $2, $3)`,
					fakeBidId, maxResults, resultOffset)

				Expect(queryErr).NotTo(HaveOccurred())
				Expect(actualBidEvents).To(ConsistOf(expectedBidEvents))
			})
		})
	})
})
