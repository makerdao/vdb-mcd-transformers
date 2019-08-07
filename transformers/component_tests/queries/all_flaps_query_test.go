package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("All flaps query", func() {
	var (
		db              *postgres.DB
		dealRepo        deal.DealRepository
		headerRepo      repositories.HeaderRepository
		contractAddress = "contract address"

		blockOne          = rand.Int()
		blockOneTimestamp = int64(111111111)

		blockTwo          = blockOne + 1
		blockTwoTimestamp = int64(222222222)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets the most recent flap for every bid id", func() {
		fakeBidIdOne := rand.Int()
		fakeBidIdTwo := fakeBidIdOne + 1

		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(blockOneTimestamp, int64(blockOne))
		_, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(blockTwoTimestamp, int64(blockTwo))
		blockTwoHeader.Hash = "blockTwoHeader"
		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidIdOne)
		test_helpers.CreateFlap(db, blockOneHeader, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidIdOne)
		test_helpers.CreateFlap(db, blockTwoHeader, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		flapStorageValuesThree := test_helpers.GetFlapStorageValues(3, fakeBidIdTwo)
		test_helpers.CreateFlap(db, blockTwoHeader, flapStorageValuesThree, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)

		var actualBids []test_helpers.FlapBid
		queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, gal, dealt, created, updated FROM api.all_flaps()`)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedBidOne := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdOne), "false", blockTwoHeader.Timestamp, blockOneHeader.Timestamp, flapStorageValuesTwo)
		expectedBidTwo := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidIdTwo), "false", blockTwoHeader.Timestamp, blockTwoHeader.Timestamp, flapStorageValuesThree)

		Expect(len(actualBids)).To(Equal(2))
		Expect(actualBids).To(ConsistOf([]test_helpers.FlapBid{
			expectedBidOne,
			expectedBidTwo,
		}))
	})
})
