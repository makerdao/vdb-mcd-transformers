package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/events/tend"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("Flap computed columns", func() {
	var (
		db              *postgres.DB
		flapKickRepo    flap_kick.FlapKickRepository
		headerRepo      repositories.HeaderRepository
		contractAddress = fakes.FakeAddress.Hex()

		fakeBidId      = rand.Int()
		blockOne       = rand.Int()
		blockOneHeader = fakes.GetFakeHeader(int64(blockOne))
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		flapKickRepo = flap_kick.FlapKickRepository{}
		flapKickRepo.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("flap_bid_events", func() {
		It("returns the bid events for a flap", func() {
			headerId, headerErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerErr).NotTo(HaveOccurred())
			flapKickLog := test_data.CreateTestLog(headerId, db)

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, blockOneHeader, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flapKickEvent := test_data.FlapKickModel
			flapKickEvent.ContractAddress = contractAddress
			flapKickEvent.BidId = strconv.Itoa(fakeBidId)
			flapKickEvent.HeaderID = headerId
			flapKickEvent.LogID = flapKickLog.ID
			flapKickErr := flapKickRepo.Create([]interface{}{flapKickEvent})
			Expect(flapKickErr).NotTo(HaveOccurred())

			expectedBidEvents := test_helpers.BidEvent{
				BidId:           strconv.Itoa(fakeBidId),
				Lot:             flapKickEvent.Lot,
				BidAmount:       flapKickEvent.Bid,
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
			headerId, headerErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerErr).NotTo(HaveOccurred())
			flapKickLog := test_data.CreateTestLog(headerId, db)

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, blockOneHeader, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flapKickEvent := test_data.FlapKickModel
			flapKickEvent.ContractAddress = contractAddress
			flapKickEvent.BidId = strconv.Itoa(fakeBidId)
			flapKickEvent.HeaderID = headerId
			flapKickEvent.LogID = flapKickLog.ID
			flapKickErr := flapKickRepo.Create([]interface{}{flapKickEvent})
			Expect(flapKickErr).NotTo(HaveOccurred())

			blockTwo := blockOne + 1
			blockTwoHeader := fakes.GetFakeHeader(int64(blockTwo))
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			irrelevantFlipKickLog := test_data.CreateTestLog(headerTwoId, db)

			irrelevantBidId := fakeBidId + 9999999999999
			irrelevantFlapStorageValues := test_helpers.GetFlapStorageValues(2, irrelevantBidId)
			test_helpers.CreateFlap(db, blockTwoHeader, irrelevantFlapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(irrelevantBidId)), contractAddress)

			irrelevantFlapKickEvent := test_data.FlapKickModel
			irrelevantFlapKickEvent.ContractAddress = contractAddress
			irrelevantFlapKickEvent.BidId = strconv.Itoa(irrelevantBidId)
			irrelevantFlapKickEvent.HeaderID = headerTwoId
			irrelevantFlapKickEvent.LogID = irrelevantFlipKickLog.ID

			flapKickErr = flapKickRepo.Create([]interface{}{irrelevantFlapKickEvent})
			Expect(flapKickErr).NotTo(HaveOccurred())

			expectedBidEvents := test_helpers.BidEvent{
				BidId:           strconv.Itoa(fakeBidId),
				Lot:             flapKickEvent.Lot,
				BidAmount:       flapKickEvent.Bid,
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
				flapKickEvent    flap_kick.FlapKickModel
			)

			BeforeEach(func() {
				headerId, headerErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
				Expect(headerErr).NotTo(HaveOccurred())
				logId := test_data.CreateTestLog(headerId, db).ID

				flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
				test_helpers.CreateFlap(db, blockOneHeader, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

				flapKickEvent = test_data.FlapKickModel
				flapKickEvent.ContractAddress = contractAddress
				flapKickEvent.BidId = strconv.Itoa(fakeBidId)
				flapKickEvent.HeaderID = headerId
				flapKickEvent.LogID = logId
				flapKickErr := flapKickRepo.Create([]interface{}{flapKickEvent})
				Expect(flapKickErr).NotTo(HaveOccurred())

				blockTwo := blockOne + 1
				blockTwoHeader := fakes.GetFakeHeader(int64(blockTwo))
				headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
				Expect(headerTwoErr).NotTo(HaveOccurred())
				logTwoId := test_data.CreateTestLog(headerTwoId, db).ID

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
					TendHeaderId:    headerTwoId,
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
					Lot:             flapKickEvent.Lot,
					BidAmount:       flapKickEvent.Bid,
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
