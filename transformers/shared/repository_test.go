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

package shared

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var _ = Describe("Create function", func() {
	const createTestEventTableQuery = `CREATE TABLE maker.testEvent(
		id        SERIAL PRIMARY KEY,
		header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
		variable1 TEXT,
		ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
		urn_id    INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
		log_idx   INTEGER NOT NULL,
		tx_idx    INTEGER NOT NULL,
		raw_log   JSONB,
		UNIQUE (header_id, tx_idx, log_idx)
		);`
	const addCheckedColumnQuery = `ALTER TABLE public.checked_headers
		ADD COLUMN testevent_checked INTEGER NOT NULL DEFAULT 0;`

	const hexIlk = "0x464b450000000000000000000000000000000000000000000000000000000000"

	var (
		headerRepository repositories.HeaderRepository
		testModel        InsertionModel
		db               *postgres.DB
		fakeLog, _       = json.Marshal("fake log")
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		_, _ = db.Exec(createTestEventTableQuery)
		_, _ = db.Exec(addCheckedColumnQuery)
		headerRepository = repositories.NewHeaderRepository(db)

		testModel = InsertionModel{
			TableName: "testEvent",
			OrderedColumns: []string{
				"header_id", "log_idx", "tx_idx", "raw_log", string(constants.IlkFK), string(constants.UrnFK), "variable1",
			},
			ColumnValues: ColumnValues{
				"log_idx":   "1",
				"tx_idx":    "2",
				"raw_log":   fakeLog,
				"variable1": "value1",
			},
			ForeignKeyValues: ForeignKeyValues{
				constants.IlkFK: hexIlk,
				constants.UrnFK: "0x12345",
			},
		}
	})

	AfterEach(func() {
		db.MustExec(`DROP TABLE maker.testEvent;`)
		db.MustExec(`ALTER TABLE public.checked_headers DROP COLUMN testevent_checked;`)
	})

	// Needs to run before the other tests, since those insert keys in map
	It("memoizes queries", func() {
		Expect(len(modelToQuery)).To(Equal(0))
		getMemoizedQuery(testModel)
		Expect(len(modelToQuery)).To(Equal(1))
		getMemoizedQuery(testModel)
		Expect(len(modelToQuery)).To(Equal(1))
	})

	It("persists a model to postgres", func() {
		header := fakes.GetFakeHeader(1)
		headerID, headerErr := headerRepository.CreateOrUpdateHeader(header)
		Expect(headerErr).NotTo(HaveOccurred())

		createErr := Create(headerID, []InsertionModel{testModel}, db)
		Expect(createErr).NotTo(HaveOccurred())

		var res TestEvent
		dbErr := db.Get(&res, `SELECT log_idx, tx_idx, raw_log, variable1
            FROM maker.testEvent;`)
		Expect(dbErr).NotTo(HaveOccurred())

		Expect(res.LogIdx).To(Equal(testModel.ColumnValues["log_idx"]))
		Expect(res.TxIdx).To(Equal(testModel.ColumnValues["tx_idx"]))
		Expect(res.Variable1).To(Equal(testModel.ColumnValues["variable1"]))
	})

	Describe("returns errors", func() {
		It("for empty model slice", func() {
			err := Create(0, []InsertionModel{}, db)
			Expect(err).To(MatchError("repository got empty model slice"))
		})

		It("for unknown foreign keys", func() {
			brokenModel := InsertionModel{
				TableName:      "testEvent",
				OrderedColumns: nil,
				ColumnValues:   nil,
				ForeignKeyValues: ForeignKeyValues{
					"unknownFK": "value",
				},
			}
			err := Create(0, []InsertionModel{brokenModel}, db)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("repository got unrecognised FK"))
			Expect(err.Error()).Should(ContainSubstring("error gettings FK ids"))
		})

		It("for failed SQL inserts", func() {
			brokenModel := InsertionModel{
				TableName: "testEvent",
				// Wrong name of last column compared to DB, will generate incorrect query
				OrderedColumns: []string{
					"header_id", "log_idx", "tx_idx", "raw_log", string(constants.IlkFK), string(constants.UrnFK), "variable2",
				},
				ColumnValues: ColumnValues{
					"log_idx":   "1",
					"tx_idx":    "2",
					"raw_log":   fakeLog,
					"variable1": "value1",
				},
				ForeignKeyValues: ForeignKeyValues{
					constants.IlkFK: hexIlk,
					constants.UrnFK: "0x12345",
				},
			}

			// Remove cached queries, or we won't generate a new (incorrect) one
			delete(modelToQuery, "testEvent")
			header := fakes.GetFakeHeader(1)
			headerID, headerErr := headerRepository.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			createErr := Create(headerID, []InsertionModel{brokenModel}, db)
			Expect(createErr).To(HaveOccurred())
			// Remove incorrect query, so other tests won't get it
			delete(modelToQuery, "testEvent")
		})
	})

	It("upserts queries with conflicting source", func() {
		header := fakes.GetFakeHeader(1)
		headerID, headerErr := headerRepository.CreateOrUpdateHeader(header)
		Expect(headerErr).NotTo(HaveOccurred())

		conflictingModel := InsertionModel{
			TableName: "testEvent",
			OrderedColumns: []string{
				"header_id", "log_idx", "tx_idx", "raw_log", string(constants.IlkFK), string(constants.UrnFK), "variable1",
			},
			ColumnValues: ColumnValues{
				"log_idx":   "1",
				"tx_idx":    "2",
				"raw_log":   fakeLog,
				"variable1": "conflictingValue",
			},
			ForeignKeyValues: ForeignKeyValues{
				constants.IlkFK: hexIlk,
				constants.UrnFK: "0x12345",
			},
		}

		createErr := Create(headerID, []InsertionModel{testModel, conflictingModel}, db)
		Expect(createErr).NotTo(HaveOccurred())

		var res TestEvent
		dbErr := db.Get(&res, `SELECT log_idx, tx_idx, raw_log, variable1
            FROM maker.testEvent;`)
		Expect(dbErr).NotTo(HaveOccurred())
		Expect(res.Variable1).To(Equal(conflictingModel.ColumnValues["variable1"]))
	})

	It("marks headers checked", func() {
		header := fakes.GetFakeHeader(1)
		headerID, headerErr := headerRepository.CreateOrUpdateHeader(header)
		Expect(headerErr).NotTo(HaveOccurred())

		createErr := Create(headerID, []InsertionModel{testModel}, db)
		Expect(createErr).NotTo(HaveOccurred())

		var checked int
		dbErr := db.Get(&checked, `SELECT testevent_checked FROM public.checked_headers;`)
		Expect(dbErr).NotTo(HaveOccurred())
		Expect(checked).To(Equal(1))
	})

	It("generates correct queries", func() {
		actualQuery := generateInsertionQuery(testModel)
		expectedQuery := `INSERT INTO maker.testEvent (header_id, log_idx, tx_idx, raw_log, ilk_id, urn_id, variable1) VALUES($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET header_id = $1, log_idx = $2, tx_idx = $3, raw_log = $4, ilk_id = $5, urn_id = $6, variable1 = $7;`
		Expect(actualQuery).To(Equal(expectedQuery))
	})

	It("looks up FK id and persists in columnToValue", func() {
		guy := "0x12345"
		foreignKeyValues := ForeignKeyValues{constants.IlkFK: hexIlk, constants.UrnFK: guy}
		columnToValue := ColumnValues{}

		tx, txErr := db.Beginx()
		Expect(txErr).NotTo(HaveOccurred())
		fkErr := populateForeignKeyIDs(foreignKeyValues, columnToValue, tx)
		Expect(fkErr).NotTo(HaveOccurred())
		commitErr := tx.Commit()
		Expect(commitErr).NotTo(HaveOccurred())

		var expectedUrnID int
		urnErr := db.Get(&expectedUrnID, `SELECT id FROM maker.urns WHERE identifier = $1`, guy)
		Expect(urnErr).NotTo(HaveOccurred())
		actualUrnID := columnToValue[string(constants.UrnFK)].(int)
		Expect(actualUrnID).To(Equal(expectedUrnID))
	})
})

type TestEvent struct {
	LogIdx    string `db:"log_idx"`
	TxIdx     string `db:"tx_idx"`
	RawLog    string `db:"raw_log"`
	Variable1 string
}
