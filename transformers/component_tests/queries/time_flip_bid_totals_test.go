// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
		addressId       int64
		bidId           int
		blockOne        int
		timestampOne    int64
		bid, lot        int
		headerOne       core.Header
		flipKickEvent   event.InsertionModel
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		bidId = rand.Int()
		bid = rand.Int() % 1000000
		lot = rand.Int() % 1000000

		blockOne = rand.Int()
		timestampOne = int64(rand.Int31())
		headerOne = createHeader(blockOne, int(timestampOne), headerRepo)

		flipKickLog := test_data.CreateTestLog(headerOne.Id, db)

		var addressErr error
		addressId, addressErr = shared.GetOrCreateAddress(contractAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		flipKickEvent = test_data.FlipKickModel()
		flipKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
		flipKickEvent.ColumnValues[event.LogFK] = flipKickLog.ID
		flipKickEvent.ColumnValues[event.AddressFK] = addressId
		flipKickEvent.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidId)
		flipKickEvent.ColumnValues[constants.BidColumn] = strconv.Itoa(0)
		flipKickEvent.ColumnValues[constants.LotColumn] = strconv.Itoa(lot)
		flipKickErr := event.PersistModels([]event.InsertionModel{flipKickEvent}, db)
		Expect(flipKickErr).NotTo(HaveOccurred())
	})

	Context("when called with an hourly 2 hour range with the range start on the first block", func() {
		It("returns the all the bid results under the first hour and 0 for the second hour", func() {
			tendLot := lot
			tendBidAmount := bid + (rand.Int() % 1000000)
			dentLot := lot / (rand.Int()%8 + 2)
			dentBidAmount := tendBidAmount

			oneHour := timestampOne + 3600
			twoHours := timestampOne + 7200
			sevenHours := timestampOne + 25200

			dateStart := time.Unix(timestampOne, 0).UTC().Format(time.RFC3339)
			dateMiddle := time.Unix(oneHour, 0).UTC().Format(time.RFC3339)
			dateEnd := time.Unix(twoHours, 0).UTC().Format(time.RFC3339)

			flipTendLog := test_data.CreateTestLog(headerOne.Id, db)
			flipTendErr := test_helpers.CreateTend(test_helpers.TendCreationInput{
				DB:              db,
				ContractAddress: contractAddress,
				BidId:           bidId,
				Lot:             tendLot,
				BidAmount:       tendBidAmount,
				TendHeaderId:    headerOne.Id,
				TendLogId:       flipTendLog.ID,
			})
			Expect(flipTendErr).NotTo(HaveOccurred())

			headerTwo := createHeader(blockOne+1, int(twoHours), headerRepo)

			tickLog := test_data.CreateTestLog(headerTwo.Id, db)
			tickErr := test_helpers.CreateTick(test_helpers.TickCreationInput{
				DB:              db,
				BidId:           bidId,
				ContractAddress: contractAddress,
				TickHeaderId:    headerTwo.Id,
				TickLogId:       tickLog.ID,
			})
			Expect(tickErr).NotTo(HaveOccurred())

			flipStorageValuesBlockTwo := test_helpers.GetFlipStorageValues(2, test_helpers.FakeIlk.Hex, bidId)
			test_helpers.CreateFlip(db, headerTwo, flipStorageValuesBlockTwo, test_helpers.GetFlipMetadatas(strconv.Itoa(bidId)), contractAddress)

			headerThree := createHeader(blockOne+2, int(sevenHours), headerRepo)

			flipDentLog := test_data.CreateTestLog(headerThree.Id, db)
			flipDentErr := test_helpers.CreateDent(test_helpers.DentCreationInput{
				DB:              db,
				BidId:           bidId,
				ContractAddress: contractAddress,
				Lot:             dentLot,
				BidAmount:       dentBidAmount,
				DentHeaderId:    headerThree.Id,
				DentLogId:       flipDentLog.ID,
			})
			Expect(flipDentErr).NotTo(HaveOccurred())

			flipDealErr := test_helpers.CreateDeal(test_helpers.DealCreationInput{
				DB:              db,
				BidId:           bidId,
				ContractAddress: contractAddress,
				DealHeaderId:    headerThree.Id,
			})
			Expect(flipDealErr).NotTo(HaveOccurred())

			var actualBidTotals []test_helpers.BucketedBidTotals
			queryErr := db.Select(&actualBidTotals, `SELECT bucket_start, bucket_end, bucket_interval, lot_start, lot_end, bid_amount_start, bid_amount_end FROM api.time_flip_bid_totals($1, $2, $3, '1 hour'::INTERVAL)`, test_helpers.FakeIlk.Identifier, dateStart, dateEnd)
			Expect(queryErr).NotTo(HaveOccurred())

			Expect(actualBidTotals).To(ConsistOf(
				test_helpers.BucketedBidTotals{BucketStart: dateStart, BucketEnd: dateMiddle, BucketInterval: "01:00:00", LotStart: flipKickEvent.ColumnValues["lot"].(string), LotEnd: strconv.Itoa(dentLot), BidAmountStart: flipKickEvent.ColumnValues["bid"].(string), BidAmountEnd: strconv.Itoa(tendBidAmount)},
				test_helpers.BucketedBidTotals{BucketStart: dateMiddle, BucketEnd: dateEnd, BucketInterval: "01:00:00", LotStart: "0", LotEnd: "0", BidAmountStart: "0", BidAmountEnd: "0"}))
		})
	})
})
