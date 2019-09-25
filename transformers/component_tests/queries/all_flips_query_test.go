package queries

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("All flips view", func() {
	var (
		db              *postgres.DB
		flipKickRepo    flip_kick.FlipKickRepository
		dealRepo        deal.DealRepository
		headerRepo      repositories.HeaderRepository
		contractAddress = "contract address"
		hash1           = common.BytesToHash([]byte{5, 4, 3, 2, 1}).String()
		hash2           = common.BytesToHash([]byte{1, 2, 3, 4, 5}).String()
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		flipKickRepo = flip_kick.FlipKickRepository{}
		flipKickRepo.SetDB(db)
		dealRepo = deal.DealRepository{}
		dealRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		rand.Seed(time.Now().UnixNano())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	It("gets the latest state of every bid on the flipper", func() {
		fakeBidId := rand.Int()
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1000

		// insert 2 records for the same bid
		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
		blockOneHeader.Hash = hash1
		headerId, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		ilkId, urnId, setupErr := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			IlkHex:           test_helpers.FakeIlk.Hex,
			UrnGuy:           test_data.FlipKickModel.ColumnValues["usr"].(string),
			FlipKickRepo:     flipKickRepo,
			FlipKickHeaderId: headerId,
		})
		Expect(setupErr).NotTo(HaveOccurred())

		flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(db, blockOneHeader, flipStorageValuesOne,
			test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
		blockTwoHeader.Hash = hash2
		headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())
		flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(db, blockTwoHeader, flipStorageValuesTwo,
			test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		// insert a separate bid with the same ilk
		fakeBidId2 := fakeBidId + 1
		flipStorageValuesThree := test_helpers.GetFlipStorageValues(3, test_helpers.FakeIlk.Hex, fakeBidId2)
		test_helpers.CreateFlip(db, blockTwoHeader, flipStorageValuesThree,
			test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId2)), contractAddress)

		expectedBid1 := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidId), strconv.FormatInt(ilkId, 10),
			strconv.FormatInt(urnId, 10), "false", blockTwoHeader.Timestamp, blockOneHeader.Timestamp, flipStorageValuesTwo)
		var actualBid1 test_helpers.FlipBid
		queryErr1 := db.Get(&actualBid1, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.all_flips($1) WHERE bid_id = $2`,
			test_helpers.FakeIlk.Identifier, fakeBidId)
		Expect(queryErr1).NotTo(HaveOccurred())
		Expect(expectedBid1).To(Equal(actualBid1))

		flipKickLog := test_data.CreateTestLog(headerTwoId, db)
		flipKickErr := test_helpers.CreateFlipKick(contractAddress, fakeBidId2, headerTwoId, flipKickLog.ID, test_data.FlipKickModel.ColumnValues["usr"].(string), flipKickRepo)
		Expect(flipKickErr).NotTo(HaveOccurred())

		expectedBid2 := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidId2), strconv.FormatInt(ilkId, 10),
			strconv.FormatInt(urnId, 10), "false", blockTwoHeader.Timestamp, blockTwoHeader.Timestamp, flipStorageValuesThree)
		var actualBid2 test_helpers.FlipBid
		queryErr2 := db.Get(&actualBid2, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.all_flips($1) WHERE bid_id = $2`,
			test_helpers.FakeIlk.Identifier, fakeBidId2)
		Expect(queryErr2).NotTo(HaveOccurred())
		Expect(expectedBid2).To(Equal(actualBid2))

		var bidCount int
		countQueryErr := db.Get(&bidCount, `SELECT COUNT(*) FROM api.all_flips($1)`, test_helpers.FakeIlk.Identifier)
		Expect(countQueryErr).NotTo(HaveOccurred())
		Expect(bidCount).To(Equal(2))
	})

	It("ignores bids from other contracts", func() {
		fakeBidId := rand.Int()
		blockNumber := rand.Int()
		timestamp := int(rand.Int31())

		header := fakes.GetFakeHeaderWithTimestamp(int64(timestamp), int64(blockNumber))
		header.Hash = hash1
		headerId, headerOneErr := headerRepo.CreateOrUpdateHeader(header)
		Expect(headerOneErr).NotTo(HaveOccurred())

		_, _, setupErr1 := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			IlkHex:           test_helpers.FakeIlk.Hex,
			UrnGuy:           test_data.FlipKickModel.ColumnValues["usr"].(string),
			FlipKickRepo:     flipKickRepo,
			FlipKickHeaderId: headerId,
		})
		Expect(setupErr1).NotTo(HaveOccurred())
		flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(
			db, header, flipStorageValues, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		irrelevantBidId := fakeBidId + 1
		irrelevantAddress := "contract address2"
		irrelevantIlkHex := test_helpers.AnotherFakeIlk.Hex
		irrelevantUrn := test_data.FlipKickModel.ColumnValues["gal"].(string)
		_, _, setupErr2 := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           irrelevantBidId,
				ContractAddress: irrelevantAddress,
			},
			Dealt:            false,
			IlkHex:           irrelevantIlkHex,
			UrnGuy:           irrelevantUrn,
			FlipKickRepo:     flipKickRepo,
			FlipKickHeaderId: headerId,
		})
		Expect(setupErr2).NotTo(HaveOccurred())
		irrelevantFlipValues := test_helpers.GetFlipStorageValues(2, irrelevantIlkHex, irrelevantBidId)
		test_helpers.CreateFlip(db, header, irrelevantFlipValues,
			test_helpers.GetFlipMetadatas(strconv.Itoa(irrelevantBidId)), irrelevantAddress)

		var bidCount int
		countQueryErr := db.Get(&bidCount, `SELECT COUNT(*) FROM api.all_flips($1)`, test_helpers.FakeIlk.Identifier)
		Expect(countQueryErr).NotTo(HaveOccurred())
		Expect(bidCount).To(Equal(1))
	})

	It("gets the all flip bids when there are multiple flip contracts", func() {
		bidIdOne := rand.Int()
		ilkOne := test_helpers.FakeIlk
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())

		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
		blockOneHeader.Hash = hash1
		headerId, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, ilkOne.Hex, bidIdOne)
		test_helpers.CreateFlip(db, blockOneHeader, flipStorageValuesOne,
			test_helpers.GetFlipMetadatas(strconv.Itoa(bidIdOne)), contractAddress)

		ilkId, urnId, setupErr := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           bidIdOne,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			IlkHex:           ilkOne.Hex,
			UrnGuy:           test_data.FlipKickModel.ColumnValues["usr"].(string),
			FlipKickRepo:     flipKickRepo,
			FlipKickHeaderId: headerId,
		})
		Expect(setupErr).NotTo(HaveOccurred())

		bidIdTwo := bidIdOne + 1
		ilkTwo := test_helpers.AnotherFakeIlk
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1000
		anotherFlipAddress := "flip2"

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
		blockTwoHeader.Hash = hash2
		headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, ilkTwo.Hex, bidIdTwo)
		test_helpers.CreateFlip(db, blockTwoHeader, flipStorageValuesTwo,
			test_helpers.GetFlipMetadatas(strconv.Itoa(bidIdTwo)), anotherFlipAddress)

		// insert a new bid associated with a different flip contract address
		ilkIdTwo, urnIdTwo, setupErr := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				Db:              db,
				BidId:           bidIdTwo,
				ContractAddress: anotherFlipAddress,
			},
			Dealt:            false,
			IlkHex:           ilkTwo.Hex,
			UrnGuy:           test_data.FlipKickModel.ColumnValues["gal"].(string),
			FlipKickRepo:     flipKickRepo,
			FlipKickHeaderId: headerTwoId,
		})
		Expect(setupErr).NotTo(HaveOccurred())

		expectedBid1 := test_helpers.FlipBidFromValues(strconv.Itoa(bidIdOne), strconv.FormatInt(ilkId, 10),
			strconv.FormatInt(urnId, 10), "false", blockOneHeader.Timestamp, blockOneHeader.Timestamp, flipStorageValuesOne)

		expectedBid2 := test_helpers.FlipBidFromValues(strconv.Itoa(bidIdTwo), strconv.Itoa(int(ilkIdTwo)),
			strconv.Itoa(int(urnIdTwo)), "false", blockTwoHeader.Timestamp, blockTwoHeader.Timestamp, flipStorageValuesTwo)

		var actualBid1 test_helpers.FlipBid
		queryErr1 := db.Get(&actualBid1, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.all_flips($1) WHERE bid_id = $2`,
			ilkOne.Identifier, bidIdOne)
		Expect(queryErr1).NotTo(HaveOccurred())
		Expect(expectedBid1).To(Equal(actualBid1))

		var actualBid2 test_helpers.FlipBid
		queryErr2 := db.Get(&actualBid2, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.all_flips($1) WHERE bid_id = $2`,
			ilkTwo.Identifier, bidIdTwo)
		Expect(queryErr2).NotTo(HaveOccurred())
		Expect(expectedBid2).To(Equal(actualBid2))
	})

	Describe("result pagination", func() {
		var (
			header                                     core.Header
			headerId, logId                            int64
			fakeBidIdOne, fakeBidIdTwo                 int
			ilkId, urnId                               int64
			flipOneStorageValues, flipTwoStorageValues map[string]interface{}
		)

		BeforeEach(func() {
			fakeBidIdOne = rand.Int()
			fakeBidIdTwo = fakeBidIdOne + 1
			blockNumber := rand.Int()
			headerTimestamp := int(rand.Int31())

			header = fakes.GetFakeHeaderWithTimestamp(int64(headerTimestamp), int64(blockNumber))
			var headerErr error
			headerId, headerErr = headerRepo.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())
			logId = test_data.CreateTestLog(headerId, db).ID

			var setupErr error
			ilkId, urnId, setupErr = test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidId:           fakeBidIdOne,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FlipKickModel.ColumnValues["usr"].(string),
				FlipKickRepo:     flipKickRepo,
				FlipKickHeaderId: headerId,
			})
			Expect(setupErr).NotTo(HaveOccurred())

			flipOneStorageValues = test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidIdOne)
			test_helpers.CreateFlip(db, header, flipOneStorageValues,
				test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

			// insert a separate bid for the same urn
			flipTwoStorageValues = test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidIdTwo)
			test_helpers.CreateFlip(db, header, flipTwoStorageValues,
				test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)
		})

		It("limits results if max_results argument is provided", func() {
			flipKickErr := test_helpers.CreateFlipKick(contractAddress, fakeBidIdTwo, headerId, logId, test_data.FlipKickModel.ColumnValues["usr"].(string), flipKickRepo)
			Expect(flipKickErr).NotTo(HaveOccurred())

			maxResults := 1
			var actualBids []test_helpers.FlipBid
			queryErr := db.Select(&actualBids, `
				SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated
				FROM api.all_flips($1, $2)`,
				test_helpers.FakeIlk.Identifier, maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidIdTwo), strconv.FormatInt(ilkId, 10),
				strconv.FormatInt(urnId, 10), "false", header.Timestamp, header.Timestamp, flipTwoStorageValues)
			Expect(actualBids).To(Equal([]test_helpers.FlipBid{expectedBid}))
		})

		It("offsets results if offset is provided", func() {
			flipKickErr := test_helpers.CreateFlipKick(contractAddress, fakeBidIdOne, headerId, logId, test_data.FlipKickModel.ColumnValues["usr"].(string), flipKickRepo)
			Expect(flipKickErr).NotTo(HaveOccurred())

			maxResults := 1
			resultOffset := 1
			var actualBids []test_helpers.FlipBid
			queryErr := db.Select(&actualBids, `
				SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated
				FROM api.all_flips($1, $2, $3)`,
				test_helpers.FakeIlk.Identifier, maxResults, resultOffset)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidIdOne), strconv.FormatInt(ilkId, 10),
				strconv.FormatInt(urnId, 10), "false", header.Timestamp, header.Timestamp, flipOneStorageValues)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})
	})
})
