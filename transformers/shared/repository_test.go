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

package shared_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"strconv"
)

var _ = Describe("Shared repository", func() {
	var db *postgres.DB
	const hexIlk = "0x464b450000000000000000000000000000000000000000000000000000000000"

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
	})

	Describe("Create", func() {
		const createTestEventTableQuery = `CREATE TABLE maker.testEvent(
		id        SERIAL PRIMARY KEY,
		header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
		log_id    BIGINT  NOT NULL REFERENCES header_sync_logs (id) ON DELETE CASCADE,
		variable1 TEXT,
		ilk_id    INTEGER NOT NULL REFERENCES maker.ilks (id) ON DELETE CASCADE,
		urn_id    INTEGER NOT NULL REFERENCES maker.urns (id) ON DELETE CASCADE,
		UNIQUE (header_id, log_id)
		);`
		const addCheckedColumnQuery = `ALTER TABLE public.checked_headers
		ADD COLUMN testevent INTEGER NOT NULL DEFAULT 0;`

		var (
			headerID, logID  int64
			headerRepository repositories.HeaderRepository
			testModel        shared.InsertionModel
		)

		BeforeEach(func() {
			db = test_config.NewTestDB(test_config.NewTestNode())
			test_config.CleanTestDB(db)
		})

		BeforeEach(func() {
			_, _ = db.Exec(createTestEventTableQuery)
			_, _ = db.Exec(addCheckedColumnQuery)
			headerRepository = repositories.NewHeaderRepository(db)
			var insertHeaderErr error
			headerID, insertHeaderErr = headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(insertHeaderErr).NotTo(HaveOccurred())
			headerSyncLog := test_data.CreateTestLog(headerID, db)
			logID = headerSyncLog.ID

			testModel = shared.InsertionModel{
				TableName: "testEvent",
				OrderedColumns: []string{
					constants.HeaderFK, constants.LogFK, string(constants.IlkFK), string(constants.UrnFK), "variable1",
				},
				ColumnValues: shared.ColumnValues{
					constants.HeaderFK: headerID,
					constants.LogFK:    strconv.FormatInt(logID, 10),
					"variable1":        "value1",
				},
				ForeignKeyValues: shared.ForeignKeyValues{
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
			Expect(len(shared.ModelToQuery)).To(Equal(0))
			shared.GetMemoizedQuery(testModel)
			Expect(len(shared.ModelToQuery)).To(Equal(1))
			shared.GetMemoizedQuery(testModel)
			Expect(len(shared.ModelToQuery)).To(Equal(1))
		})

		It("persists a model to postgres", func() {
			createErr := shared.Create([]shared.InsertionModel{testModel}, db)
			Expect(createErr).NotTo(HaveOccurred())

			var res TestEvent
			dbErr := db.Get(&res, `SELECT log_id, variable1 FROM maker.testEvent;`)
			Expect(dbErr).NotTo(HaveOccurred())

			Expect(res.LogID).To(Equal(testModel.ColumnValues[constants.LogFK]))
			Expect(res.Variable1).To(Equal(testModel.ColumnValues["variable1"]))
		})

		Describe("returns errors", func() {
			It("for empty model slice", func() {
				err := shared.Create([]shared.InsertionModel{}, db)
				Expect(err).To(MatchError("repository got empty model slice"))
			})

			It("for unknown foreign keys", func() {
				brokenModel := shared.InsertionModel{
					TableName:      "testEvent",
					OrderedColumns: nil,
					ColumnValues:   shared.ColumnValues{constants.HeaderFK: 0},
					ForeignKeyValues: shared.ForeignKeyValues{
						"unknownFK": "value",
					},
				}
				err := shared.Create([]shared.InsertionModel{brokenModel}, db)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("repository got unrecognised FK"))
				Expect(err.Error()).Should(ContainSubstring("error gettings FK ids"))
			})

			It("upserts queries with conflicting source", func() {
				header := fakes.GetFakeHeader(1)
				headerID, headerErr := headerRepository.CreateOrUpdateHeader(header)
				Expect(headerErr).NotTo(HaveOccurred())

				conflictingModel := shared.InsertionModel{
					TableName: "testEvent",
					OrderedColumns: []string{
						constants.HeaderFK, constants.LogFK, string(constants.IlkFK), string(constants.UrnFK), "variable2",
					},
					ColumnValues: shared.ColumnValues{
						constants.HeaderFK: headerID,
						constants.LogFK:    logID,
						"variable1":        "value1",
					},
					ForeignKeyValues: shared.ForeignKeyValues{
						constants.IlkFK: hexIlk,
						constants.UrnFK: "0x12345",
					},
				}

				// Remove cached queries, or we won't generate a new (incorrect) one
				delete(shared.ModelToQuery, "testEvent")
				conflictingModel.ColumnValues[constants.HeaderFK] = headerID

				createErr := shared.Create([]shared.InsertionModel{conflictingModel}, db)
				Expect(createErr).To(HaveOccurred())
				// Remove incorrect query, so other tests won't get it
				delete(shared.ModelToQuery, "testEvent")
			})
		})

		It("upserts queries with conflicting source", func() {
			conflictingModel := shared.InsertionModel{
				TableName: "testEvent",
				OrderedColumns: []string{
					constants.HeaderFK, constants.LogFK, string(constants.IlkFK), string(constants.UrnFK), "variable1",
				},
				ColumnValues: shared.ColumnValues{
					constants.HeaderFK: headerID,
					constants.LogFK:    logID,
					"variable1":        "conflictingValue",
				},
				ForeignKeyValues: shared.ForeignKeyValues{
					constants.IlkFK: hexIlk,
					constants.UrnFK: "0x12345",
				},
			}

			createErr := shared.Create([]shared.InsertionModel{testModel, conflictingModel}, db)
			Expect(createErr).NotTo(HaveOccurred())

			var res TestEvent
			dbErr := db.Get(&res, `SELECT log_id, variable1 FROM maker.testEvent;`)
			Expect(dbErr).NotTo(HaveOccurred())
			Expect(res.Variable1).To(Equal(conflictingModel.ColumnValues["variable1"]))
		})

		It("generates correct queries", func() {
			actualQuery := shared.GenerateInsertionQuery(testModel)
			expectedQuery := `INSERT INTO maker.testEvent (header_id, log_id, ilk_id, urn_id, variable1) VALUES($1, $2, $3, $4, $5)
		ON CONFLICT (header_id, log_id) DO UPDATE SET header_id = $1, log_id = $2, ilk_id = $3, urn_id = $4, variable1 = $5;`
			Expect(actualQuery).To(Equal(expectedQuery))
		})

		It("marks log transformed", func() {
			createErr := shared.Create([]shared.InsertionModel{testModel}, db)
			Expect(createErr).NotTo(HaveOccurred())

			var logTransformed bool
			getErr := db.Get(&logTransformed, `SELECT transformed FROM public.header_sync_logs WHERE id = $1`, logID)
			Expect(getErr).NotTo(HaveOccurred())
			Expect(logTransformed).To(BeTrue())
		})

		Describe("FK columns", func() {
			It("looks up FK id and persists in columnToValue for IlkFK", func() {
				foreignKeyValues := shared.ForeignKeyValues{constants.IlkFK: hexIlk}
				columnToValue := shared.ColumnValues{}

				tx, txErr := db.Beginx()
				Expect(txErr).NotTo(HaveOccurred())
				fkErr := shared.PopulateForeignKeyIDs(foreignKeyValues, columnToValue, tx)
				Expect(fkErr).NotTo(HaveOccurred())
				commitErr := tx.Commit()
				Expect(commitErr).NotTo(HaveOccurred())

				ilkIdentifier := shared.DecodeHexToText(hexIlk)
				var expectedIlkID int64
				ilkErr := db.Get(&expectedIlkID, `SELECT id FROM maker.ilks WHERE identifier = $1`, ilkIdentifier)
				Expect(ilkErr).NotTo(HaveOccurred())
				actualIlkID := columnToValue[string(constants.IlkFK)].(int64)
				Expect(actualIlkID).To(Equal(expectedIlkID))
			})

			It("looks up FK id and persists in columnToValue for UrnFK", func() {
				guy := "0x12345"
				foreignKeyValues := shared.ForeignKeyValues{constants.UrnFK: guy}
				columnToValue := shared.ColumnValues{}

				tx, txErr := db.Beginx()
				Expect(txErr).NotTo(HaveOccurred())
				fkErr := shared.PopulateForeignKeyIDs(foreignKeyValues, columnToValue, tx)
				Expect(fkErr).NotTo(HaveOccurred())
				commitErr := tx.Commit()
				Expect(commitErr).NotTo(HaveOccurred())

				var expectedUrnID int64
				urnErr := db.Get(&expectedUrnID, `SELECT id FROM maker.urns WHERE identifier = $1`, guy)
				Expect(urnErr).NotTo(HaveOccurred())
				actualUrnID := columnToValue[string(constants.UrnFK)].(int64)
				Expect(actualUrnID).To(Equal(expectedUrnID))
			})

			It("looks up FK id and persists in columnToValue for AddressFK", func() {
				foreignKeyValues := shared.ForeignKeyValues{constants.AddressFK: fakes.FakeAddress.Hex()}
				columnToValue := shared.ColumnValues{}

				tx, txErr := db.Beginx()
				Expect(txErr).NotTo(HaveOccurred())
				fkErr := shared.PopulateForeignKeyIDs(foreignKeyValues, columnToValue, tx)
				Expect(fkErr).NotTo(HaveOccurred())
				commitErr := tx.Commit()
				Expect(commitErr).NotTo(HaveOccurred())

				var expectedAddressID int64
				addressErr := db.Get(&expectedAddressID, `SELECT id FROM public.addresses WHERE address = $1`, fakes.FakeAddress.Hex())
				Expect(addressErr).NotTo(HaveOccurred())
				actualAddressID := columnToValue[string(constants.AddressFK)].(int64)
				Expect(actualAddressID).To(Equal(expectedAddressID))
			})
		})
	})

	Describe("GetOrCreateIlk", func() {
		It("returns ID for same ilk with or without hex prefix", func() {
			ilkIDOne, insertErrOne := shared.GetOrCreateIlk(hexIlk, db)
			Expect(insertErrOne).NotTo(HaveOccurred())

			ilkIDTwo, insertErrTwo := shared.GetOrCreateIlk(hexIlk[2:], db)
			Expect(insertErrTwo).NotTo(HaveOccurred())

			Expect(ilkIDOne).To(Equal(ilkIDTwo))
		})
	})

	Describe("GetOrCreateIlkInTransaction", func() {
		It("returns ID for same ilk with or without hex prefix", func() {
			tx, txErr := db.Beginx()
			Expect(txErr).NotTo(HaveOccurred())

			ilkIDOne, insertErrOne := shared.GetOrCreateIlkInTransaction(hexIlk, tx)
			Expect(insertErrOne).NotTo(HaveOccurred())

			ilkIDTwo, insertErrTwo := shared.GetOrCreateIlkInTransaction(hexIlk[2:], tx)
			Expect(insertErrTwo).NotTo(HaveOccurred())

			commitErr := tx.Commit()
			Expect(commitErr).NotTo(HaveOccurred())

			Expect(ilkIDOne).NotTo(BeZero())
			Expect(ilkIDOne).To(Equal(ilkIDTwo))
		})
	})

	Describe("GetOrCreateAddress", func() {
		It("creates an address record", func() {
			_, err := shared.GetOrCreateAddress(fakes.FakeAddress.Hex(), db)
			Expect(err).NotTo(HaveOccurred())

			var address string
			db.Get(&address, `SELECT address from addresses LIMIT 1`)
			Expect(address).To(Equal(fakes.FakeAddress.Hex()))
		})

		It("returns the id for an address that already exists", func() {
			//create the address record
			createAddressId, createErr := shared.GetOrCreateAddress(fakes.FakeAddress.Hex(), db)
			Expect(createErr).NotTo(HaveOccurred())

			//get the address record
			getAddressId, getErr := shared.GetOrCreateAddress(fakes.FakeAddress.Hex(), db)
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

			_, createErr := shared.GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
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
			createAddressId, createErr := shared.GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
			Expect(createErr).NotTo(HaveOccurred())

			//get the address record
			getAddressId, getErr := shared.GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
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

			_, createErr := shared.GetOrCreateAddressInTransaction(fakes.FakeAddress.Hex(), tx)
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
	LogID     string `db:"log_id"`
	Variable1 string
}
