package queries

import (
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Single flip view", func() {
	var (
		blockOne, timestampOne int
		contractAddress        = fakes.FakeAddress.Hex()
		fakeBidId              = rand.Int()
		headerOne              core.Header
		headerRepo             datastore.HeaderRepository
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)

		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	Context("get_flip_with_address", func() {
		It("gets only the specified flip", func() {
			flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
			test_helpers.CreateFlip(db, headerOne, flipStorageValuesOne, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			ilkId, urnId, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
					DealHeaderId:    headerOne.Id,
				},
				Dealt:            true,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FlipKickModel().ColumnValues["usr"].(string),
				FlipKickHeaderId: headerOne.Id,
			})
			Expect(err).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidId), strconv.FormatInt(ilkId, 10),
				strconv.FormatInt(urnId, 10), "true", headerOne.Timestamp, headerOne.Timestamp, flipStorageValuesOne)
			expectedBid.BlockHeight = strconv.FormatInt(headerOne.BlockNumber, 10)

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidId)
			test_helpers.CreateFlip(db, headerTwo, flipStorageValuesTwo, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			var actualBid test_helpers.FlipBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, flip_address, created, updated, block_height FROM api.get_flip_with_address($1, $2, $3, $4)`,
				fakeBidId, contractAddress, test_helpers.FakeIlk.Identifier, blockOne)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid).To(Equal(actualBid))
		})

		Describe("dealt", func() {
			It("is false if no deal events", func() {
				flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
				test_helpers.CreateFlip(db, headerOne, flipStorageValues, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

				ilkId, urnId, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
					DealCreationInput: test_helpers.DealCreationInput{
						DB:              db,
						BidId:           fakeBidId,
						ContractAddress: contractAddress,
					},
					Dealt:            false,
					IlkHex:           test_helpers.FakeIlk.Hex,
					UrnGuy:           test_data.FlipKickModel().ColumnValues["usr"].(string),
					FlipKickHeaderId: headerOne.Id,
				})
				Expect(err).NotTo(HaveOccurred())

				expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidId), strconv.FormatInt(ilkId, 10),
					strconv.FormatInt(urnId, 10), "false", headerOne.Timestamp, headerOne.Timestamp, flipStorageValues)

				var actualBid test_helpers.FlipBid
				queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, flip_address, created, updated FROM api.get_flip_with_address($1, $2, $3, $4)`,
					fakeBidId, contractAddress, test_helpers.FakeIlk.Identifier, blockOne)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(expectedBid).To(Equal(actualBid))
			})

			It("is false if deal event in later block", func() {
				flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
				test_helpers.CreateFlip(db, headerOne, flipStorageValuesOne, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

				headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)

				flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidId)
				test_helpers.CreateFlip(db, headerTwo, flipStorageValuesTwo, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

				ilkId, urnId, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
					DealCreationInput: test_helpers.DealCreationInput{
						DB:              db,
						BidId:           fakeBidId,
						ContractAddress: contractAddress,
						DealHeaderId:    headerTwo.Id,
					},
					Dealt:            true,
					IlkHex:           test_helpers.FakeIlk.Hex,
					UrnGuy:           test_data.FlipKickModel().ColumnValues["usr"].(string),
					FlipKickHeaderId: headerOne.Id,
				})
				Expect(err).NotTo(HaveOccurred())

				expectedBid := test_helpers.FlipBidFromValues(
					strconv.Itoa(fakeBidId), strconv.FormatInt(ilkId, 10), strconv.FormatInt(urnId, 10), "false",
					headerOne.Timestamp, headerOne.Timestamp, flipStorageValuesOne)

				var actualBid test_helpers.FlipBid
				queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, flip_address, created, updated FROM api.get_flip_with_address($1, $2, $3, $4)`,
					fakeBidId, contractAddress, test_helpers.FakeIlk.Identifier, blockOne)
				Expect(queryErr).NotTo(HaveOccurred())

				Expect(expectedBid).To(Equal(actualBid))
			})
		})

		It("gets created and updated blocks", func() {
			flipStorageValuesOne := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, fakeBidId)
			test_helpers.CreateFlip(db, headerOne, flipStorageValuesOne, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			ilkId, urnId, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           fakeBidId,
					ContractAddress: contractAddress,
					DealHeaderId:    headerOne.Id,
				},
				Dealt:            true,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FlipKickModel().ColumnValues["usr"].(string),
				FlipKickHeaderId: headerOne.Id,
			})
			Expect(err).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, timestampOne+1, headerRepo)
			flipStorageValuesTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, fakeBidId)
			test_helpers.CreateFlip(db, headerTwo, flipStorageValuesTwo, test_helpers.GetFlipMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(fakeBidId), strconv.FormatInt(ilkId, 10),
				strconv.FormatInt(urnId, 10), "true", headerTwo.Timestamp, headerOne.Timestamp, flipStorageValuesOne)

			var actualBid test_helpers.FlipBid
			queryErr := db.Get(&actualBid, `SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated FROM api.get_flip_with_address($1, $2, $3, $4)`,
				fakeBidId, contractAddress, test_helpers.FakeIlk.Identifier, headerTwo.BlockNumber)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(expectedBid.Created).To(Equal(actualBid.Created))
			Expect(expectedBid.Updated).To(Equal(actualBid.Updated))
		})
	})
})
