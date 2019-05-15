package queries

import (
	"database/sql"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("File event computed columns", func() {
	var (
		db         *postgres.DB
		fakeBlock  int
		fakeHeader core.Header
		fileEvent  ilk.VatFileIlkModel
		fileRepo   ilk.VatFileIlkRepository
		headerId   int64
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

		fileRepo = ilk.VatFileIlkRepository{}
		fileRepo.SetDB(db)
		fileEvent = test_data.VatFileIlkDustModel
		fileEvent.Ilk = test_helpers.FakeIlk.Hex
		insertFileErr := fileRepo.Create(headerId, []interface{}{fileEvent})
		Expect(insertFileErr).NotTo(HaveOccurred())
	})

	Describe("file_event_ilk", func() {
		It("returns ilk_state for a file_event", func() {
			vatRepository.SetDB(db)
			catRepository.SetDB(db)
			jugRepository.SetDB(db)

			ilkValues := test_helpers.GetIlkValues(0)
			createIlkAtBlock(fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

			var result test_helpers.IlkState
			err := db.Get(&result,
				`SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated
                    FROM api.file_event_ilk(
                        (SELECT (id, ilk_name, what, data, block_height, tx_idx)::api.file_event FROM api.ilk_files($1))
                    )`, test_helpers.FakeIlk.Name)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("file_event_tx", func() {
		It("returns transaction for a file event", func() {
			expectedTx := Tx{
				TransactionHash: sql.NullString{String: "txHash", Valid: true},
				TransactionIndex: sql.NullInt64{
					Int64: int64(fileEvent.TransactionIndex),
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
				BlockHash:   sql.NullString{String: fakeHeader.Hash, Valid: true},
				TxFrom:      sql.NullString{String: "fromAddress", Valid: true},
				TxTo:        sql.NullString{String: "toAddress", Valid: true},
			}

			_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			err := db.Get(&actualTx, `SELECT * FROM api.file_event_tx(
			    (SELECT (id, ilk_name, what, data, block_height, tx_idx)::api.file_event FROM api.ilk_files($1)))`,
				test_helpers.FakeIlk.Name)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: sql.NullString{String: "wrongTxHash", Valid: true},
				TransactionIndex: sql.NullInt64{
					Int64: int64(fileEvent.TransactionIndex) + 1,
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
				BlockHash:   sql.NullString{String: fakeHeader.Hash, Valid: true},
				TxFrom:      sql.NullString{String: "wrongFromAddress", Valid: true},
				TxTo:        sql.NullString{String: "wrongToAddress", Valid: true},
			}

			_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			err := db.Get(&actualTx, `SELECT * FROM api.file_event_tx(
			    (SELECT (id, ilk_name, what, data, block_height, tx_idx)::api.file_event FROM api.ilk_files($1)))`,
				test_helpers.FakeIlk.Name)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
