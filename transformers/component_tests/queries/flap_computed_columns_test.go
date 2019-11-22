package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flap_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/tend"
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

var _ = Describe("Flap computed columns", func() {
	var (
		db                     *postgres.DB
		flapKickRepo           flap_kick.FlapKickRepository
		headerRepo             repositories.HeaderRepository
		contractAddress        = fakes.FakeAddress.Hex()
		fakeBidId              = rand.Int()
		blockOne, timestampOne int
		headerOne              core.Header
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flapKickRepo = flap_kick.FlapKickRepository{}
		flapKickRepo.SetDB(db)

		blockOne = rand.Int()
		timestampOne := int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("flap_bid_events", func() {
		It("returns the bid events for a flap", func() {
			flapKickLog := test_data.CreateTestLog(headerOne.Id, db)

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flapKickEvent := test_data.FlapKickModel()
			flapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			flapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
			flapKickEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			flapKickEvent.ColumnValues[constants.LogFK] = flapKickLog.ID
			flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEvent})
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
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flapKickEvent := test_data.FlapKickModel()
			flapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			flapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
			flapKickEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
			flapKickEvent.ColumnValues[constants.LogFK] = flapKickLog.ID
			flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEvent})
			Expect(flapKickErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			irrelevantFlipKickLog := test_data.CreateTestLog(headerTwo.Id, db)

			irrelevantBidId := fakeBidId + 9999999999999
			irrelevantFlapStorageValues := test_helpers.GetFlapStorageValues(2, irrelevantBidId)
			test_helpers.CreateFlap(db, headerTwo, irrelevantFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(irrelevantBidId)), contractAddress)

			irrelevantFlapKickEvent := test_data.FlapKickModel()
			irrelevantFlapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
			irrelevantFlapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(irrelevantBidId)
			irrelevantFlapKickEvent.ColumnValues[constants.HeaderFK] = headerTwo.Id
			irrelevantFlapKickEvent.ColumnValues[constants.LogFK] = irrelevantFlipKickLog.ID

			flapKickErr = flapKickRepo.Create([]shared.InsertionModel{irrelevantFlapKickEvent})
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
				flapKickEvent    shared.InsertionModel
			)

			BeforeEach(func() {
				logId := test_data.CreateTestLog(headerOne.Id, db).ID

				flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
				test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

				flapKickEvent = test_data.FlapKickModel()
				flapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
				flapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
				flapKickEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
				flapKickEvent.ColumnValues[constants.LogFK] = logId
				flapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEvent})
				Expect(flapKickErr).NotTo(HaveOccurred())

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
				logTwoId := test_data.CreateTestLog(headerTwo.Id, db).ID

				tendBid = rand.Int()
				tendLot = rand.Int()
				tendRepo := tend.TendRepository{}
				tendRepo.SetDB(db)
				flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
					Lot:             tendLot,
					BidAmount:       tendBid,
					TendRepo:        tendRepo,
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
					Lot:             flapKickEvent.ColumnValues["lot"].(string),
					BidAmount:       flapKickEvent.ColumnValues["bid"].(string),
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
