package queries

import (
	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get flap query", func() {
	var (
		headerRepo                 repositories.HeaderRepository
		contractAddress            = fakes.RandomString(42)
		fakeBidId                  = rand.Int()
		blockOne, blockTwo         int
		timestampOne, timestampTwo int
		headerOne, headerTwo       core.Header
		diffID int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		blockTwo = blockOne + 1
		timestampOne = int(rand.Int31())
		timestampTwo = timestampOne + 1
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
		headerTwo = createHeader(blockTwo, timestampTwo, headerRepo)

		diffID = storage_helper.CreateDiffRecord(db)
	})

	It("gets the specified flap", func() {
		err := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealHeaderId:    headerOne.Id,
			},
			Dealt:            true,
			FlapKickHeaderId: headerOne.Id,
		})
		Expect(err).NotTo(HaveOccurred())

		flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidId)
		test_helpers.CreateFlap(db, diffID, headerOne.Id, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidId)
		test_helpers.CreateFlap(db, diffID, headerTwo.Id, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidId), "true", headerOne.Timestamp, headerOne.Timestamp, flapStorageValuesOne)

		var actualBid test_helpers.FlapBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockOne)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	It("gets the correct created and updated timestamps based on the requested block", func() {
		err := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealHeaderId:    headerTwo.Id,
			},
			Dealt:            true,
			FlapKickHeaderId: headerOne.Id,
		})
		Expect(err).NotTo(HaveOccurred())

		flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidId)
		test_helpers.CreateFlap(db, diffID, headerTwo.Id, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidId)
		test_helpers.CreateFlap(db, diffID, headerOne.Id, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		// creating another block + updated storage values to ensure that get_flap uses the specified block
		headerThree := createHeader(blockTwo+1, timestampTwo+1, headerRepo)
		flapStorageValuesThree := test_helpers.GetFlapStorageValues(3, fakeBidId)
		test_helpers.CreateFlap(db, diffID, headerThree.Id, flapStorageValuesThree, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidId), "true", headerTwo.Timestamp, headerOne.Timestamp, flapStorageValuesTwo)

		var actualBid test_helpers.FlapBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockTwo)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	Describe("Dealt", func() {
		It("is false if no deal events", func() {
			err := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlapKickHeaderId: headerOne.Id,
			})
			Expect(err).NotTo(HaveOccurred())

			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, diffID, headerOne.Id, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidId), "false", headerOne.Timestamp, headerOne.Timestamp, flapStorageValues)

			var actualBid test_helpers.FlapBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})

		It("is false if deal event in later block", func() {
			// todo: change how created timestamp is retrieved so this test can pass if we set up flap bid context after storage vals are created
			err := test_helpers.SetUpFlapBidContext(test_helpers.FlapBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
					DealHeaderId:    headerTwo.Id,
				},
				Dealt:            true,
				FlapKickHeaderId: headerOne.Id,
			})
			Expect(err).NotTo(HaveOccurred())

			flapStorageValuesOne := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, diffID, headerOne.Id, flapStorageValuesOne, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			flapStorageValuesTwo := test_helpers.GetFlapStorageValues(2, fakeBidId)
			test_helpers.CreateFlap(db, diffID, headerTwo.Id, flapStorageValuesTwo, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidId), "false", headerOne.Timestamp, headerOne.Timestamp, flapStorageValuesOne)
			var actualBid test_helpers.FlapBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flap($1, $2)`, fakeBidId, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})
	})
})
