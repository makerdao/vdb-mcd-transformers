package queries

import (
	"database/sql"
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_poke"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
)

var _ = Describe("all poke events query", func() {
	var (
		db               *postgres.DB
		fakeBlock        int
		fakeHeader       core.Header
		fakeGethLog      types.Log
		spotPokeEvent    shared.InsertionModel
		spotPokeRepo     spot_poke.SpotPokeRepository
		headerId         int64
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		fakeBlock = rand.Int()
		fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
		var insertHeaderErr error
		headerId, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())

		fakeHeaderSyncLog := test_data.CreateTestLog(headerId, db)
		fakeGethLog = fakeHeaderSyncLog.Log

		spotPokeRepo = spot_poke.SpotPokeRepository{}
		spotPokeRepo.SetDB(db)
		spotPokeEvent = test_data.SpotPokeModel()
		spotPokeEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		spotPokeEvent.ColumnValues[constants.HeaderFK] = headerId
		spotPokeEvent.ColumnValues[constants.LogFK] = fakeHeaderSyncLog.ID
		insertSpotPokeErr := spotPokeRepo.Create([]shared.InsertionModel{spotPokeEvent})
		Expect(insertSpotPokeErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("poke_event_ilk", func() {
		It("returns ilk_state for a poke_event", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)
			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

			var result test_helpers.IlkState
			err := db.Get(&result, `
				SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated
				FROM api.poke_event_ilk(
					(SELECT (ilk_id, val, spot, block_height, log_id)::api.poke_event FROM api.all_poke_events()))`)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("poke_event_tx", func() {
		It("returns transaction for a poke_event", func() {
			expectedTx := Tx{
				TransactionHash:  test_helpers.GetValidNullString("txHash"),
				TransactionIndex: sql.NullInt64{Int64: int64(fakeGethLog.TxIndex), Valid: true},
				BlockHeight:      sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
				BlockHash:        test_helpers.GetValidNullString(fakeHeader.Hash),
				TxFrom:           test_helpers.GetValidNullString("fromAddress"),
				TxTo:             test_helpers.GetValidNullString("toAddress"),
			}

			_, err := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
		        VALUES ($1, $2, $3, $4, $5)`, headerId, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(err).NotTo(HaveOccurred())

			var actualTx Tx
			err = db.Get(&actualTx, `
				SELECT * FROM api.poke_event_tx(
					(SELECT (ilk_id, val, spot, block_height, log_id)::api.poke_event FROM api.all_poke_events()))`)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(fakeGethLog.TxIndex) + 1,
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(fakeHeader.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			err := db.Get(&actualTx, `
				SELECT * FROM api.poke_event_tx(
					(SELECT (ilk_id, val, spot, block_height, log_id)::api.poke_event FROM api.all_poke_events()))`)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})

		It("does not return transaction from different block with same index", func() {
			lowerBlockNumber := fakeBlock - 1
			anotherHeader := fakes.GetFakeHeader(int64(lowerBlockNumber))
			anotherHeaderID, insertHeaderErr := headerRepository.CreateOrUpdateHeader(anotherHeader)
			Expect(insertHeaderErr).NotTo(HaveOccurred())
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(fakeGethLog.TxIndex),
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(lowerBlockNumber), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(fakeHeader.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, anotherHeaderID, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			err := db.Get(&actualTx, `
				SELECT * FROM api.poke_event_tx(
					(SELECT (ilk_id, val, spot, block_height, log_id)::api.poke_event FROM api.all_poke_events()))`)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
