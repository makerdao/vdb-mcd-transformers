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
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vow_fess"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
)

var _ = Describe("Sin queue event computed columns", func() {
	var (
		db               *postgres.DB
		fakeBlock        int
		fakeEra          string
		fakeHeader       core.Header
		vowFessEvent     shared.InsertionModel
		vowFessRepo      vow_fess.VowFessRepository
		headerId         int64
		headerRepository repositories.HeaderRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)

		headerRepository = repositories.NewHeaderRepository(db)
		fakeBlock = rand.Int()
		fakeEra = strconv.Itoa(int(rand.Int31()))
		fakeHeader = fakes.GetFakeHeader(int64(fakeBlock))
		fakeHeader.Timestamp = fakeEra
		var insertHeaderErr error
		headerId, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakeHeader)
		Expect(insertHeaderErr).NotTo(HaveOccurred())

		vowFessRepo = vow_fess.VowFessRepository{}
		vowFessRepo.SetDB(db)
		vowFessEvent = test_data.VowFessModel
		insertErr := vowFessRepo.Create(headerId, []shared.InsertionModel{vowFessEvent})
		Expect(insertErr).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("sin_queue_event_tx", func() {
		It("returns transaction for a sin_queue_event", func() {
			expectedTx := Tx{
				TransactionHash:  test_helpers.GetValidNullString("txHash"),
				TransactionIndex: sql.NullInt64{Int64: int64(vowFessEvent.ColumnValues["tx_idx"].(uint)), Valid: true},
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
				SELECT * FROM api.sin_queue_event_tx(
					(SELECT (era, act, block_height, tx_idx)::api.sin_queue_event FROM api.all_sin_queue_events($1)))`,
				fakeEra)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(Equal(expectedTx))
		})

		It("does not return transaction from same block with different index", func() {
			wrongTx := Tx{
				TransactionHash: test_helpers.GetValidNullString("wrongTxHash"),
				TransactionIndex: sql.NullInt64{
					Int64: int64(vowFessEvent.ColumnValues["tx_idx"].(uint)) + 1,
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
				SELECT * FROM api.sin_queue_event_tx(
					(SELECT (era, act, block_height, tx_idx)::api.sin_queue_event FROM api.all_sin_queue_events($1)))`,
				fakeEra)

			Expect(err).NotTo(HaveOccurred())
			Expect(actualTx).To(BeZero())
		})
	})
})
