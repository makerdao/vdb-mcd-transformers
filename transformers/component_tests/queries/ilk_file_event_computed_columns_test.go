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

	storage_helper "github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vat_file/ilk"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ilk file event computed columns", func() {
	var (
		blockOne, timestampOne int
		fakeGethLog            types.Log
		fileEvent              shared.InsertionModel
		fileRepo               ilk.VatFileIlkRepository
		headerOne              core.Header
		headerRepository       repositories.HeaderRepository
		diffID                 int64
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepository)
		fakeHeaderSyncLog := test_data.CreateTestLog(headerOne.Id, db)
		fakeGethLog = fakeHeaderSyncLog.Log

		fileRepo = ilk.VatFileIlkRepository{}
		fileRepo.SetDB(db)
		fileEvent = test_data.VatFileIlkDustModel()
		fileEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		fileEvent.ColumnValues[constants.HeaderFK] = headerOne.Id
		fileEvent.ColumnValues[constants.LogFK] = fakeHeaderSyncLog.ID
		insertFileErr := fileRepo.Create([]shared.InsertionModel{fileEvent})
		Expect(insertFileErr).NotTo(HaveOccurred())

		diffID = storage_helper.CreateFakeDiffRecord(db)
	})

	Describe("ilk_file_event_ilk", func() {
		It("returns ilk_state for an ilk_file_event", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, diffID, headerOne, ilkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, ilkValues)

			var result test_helpers.IlkState
			err := db.Get(&result,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated
                    FROM api.ilk_file_event_ilk(
                        (SELECT (ilk_identifier, what, data, block_height, log_id)::api.ilk_file_event FROM api.all_ilk_file_events($1))
                    )`, test_helpers.FakeIlk.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("ilk_file_event_tx", func() {
		It("returns transaction for an ilk_file_event", func() {
			expectedTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("txHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(fakeGethLog.TxIndex),
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: headerOne.BlockNumber, Valid: true},
				BlockHash:   test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:      test_helpers.GetValidNullString("fromAddress"),
				TxTo:        test_helpers.GetValidNullString("toAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerOne.Id, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			err := db.Get(&actualTx, `SELECT * FROM api.ilk_file_event_tx(
			    (SELECT (ilk_identifier, what, data, block_height, log_id)::api.ilk_file_event FROM api.all_ilk_file_events($1)))`,
				test_helpers.FakeIlk.Identifier)

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
				BlockHeight: sql.NullInt64{Int64: headerOne.BlockNumber, Valid: true},
				BlockHash:   test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerOne.Id, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			err := db.Get(&actualTx, `SELECT * FROM api.ilk_file_event_tx(
			    (SELECT (ilk_identifier, what, data, block_height, log_id)::api.ilk_file_event FROM api.all_ilk_file_events($1)))`,
				test_helpers.FakeIlk.Identifier)

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

			_, insertErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerZero.Id, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertErr).NotTo(HaveOccurred())

			var actualTx Tx
			err := db.Get(&actualTx, `SELECT * FROM api.ilk_file_event_tx(
			    (SELECT (ilk_identifier, what, data, block_height, log_id)::api.ilk_file_event FROM api.all_ilk_file_events($1)))`,
				test_helpers.FakeIlk.Identifier)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
