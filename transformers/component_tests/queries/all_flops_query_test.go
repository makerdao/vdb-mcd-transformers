package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/deal"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flop_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All flops query", func() {
	var (
		db              *postgres.DB
		flopKickRepo    flop_kick.FlopKickRepository
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
		flopKickRepo = flop_kick.FlopKickRepository{}
		flopKickRepo.SetDB(db)
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
		headerOneId, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(blockTwoTimestamp, int64(blockTwo))
		blockTwoHeader.Hash = "blockTwoHeader"
		headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidIdOne,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlopKickRepo:     flopKickRepo,
			FlopKickHeaderId: headerOneId,
		})
		Expect(contextErr).NotTo(HaveOccurred())

		initialFlopOneStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidIdOne)
		test_helpers.CreateFlop(db, blockOneHeader, initialFlopOneStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		updatedFlopOneStorageValues := test_helpers.GetFlopStorageValues(2, fakeBidIdOne)
		test_helpers.CreateFlop(db, blockTwoHeader, updatedFlopOneStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		flopStorageValuesTwo := test_helpers.GetFlopStorageValues(3, fakeBidIdTwo)
		test_helpers.CreateFlop(db, blockTwoHeader, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)

		contextErr = test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidIdTwo,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlopKickRepo:     flopKickRepo,
			FlopKickHeaderId: headerTwoId,
		})
		Expect(contextErr).NotTo(HaveOccurred())

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

	Describe("result pagination", func() {
		var (
			headerId                                   int64
			header                                     core.Header
			fakeBidIdOne, fakeBidIdTwo                 int
			flopStorageValuesOne, flopStorageValuesTwo map[string]interface{}
		)

		BeforeEach(func() {
			fakeBidIdOne = rand.Int()
			fakeBidIdTwo = fakeBidIdOne + 1

			header = fakes.GetFakeHeaderWithTimestamp(blockOneTimestamp, int64(blockOne))
			var headerErr error
			headerId, headerErr = headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			flopStorageValuesOne = test_helpers.GetFlopStorageValues(1, fakeBidIdOne)
			test_helpers.CreateFlop(db, header, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

			flopStorageValuesTwo = test_helpers.GetFlopStorageValues(2, fakeBidIdTwo)
			test_helpers.CreateFlop(db, header, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)
		})

		It("limits results if max_results argument is provided", func() {
			contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidId:           fakeBidIdTwo,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlopKickRepo:     flopKickRepo,
				FlopKickHeaderId: headerId,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			maxResults := 1
			var actualBids []test_helpers.FlopBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flops($1)`,
				maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdTwo), "false", header.Timestamp,
				header.Timestamp, flopStorageValuesTwo)
			Expect(actualBids).To(Equal([]test_helpers.FlopBid{expectedBid}))
		})

		It("offsets results if offset is provided", func() {
			contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidId:           fakeBidIdOne,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlopKickRepo:     flopKickRepo,
				FlopKickHeaderId: headerId,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			maxResults := 1
			resultOffset := 1
			var actualBids []test_helpers.FlopBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated FROM api.all_flops($1, $2)`,
				maxResults, resultOffset)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdOne), "false", header.Timestamp,
				header.Timestamp, flopStorageValuesOne)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})
	})
})
