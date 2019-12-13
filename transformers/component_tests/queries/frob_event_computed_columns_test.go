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

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/queries/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Frob event computed columns", func() {
	var (
		blockOne, timestampOne int
		fakeGuy                = fakes.RandomString(42)
		headerOne              core.Header
		frobGethLog            types.Log
		frobEvent              event.InsertionModel
		vatRepository          vat.VatStorageRepository
		headerRepository       repositories.HeaderRepository
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		blockOne = rand.Int()
		timestampOne = int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepository)

		frobHeaderSyncLog := test_data.CreateTestLog(headerOne.Id, db)
		frobGethLog = frobHeaderSyncLog.Log

		frobEvent = test_data.VatFrobModelWithPositiveDart()
		urnID, urnErr := shared.GetOrCreateUrn(fakeGuy, test_helpers.FakeIlk.Hex, db)
		Expect(urnErr).NotTo(HaveOccurred())
		frobEvent.ColumnValues[constants.UrnColumn] = urnID
		frobEvent.ColumnValues[event.HeaderFK] = headerOne.Id
		frobEvent.ColumnValues[event.LogFK] = frobHeaderSyncLog.ID
		insertFrobErr := event.PersistModels([]event.InsertionModel{frobEvent}, db)
		Expect(insertFrobErr).NotTo(HaveOccurred())
	})

	Describe("frob_event_ilk", func() {
		It("returns ilk_state for a frob_event", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, headerOne, ilkValues, test_helpers.FakeIlkVatMetadatas, test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, headerOne.Timestamp, headerOne.Timestamp, ilkValues)

			var result test_helpers.IlkState
			getIlkErr := db.Get(&result,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated
                    FROM api.frob_event_ilk(
                        (SELECT (ilk_identifier, urn_identifier, dink, dart, ilk_rate, block_height, log_id)::api.frob_event FROM api.all_frobs($1))
                    )`, test_helpers.FakeIlk.Identifier)

			Expect(getIlkErr).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("frob_event_urn", func() {
		It("returns urn_state for a frob_event", func() {
			urnSetupData := test_helpers.GetUrnSetupData()
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			vatRepository.SetDB(db)
			test_helpers.CreateUrn(db, urnSetupData, headerOne, urnMetadata, vatRepository)

			var actualUrn test_helpers.UrnState
			getUrnErr := db.Get(&actualUrn,
				`SELECT urn_identifier, ilk_identifier FROM api.frob_event_urn(
                        (SELECT (ilk_identifier, urn_identifier, dink, dart, ilk_rate, block_height, log_id)::api.frob_event FROM api.all_frobs($1)))`,
				test_helpers.FakeIlk.Identifier)
			Expect(getUrnErr).NotTo(HaveOccurred())

			expectedUrn := test_helpers.UrnState{
				UrnIdentifier: fakeGuy,
				IlkIdentifier: test_helpers.FakeIlk.Identifier,
			}

			test_helpers.AssertUrn(actualUrn, expectedUrn)
		})
	})

	Describe("frob_event_tx", func() {
		It("returns transaction for a frob_event", func() {
			expectedTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("txHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(frobGethLog.TxIndex),
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(blockOne), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:      test_helpers.GetValidNullString("fromAddress"),
				TxTo:        test_helpers.GetValidNullString("toAddress"),
			}

			_, insertTxErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerOne.Id, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(insertTxErr).NotTo(HaveOccurred())

			var actualTx Tx
			getTxErr := db.Get(&actualTx, `SELECT * FROM api.frob_event_tx(
			    (SELECT (ilk_identifier, urn_identifier, dink, dart, ilk_rate, block_height, log_id)::api.frob_event FROM api.all_frobs($1)))`,
				test_helpers.FakeIlk.Identifier)

			Expect(getTxErr).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(frobGethLog.TxIndex) + 1,
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(blockOne), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertTxErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerOne.Id, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertTxErr).NotTo(HaveOccurred())

			var actualTx Tx
			getTxErr := db.Get(&actualTx, `SELECT * FROM api.frob_event_tx(
			    (SELECT (ilk_identifier, urn_identifier, dink, dart, ilk_rate, block_height, log_id)::api.frob_event FROM api.all_frobs($1)))`,
				test_helpers.FakeIlk.Identifier)

			Expect(getTxErr).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})

		It("does not return transaction from different block with same index", func() {
			headerZero := createHeader(blockOne-1, timestampOne-1, headerRepository)
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(frobGethLog.TxIndex),
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: headerZero.BlockNumber, Valid: true},
				BlockHash:   test_helpers.GetValidNullString(headerOne.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertTxErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerZero.Id, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertTxErr).NotTo(HaveOccurred())

			var actualTx Tx
			getTxErr := db.Get(&actualTx, `SELECT * FROM api.frob_event_tx(
			    (SELECT (ilk_identifier, urn_identifier, dink, dart, ilk_rate, block_height, log_id)::api.frob_event FROM api.all_frobs($1)))`,
				test_helpers.FakeIlk.Identifier)

			Expect(getTxErr).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
