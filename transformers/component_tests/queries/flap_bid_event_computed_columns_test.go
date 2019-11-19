package queries

import (
	"database/sql"
	"math/rand"
	"strconv"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/flap_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("flap_bid_event computed columns", func() {
	var (
		db              *postgres.DB
		blockNumber     = rand.Int()
		header          core.Header
		flapKickLog     core.HeaderSyncLog
		headerId        int64
		headerRepo      repositories.HeaderRepository
		flapKickRepo    flap_kick.FlapKickRepository
		flapKickEvent   shared.InsertionModel
		contractAddress = "FlapAddress"
		fakeBidId       = rand.Int()
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		headerRepo = repositories.NewHeaderRepository(db)
		header = fakes.GetFakeHeader(int64(blockNumber))
		var insertHeaderErr error
		headerId, insertHeaderErr = headerRepo.CreateOrUpdateHeader(header)
		Expect(insertHeaderErr).NotTo(HaveOccurred())
		flapKickLog = test_data.CreateTestLog(headerId, db)

		flapKickRepo = flap_kick.FlapKickRepository{}
		flapKickRepo.SetDB(db)

		flapKickEvent = test_data.FlapKickModel()
		flapKickEvent.ColumnValues["bid_id"] = strconv.Itoa(fakeBidId)
		flapKickEvent.ForeignKeyValues[constants.AddressFK] = contractAddress
		flapKickEvent.ColumnValues[constants.HeaderFK] = headerId
		flapKickEvent.ColumnValues[constants.LogFK] = flapKickLog.ID
		insertFlapKickErr := flapKickRepo.Create([]shared.InsertionModel{flapKickEvent})
		Expect(insertFlapKickErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("flap_bid_event_bid", func() {
		It("returns flap_bid for a flap_bid_event", func() {
			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, header, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidId), "false", header.Timestamp, header.Timestamp, flapStorageValues)

			var actualBid test_helpers.FlapBid
			err := db.Get(&actualBid, `
				SELECT bid_id, guy, tic, "end", lot, bid, dealt, created, updated
				FROM api.flap_bid_event_bid(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flap_bid_event
					FROM api.all_flap_bid_events())
				)`)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualBid).To(Equal(expectedBid))
		})
	})

	Describe("flap_bid_event_tx", func() {
		It("returns transaction for a flap_bid_event", func() {
			expectedTx := Tx{
				TransactionHash:  test_helpers.GetValidNullString("txHash"),
				TransactionIndex: sql.NullInt64{Int64: int64(flapKickLog.Log.TxIndex), Valid: true},
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
				SELECT * FROM api.flap_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flap_bid_event
					FROM api.all_flap_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(flapKickLog.Log.TxIndex) + 1,
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
				SELECT * FROM api.flap_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flap_bid_event
					FROM api.all_flap_bid_events()))`)

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
					Int64: int64(flapKickLog.Log.TxIndex),
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
				SELECT * FROM api.flap_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flap_bid_event
					FROM api.all_flap_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
