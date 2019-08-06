package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
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
		flapKickRepo    flap_kick.FlapKickRepository
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
		flapKickRepo = flap_kick.FlapKickRepository{}
		flapKickRepo.SetDB(db)
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
		dealLog := test_data.CreateTestLog(headerId, db)

		err := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderId:    headerId,
				DealLogID:       dealLog.ID,
			},
			Dealt:            true,
			FlapKickRepo:     flapKickRepo,
			FlapKickHeaderID: headerId,
		})
		Expect(err).NotTo(HaveOccurred())

		flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidId)
		test_helpers.CreateFlap(db, blockOneHeader, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())
		flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidId)
		test_helpers.CreateFlap(db, blockTwoHeader, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidId), "true", blockOneHeader.Timestamp, blockOneHeader.Timestamp, flapStorageValuesOne)

		var actualBid test_helpers.FlapBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockOne)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	It("gets the correct created and updated timestamps based on the requested block", func() {
		headerOneId, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())
		dealLog := test_data.CreateTestLog(headerTwoId, db)

		err := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderId:    headerTwoId,
				DealLogID:       dealLog.ID,
			},
			Dealt:            true,
			FlapKickRepo:     flapKickRepo,
			FlapKickHeaderID: headerOneId,
		})
		Expect(err).NotTo(HaveOccurred())

		flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidId)
		test_helpers.CreateFlap(db, blockTwoHeader, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidId)
		test_helpers.CreateFlap(db, blockOneHeader, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		// creating another block + updated storage values to ensure that get_flap uses the specified block
		blockThree := blockTwo + 1
		timestampThree := timestampTwo + 1000
		blockThreeHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampThree), int64(blockThree))
		flapStorageValuesThree := test_helpers.GetFlapStorageValues(3, fakeBidId)
		test_helpers.CreateFlap(db, blockThreeHeader, flapStorageValuesThree, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidId), "true", blockTwoHeader.Timestamp, blockOneHeader.Timestamp, flapStorageValuesTwo)

		var actualBid test_helpers.FlapBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockTwo)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	Describe("Dealt", func() {
		It("is false if no deal events", func() {
			header := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
			headerId, headerErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			err := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlapKickRepo:     flapKickRepo,
				FlapKickHeaderID: headerId,
			})
			Expect(err).NotTo(HaveOccurred())

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, header, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidId), "false", header.Timestamp, header.Timestamp, flapStorageValues)

			var actualBid test_helpers.FlapBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})

		It("is false if deal event in later block", func() {
			headerId, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerOneErr).NotTo(HaveOccurred())

			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
			Expect(headerTwoErr).NotTo(HaveOccurred())

			dealLog := test_data.CreateTestLog(headerTwoId, db)
			// todo: change how created timestamp is retrieved so this test can pass if we set up flap bid context after storage vals are created
			err := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
					DealRepo:        dealRepo,
					DealHeaderId:    headerTwoId,
					DealLogID:       dealLog.ID,
				},
				Dealt:            true,
				FlapKickRepo:     flapKickRepo,
				FlapKickHeaderID: headerId,
			})
			Expect(err).NotTo(HaveOccurred())

			flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, blockOneHeader, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidId)
			test_helpers.CreateFlap(db, blockTwoHeader, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidId), "false", blockOneHeader.Timestamp, blockOneHeader.Timestamp, flapStorageValuesOne)
			var actualBid test_helpers.FlapBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})
	})
})
