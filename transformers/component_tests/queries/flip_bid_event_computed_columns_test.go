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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/flip_kick"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("Flip bid event computed columns", func() {
	var (
		db              *postgres.DB
		blockNumber     = rand.Int()
		header          core.Header
		contractAddress = fakes.FakeAddress.Hex()
		bidId           int
		flipKickRepo    flip_kick.FlipKickRepository
		headerId        int64
		headerRepo      repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		bidId = rand.Int()
		blockNumber = rand.Int()

		headerRepo = repositories.NewHeaderRepository(db)
		header = fakes.GetFakeHeader(int64(blockNumber))
		var insertHeaderErr error
		headerId, insertHeaderErr = headerRepo.CreateOrUpdateHeader(header)
		Expect(insertHeaderErr).NotTo(HaveOccurred())

		flipKickRepo = flip_kick.FlipKickRepository{}
		flipKickRepo.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("flip_bid_event_bid", func() {
		It("returns flip bid for a flip_bid_event", func() {
			flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, bidId)
			flipMetadatas := test_helpers.GetFlipMetadatas(strconv.Itoa(bidId))
			test_helpers.CreateFlip(db, header, flipStorageValues, flipMetadatas, contractAddress)

			ilkId, urnId, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidId:           bidId,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FlipKickModel.Usr,
				FlipKickRepo:     flipKickRepo,
				FlipKickHeaderId: headerId,
			})
			Expect(err).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(bidId), strconv.FormatInt(ilkId, 10),
				strconv.FormatInt(urnId, 10), "false", header.Timestamp, header.Timestamp, flipStorageValues)

			var actualBid test_helpers.FlipBid
			queryErr := db.Get(&actualBid, `
				SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated
				FROM api.flip_bid_event_bid(
					(SELECT (bid_id, lot, bid_amount, act, block_height, tx_idx, contract_address)::api.flip_bid_event FROM api.all_flip_bid_events())
				)`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBid).To(Equal(expectedBid))
		})

		It("gets the correct flipper for the event (using the contract address that matches the event)", func() {
			irrelevantContractAddress := "different flipper"
			irrelevantFlipStorageValues := test_helpers.GetFlipStorageValues(0, test_helpers.AnotherFakeIlk.Hex, bidId)
			irrelevantFlipMetadatas := test_helpers.GetFlipMetadatas(strconv.Itoa(bidId))
			test_helpers.CreateFlip(db, header, irrelevantFlipStorageValues, irrelevantFlipMetadatas, irrelevantContractAddress)

			_, _, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidId:           bidId,
					ContractAddress: irrelevantContractAddress,
				},
				Dealt:            false,
				IlkHex:           test_helpers.AnotherFakeIlk.Hex,
				UrnGuy:           test_data.FakeUrn,
				FlipKickRepo:     flipKickRepo,
				FlipKickHeaderId: headerId,
			})
			Expect(err).NotTo(HaveOccurred())

			flipStorageValues := test_helpers.GetFlipStorageValues(1, test_helpers.FakeIlk.Hex, bidId)
			flipMetadatas := test_helpers.GetFlipMetadatas(strconv.Itoa(bidId))
			test_helpers.CreateFlip(db, header, flipStorageValues, flipMetadatas, contractAddress)

			ilkId, urnId, err := test_helpers.SetUpFlipBidContext(test_helpers.FlipBidContextInput{
				DealCreationInput: test_helpers.DealCreationInput{
					Db:              db,
					BidId:           bidId,
					ContractAddress: contractAddress,
				},
				Dealt:            false,
				IlkHex:           test_helpers.FakeIlk.Hex,
				UrnGuy:           test_data.FakeUrn,
				FlipKickRepo:     flipKickRepo,
				FlipKickHeaderId: headerId,
			})
			Expect(err).NotTo(HaveOccurred())

			expectedBid := test_helpers.FlipBidFromValues(strconv.Itoa(bidId), strconv.FormatInt(ilkId, 10),
				strconv.FormatInt(urnId, 10), "false", header.Timestamp, header.Timestamp, flipStorageValues)

			var actualBid test_helpers.FlipBid
			queryErr := db.Get(&actualBid, `
				SELECT bid_id, ilk_id, urn_id, guy, tic, "end", lot, bid, gal, dealt, tab, created, updated
				FROM api.flip_bid_event_bid(
					(SELECT (bid_id, lot, bid_amount, act, block_height, tx_idx, contract_address)::api.flip_bid_event FROM api.all_flip_bid_events() WHERE contract_address = $1)
				)`, contractAddress)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualBid).To(Equal(expectedBid))
		})
	})

	Describe("flip_bid_event_tx", func() {
		BeforeEach(func() {
			flipKickEvent := test_data.FlipKickModel
			flipKickEvent.ContractAddress = contractAddress
			flipKickEvent.BidId = strconv.Itoa(bidId)
			flipKickErr := flipKickRepo.Create(headerId, []interface{}{flipKickEvent})
			Expect(flipKickErr).NotTo(HaveOccurred())
		})

		It("returns transaction for a flip bid event", func() {
			expectedTx := Tx{
				TransactionHash:  test_helpers.GetValidNullString("txHash"),
				TransactionIndex: sql.NullInt64{Int64: int64(test_data.FlipKickModel.TransactionIndex), Valid: true},
				BlockHeight:      sql.NullInt64{Int64: int64(blockNumber), Valid: true},
				BlockHash:        test_helpers.GetValidNullString(header.Hash),
				TxFrom:           test_helpers.GetValidNullString("fromAddress"),
				TxTo:             test_helpers.GetValidNullString("toAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			queryErr := db.Get(&actualTx, `
				SELECT * FROM api.flip_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, tx_idx, contract_address)::api.flip_bid_event FROM api.all_flip_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(test_data.FlipKickModel.TransactionIndex) + 1,
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(blockNumber), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(header.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx []Tx
			queryErr := db.Select(&actualTx, `
				SELECT * FROM api.flip_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, tx_idx, contract_address)::api.flip_bid_event FROM api.all_flip_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
