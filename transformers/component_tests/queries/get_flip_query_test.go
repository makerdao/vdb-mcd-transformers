package queries

import (
	"github.com/vulcanize/mcd_transformers/transformers/storage"
	"math/rand"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/deal"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flip"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Single flip view", func() {
	var (
		db              *postgres.DB
		flipRepo        flip.FlipStorageRepository
		flipKickRepo    flip_kick.FlipKickRepository
		dealRepo        deal.DealRepository
		headerRepo      repositories.HeaderRepository
		contractAddress = "contract address"
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		flipRepo = flip.FlipStorageRepository{ContractAddress: contractAddress}
		flipRepo.SetDB(db)
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

	It("gets only the specified flip", func() {
		fakeBidId := rand.Int()
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1000

		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
		headerId, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(db, blockOneHeader, flipStorageValuesOne, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		ilkId, urnId, err := test_helpers.SetUpFlipBidBackgroundData(test_helpers.FlipBidSetupData{
			Db:               db,
			BidId:            fakeBidId,
			IlkHex:           test_helpers.FakeIlk.Hex,
			UrnGuy:           test_data.FlipKickModel.Usr,
			ContractAddress:  contractAddress,
			FlipKickRepo:     flipKickRepo,
			FlipKickHeaderId: headerId,
			DealRepo:         dealRepo,
			DealHeaderId:     headerId,
			Dealt:            true,
		})
		Expect(err).NotTo(HaveOccurred())

		expectedBid := test_helpers.FlipBidStateFromValues(
			strconv.Itoa(fakeBidId), strconv.Itoa(ilkId), strconv.Itoa(urnId), "true", blockOneHeader.Timestamp, blockOneHeader.Timestamp, flipStorageValuesOne)

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
		blockTwoHeader.Hash = common.BytesToHash([]byte{5, 4, 3, 2, 1}).String()
		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())
		flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(db, blockTwoHeader, flipStorageValuesTwo, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		var actualBid test_helpers.FlipBidState
		queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.get_flip($1, $2, $3)`, fakeBidId, test_helpers.FakeIlk.Hex, blockOne)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	Describe("dealt", func() {
		It("is false if no deal events", func() {
			fakeBidId := rand.Int()
			blockNumber := rand.Int()
			timestamp := int(rand.Int31())

			header := fakes.GetFakeHeaderWithTimestamp(int64(timestamp), int64(blockNumber))
			headerId, headerOneErr := headerRepo.CreateOrUpdateHeader(header)
			Expect(headerOneErr).NotTo(HaveOccurred())

			flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
			test_helpers.CreateFlip(db, header, flipStorageValues, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			ilkId, urnId, err := test_helpers.SetUpFlipBidBackgroundData(test_helpers.FlipBidSetupData{
				Db:               db,
				BidId:            fakeBidId,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FlipKickModel.Usr,
				ContractAddress:  contractAddress,
				FlipKickRepo:     flipKickRepo,
				FlipKickHeaderId: headerId,
				Dealt:            false,
			})
			Expect(err).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidStateFromValues(
				strconv.Itoa(fakeBidId), strconv.Itoa(ilkId), strconv.Itoa(urnId), "false", header.Timestamp, header.Timestamp, flipStorageValues)

			var actualBid test_helpers.FlipBidState
			queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.get_flip($1, $2, $3)`, fakeBidId, test_helpers.FakeIlk.Hex, blockNumber)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})

		It("is false if deal event in later block", func() {
			fakeBidId := rand.Int()
			blockOne := rand.Int()
			timestampOne := int(rand.Int31())
			blockTwo := blockOne + 1
			timestampTwo := timestampOne + 1000

			blockOneHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
			headerOneId, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
			Expect(headerOneErr).NotTo(HaveOccurred())

			flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
			test_helpers.CreateFlip(db, blockOneHeader, flipStorageValuesOne, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
			blockTwoHeader.Hash = common.BytesToHash([]byte{5, 4, 3, 2, 1}).String()
			headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
			Expect(headerTwoErr).NotTo(HaveOccurred())
			flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidId)
			test_helpers.CreateFlip(db, blockTwoHeader, flipStorageValuesTwo, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			ilkId, urnId, err := test_helpers.SetUpFlipBidBackgroundData(test_helpers.FlipBidSetupData{
				Db:               db,
				BidId:            fakeBidId,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FlipKickModel.Usr,
				ContractAddress:  contractAddress,
				FlipKickRepo:     flipKickRepo,
				FlipKickHeaderId: headerOneId,
				DealRepo:         dealRepo,
				DealHeaderId:     headerTwoId,
				Dealt:            true,
			})
			Expect(err).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidStateFromValues(
				strconv.Itoa(fakeBidId), strconv.Itoa(ilkId), strconv.Itoa(urnId), "false", blockOneHeader.Timestamp, blockOneHeader.Timestamp, flipStorageValuesOne)

			var actualBid test_helpers.FlipBidState
			queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.get_flip($1, $2, $3)`, fakeBidId, test_helpers.FakeIlk.Hex, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})
	})

	It("gets created and updated blocks", func() {
		fakeBidId := rand.Int()
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1000

		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
		headerOneId, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
		blockTwoHeader.Hash = common.BytesToHash([]byte{5, 4, 3, 2, 1}).String()
		_, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())

		flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(db, blockTwoHeader, flipStorageValues, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		ilkId, urnId, err := test_helpers.SetUpFlipBidBackgroundData(test_helpers.FlipBidSetupData{
			Db:               db,
			BidId:            fakeBidId,
			IlkHex:           test_helpers.FakeIlk.Hex,
			UrnGuy:           test_data.FlipKickModel.Usr,
			ContractAddress:  contractAddress,
			FlipKickRepo:     flipKickRepo,
			FlipKickHeaderId: headerOneId,
			Dealt:            false,
		})
		Expect(err).NotTo(HaveOccurred())

		expectedBid := test_helpers.FlipBidStateFromValues(
			strconv.Itoa(fakeBidId), strconv.Itoa(ilkId), strconv.Itoa(urnId), "false", blockTwoHeader.Timestamp, blockOneHeader.Timestamp, flipStorageValues)

		var actualBid test_helpers.FlipBidState
		queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.get_flip($1, $2, $3)`, fakeBidId, test_helpers.FakeIlk.Hex, blockTwo)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})

	It("returns bid state prior to deletion if bid is deleted", func() {
		fakeBidId := rand.Int()
		blockOne := rand.Int()
		timestampOne := int(rand.Int31())
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1000

		blockOneHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampOne), int64(blockOne))
		headerOneId, headerOneErr := headerRepo.CreateOrUpdateHeader(blockOneHeader)
		Expect(headerOneErr).NotTo(HaveOccurred())

		flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(db, blockOneHeader, flipStorageValuesOne, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		blockTwoHeader := fakes.GetFakeHeaderWithTimestamp(int64(timestampTwo), int64(blockTwo))
		blockTwoHeader.Hash = common.BytesToHash([]byte{5, 4, 3, 2, 1}).String()
		headerTwoId, headerTwoErr := headerRepo.CreateOrUpdateHeader(blockTwoHeader)
		Expect(headerTwoErr).NotTo(HaveOccurred())
		flipStorageValuesTwo := getDealtFlipStorageValues(test_helpers.FakeIlk.Hex, fakeBidId)
		test_helpers.CreateFlip(db, blockTwoHeader, flipStorageValuesTwo, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

		ilkId, urnId, err := test_helpers.SetUpFlipBidBackgroundData(test_helpers.FlipBidSetupData{
			Db:               db,
			BidId:            fakeBidId,
			IlkHex:           test_helpers.FakeIlk.Hex,
			UrnGuy:           test_data.FlipKickModel.Usr,
			ContractAddress:  contractAddress,
			FlipKickRepo:     flipKickRepo,
			FlipKickHeaderId: headerOneId,
			DealRepo:         dealRepo,
			DealHeaderId:     headerTwoId,
			Dealt:            true,
		})
		Expect(err).NotTo(HaveOccurred())

		expectedBid := test_helpers.FlipBidStateFromValues(
			strconv.Itoa(fakeBidId), strconv.Itoa(ilkId), strconv.Itoa(urnId), "true", blockTwoHeader.Timestamp, blockOneHeader.Timestamp, flipStorageValuesOne)

		var actualBid test_helpers.FlipBidState
		queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.get_flip($1, $2, $3)`, fakeBidId, test_helpers.FakeIlk.Hex, blockTwo)
		Expect(queryErr).NotTo(HaveOccurred())

		Expect(expectedBid).To(Equal(actualBid))
	})
})

func getDealtFlipStorageValues(ilk string, bidId int) map[string]interface{} {
	emptyAddress := "0x0000000000000000000000000000000000000000"
	zeroStr := strconv.Itoa(0)
	packedValues := map[int]string{0: emptyAddress, 1: zeroStr, 2: zeroStr}

	valuesMap := make(map[string]interface{})
	valuesMap[storage.Ilk] = ilk
	valuesMap[storage.Kicks] = strconv.Itoa(bidId)
	valuesMap[storage.BidBid] = zeroStr
	valuesMap[storage.BidLot] = zeroStr
	valuesMap[storage.BidUsr] = emptyAddress
	valuesMap[storage.BidGal] = emptyAddress
	valuesMap[storage.BidTab] = zeroStr
	valuesMap[storage.Packed] = packedValues

	return valuesMap
}
