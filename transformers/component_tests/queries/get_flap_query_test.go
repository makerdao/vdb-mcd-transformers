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

var _ = Describe("Get flap query", func() {
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

	It("gets the specified flap", func() {
		headerId, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidId)
		test_helpers.CreateFlap(db, blockOneHeader, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		err := test_helpers.SetUpBidContext(test_helpers.BidSetupData{
			Db:              db,
			BidId:           fakeBidId,
			ContractAddress: contractAddress,
			DealRepo:        dealRepo,
			DealHeaderId:    headerId,
			Dealt:           true,
		})
		Expect(err).NotTo(HaveOccurred())

		expectedBid := test_helpers.BidFromValues(strconv.Itoa(fakeBidId), "true", blockOneHeader.Timestamp, blockOneHeader.Timestamp, flapStorageValuesOne)

		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())
		flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidId)
		test_helpers.CreateFlap(db, blockTwoHeader, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		var actualBid test_helpers.Bid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, gal, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockOne)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	It("gets created and updated blocks", func() {
		_, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		err := test_helpers.SetUpBidContext(test_helpers.BidSetupData{
			Db:              db,
			BidId:           fakeBidId,
			ContractAddress: contractAddress,
			DealRepo:        dealRepo,
			DealHeaderId:    headerTwoId,
			Dealt:           true,
		})
		Expect(err).NotTo(HaveOccurred())

		flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidId)
		test_helpers.CreateFlap(db, blockOneHeader, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidId)
		test_helpers.CreateFlap(db, blockTwoHeader, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		expectedBid := test_helpers.BidFromValues(strconv.Itoa(fakeBidId), "true", blockTwoHeader.Timestamp, blockOneHeader.Timestamp, flapStorageValuesTwo)

		var actualBid test_helpers.Bid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, gal, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockTwo)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	Describe("Dealt", func() {
		It("is false if no deal events", func() {
			header := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
			headerId, headerErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, header, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flapBidContextErr := test_helpers.SetUpBidContext(test_helpers.BidSetupData{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderId:    headerId,
				Dealt:           false,
			})
			Expect(flapBidContextErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.BidFromValues(strconv.Itoa(fakeBidId), "false", header.Timestamp, header.Timestamp, flapStorageValues)

			var actualBid test_helpers.Bid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, gal, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})

		It("is false if deal event in later block", func() {
			_, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerOneErr).NotTo(HaveOccurred())
			flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, blockOneHeader, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidId)
			test_helpers.CreateFlap(db, blockTwoHeader, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flapBidContextErr := test_helpers.SetUpBidContext(test_helpers.BidSetupData{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderId:    headerTwoId,
				Dealt:           true,
			})
			Expect(flapBidContextErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.BidFromValues(strconv.Itoa(fakeBidId), "false", blockOneHeader.Timestamp, blockOneHeader.Timestamp, flapStorageValuesOne)
			var actualBid test_helpers.Bid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, gal, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})
	})
})
