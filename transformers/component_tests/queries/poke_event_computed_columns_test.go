package queries

import (
	"database/sql"
	"math/rand"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("all poke events query", func() {
	var (
		blockOne, timestampOne int
		headerOne              core.Header
		fakeGethLog            types.Log
		spotPokeEvent          event.InsertionModel
		headerRepository       repositories.HeaderRepository
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepository)
		fakeEventLog := test_data.CreateTestLog(headerOne.Id, db)
		fakeGethLog = fakeEventLog.Log

		ilkID, ilkErr := shared.GetOrCreateIlk(test_helpers.FakeIlk.Hex, db)
		Expect(ilkErr).NotTo(HaveOccurred())

		spotPokeEvent = test_data.SpotPokeModel()
		spotPokeEvent.ColumnValues[event.HeaderFK] = headerOne.Id
		spotPokeEvent.ColumnValues[event.LogFK] = fakeEventLog.ID
		spotPokeEvent.ColumnValues[constants.IlkColumn] = ilkID
		insertSpotPokeErr := event.PersistModels([]event.InsertionModel{spotPokeEvent}, db)
		Expect(insertSpotPokeErr).NotTo(HaveOccurred())
	})

	Describe("poke_event_ilk", func() {
		It("returns ilk_state for a poke_event", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, headerOne, ilkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)
			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, ilkValues)

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
				BlockHeight:      sql.NullInt64{Int64: int64(blockOne), Valid: true},
				BlockHash:        test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:           test_helpers.GetValidNullString("fromAddress"),
				TxTo:             test_helpers.GetValidNullString("toAddress"),
			}

			_, err := db.Exec(`INSERT INTO public.transactions (header_id, hash, tx_from, tx_index, tx_to)
		        VALUES ($1, $2, $3, $4, $5)`, headerOne.Id, expectedTx.TransactionHash, expectedTx.TxFrom,
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
				BlockHeight: sql.NullInt64{Int64: int64(blockOne), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO public.transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerOne.Id, wrongTx.TransactionHash, wrongTx.TxFrom,
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
			headerZero := createHeader(blockOne-1, timestampOne-1, headerRepository)
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(fakeGethLog.TxIndex),
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

			var actualTx Tx
			err := db.Get(&actualTx, `
				SELECT * FROM api.poke_event_tx(
					(SELECT (ilk_id, val, spot, block_height, log_id)::api.poke_event FROM api.all_poke_events()))`)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
