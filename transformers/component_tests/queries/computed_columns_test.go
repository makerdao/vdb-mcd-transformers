package queries

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_file/ilk"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strings"
)

var _ = Describe("Extension function", func() {
	var (
		db         *postgres.DB
		fakeBlock  int
		fakeHeader core.Header
		headerId   int64
		err        error
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepository = repositories.NewHeaderRepository(db)
		fakeBlock = rand.Int()
		fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
		headerId, err = headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("using ilk and frob", func() {
		var (
			frobRepo  vat_frob.VatFrobRepository
			frobEvent vat_frob.VatFrobModel
			ilkValues map[string]string
		)

		BeforeEach(func() {
			frobRepo = vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			vatRepository.SetDB(db)
			catRepository.SetDB(db)
			jugRepository.SetDB(db)

			// Create an ilk
			ilkValues = test_helpers.GetIlkValues(0)
			createIlkAtBlock(fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)

			// Create a frob event
			frobEvent = test_data.VatFrobModelWithPositiveDart
			frobEvent.Ilk = test_helpers.FakeIlk.Hex
			err = frobRepo.Create(headerId, []interface{}{frobEvent})
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("ilk_state_frobs", func() {
			It("returns relevant frobs for an ilk_state", func() {
				irrelevantFrob := test_data.VatFrobModelWithPositiveDart
				irrelevantFrob.Ilk = test_helpers.AnotherFakeIlk.Hex
				irrelevantFrob.Urn = "anotherGuy"
				irrelevantFrob.TransactionIndex = frobEvent.TransactionIndex + 1
				err = frobRepo.Create(headerId, []interface{}{irrelevantFrob})
				Expect(err).NotTo(HaveOccurred())

				var ilkId int
				err = db.Get(&ilkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, test_helpers.FakeIlk.Hex)
				Expect(err).NotTo(HaveOccurred())

				var actualFrobs []test_helpers.FrobEvent
				err = db.Select(&actualFrobs,
					`SELECT ilk_name, urn_id, dink, dart FROM maker.ilk_state_frobs(
                        (SELECT (ilk_id, ilk_name, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated)::maker.ilk_state
                         FROM maker.get_ilk($1, $2))
                    )`, fakeBlock, ilkId)
				Expect(err).NotTo(HaveOccurred())

				expectedFrobs := []test_helpers.FrobEvent{{
					IlkName: test_helpers.FakeIlk.Name,
					UrnId:   frobEvent.Urn,
					Dink:    frobEvent.Dink,
					Dart:    frobEvent.Dart,
				}}

				Expect(actualFrobs).To(Equal(expectedFrobs))
			})
		})

		Describe("frob_event_ilk", func() {
			It("returns ilk_state for a frob_event", func() {
				expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

				var result test_helpers.IlkState
				err = db.Get(&result,
					`SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated
                    FROM maker.frob_event_ilk(
                        (SELECT (ilk_name, urn_id, dink, dart, block_height, tx_idx)::maker.frob_event FROM maker.all_frobs($1))
                    )`, test_helpers.FakeIlk.Name)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(expectedIlk))
			})
		})
	})

	Describe("using ilk and file", func() {
		var (
			fileRepo  ilk.VatFileIlkRepository
			fileEvent ilk.VatFileIlkModel
			ilkValues map[string]string
		)

		BeforeEach(func() {
			fileRepo = ilk.VatFileIlkRepository{}
			fileRepo.SetDB(db)
			vatRepository.SetDB(db)
			catRepository.SetDB(db)
			jugRepository.SetDB(db)

			// Create an ilk
			ilkValues = test_helpers.GetIlkValues(0)
			createIlkAtBlock(fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas)

			// Create a file event
			fileEvent = test_data.VatFileIlkDustModel
			fileEvent.Ilk = test_helpers.FakeIlk.Hex
			err = fileRepo.Create(headerId, []interface{}{fileEvent})
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("ilks_state_files", func() {
			It("returns file event for an ilk state", func() {
				var ilkId int
				err = db.Get(&ilkId, `SELECT id FROM maker.ilks WHERE ilk = $1`, test_helpers.FakeIlk.Hex)
				Expect(err).NotTo(HaveOccurred())

				var actualFiles []test_helpers.FileEvent
				err = db.Select(&actualFiles,
					`SELECT id, ilk_name, what, data FROM maker.ilk_state_files(
                        (SELECT (ilk_id, ilk_name, block_height, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated)::maker.ilk_state
                         FROM maker.get_ilk($1, $2))
                    )`, fakeBlock, ilkId)
				Expect(err).NotTo(HaveOccurred())

				expectedFiles := []test_helpers.FileEvent{{
					Id: strings.ToLower(test_data.EthVatFileIlkDustLog.Address.Hex()),
					IlkName: sql.NullString{
						String: test_helpers.FakeIlk.Name,
						Valid:  true,
					},
					What: fileEvent.What,
					Data: fileEvent.Data,
				}}

				Expect(actualFiles).To(Equal(expectedFiles))
			})
		})

		Describe("file_event_ilk", func() {
			It("returns ilk_state for a file_event", func() {
				expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

				var result test_helpers.IlkState
				err = db.Get(&result,
					`SELECT ilk_name, rate, art, spot, line, dust, chop, lump, flip, rho, duty, created, updated
                    FROM maker.file_event_ilk(
                        (SELECT (id, ilk_name, what, data, block_height, tx_idx)::maker.file_event FROM maker.ilk_files($1))
                    )`, test_helpers.FakeIlk.Name)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(expectedIlk))
			})
		})
	})

	Describe("using urn and frob", func() {
		var (
			frobRepo  vat_frob.VatFrobRepository
			frobEvent vat_frob.VatFrobModel
			fakeGuy   = "fakeAddress"
		)

		BeforeEach(func() {
			frobRepo = vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			vatRepository.SetDB(db)

			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			frobEvent = test_data.VatFrobModelWithPositiveDart
			frobEvent.Urn = fakeGuy
			frobEvent.Ilk = test_helpers.FakeIlk.Hex
			err = frobRepo.Create(headerId, []interface{}{frobEvent})
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("frob_event_urn", func() {
			It("returns urn_state for a frob_event", func() {
				var actualUrn test_helpers.UrnState
				err = db.Get(&actualUrn,
					`SELECT urn_id, ilk_name FROM maker.frob_event_urn(
                        (SELECT (ilk_name, urn_id, dink, dart, block_height, tx_idx)::maker.frob_event FROM maker.all_frobs($1)))`,
					test_helpers.FakeIlk.Name)
				Expect(err).NotTo(HaveOccurred())

				expectedUrn := test_helpers.UrnState{
					UrnId:   fakeGuy,
					IlkName: test_helpers.FakeIlk.Name,
				}

				test_helpers.AssertUrn(actualUrn, expectedUrn)
			})
		})

		Describe("urn_state_frobs", func() {
			It("returns relevant frobs for an urn_state", func() {
				var actualFrobs test_helpers.FrobEvent
				err = db.Get(&actualFrobs,
					`SELECT ilk_name, urn_id, dink, dart FROM maker.urn_state_frobs(
                        (SELECT (urn_id, ilk_name, block_height, ink, art, ratio, safe, created, updated)::maker.urn_state
                         FROM maker.all_urns($1))
                    )`, fakeBlock)
				Expect(err).NotTo(HaveOccurred())

				expectedFrobs := test_helpers.FrobEvent{
					IlkName: test_helpers.FakeIlk.Name,
					UrnId:   fakeGuy,
					Dink:    frobEvent.Dink,
					Dart:    frobEvent.Dart,
				}

				Expect(actualFrobs).To(Equal(expectedFrobs))
			})
		})
	})

	Describe("tx_era", func() {
		It("returns an era object for a transaction", func() {
			txFrom := "fromAddress"
			txTo := "toAddress"
			txIndex := rand.Intn(10)
			_, err = db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
                VALUES ($1, $2, $3, $4, $5)`, headerId, fakeHeader.Hash, txFrom, txIndex, txTo)
			Expect(err).NotTo(HaveOccurred())

			var actualEra Era
			err = db.Get(&actualEra, `SELECT * FROM maker.tx_era(
                    (SELECT (txs.hash, txs.tx_index, h.block_number, h.hash, txs.tx_from, txs.tx_to)::maker.tx
			        FROM header_sync_transactions txs
			        LEFT JOIN headers h ON h.id = txs.header_id)
			    )`)
			Expect(err).NotTo(HaveOccurred())

			expectedEra := Era{
				Epoch: fakeHeader.Timestamp,
				Iso:   "1973-07-10T00:11:51Z", // Z for Zulu, meaning UTC
			}
			Expect(actualEra).To(Equal(expectedEra))
		})
	})

	Describe("frob_event_tx", func() {
		var (
			frobRepo  vat_frob.VatFrobRepository
			frobEvent vat_frob.VatFrobModel
		)

		BeforeEach(func() {
			frobRepo = vat_frob.VatFrobRepository{}
			frobRepo.SetDB(db)
			frobEvent = test_data.VatFrobModelWithPositiveDart
			frobEvent.Ilk = test_helpers.FakeIlk.Hex
			err = frobRepo.Create(headerId, []interface{}{frobEvent})
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns transaction for a frob_event", func() {
			expectedTx := Tx{
				TransactionHash: sql.NullString{String: "txHash", Valid: true},
				TransactionIndex: sql.NullInt64{
					Int64: int64(frobEvent.TransactionIndex),
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
				BlockHash:   sql.NullString{String: fakeHeader.Hash, Valid: true},
				TxFrom:      sql.NullString{String: "fromAddress", Valid: true},
				TxTo:        sql.NullString{String: "toAddress", Valid: true},
			}

			_, err = db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(err).NotTo(HaveOccurred())

			var actualTx Tx
			err = db.Get(&actualTx, `SELECT * FROM maker.frob_event_tx(
			    (SELECT (ilk_name, urn_id, dink, dart, block_height, tx_idx)::maker.frob_event FROM maker.all_frobs($1)))`,
				test_helpers.FakeIlk.Name)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: sql.NullString{String: "wrongTxHash", Valid: true},
				TransactionIndex: sql.NullInt64{
					Int64: int64(frobEvent.TransactionIndex) + 1,
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
				BlockHash:   sql.NullString{String: fakeHeader.Hash, Valid: true},
				TxFrom:      sql.NullString{String: "wrongFromAddress", Valid: true},
				TxTo:        sql.NullString{String: "wrongToAddress", Valid: true},
			}

			_, err = db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(err).NotTo(HaveOccurred())

			var actualTx Tx
			err = db.Get(&actualTx, `SELECT * FROM maker.frob_event_tx(
			    (SELECT (ilk_name, urn_id, dink, dart, block_height, tx_idx)::maker.frob_event FROM maker.all_frobs($1)))`,
				test_helpers.FakeIlk.Name)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})

	Describe("file_event_tx", func() {
		var (
			fileRepo  ilk.VatFileIlkRepository
			fileEvent ilk.VatFileIlkModel
		)

		BeforeEach(func() {
			fileRepo = ilk.VatFileIlkRepository{}
			fileRepo.SetDB(db)
			fileEvent = test_data.VatFileIlkDustModel
			fileEvent.Ilk = test_helpers.FakeIlk.Hex
			err = fileRepo.Create(headerId, []interface{}{fileEvent})
			Expect(err).NotTo(HaveOccurred())
		})

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

			_, err = db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(err).NotTo(HaveOccurred())

			var actualTx Tx
			err = db.Get(&actualTx, `SELECT * FROM maker.file_event_tx(
			    (SELECT (id, ilk_name, what, data, block_height, tx_idx)::maker.file_event FROM maker.ilk_files($1)))`,
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

			_, err = db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(err).NotTo(HaveOccurred())

			var actualTx Tx
			err = db.Get(&actualTx, `SELECT * FROM maker.file_event_tx(
			    (SELECT (id, ilk_name, what, data, block_height, tx_idx)::maker.file_event FROM maker.ilk_files($1)))`,
				test_helpers.FakeIlk.Name)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})

type Era struct {
	Epoch string
	Iso   string
}

type Tx struct {
	TransactionHash  sql.NullString `db:"transaction_hash"`
	TransactionIndex sql.NullInt64  `db:"transaction_index"`
	BlockHeight      sql.NullInt64  `db:"block_height"`
	BlockHash        sql.NullString `db:"block_hash"`
	TxFrom           sql.NullString `db:"tx_from"`
	TxTo             sql.NullString `db:"tx_to"`
}
