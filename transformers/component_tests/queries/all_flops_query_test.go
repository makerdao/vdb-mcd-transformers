package queries

import (
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flop"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All flops query", func() {
	var (
		flopRepo               flop.StorageRepository
		headerRepo             datastore.HeaderRepository
		contractAddress        = fakes.FakeAddress.Hex()
		blockOne, timestampOne int
		headerOne              core.Header
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		flopRepo = flop.StorageRepository{}
		flopRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	It("gets the most recent flop for every bid id", func() {
		fakeBidIdOne := rand.Int()
		fakeBidIdTwo := fakeBidIdOne + 1

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

		contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidIdOne,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlopKickHeaderId: headerOne.Id,
		})
		Expect(contextErr).NotTo(HaveOccurred())

		initialFlopOneStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidIdOne)
		test_helpers.CreateFlop(db, headerOne, initialFlopOneStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		updatedFlopOneStorageValues := test_helpers.GetFlopStorageValues(2, fakeBidIdOne)
		test_helpers.CreateFlop(db, headerTwo, updatedFlopOneStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

		flopStorageValuesTwo := test_helpers.GetFlopStorageValues(3, fakeBidIdTwo)
		test_helpers.CreateFlop(db, headerTwo, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)

		contextErr = test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidIdTwo,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlopKickHeaderId: headerTwo.Id,
		})
		Expect(contextErr).NotTo(HaveOccurred())

		var actualBids []test_helpers.FlopBid
		queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated, flop_address FROM api.all_flops()`)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedBidOne := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdOne), "false", headerTwo.Timestamp, headerOne.Timestamp, contractAddress, updatedFlopOneStorageValues)
		expectedBidTwo := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdTwo), "false", headerTwo.Timestamp, headerTwo.Timestamp, contractAddress, flopStorageValuesTwo)

		Expect(len(actualBids)).To(Equal(2))
		Expect(actualBids).To(ConsistOf([]test_helpers.FlopBid{
			expectedBidOne,
			expectedBidTwo,
		}))
	})

	It("returns flops from every flop contract", func() {
		fakeBidID := rand.Int()

		contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidID,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			FlopKickHeaderId: headerOne.Id,
		})
		Expect(contextErr).NotTo(HaveOccurred())
		flopStorageValuesOne := test_helpers.GetFlopStorageValues(1, fakeBidID)
		test_helpers.CreateFlop(db, headerOne, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddress)

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
		contractAddressTwo := common.HexToAddress(fakes.RandomString(40)).Hex()
		contextErr = test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidID,
				ContractAddress: contractAddressTwo,
			},
			Dealt:            false,
			FlopKickHeaderId: headerTwo.Id,
		})
		Expect(contextErr).NotTo(HaveOccurred())
		flopStorageValuesTwo := test_helpers.GetFlopStorageValues(3, fakeBidID)
		test_helpers.CreateFlop(db, headerTwo, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidID)), contractAddressTwo)

		var actualBids []test_helpers.FlopBid
		queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated, flop_address FROM api.all_flops()`)
		Expect(queryErr).NotTo(HaveOccurred())

		expectedBidOne := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidID), "false", headerOne.Timestamp, headerOne.Timestamp, contractAddress, flopStorageValuesOne)
		expectedBidTwo := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidID), "false", headerTwo.Timestamp, headerTwo.Timestamp, contractAddressTwo, flopStorageValuesTwo)

		Expect(len(actualBids)).To(Equal(2))
		Expect(actualBids).To(ConsistOf([]test_helpers.FlopBid{
			expectedBidOne,
			expectedBidTwo,
		}))
	})

	Describe("result pagination", func() {
		var (
			fakeBidIdOne, fakeBidIdTwo                 int
			flopStorageValuesOne, flopStorageValuesTwo map[string]interface{}
		)

		BeforeEach(func() {
			fakeBidIdOne = rand.Int()
			fakeBidIdTwo = fakeBidIdOne + 1

			flopStorageValuesOne = test_helpers.GetFlopStorageValues(1, fakeBidIdOne)
			test_helpers.CreateFlop(db, headerOne, flopStorageValuesOne, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

			flopStorageValuesTwo = test_helpers.GetFlopStorageValues(2, fakeBidIdTwo)
			test_helpers.CreateFlop(db, headerOne, flopStorageValuesTwo, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)
		})

		It("limits results if max_results argument is provided", func() {
			contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidIdTwo,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlopKickHeaderId: headerOne.Id,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			maxResults := 1
			var actualBids []test_helpers.FlopBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated, flop_address FROM api.all_flops($1)`,
				maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdTwo), "false", headerOne.Timestamp,
				headerOne.Timestamp, contractAddress, flopStorageValuesTwo)
			Expect(actualBids).To(Equal([]test_helpers.FlopBid{expectedBid}))
		})

		It("offsets results if offset is provided", func() {
			contextErr := test_helpers.SetUpFlopBidContext(test_helpers.FlopBidCreationInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidIdOne,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				FlopKickHeaderId: headerOne.Id,
			})
			Expect(contextErr).NotTo(HaveOccurred())

			maxResults := 1
			resultOffset := 1
			var actualBids []test_helpers.FlopBid
			queryErr := db.Select(&actualBids, `SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated, flop_address FROM api.all_flops($1, $2)`,
				maxResults, resultOffset)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidIdOne), "false", headerOne.Timestamp,
				headerOne.Timestamp, contractAddress, flopStorageValuesOne)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})
	})
})
