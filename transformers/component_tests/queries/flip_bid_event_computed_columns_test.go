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
	"database/sql"
	"math/rand"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flip bid event computed columns", func() {
	var (
		blockOne, timestampOne int
		headerOne              core.Header
		contractAddress        = fakes.FakeAddress.Hex()
		bidId                  int
		headerRepo             datastore.HeaderRepository
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		bidId = rand.Int()

		headerRepo = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
	})

	Describe("flip_bid_event_bid", func() {
		It("returns flip bid for a flip_bid_event", func() {
			flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, bidId)
			flipMetadatas := test_helpers.GetFlipMetadatas(strconv.Itoa(bidId))
			test_helpers.CreateFlip(db, headerOne, flipStorageValues, flipMetadatas, contractAddress)

			ilkId, urnId, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           bidId,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FlipKickModel().ColumnValues["usr"].(string),
				FlipKickHeaderId: headerOne.Id,
			})
			Expect(err).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(bidId), strconv.FormatInt(ilkId, 10),
				strconv.FormatInt(urnId, 10), "false", headerOne.Timestamp, headerOne.Timestamp, flipStorageValues)

			var actualBid test_helpers.FlipBid
			queryErr := db.Get(&actualBid, `
				SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, flip_address, created, updated
				FROM api.flip_bid_event_bid(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flip_bid_event FROM api.all_flip_bid_events())
				)`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBid).To(Equal(expectedBid))
		})

		It("gets the correct flipper for the event (using the contract address that matches the event)", func() {
			irrelevantContractAddress := "different flipper"
			irrelevantFlipStorageValues := test_helpers.GetFlipStorageValues(0, test_helpers.AnotherFakeIlk.Hex, bidId)
			irrelevantFlipMetadatas := test_helpers.GetFlipMetadatas(strconv.Itoa(bidId))
			test_helpers.CreateFlip(db, headerOne, irrelevantFlipStorageValues, irrelevantFlipMetadatas, irrelevantContractAddress)

			_, _, irrelevantFlipContextErr := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           bidId,
					ContractAddress: irrelevantContractAddress,
				},
				Dealt:            false,
				IlkHex:           test_helpers.AnotherFakeIlk.Hex,
				UrnGuy:           test_data.FakeUrn,
				FlipKickHeaderId: headerOne.Id,
			})
			Expect(irrelevantFlipContextErr).NotTo(HaveOccurred())

			flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, bidId)
			flipMetadatas := test_helpers.GetFlipMetadatas(strconv.Itoa(bidId))
			test_helpers.CreateFlip(db, headerOne, flipStorageValues, flipMetadatas, contractAddress)

			ilkId, urnId, flipContextErr := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					DB:              db,
					BidId:           bidId,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FakeUrn,
				FlipKickHeaderId: headerOne.Id,
			})
			Expect(flipContextErr).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(bidId), strconv.FormatInt(ilkId, 10),
				strconv.FormatInt(urnId, 10), "false", headerOne.Timestamp, headerOne.Timestamp, flipStorageValues)

			var actualBid test_helpers.FlipBid
			queryErr := db.Get(&actualBid, `
				SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, flip_address, created, updated
				FROM api.flip_bid_event_bid(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flip_bid_event FROM api.all_flip_bid_events() WHERE contract_address = $1)
				)`, contractAddress)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBid).To(Equal(expectedBid))
		})
	})

	Describe("flip_bid_event_tx", func() {
		var flipKickGethLog types.Log

		BeforeEach(func() {
			flipKickEventLog := test_data.CreateTestLog(headerOne.Id, db)
			flipKickGethLog = flipKickEventLog.Log

			addressId, addressErr := repository.GetOrCreateAddress(db, contractAddress)
			Expect(addressErr).NotTo(HaveOccurred())

			flipKickEvent := test_data.FlipKickModel()
			flipKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
			flipKickEvent.ColumnValues[event.LogFK] = flipKickEventLog.ID
			flipKickEvent.ColumnValues[event.AddressFK] = addressId
			flipKickEvent.ColumnValues[constants.BidIDColumn] = strconv.Itoa(bidId)
			flipKickErr := event.PersistModels([]event.InsertionModel{flipKickEvent}, db)
			Expect(flipKickErr).NotTo(HaveOccurred())
		})

		It("returns transaction for a flip bid event", func() {
			expectedTx := Tx{
				TransactionHash:  test_helpers.GetValidNullString("txHash"),
				TransactionIndex: sql.NullInt64{Int64: int64(flipKickGethLog.TxIndex), Valid: true},
				BlockHeight:      sql.NullInt64{Int64: int64(blockOne), Valid: true},
				BlockHash:        test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:           test_helpers.GetValidNullString("fromAddress"),
				TxTo:             test_helpers.GetValidNullString("toAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO public.transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerOne.Id, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			queryErr := db.Get(&actualTx, `
				SELECT * FROM api.flip_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flip_bid_event FROM api.all_flip_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(flipKickGethLog.TxIndex) + 1,
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(blockOne), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO public.transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerOne.Id, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx []Tx
			queryErr := db.Select(&actualTx, `
				SELECT * FROM api.flip_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flip_bid_event FROM api.all_flip_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})

		It("does not return transaction from different block with same index", func() {
			headerZero := createHeader(blockOne-1, timestampOne-1, headerRepo)
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(flipKickGethLog.TxIndex),
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: headerZero.BlockNumber, Valid: true},
				BlockHash:   test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO public.transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerZero.Id, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx []Tx
			queryErr := db.Select(&actualTx, `
				SELECT * FROM api.flip_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flip_bid_event FROM api.all_flip_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
