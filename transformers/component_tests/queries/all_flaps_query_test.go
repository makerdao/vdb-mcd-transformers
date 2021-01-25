package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All flaps query", func() {
	var (
		headerRepo             datastore.HeaderRepository
		contractAddress        = fakes.FakeAddress.Hex()
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
		test_helpers.CreateFlap(db, headerOne, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidIdOne)
		test_helpers.CreateFlap(db, headerTwo, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

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
		test_helpers.CreateFlap(db, headerTwo, flapStorageValuesThree, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)

		var actualBids []test_helpers.FlapBid
		queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated, flap_address FROM api.all_flaps()`)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedBidOne := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdOne), "false", headerTwo.Timestamp, headerOne.Timestamp, contractAddress, flapStorageValuesTwo)
		expectedBidTwo := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdTwo), "false", headerTwo.Timestamp, headerTwo.Timestamp, contractAddress, flapStorageValuesThree)

		Expect(len(actualBids)).To(Equal(2))
		Expect(actualBids).To(ConsistOf([]test_helpers.FlapBid{
			expectedBidOne,
			expectedBidTwo,
		}))
	})

	Context("when there are multiple flap contracts", func() {
		var (
			expectedBidOne         test_helpers.FlapBid
			anotherContractAddress = fakes.AnotherFakeAddress.Hex()
			fakeBidIdOne           int
		)

		BeforeEach(func() {
			fakeBidIdOne = rand.Int()
			flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidIdOne)
			test_helpers.CreateFlap(db, headerOne, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

			var setupErr error
			setupErr = test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidIdOne,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlapKickHeaderId: headerOne.Id,
			})
			Expect(setupErr).NotTo(HaveOccurred())

			expectedBidOne = test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdOne), "false", headerOne.Timestamp, headerOne.Timestamp, contractAddress, flapStorageValuesOne)
		})

		It("gets all the flap bids when they have different bid ids", func() {
			fakeBidIdTwo := fakeBidIdOne + 1
			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

			flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidIdTwo)
			test_helpers.CreateFlap(db, headerTwo, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdTwo)), anotherContractAddress)

			contextErr := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidIdTwo,
					ContractAddress: anotherContractAddress,
				},
				Dealt:            false,
				FlapKickHeaderId: headerTwo.Id,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			expectedBidTwo := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdTwo), "false", headerTwo.Timestamp, headerTwo.Timestamp, anotherContractAddress, flapStorageValuesTwo)

			var actualBids []test_helpers.FlapBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated, flap_address FROM api.all_flaps()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(len(actualBids)).To(Equal(2))
			Expect(actualBids).To(ContainElement(expectedBidOne))
			Expect(actualBids).To(ContainElement(expectedBidTwo))
		})

		It("gets the all flap bids when they have the same bid ids on different contracts", func() {
			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

			flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidIdOne)
			test_helpers.CreateFlap(db, headerTwo, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), anotherContractAddress)

			contextErr := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidIdOne,
					ContractAddress: anotherContractAddress,
				},
				Dealt:            false,
				FlapKickHeaderId: headerTwo.Id,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			expectedBidTwo := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdOne), "false", headerTwo.Timestamp, headerTwo.Timestamp, anotherContractAddress, flapStorageValuesTwo)

			var actualBids []test_helpers.FlapBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated, flap_address FROM api.all_flaps()`)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(len(actualBids)).To(Equal(2))
			Expect(actualBids).To(ContainElement(expectedBidOne))
			Expect(actualBids).To(ContainElement(expectedBidTwo))
		})
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
			test_helpers.CreateFlap(db, headerOne, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

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
			test_helpers.CreateFlap(db, headerTwo, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)
		})

		It("limits results if max_results argument is provided", func() {
			maxResults := 1
			var actualBids []test_helpers.FlapBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated, flap_address FROM api.all_flaps($1)`, maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdTwo), "false",
				headerTwo.Timestamp, headerTwo.Timestamp, contractAddress, flapStorageValuesTwo)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})

		It("offsets results if offset is provided", func() {
			maxResults := 1
			resultOffset := 1
			var actualBids []test_helpers.FlapBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated, flap_address FROM api.all_flaps($1, $2)`,
				maxResults, resultOffset)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdOne), "false",
				headerOne.Timestamp, headerOne.Timestamp, contractAddress, flapStorageValuesOne)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})
	})
})
