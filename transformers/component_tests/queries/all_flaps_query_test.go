package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All flaps query", func() {
	var (
		headerRepo      repositories.HeaderRepository
		contractAddress = "contract address"

		blockOne, timestampOne int
		headerOne              core.Header
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	It("gets the most recent flap for every bid id", func() {
		fakeBidIdOne := rand.Int()
		fakeBidIdTwo := fakeBidIdOne + 1

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

		contextErr := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidIdOne,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlapKickHeaderId: headerOne.Id,
		})
		Expect(contextErr).NotTo(HaveOccurred())

		flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidIdOne)
		test_helpers.CreateFlap(db, 0, headerOne, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidIdOne)
		test_helpers.CreateFlap(db, 0, headerTwo, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		contextErr = test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidIdTwo,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlapKickHeaderId: headerTwo.Id,
		})
		Expect(contextErr).NotTo(HaveOccurred())
		flapStorageValuesThree := test_helpers.GetFlapStorageValues(3, fakeBidIdTwo)
		test_helpers.CreateFlap(db, 0, headerTwo, flapStorageValuesThree, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)

		var actualBids []test_helpers.FlapBid
		queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flaps()`)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedBidOne := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdOne), "false", headerTwo.Timestamp, headerOne.Timestamp, flapStorageValuesTwo)
		expectedBidTwo := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdTwo), "false", headerTwo.Timestamp, headerTwo.Timestamp, flapStorageValuesThree)

		Expect(len(actualBids)).To(Equal(2))
		Expect(actualBids).To(ConsistOf([]test_helpers.FlapBid{
			expectedBidOne,
			expectedBidTwo,
		}))
	})

	Describe("result pagination", func() {
		var (
			fakeBidIdOne, fakeBidIdTwo                 int
			headerTwo                                  core.Header
			flapStorageValuesOne, flapStorageValuesTwo map[string]interface{}
		)

		BeforeEach(func() {
			fakeBidIdOne = rand.Int()
			fakeBidIdTwo = fakeBidIdOne + 1

			contextErr := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidIdOne,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlapKickHeaderId: headerOne.Id,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			flapStorageValuesOne = test_helpers.GetFlapStorageValues(2, fakeBidIdOne)
			test_helpers.CreateFlap(db, 0, headerOne, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

			headerTwo = createHeader(blockOne+1, timestampOne+1, headerRepo)

			contextErr = test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidIdTwo,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlapKickHeaderId: headerTwo.Id,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			flapStorageValuesTwo = test_helpers.GetFlapStorageValues(3, fakeBidIdTwo)
			test_helpers.CreateFlap(db, 0, headerTwo, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)
		})

		It("limits results if max_results argument is provided", func() {
			maxResults := 1
			var actualBids []test_helpers.FlapBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flaps($1)`, maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdTwo), "false",
				headerTwo.Timestamp, headerTwo.Timestamp, flapStorageValuesTwo)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})

		It("offsets results if offset is provided", func() {
			maxResults := 1
			resultOffset := 1
			var actualBids []test_helpers.FlapBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flaps($1, $2)`,
				maxResults, resultOffset)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdOne), "false",
				headerOne.Timestamp, headerOne.Timestamp, flapStorageValuesOne)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})
	})
})
