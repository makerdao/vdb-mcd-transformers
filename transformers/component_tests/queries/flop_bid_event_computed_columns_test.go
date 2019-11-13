package queries

import (
	"database/sql"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/flop_kick"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Flop bid event computed columns", func() {
	var (
		db              *postgres.DB
		blockNumber     = rand.Int()
		timestamp       = int(rand.Int31())
		header          core.Header
		contractAddress = "0x763ztv6x68exwqrgtl325e7hrcvavid4e3fcb4g"
		fakeBidId       = rand.Int()
		flopKickGethLog types.Log
		flopKickRepo    flop_kick.FlopKickRepository
		flopKickEvent   shared.InsertionModel
		headerId        int64
		headerRepo      repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		headerRepo = repositories.NewHeaderRepository(db)
		header = fakes.GetFakeHeaderWithTimestamp(int64(timestamp), int64(blockNumber))
		var insertHeaderErr error
		headerId, insertHeaderErr = headerRepo.CreateOrUpdateHeader(header)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		flopKickHeaderSyncLog := test_data.CreateTestLog(headerId, db)
		flopKickGethLog = flopKickHeaderSyncLog.Log

		flopKickRepo = flop_kick.FlopKickRepository{}
		flopKickRepo.SetDB(db)

		flopKickEvent = test_data.FlopKickModel()
		flopKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
		flopKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
		flopKickEvent.ColumnValues[constants.HeaderFK] = headerId
		flopKickEvent.ColumnValues[constants.LogFK] = flopKickHeaderSyncLog.ID
		insertFlopKickErr := flopKickRepo.Create([]shared.InsertionModel{flopKickEvent})

		Expect(insertFlopKickErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("flop_bid_event_bid", func() {
		It("returns flop bid for a flop_bid_event", func() {
			flopStorageValues := test_helpers.GetFlopStorageValues(1, fakeBidId)
			test_helpers.CreateFlop(db, header, flopStorageValues, test_helpers.GetFlopMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			expectedBid := test_helpers.FlopBidFromValues(strconv.Itoa(fakeBidId), "false", header.Timestamp, header.Timestamp, flopStorageValues)

			var actualBid test_helpers.FlopBid
			err := db.Get(&actualBid, `
				SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated
				FROM api.flop_bid_event_bid(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flop_bid_event FROM api.all_flop_bid_events())
				)`)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualBid).To(Equal(expectedBid))
		})
	})

	Describe("flop_bid_event_tx", func() {
		It("returns transaction for a flop bid event", func() {
			expectedTx := Tx{
				TransactionHash:  test_helpers.GetValidNullString("txHash"),
				TransactionIndex: sql.NullInt64{Int64: int64(flopKickGethLog.TxIndex), Valid: true},
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
				SELECT * FROM api.flop_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flop_bid_event FROM api.all_flop_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(flopKickGethLog.TxIndex + 1),
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
				SELECT * FROM api.flop_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flop_bid_event FROM api.all_flop_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})

		It("does not return transaction from different block with same index", func() {
			lowerBlockNumber := blockNumber - 1
			anotherHeader := fakes.GetFakeHeader(int64(lowerBlockNumber))
			anotherHeaderID, insertHeaderErr := headerRepo.CreateOrUpdateHeader(anotherHeader)
			Expect(insertHeaderErr).NotTo(HaveOccurred())
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(flopKickGethLog.TxIndex),
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(lowerBlockNumber), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(header.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, anotherHeaderID, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx []Tx
			queryErr := db.Select(&actualTx, `
				SELECT * FROM api.flop_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flop_bid_event FROM api.all_flop_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
