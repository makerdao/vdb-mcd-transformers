package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dent"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flop_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flop computed columns", func() {
	var (
		db                     *postgres.DB
		flopKickRepo           flop_kick.FlopKickRepository
		headerRepo             repositories.HeaderRepository
		contractAddress        = fakes.RandomString(42)
		fakeBidId              = rand.Int()
		blockOne, timestampOne int
		headerOne              core.Header
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flopKickRepo = flop_kick.FlopKickRepository{}
		flopKickRepo.SetDB(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("flop_bid_events", func() {
		It("returns the bid events for flop", func() {
			flopKickLog := test_data.CreateTestLog(headerOne.Id, db)

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flopKickEvent := test_data.FlopKickModel()
			flopKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			flopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
			flopKickEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
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
			flopKickLog := test_data.CreateTestLog(headerOne.Id, db)

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flopKickEvent := test_data.FlopKickModel()
			flopKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			flopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
			flopKickEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			flopKickEvent.ColumnValues[constants.LogFK] = flopKickLog.ID
			flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEvent})
			Expect(flopKickErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			irrelevantFlopKickLog := test_data.CreateTestLog(headerTwo.Id, db)

			irrelevantBidId := fakeBidId + 9999999999999
			irrelevantFlopStorageValues := test_helpers.GetFlopStorageValues(2, irrelevantBidId)
			test_helpers.CreateFlop(db, headerTwo, irrelevantFlopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(irrelevantBidId)), contractAddress)

			irrelevantFlopKickEvent := test_data.FlopKickModel()
			irrelevantFlopKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			irrelevantFlopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(irrelevantBidId)
			irrelevantFlopKickEvent.ColumnValues[constants.HeaderFK] = headerTwo.Id
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
				logId := test_data.CreateTestLog(headerOne.Id, db).ID

				flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
				test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

				flopKickEvent = test_data.FlopKickModel()
				flopKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
				flopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
				flopKickEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
				flopKickEvent.ColumnValues[constants.LogFK] = logId
				flopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEvent})

				Expect(flopKickErr).NotTo(HaveOccurred())

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
				logTwoId := test_data.CreateTestLog(headerTwo.Id, db).ID

				dentLot = rand.Int()
				dentBid = rand.Int()
				dentRepo := dent.Repository{}
				dentRepo.SetDB(db)
				flopDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
					Db:              db,
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
					Lot:             dentLot,
					BidAmount:       dentBid,
					DentRepo:        dentRepo,
					DentHeaderId:    headerTwo.Id,
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
