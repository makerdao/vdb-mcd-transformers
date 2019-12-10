package queries

import (
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

var _ = Describe("get flop query", func() {
	var (
		headerRepo                 repositories.HeaderRepository
		contractAddress            = fakes.RandomString(42)
		fakeBidId                  = rand.Int()
		blockOne, blockTwo         int
		timestampOne, timestampTwo int
		headerOne, headerTwo       core.Header
		flopStorageValuesOne       = test_helpers.GetFlopStorageValues(1, fakeBidId)
		flopStorageValuesTwo       = test_helpers.GetFlopStorageValues(2, fakeBidId)
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
	})

	It("gets the specified flop", func() {
		err := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealHeaderId:    headerOne.Id,
			},
			Dealt:            true,
			FlopKickHeaderId: headerOne.Id,
		})
		Expect(err).NotTo(HaveOccurred())

		test_helpers.CreateFlop(db, headerOne, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)
		test_helpers.CreateFlop(db, headerTwo, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidId), "true", headerOne.Timestamp, headerOne.Timestamp, flopStorageValuesOne)

		var actualBid test_helpers.FlopBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidId, blockOne)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	It("gets created and updated blocks", func() {
		err := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealHeaderId:    headerTwo.Id,
			},
			Dealt:            true,
			FlopKickHeaderId: headerOne.Id,
		})
		Expect(err).NotTo(HaveOccurred())

		test_helpers.CreateFlop(db, headerOne, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)
		test_helpers.CreateFlop(db, headerTwo, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		headerThree := createHeader(blockTwo+1, timestampTwo+1, headerRepo)
		flopStorageValuesThree := test_helpers.GetFlopStorageValues(3, fakeBidId)
		test_helpers.CreateFlop(db, headerThree, flopStorageValuesThree, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidId), "true", headerTwo.Timestamp, headerOne.Timestamp, flopStorageValuesTwo)

		var actualBid test_helpers.FlopBid
		queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidId, blockTwo)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	Describe("dealt", func() {
		It("is false if no deal events", func() {
			err := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlopKickHeaderId: headerOne.Id,
			})
			Expect(err).NotTo(HaveOccurred())

			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, headerOne, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidId), "false", headerOne.Timestamp, headerOne.Timestamp, flopStorageValues)

			var actualBid test_helpers.FlopBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidId, headerOne.BlockNumber)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})

		It("is false if deal event in later block", func() {
			err := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
					DealHeaderId:    headerTwo.Id,
				},
				Dealt:            true,
				FlopKickHeaderId: headerOne.Id,
			})
			Expect(err).NotTo(HaveOccurred())

			test_helpers.CreateFlop(db, headerOne, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)
			test_helpers.CreateFlop(db, headerTwo, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			expectedBid := test_helpers.FlopBidFromValues(
				strconv.Itoa(fakeBidId), "false", headerOne.Timestamp, headerOne.Timestamp, flopStorageValuesOne)

			var actualBid test_helpers.FlopBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.get_flop($1, $2)`, fakeBidId, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})
	})
})
