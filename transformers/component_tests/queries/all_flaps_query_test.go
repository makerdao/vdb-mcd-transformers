package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/deal"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flap_kick"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All flaps query", func() {
	var (
		db              *postgres.DB
		flapKickRepo    flap_kick.Repository
		dealRepo        deal.Repository
		headerRepo      repositories.HeaderRepository
		contractAddress = "contract address"

		blockOne, timestampOne int
		headerOne              core.Header
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		flapKickRepo = flap_kick.Repository{}
		flapKickRepo.SetDB(db)
		dealRepo = deal.Repository{}
		dealRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets the most recent flap for every bid id", func() {
		fakeBidIdOne := rand.Int()
		fakeBidIdTwo := fakeBidIdOne + 1

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

		contextErr := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidIdOne,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlapKickRepo:     flapKickRepo,
			FlapKickHeaderId: headerOne.Id,
		})
		Expect(contextErr).NotTo(HaveOccurred())

		flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidIdOne)
		test_helpers.CreateFlap(db, headerOne, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidIdOne)
		test_helpers.CreateFlap(db, headerTwo, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		contextErr = test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidIdTwo,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlapKickRepo:     flapKickRepo,
			FlapKickHeaderId: headerTwo.Id,
		})
		Expect(contextErr).NotTo(HaveOccurred())
		flapStorageValuesThree := test_helpers.GetFlapStorageValues(3, fakeBidIdTwo)
		test_helpers.CreateFlap(db, headerTwo, flapStorageValuesThree, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)

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
					Db:              db,
					BidId:           fakeBidIdOne,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlapKickRepo:     flapKickRepo,
				FlapKickHeaderId: headerOne.Id,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			flapStorageValuesOne = test_helpers.GetFlapStorageValues(2, fakeBidIdOne)
			test_helpers.CreateFlap(db, headerOne, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

			headerTwo = createHeader(blockOne+1, timestampOne+1, headerRepo)

			contextErr = test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidId:           fakeBidIdTwo,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlapKickRepo:     flapKickRepo,
				FlapKickHeaderId: headerTwo.Id,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			flapStorageValuesTwo = test_helpers.GetFlapStorageValues(3, fakeBidIdTwo)
			test_helpers.CreateFlap(db, headerTwo, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)
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
