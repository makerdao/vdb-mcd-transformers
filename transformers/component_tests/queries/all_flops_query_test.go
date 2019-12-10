package queries

import (
	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All flops query", func() {
	var (
		flopRepo               flop.FlopStorageRepository
		headerRepo             repositories.HeaderRepository
		contractAddress        = fakes.RandomString(42)
		blockOne, timestampOne int
		headerOne              core.Header
		diffID                 int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		flopRepo = flop.FlopStorageRepository{}
		flopRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)

		diffID = storage_helper.CreateFakeDiffRecord(db)
	})

	It("gets the most recent flop for every bid id", func() {
		fakeBidIdOne := rand.Int()
		fakeBidIdTwo := fakeBidIdOne + 1

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

		contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidIdOne,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlopKickHeaderId: headerOne.Id,
		})
		Expect(contextErr).NotTo(HaveOccurred())

		initialFlopOneStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidIdOne)
		test_helpers.CreateFlop(db, diffID, headerOne.Id, initialFlopOneStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		updatedFlopOneStorageValues := test_helpers.GetFlopStorageValues(2, fakeBidIdOne)
		test_helpers.CreateFlop(db, diffID, headerTwo.Id, updatedFlopOneStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		flopStorageValuesTwo := test_helpers.GetFlopStorageValues(3, fakeBidIdTwo)
		test_helpers.CreateFlop(db, diffID, headerTwo.Id, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)

		contextErr = test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidIdTwo,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlopKickHeaderId: headerTwo.Id,
		})
		Expect(contextErr).NotTo(HaveOccurred())

		var actualBids []test_helpers.FlopBid
		queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flops()`)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedBidOne := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdOne), "false", headerTwo.Timestamp, headerOne.Timestamp, updatedFlopOneStorageValues)
		expectedBidTwo := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdTwo), "false", headerTwo.Timestamp, headerTwo.Timestamp, flopStorageValuesTwo)

		Expect(len(actualBids)).To(Equal(2))
		Expect(actualBids).To(ConsistOf([]test_helpers.FlopBid{
			expectedBidOne,
			expectedBidTwo,
		}))
	})

	Describe("result pagination", func() {
		var (
			fakeBidIdOne, fakeBidIdTwo                 int
			flopStorageValuesOne, flopStorageValuesTwo map[string]interface{}
		)

		BeforeEach(func() {
			fakeBidIdOne = rand.Int()
			fakeBidIdTwo = fakeBidIdOne + 1

			flopStorageValuesOne = test_helpers.GetFlopStorageValues(1, fakeBidIdOne)
			test_helpers.CreateFlop(db, diffID, headerOne.Id, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

			flopStorageValuesTwo = test_helpers.GetFlopStorageValues(2, fakeBidIdTwo)
			test_helpers.CreateFlop(db, diffID, headerOne.Id, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)
		})

		It("limits results if max_results argument is provided", func() {
			contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidIdTwo,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlopKickHeaderId: headerOne.Id,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			maxResults := 1
			var actualBids []test_helpers.FlopBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flops($1)`,
				maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdTwo), "false", headerOne.Timestamp,
				headerOne.Timestamp, flopStorageValuesTwo)
			Expect(actualBids).To(Equal([]test_helpers.FlopBid{expectedBid}))
		})

		It("offsets results if offset is provided", func() {
			contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidIdOne,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlopKickHeaderId: headerOne.Id,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			maxResults := 1
			resultOffset := 1
			var actualBids []test_helpers.FlopBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flops($1, $2)`,
				maxResults, resultOffset)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdOne), "false", headerOne.Timestamp,
				headerOne.Timestamp, flopStorageValuesOne)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})
	})
})
