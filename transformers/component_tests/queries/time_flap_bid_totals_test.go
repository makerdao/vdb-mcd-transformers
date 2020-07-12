package queries

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Time Flip Bid Totals query", func() {
	var (
		headerRepo      datastore.HeaderRepository
		contractAddress = fakes.FakeAddress.Hex()
		blockOne        int
		timestampOne    int64
		bid, lot        int
		headerOne       core.Header
		fakeBidId       int
		flapKickEvent   event.InsertionModel
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		fakeBidId = rand.Int()
		bid = rand.Int() % 1000000
		lot = rand.Int() % 1000000

		blockOne = rand.Int()
		timestampOne = int64(rand.Int31())
		headerOne = createHeader(blockOne, int(timestampOne), headerRepo)
		flapKickLog := test_data.CreateTestLog(headerOne.Id, db)
		addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		flapKickEvent = test_data.FlapKickModel()
		flapKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
		flapKickEvent.ColumnValues[event.LogFK] = flapKickLog.ID
		flapKickEvent.ColumnValues[event.AddressFK] = addressId
		flapKickEvent.ColumnValues[constants.BidIDColumn] = strconv.Itoa(fakeBidId)
		flapKickEvent.ColumnValues[constants.BidColumn] = strconv.Itoa(0)
		flapKickEvent.ColumnValues[constants.LotColumn] = strconv.Itoa(lot)
		flapKickErr := event.PersistModels([]event.InsertionModel{flapKickEvent}, db)
		Expect(flapKickErr).NotTo(HaveOccurred())
	})

	Context("when called with an hourly 2 hour range with the range start on the first block", func() {
		It("returns the all the bid results under the first hour and 0 for the second hour", func() {
			tendLog := test_data.CreateTestLog(headerOne.Id, db)

			oneHour := timestampOne + 3600
			twoHours := timestampOne + 7200

			dateStart := time.Unix(timestampOne, 0).UTC().Format(time.RFC3339)
			dateMiddle := time.Unix(oneHour, 0).UTC().Format(time.RFC3339)
			dateEnd := time.Unix(twoHours, 0).UTC().Format(time.RFC3339)

			flapTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				BidId:           fakeBidId,
				Lot:             lot,
				BidAmount:       bid,
				TendHeaderId:    headerOne.Id,
				TendLogId:       tendLog.ID,
			})
			Expect(flapTendErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, int(twoHours), headerRepo)

			tickLog := test_data.CreateTestLog(headerOne.Id, db)
			flapTickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				TickHeaderId:    headerTwo.Id,
				TickLogId:       tickLog.ID,
			})
			Expect(flapTickErr).NotTo(HaveOccurred())

			flapDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				DB:              db,
				BidId:           fakeBidId,
				ContractAddress: contractAddress,
				DealHeaderId:    headerTwo.Id,
			})
			Expect(flapDealErr).NotTo(HaveOccurred())

			var actualBidTotals []test_helpers.BucketedBidTotals
			queryErr := db.Select(&actualBidTotals, `SELECT bucket_start, bucket_end, bucket_interval, lot_start, lot_end, bid_amount_start, bid_amount_end FROM api.time_flap_bid_totals($1, $2, '1 hour'::INTERVAL)`, dateStart, dateEnd)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidTotals).To(ConsistOf(
				test_helpers.BucketedBidTotals{BucketStart: dateStart, BucketEnd: dateMiddle, BucketInterval: "01:00:00", LotStart: flapKickEvent.ColumnValues["lot"].(string), LotEnd: strconv.Itoa(lot), BidAmountStart: flapKickEvent.ColumnValues["bid"].(string), BidAmountEnd: strconv.Itoa(bid)},
				test_helpers.BucketedBidTotals{BucketStart: dateMiddle, BucketEnd: dateEnd, BucketInterval: "01:00:00", LotStart: "0", LotEnd: "0", BidAmountStart: "0", BidAmountEnd: "0"}))
		})
	})
})
