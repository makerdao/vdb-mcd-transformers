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
	"time"
)

var _ = Describe("get flop query", func() {
	var (
		db              *postgres.DB
		dealRepo        deal.DealRepository
		headerRepo      repositories.HeaderRepository
		contractAddress = "contract address"

		fakeBidId      = rand.Int()
		blockOne       = rand.Int()
		timestampOne   = int(rand.Int31())
		hashOne        = "hashOne"
		blockOneHeader = fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))

		blockTwo       = blockOne + 1
		timestampTwo   = timestampOne + 1000
		hashTwo        = "hashTwo"
		blockTwoHeader = fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))

		flopStorageValuesOne = test_helpers.GetFlopStorageValues(1, fakeBidId)
		flopStorageValuesTwo = test_helpers.GetFlopStorageValues(2, fakeBidId)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		blockOneHeader.Hash = hashOne
		blockTwoHeader.Hash = hashTwo
		rand.Seed(time.Now().UnixNano())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets the specified flop", func() {
		headerId, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		test_helpers.CreateFlop(db, blockOneHeader, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		err := test_helpers.CreateDeal(test_helpers.DealCreationInput{
			Db:              db,
			BidId:           fakeBidId,
			ContractAddress: contractAddress,
			DealRepo:        dealRepo,
			DealHeaderId:    headerId,
		})
		Expect(err).NotTo(HaveOccurred())

		dealt := "true"
		expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidId), dealt, blockOneHeader.Timestamp, blockOneHeader.Timestamp, flopStorageValuesOne)

		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())
		test_helpers.CreateFlop(db, blockTwoHeader, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		var actualBid test_helpers.FlopBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidId, blockOne)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	It("gets created and updated blocks", func() {
		_, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		flopBidContextErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
			Db:              db,
			BidId:           fakeBidId,
			ContractAddress: contractAddress,
			DealRepo:        dealRepo,
			DealHeaderId:    headerTwoId,
		})
		Expect(flopBidContextErr).NotTo(HaveOccurred())

		test_helpers.CreateFlop(db, blockOneHeader, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)
		test_helpers.CreateFlop(db, blockTwoHeader, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		dealt := "true"
		expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidId), dealt, blockTwoHeader.Timestamp, blockOneHeader.Timestamp, flopStorageValuesTwo)

		var actualBid test_helpers.FlopBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidId, blockTwo)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	Describe("dealt", func() {
		It("is false if no deal events", func() {
			blockNumber := rand.Int()
			timestamp := int(rand.Int31())

			header := fakes.GetFakeHeaderWithTimestamp(int64(timestamp), int64(blockNumber))
			_, headerErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, header, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			dealt := "false"
			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidId), dealt, header.Timestamp, header.Timestamp, flopStorageValues)

			var actualBid test_helpers.FlopBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidId, blockNumber)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})

		It("is false if deal event in later block", func() {
			_, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerOneErr).NotTo(HaveOccurred())

			test_helpers.CreateFlop(db, blockOneHeader, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			test_helpers.CreateFlop(db, blockTwoHeader, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flopBidContextErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderId:    headerTwoId,
			})
			Expect(flopBidContextErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlopBidFromValues(
				strconv.Itoa(fakeBidId), "false", blockOneHeader.Timestamp, blockOneHeader.Timestamp, flopStorageValuesOne)

			var actualBid test_helpers.FlopBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidId, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})
	})
})
