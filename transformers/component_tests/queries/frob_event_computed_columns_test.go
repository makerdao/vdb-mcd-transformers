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

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Frob event computed columns", func() {
	var (
		db               *postgres.DB
		fakeBlock        int
		fakeGuy          = "fakeAddress"
		fakeHeader       core.Header
		frobRepo         vat_frob.VatFrobRepository
		frobEvent        shared.InsertionModel
		headerId         int64
		vatRepository    vat.VatStorageRepository
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

		frobRepo = vat_frob.VatFrobRepository{}
		frobRepo.SetDB(db)
		frobEvent = test_data.CopyModel(test_data.VatFrobModelWithPositiveDart)
		frobEvent.ForeignKeyValues[constants.UrnFK] = fakeGuy
		frobEvent.ForeignKeyValues[constants.IlkFK] = test_helpers.FakeIlk.Hex
		insertFrobErr := frobRepo.Create(headerId, []shared.InsertionModel{frobEvent})
		Expect(insertFrobErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("frob_event_ilk", func() {
		It("returns ilk_state for a frob_event", func() {
			ilkValues := test_helpers.GetIlkValues(0)
			test_helpers.CreateIlk(db, fakeHeader, ilkValues, test_helpers.FakeIlkVatMetadatas,
				test_helpers.FakeIlkCatMetadatas, test_helpers.FakeIlkJugMetadatas, test_helpers.FakeIlkSpotMetadatas)

			expectedIlk := test_helpers.IlkStateFromValues(test_helpers.FakeIlk.Hex, fakeHeader.Timestamp, fakeHeader.Timestamp, ilkValues)

			var result test_helpers.IlkState
			getIlkErr := db.Get(&result,
				`SELECT ilk_identifier, rate, art, spot, line, dust, chop, lump, flip, rho, duty, pip, mat, created, updated
                    FROM api.frob_event_ilk(
                        (SELECT (ilk_identifier, urn_identifier, dink, dart, block_height, tx_idx)::api.frob_event FROM api.all_frobs($1))
                    )`, test_helpers.FakeIlk.Identifier)

			Expect(getIlkErr).NotTo(HaveOccurred())
			Expect(result).To(Equal(expectedIlk))
		})
	})

	Describe("frob_event_urn", func() {
		It("returns urn_state for a frob_event", func() {
			urnSetupData := test_helpers.GetUrnSetupData(fakeBlock, 1)
			urnSetupData.Header.Hash = fakeHeader.Hash
			urnMetadata := test_helpers.GetUrnMetadata(test_helpers.FakeIlk.Hex, fakeGuy)
			vatRepository.SetDB(db)
			test_helpers.CreateUrn(urnSetupData, urnMetadata, vatRepository, headerRepository)

			var actualUrn test_helpers.UrnState
			getUrnErr := db.Get(&actualUrn,
				`SELECT urn_identifier, ilk_identifier FROM api.frob_event_urn(
                        (SELECT (ilk_identifier, urn_identifier, dink, dart, block_height, tx_idx)::api.frob_event FROM api.all_frobs($1)))`,
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
					Int64: int64(frobEvent.ColumnValues["tx_idx"].(uint)),
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(fakeHeader.Hash),
				TxFrom:      test_helpers.GetValidNullString("fromAddress"),
				TxTo:        test_helpers.GetValidNullString("toAddress"),
			}

			_, insertTxErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, expectedTx.TransactionHash, expectedTx.TxFrom,
				expectedTx.TransactionIndex, expectedTx.TxTo)
			Expect(insertTxErr).NotTo(HaveOccurred())

			var actualTx Tx
			getTxErr := db.Get(&actualTx, `SELECT * FROM api.frob_event_tx(
			    (SELECT (ilk_identifier, urn_identifier, dink, dart, block_height, tx_idx)::api.frob_event FROM api.all_frobs($1)))`,
				test_helpers.FakeIlk.Identifier)

			Expect(getTxErr).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(frobEvent.ColumnValues["tx_idx"].(uint)) + 1,
					Valid: true,
				},
				BlockHeight: sql.NullInt64{Int64: int64(fakeBlock), Valid: true},
				BlockHash:   test_helpers.GetValidNullString(fakeHeader.Hash),
				TxFrom:      test_helpers.GetValidNullString("wrongFromAddress"),
				TxTo:        test_helpers.GetValidNullString("wrongToAddress"),
			}

			_, insertTxErr := db.Exec(`INSERT INTO header_sync_transactions (header_id, hash, tx_from, tx_index, tx_to)
				VALUES ($1, $2, $3, $4, $5)`, headerId, wrongTx.TransactionHash, wrongTx.TxFrom,
				wrongTx.TransactionIndex, wrongTx.TxTo)
			Expect(insertTxErr).NotTo(HaveOccurred())

			var actualTx Tx
			getTxErr := db.Get(&actualTx, `SELECT * FROM api.frob_event_tx(
			    (SELECT (ilk_identifier, urn_identifier, dink, dart, block_height, tx_idx)::api.frob_event FROM api.all_frobs($1)))`,
				test_helpers.FakeIlk.Identifier)

			Expect(getTxErr).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
