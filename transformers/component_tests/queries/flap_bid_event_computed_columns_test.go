package queries

import (
	"database/sql"
	"math/rand"
	"strconv"

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

var _ = Describe("flap_bid_event computed columns", func() {
	var (
		blockOne, timestampOne int
		headerOne              core.Header
		flapKickLog            core.EventLog
		headerRepo             datastore.HeaderRepository
		flapKickEvent          event.InsertionModel
		contractAddress        = fakes.RandomString(42)
		fakeBidId              = rand.Int()
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		headerRepo = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
		flapKickLog = test_data.CreateTestLog(headerOne.Id, db)

		addressId, addressErr := shared.GetOrCreateAddress(contractAddress, db)
		Expect(addressErr).NotTo(HaveOccurred())

		flapKickEvent = test_data.FlapKickModel()
		flapKickEvent.ColumnValues[event.HeaderFK] = headerOne.Id
		flapKickEvent.ColumnValues[event.LogFK] = flapKickLog.ID
		flapKickEvent.ColumnValues[event.AddressFK] = addressId
		flapKickEvent.ColumnValues[constants.BidIDColumn] = strconv.Itoa(fakeBidId)
		insertFlapKickErr := event.PersistModels([]event.InsertionModel{flapKickEvent}, db)
		Expect(insertFlapKickErr).NotTo(HaveOccurred())
	})

	Describe("flap_bid_event_bid", func() {
		It("returns flap_bid for a flap_bid_event", func() {
			flapStorageValues := test_helpers.GetFlapStorageValues(1, fakeBidId)
			test_helpers.CreateFlap(db, headerOne, flapStorageValues, test_helpers.GetFlapMetadatas(strconv.Itoa(fakeBidId)), contractAddress)

			expectedBid := test_helpers.FlapBidFromValues(strconv.Itoa(fakeBidId), "false", headerOne.Timestamp, headerOne.Timestamp, flapStorageValues)

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
				SELECT * FROM api.flap_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flap_bid_event
					FROM api.all_flap_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})

		It("does not return transaction from different block with same index", func() {
			headerZero := createHeader(blockOne-1, timestampOne-1, headerRepo)
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(flapKickLog.Log.TxIndex),
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
				SELECT * FROM api.flap_bid_event_tx(
					(SELECT (bid_id, lot, bid_amount, act, block_height, log_id, contract_address)::api.flap_bid_event
					FROM api.all_flap_bid_events()))`)

			Expect(queryErr).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
