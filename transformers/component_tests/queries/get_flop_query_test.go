package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/deal"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flop_kick"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("get flop query", func() {
	var (
		db                         *postgres.DB
		flopKickRepo               flop_kick.FlopKickRepository
		dealRepo                   deal.DealRepository
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
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		flopKickRepo = flop_kick.FlopKickRepository{}
		flopKickRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		blockTwo = blockOne + 1
		timestampOne = int(rand.Int31())
		timestampTwo = timestampOne + 1
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
		headerTwo = createHeader(blockTwo, timestampTwo, headerRepo)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets the specified flop", func() {
		err := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderId:    headerOne.Id,
			},
			Dealt:            true,
			FlopKickRepo:     flopKickRepo,
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
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealRepo:        dealRepo,
				DealHeaderId:    headerTwo.Id,
			},
			Dealt:            true,
			FlopKickRepo:     flopKickRepo,
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
					Db:              db,
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlopKickRepo:     flopKickRepo,
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
					Db:              db,
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
					DealRepo:        dealRepo,
					DealHeaderId:    headerTwo.Id,
				},
				Dealt:            true,
				FlopKickRepo:     flopKickRepo,
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
