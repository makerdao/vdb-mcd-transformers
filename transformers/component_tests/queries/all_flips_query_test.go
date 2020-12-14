package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("All flips view", func() {
	var (
		headerRepo             datastore.HeaderRepository
		contractAddress        = fakes.FakeAddress.Hex()
		blockOne, timestampOne int
		headerOne              core.Header
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	It("gets the latest state of every bid on the flipper", func() {
		fakeBidId := rand.Int()

		// insert 2 records for the same bid
		ilkId, urnId, setupErr := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			IlkHex:           test_helpers.FakeIlk.Hex,
			UrnGuy:           test_data.FlipKickModel().ColumnValues[constants.UsrColumn].(string),
			FlipKickHeaderId: headerOne.Id,
		})
		Expect(setupErr).NotTo(HaveOccurred())

		flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(db, headerOne, flipStorageValuesOne, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
		flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(db, headerTwo, flipStorageValuesTwo, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		// insert a separate bid with the same ilk
		fakeBidId2 := fakeBidId + 1
		flipStorageValuesThree := test_helpers.GetFlipStorageValues(3, test_helpers.FakeIlk.Hex, fakeBidId2)
		test_helpers.CreateFlip(db, headerTwo, flipStorageValuesThree, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId2)), contractAddress)

		expectedBid1 := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidId), strconv.FormatInt(ilkId, 10),
			strconv.FormatInt(urnId, 10), "false", headerTwo.Timestamp, headerOne.Timestamp, flipStorageValuesTwo)
		var actualBid1 test_helpers.FlipBid
		queryErr1 := db.Get(&actualBid1, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, flip_address, created, updated FROM api.all_flips($1) WHERE bid_id = $2`,
			test_helpers.FakeIlk.Identifier, fakeBidId)
		Expect(queryErr1).NotTo(HaveOccurred())
		Expect(expectedBid1).To(Equal(actualBid1))

		flipKickLog := test_data.CreateTestLog(headerTwo.Id, db)
		flipKickErr := test_helpers.CreateFlipKick(contractAddress, fakeBidId2, headerTwo.Id, flipKickLog.ID, test_data.FlipKickModel().ColumnValues[constants.UsrColumn].(string), db)
		Expect(flipKickErr).NotTo(HaveOccurred())

		expectedBid2 := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidId2), strconv.FormatInt(ilkId, 10),
			strconv.FormatInt(urnId, 10), "false", headerTwo.Timestamp, headerTwo.Timestamp, flipStorageValuesThree)
		var actualBid2 test_helpers.FlipBid
		queryErr2 := db.Get(&actualBid2, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, flip_address, created, updated FROM api.all_flips($1) WHERE bid_id = $2`,
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

		_, _, setupErr1 := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			IlkHex:           test_helpers.FakeIlk.Hex,
			UrnGuy:           test_data.FlipKickModel().ColumnValues[constants.UsrColumn].(string),
			FlipKickHeaderId: headerOne.Id,
		})
		Expect(setupErr1).NotTo(HaveOccurred())
		flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(db, headerOne, flipStorageValues, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		irrelevantBidId := fakeBidId + 1
		irrelevantAddress := "contract address2"
		irrelevantIlkHex := test_helpers.AnotherFakeIlk.Hex
		irrelevantUrn := test_data.FlipKickModel().ColumnValues[constants.GalColumn].(string)
		_, _, setupErr2 := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           irrelevantBidId,
				ContractAddress: irrelevantAddress,
			},
			Dealt:            false,
			IlkHex:           irrelevantIlkHex,
			UrnGuy:           irrelevantUrn,
			FlipKickHeaderId: headerOne.Id,
		})
		Expect(setupErr2).NotTo(HaveOccurred())
		irrelevantFlipValues := test_helpers.GetFlipStorageValues(2, irrelevantIlkHex, irrelevantBidId)
		test_helpers.CreateFlip(db, headerOne, irrelevantFlipValues, test_helpers.GetFlipMetadatas(strconv.Itoa(irrelevantBidId)), irrelevantAddress)

		var bidCount int
		countQueryErr := db.Get(&bidCount, `SELECT COUNT(*) FROM api.all_flips($1)`, test_helpers.FakeIlk.Identifier)
		Expect(countQueryErr).NotTo(HaveOccurred())
		Expect(bidCount).To(Equal(1))
	})

	It("gets the all flip bids when there are multiple flip contracts", func() {
		bidIdOne := rand.Int()
		ilkOne := test_helpers.FakeIlk

		flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, ilkOne.Hex, bidIdOne)
		test_helpers.CreateFlip(db, headerOne, flipStorageValuesOne, test_helpers.GetFlipMetadatas(strconv.Itoa(bidIdOne)), contractAddress)

		ilkId, urnId, setupErr := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           bidIdOne,
				ContractAddress: contractAddress,
			},
			Dealt:            false,
			IlkHex:           ilkOne.Hex,
			UrnGuy:           test_data.FlipKickModel().ColumnValues[constants.UsrColumn].(string),
			FlipKickHeaderId: headerOne.Id,
		})
		Expect(setupErr).NotTo(HaveOccurred())

		bidIdTwo := bidIdOne + 1
		ilkTwo := test_helpers.AnotherFakeIlk
		anotherFlipAddress := fakes.AnotherFakeAddress.Hex()

		headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

		flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, ilkTwo.Hex, bidIdTwo)
		test_helpers.CreateFlip(db, headerTwo, flipStorageValuesTwo, test_helpers.GetFlipMetadatas(strconv.Itoa(bidIdTwo)), anotherFlipAddress)

		// insert a new bid associated with a different flip contract address
		ilkIdTwo, urnIdTwo, setupErr := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
			DealCreationInput: test_helpers.DealCreationInput{
				DB:              db,
				BidId:           bidIdTwo,
				ContractAddress: anotherFlipAddress,
			},
			Dealt:            false,
			IlkHex:           ilkTwo.Hex,
			UrnGuy:           test_data.FlipKickModel().ColumnValues[constants.GalColumn].(string),
			FlipKickHeaderId: headerTwo.Id,
		})
		Expect(setupErr).NotTo(HaveOccurred())

		expectedBid1 := test_helpers.FlipBidFromValues(strconv.Itoa(bidIdOne), strconv.FormatInt(ilkId, 10),
			strconv.FormatInt(urnId, 10), "false", headerOne.Timestamp, headerOne.Timestamp, flipStorageValuesOne)

		expectedBid2 := test_helpers.FlipBidFromValues(strconv.Itoa(bidIdTwo), strconv.Itoa(int(ilkIdTwo)),
			strconv.Itoa(int(urnIdTwo)), "false", headerTwo.Timestamp, headerTwo.Timestamp, flipStorageValuesTwo)
		expectedBid2.FlipAddress = anotherFlipAddress

		var actualBid1 test_helpers.FlipBid
		queryErr1 := db.Get(&actualBid1, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, flip_address, created, updated FROM api.all_flips($1) WHERE bid_id = $2`,
			ilkOne.Identifier, bidIdOne)
		Expect(queryErr1).NotTo(HaveOccurred())
		Expect(expectedBid1).To(Equal(actualBid1))

		var actualBid2 test_helpers.FlipBid
		queryErr2 := db.Get(&actualBid2, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, flip_address, created, updated FROM api.all_flips($1) WHERE bid_id = $2`,
			ilkTwo.Identifier, bidIdTwo)
		Expect(queryErr2).NotTo(HaveOccurred())
		Expect(expectedBid2).To(Equal(actualBid2))
	})

	Describe("result pagination", func() {
		var (
			logId                                      int64
			fakeBidIdOne, fakeBidIdTwo                 int
			ilkId, urnId                               int64
			flipOneStorageValues, flipTwoStorageValues map[string]interface{}
		)

		BeforeEach(func() {
			fakeBidIdOne = rand.Int()
			fakeBidIdTwo = fakeBidIdOne + 1

			logId = test_data.CreateTestLog(headerOne.Id, db).ID

			var setupErr error
			ilkId, urnId, setupErr = test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidIdOne,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FlipKickModel().ColumnValues[constants.UsrColumn].(string),
				FlipKickHeaderId: headerOne.Id,
			})
			Expect(setupErr).NotTo(HaveOccurred())

			flipOneStorageValues = test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidIdOne)
			test_helpers.CreateFlip(db, headerOne, flipOneStorageValues, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidIdOne)), contractAddress)

			// insert a separate bid for the same urn
			flipTwoStorageValues = test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidIdTwo)
			test_helpers.CreateFlip(db, headerOne, flipTwoStorageValues, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidIdTwo)), contractAddress)
		})

		It("limits results if max_results argument is provided", func() {
			flipKickErr := test_helpers.CreateFlipKick(contractAddress, fakeBidIdTwo, headerOne.Id, logId, test_data.FlipKickModel().ColumnValues[constants.UsrColumn].(string), db)
			Expect(flipKickErr).NotTo(HaveOccurred())

			maxResults := 1
			var actualBids []test_helpers.FlipBid
			queryErr := db.Select(&actualBids, `
				SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, flip_address, created, updated
				FROM api.all_flips($1, $2)`,
				test_helpers.FakeIlk.Identifier, maxResults)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidIdTwo), strconv.FormatInt(ilkId, 10),
				strconv.FormatInt(urnId, 10), "false", headerOne.Timestamp, headerOne.Timestamp, flipTwoStorageValues)
			Expect(actualBids).To(Equal([]test_helpers.FlipBid{expectedBid}))
		})

		It("offsets results if offset is provided", func() {
			flipKickErr := test_helpers.CreateFlipKick(contractAddress, fakeBidIdOne, headerOne.Id, logId, test_data.FlipKickModel().ColumnValues[constants.UsrColumn].(string), db)
			Expect(flipKickErr).NotTo(HaveOccurred())

			maxResults := 1
			resultOffset := 1
			var actualBids []test_helpers.FlipBid
			queryErr := db.Select(&actualBids, `
				SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, flip_address, created, updated
				FROM api.all_flips($1, $2, $3)`,
				test_helpers.FakeIlk.Identifier, maxResults, resultOffset)
			Expect(queryErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidIdOne), strconv.FormatInt(ilkId, 10),
				strconv.FormatInt(urnId, 10), "false", headerOne.Timestamp, headerOne.Timestamp, flipOneStorageValues)
			Expect(actualBids).To(ConsistOf(expectedBid))
		})
	})
})
