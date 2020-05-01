package queries

import (
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip"
	storageHelpers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("urn bid events query", func() {
	var (
		bidOneID,
		bidTwoID,
		bidThreeID int
		usrOne, usrTwo string
		ethFlipAddress = test_data.FlipEthAddress()
		batFlipAddress = common.HexToAddress("0x" + test_data.RandomString(40)).Hex()
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		timestamp := int64(rand.Int31())
		blockNumber := rand.Int()
		headerOneID := storageHelpers.CreateHeader(timestamp, blockNumber, db).Id
		headerTwoID := storageHelpers.CreateHeader(timestamp+1, blockNumber+1, db).Id
		headerThreeID := storageHelpers.CreateHeader(timestamp+2, blockNumber+2, db).Id
		headerFourID := storageHelpers.CreateHeader(timestamp+3, blockNumber+3, db).Id

		bidOneID = rand.Int()
		bidTwoID = bidOneID + 1
		bidThreeID = bidTwoID + 1

		diffID := storageHelpers.CreateFakeDiffRecord(db)
		usrOne = "0x" + test_data.RandomString(40)
		usrTwo = "0x" + test_data.RandomString(40)
		logOneID := test_data.CreateTestLog(headerOneID, db).ID
		logTwoID := test_data.CreateTestLog(headerOneID, db).ID
		logThreeID := test_data.CreateTestLog(headerOneID, db).ID
		logFourID := test_data.CreateTestLog(headerOneID, db).ID

		// insert records used in join
		repo := flip.FlipStorageRepository{ContractAddress: ethFlipAddress}
		repo.SetDB(db)

		ethIlkErr := repo.Create(diffID, headerOneID, flip.IlkMetadata, test_helpers.FakeIlk.Hex)
		Expect(ethIlkErr).NotTo(HaveOccurred())

		ethBidOneKeys := map[types.Key]string{constants.BidId: strconv.Itoa(bidOneID)}
		ethBidUsrMetadataOne := types.GetValueMetadata(storage.BidUsr, ethBidOneKeys, types.Address)
		ethBidUsrErrOne := repo.Create(diffID, headerOneID, ethBidUsrMetadataOne, usrOne)
		Expect(ethBidUsrErrOne).NotTo(HaveOccurred())

		ethBidTwoKeys := map[types.Key]string{constants.BidId: strconv.Itoa(bidTwoID)}
		ethBidUsrMetadataTwo := types.GetValueMetadata(storage.BidUsr, ethBidTwoKeys, types.Address)
		ethBidUsrErrTwo := repo.Create(diffID, headerOneID, ethBidUsrMetadataTwo, usrTwo)
		Expect(ethBidUsrErrTwo).NotTo(HaveOccurred())

		ethBidThreeKeys := map[types.Key]string{constants.BidId: strconv.Itoa(bidThreeID)}
		ethBidUsrMetadataThree := types.GetValueMetadata(storage.BidUsr, ethBidThreeKeys, types.Address)
		ethBidUsrErrThree := repo.Create(diffID, headerOneID, ethBidUsrMetadataThree, usrTwo)
		Expect(ethBidUsrErrThree).NotTo(HaveOccurred())

		repo.ContractAddress = batFlipAddress
		batIlkErr := repo.Create(diffID, headerOneID, flip.IlkMetadata, test_helpers.AnotherFakeIlk.Hex)
		Expect(batIlkErr).NotTo(HaveOccurred())
		batBidUsrErr := repo.Create(diffID, headerOneID, ethBidUsrMetadataOne, usrTwo)
		Expect(batBidUsrErr).NotTo(HaveOccurred())

		// insert 4 kicks for 3 different urns
		kickErrOne := test_helpers.CreateFlipKick(ethFlipAddress, bidOneID, headerOneID, logOneID, usrOne, db)
		Expect(kickErrOne).NotTo(HaveOccurred())
		kickErrTwo := test_helpers.CreateFlipKick(ethFlipAddress, bidTwoID, headerTwoID, logTwoID, usrTwo, db)
		Expect(kickErrTwo).NotTo(HaveOccurred())
		kickErrThree := test_helpers.CreateFlipKick(ethFlipAddress, bidThreeID, headerThreeID, logThreeID, usrTwo, db)
		Expect(kickErrThree).NotTo(HaveOccurred())
		kickErrFour := test_helpers.CreateFlipKick(batFlipAddress, bidOneID, headerFourID, logFourID, usrTwo, db)
		Expect(kickErrFour).NotTo(HaveOccurred())
	})

	It("only gets bid events related to the passed-in urn", func() {
		expectedFlip := test_helpers.BidEvent{ContractAddress: ethFlipAddress, BidId: strconv.Itoa(bidOneID)}

		var actualEvents []test_helpers.BidEvent
		err := db.Select(&actualEvents, `SELECT contract_address, bid_id FROM api.urn_bid_events($1, $2)`,
			usrOne, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())

		Expect(actualEvents).To(ConsistOf(expectedFlip))
	})

	It("only gets bid events from the flip contract associated with passed-in ilk", func() {
		expectedFlip := test_helpers.BidEvent{ContractAddress: batFlipAddress, BidId: strconv.Itoa(bidOneID)}

		var actualEvents []test_helpers.BidEvent
		err := db.Select(&actualEvents, `SELECT contract_address, bid_id FROM api.urn_bid_events($1, $2)`,
			usrTwo, test_helpers.AnotherFakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())

		Expect(actualEvents).To(ConsistOf(expectedFlip))
	})

	It("gets bid events from all of an urn's bids", func() {
		expectedFlips := []test_helpers.BidEvent{
			{ContractAddress: ethFlipAddress, BidId: strconv.Itoa(bidTwoID)},
			{ContractAddress: ethFlipAddress, BidId: strconv.Itoa(bidThreeID)},
		}

		var actualEvents []test_helpers.BidEvent
		err := db.Select(&actualEvents, `SELECT contract_address, bid_id FROM api.urn_bid_events($1, $2)`,
			usrTwo, test_helpers.FakeIlk.Identifier)
		Expect(err).NotTo(HaveOccurred())

		Expect(actualEvents).To(ConsistOf(expectedFlips))
	})
})
