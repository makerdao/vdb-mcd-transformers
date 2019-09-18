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

var _ = Describe("Shared repository", func() {
	var db *postgres.DB

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
	})

	const hexIlk = "0x464b450000000000000000000000000000000000000000000000000000000000"

	Describe("Create function", func() {
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
		ADD COLUMN testevent INTEGER NOT NULL DEFAULT 0;`

		var (
			headerRepository repositories.HeaderRepository
			testModel        InsertionModel
			fakeLog, _       = json.Marshal("fake log")
		)

		BeforeEach(func() {
			_, _ = db.Exec(createTestEventTableQuery)
			_, _ = db.Exec(addCheckedColumnQuery)
			headerRepository = repositories.NewHeaderRepository(db)

			testModel = InsertionModel{
				SchemaName: "maker",
				TableName:  "testEvent",
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
			db.MustExec(`ALTER TABLE public.checked_headers DROP COLUMN testevent;`)
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
					SchemaName:     "maker",
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
					SchemaName: "maker",
					TableName:  "testEvent",
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
				delete(modelToQuery, "makertestEvent")
				header := fakes.GetFakeHeader(1)
				headerID, headerErr := headerRepository.CreateOrUpdateHeader(header)
				Expect(headerErr).NotTo(HaveOccurred())

				createErr := Create(headerID, []InsertionModel{brokenModel}, db)
				Expect(createErr).To(HaveOccurred())
				// Remove incorrect query, so other tests won't get it
				delete(modelToQuery, "makertestEvent")
			})
		})

		It("upserts queries with conflicting source", func() {
			header := fakes.GetFakeHeader(1)
			headerID, headerErr := headerRepository.CreateOrUpdateHeader(header)
			Expect(headerErr).NotTo(HaveOccurred())

			conflictingModel := InsertionModel{
				SchemaName: "maker",
				TableName:  "testEvent",
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
			dbErr := db.Get(&checked, `SELECT testevent FROM public.checked_headers;`)
			Expect(dbErr).NotTo(HaveOccurred())
			Expect(checked).To(Equal(1))
		})

		It("generates correct queries", func() {
			actualQuery := generateInsertionQuery(testModel)
			expectedQuery := `INSERT INTO maker.testEvent (header_id, log_idx, tx_idx, raw_log, ilk_id, urn_id, variable1) VALUES($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET header_id = $1, log_idx = $2, tx_idx = $3, raw_log = $4, ilk_id = $5, urn_id = $6, variable1 = $7;`
			Expect(actualQuery).To(Equal(expectedQuery))
		})

		Describe("FK columns", func() {
			It("looks up FK id and persists in columnToValue for IlkFK", func() {
				foreignKeyValues := ForeignKeyValues{constants.IlkFK: hexIlk}
				columnToValue := ColumnValues{}

				tx, txErr := db.Beginx()
				Expect(txErr).NotTo(HaveOccurred())
				fkErr := populateForeignKeyIDs(foreignKeyValues, columnToValue, tx)
				Expect(fkErr).NotTo(HaveOccurred())
				commitErr := tx.Commit()
				Expect(commitErr).NotTo(HaveOccurred())

				ilkIdentifier := DecodeHexToText(hexIlk)
				var expectedIlkID int
				ilkErr := db.Get(&expectedIlkID, `SELECT id FROM maker.ilks WHERE identifier = $1`, ilkIdentifier)
				Expect(ilkErr).NotTo(HaveOccurred())
				actualIlkID := columnToValue[string(constants.IlkFK)].(int)
				Expect(actualIlkID).To(Equal(expectedIlkID))
			})

			It("looks up FK id and persists in columnToValue for UrnFK", func() {
				guy := "0x12345"
				foreignKeyValues := ForeignKeyValues{constants.UrnFK: guy}
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

			It("looks up FK id and persists in columnToValue for AddressFK", func() {
				foreignKeyValues := ForeignKeyValues{constants.AddressFK: fakes.FakeAddress.Hex()}
				columnToValue := ColumnValues{}

				tx, txErr := db.Beginx()
				Expect(txErr).NotTo(HaveOccurred())
				fkErr := populateForeignKeyIDs(foreignKeyValues, columnToValue, tx)
				Expect(fkErr).NotTo(HaveOccurred())
				commitErr := tx.Commit()
				Expect(commitErr).NotTo(HaveOccurred())

				var expectedAddressID int
				addressErr := db.Get(&expectedAddressID, `SELECT id FROM public.addresses WHERE address = $1`, fakes.FakeAddress.Hex())
				Expect(addressErr).NotTo(HaveOccurred())
				actualAddressID := columnToValue[string(constants.AddressFK)].(int)
				Expect(actualAddressID).To(Equal(expectedAddressID))
			})
		})
	})

	Describe("GetOrCreateIlk", func() {
		It("returns same ilk id for value with and without hex prefix", func() {
			ilkWithPrefix := hexIlk
			ilkWithoutPrefix := ilkWithPrefix[2:]

			ilkIdOne, ilkErrOne := GetOrCreateIlk(ilkWithPrefix, db)
			Expect(ilkErrOne).NotTo(HaveOccurred())

			ilkIdTwo, ilkErrTwo := GetOrCreateIlk(ilkWithoutPrefix, db)
			Expect(ilkErrTwo).NotTo(HaveOccurred())

			Expect(ilkIdOne).NotTo(BeZero())
			Expect(ilkIdOne).To(Equal(ilkIdTwo))
		})
	})

	Describe("GetOrCreateIlkInTransaction", func() {
		It("returns same ilk id for value with and without hex prefix", func() {
			ilkWithPrefix := hexIlk
			ilkWithoutPrefix := ilkWithPrefix[2:]

			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			ilkIdOne, ilkErrOne := GetOrCreateIlkInTransaction(ilkWithPrefix, tx)
			Expect(ilkErrOne).NotTo(HaveOccurred())

			ilkIdTwo, ilkErrTwo := GetOrCreateIlkInTransaction(ilkWithoutPrefix, tx)
			Expect(ilkErrTwo).NotTo(HaveOccurred())

			commitErr := tx.Commit()
			Expect(commitErr).NotTo(HaveOccurred())

			Expect(ilkIdOne).NotTo(BeZero())
			Expect(ilkIdOne).To(Equal(ilkIdTwo))
		})
	})

	Describe("GetOrCreateAddress", func() {
		It("creates an address record", func() {
			_, err := GetOrCreateAddress(fakes.FakeAddress.Hex(), db)
			Expect(err).NotTo(HaveOccurred())

			var address string
			db.Get(&address, `SELECT address from addresses LIMIT 1`)
			Expect(address).To(Equal(fakes.FakeAddress.Hex()))
		})

		It("returns the id for an address that already exists", func() {
			//create the address record
			createAddressId, createErr := GetOrCreateAddress(fakes.FakeAddress.Hex(), db)
			Expect(createErr).NotTo(HaveOccurred())

			//get the address record
			getAddressId, getErr := GetOrCreateAddress(fakes.FakeAddress.Hex(), db)
			Expect(getErr).NotTo(HaveOccurred())

			Expect(createAddressId).To(Equal(getAddressId))

			var addressCount int
			db.Get(&addressCount, `SELECT count(*) from addresses`)
			Expect(addressCount).To(Equal(1))
		})
	})

	Describe("GetOrCreateAddressInTransaction", func() {
		It("creates an address record", func() {
			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			_, createErr := GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
			Expect(createErr).NotTo(HaveOccurred())

			commitErr := tx.Commit()
			Expect(commitErr).NotTo(HaveOccurred())

			var address string
			db.Get(&address, `SELECT address from addresses LIMIT 1`)
			Expect(address).To(Equal(fakes.FakeAddress.Hex()))
		})

		It("returns the id for an address that already exists", func() {
			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			//create the address record
			createAddressId, createErr := GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
			Expect(createErr).NotTo(HaveOccurred())

			//get the address record
			getAddressId, getErr := GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
			Expect(getErr).NotTo(HaveOccurred())

			commitErr := tx.Commit()
			Expect(commitErr).NotTo(HaveOccurred())

			Expect(createAddressId).To(Equal(getAddressId))

			var addressCount int
			db.Get(&addressCount, `SELECT count(*) from addresses`)
			Expect(addressCount).To(Equal(1))
		})

		It("doesn't persist the address if the transaction is rolled back", func() {
			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			_, createErr := GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
			Expect(createErr).NotTo(HaveOccurred())

			commitErr := tx.Rollback()
			Expect(commitErr).NotTo(HaveOccurred())

			var addressCount int
			db.Get(&addressCount, `SELECT count(*) from addresses`)
			Expect(addressCount).To(Equal(0))
		})
	})
})

type TestEvent struct {
	LogIdx    string `db:"log_idx"`
	TxIdx     string `db:"tx_idx"`
	RawLog    string `db:"raw_log"`
	Variable1 string
}
