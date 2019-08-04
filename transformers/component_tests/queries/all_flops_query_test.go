package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flop"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("All flops query", func() {
	var (
		db              *postgres.DB
		flopRepo        flop.FlopStorageRepository
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
		flopRepo = flop.FlopStorageRepository{}
		flopRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets the most recent flop for every bid id", func() {
		fakeBidIdOne := rand.Int()
		fakeBidIdTwo := fakeBidIdOne + 1

		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(blockOneTimestamp, int64(blockOne))
		_, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(blockTwoTimestamp, int64(blockTwo))
		blockTwoHeader.Hash = "blockTwoHeader"
		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		initialFlopOneStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidIdOne)
		test_helpers.CreateFlop(db, blockOneHeader, initialFlopOneStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		updatedFlopOneStorageValues := test_helpers.GetFlopStorageValues(2, fakeBidIdOne)
		test_helpers.CreateFlop(db, blockTwoHeader, updatedFlopOneStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		flopStorageValuesTwo := test_helpers.GetFlopStorageValues(3, fakeBidIdTwo)
		test_helpers.CreateFlop(db, blockTwoHeader, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)

		var actualBids []test_helpers.FlopBid
		queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flops()`)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedBidOne := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdOne), "false", blockTwoHeader.Timestamp, blockOneHeader.Timestamp, updatedFlopOneStorageValues)
		expectedBidTwo := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdTwo), "false", blockTwoHeader.Timestamp, blockTwoHeader.Timestamp, flopStorageValuesTwo)

		Expect(len(actualBids)).To(Equal(2))
		Expect(actualBids).To(ConsistOf([]test_helpers.FlopBid{
			expectedBidOne,
			expectedBidTwo,
		}))
	})

})
