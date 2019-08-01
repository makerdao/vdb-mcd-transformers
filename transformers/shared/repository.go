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
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"strconv"
	"strings"
)

const (
	getIlkIdQuery  = `SELECT id FROM maker.ilks WHERE ilk = $1`
	getUrnIdQuery  = `SELECT id FROM maker.urns WHERE identifier = $1 AND ilk_id = $2`
	InsertIlkQuery = `INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2) RETURNING id`
	InsertUrnQuery = `INSERT INTO maker.urns (identifier, ilk_id) VALUES ($1, $2) RETURNING id`
)

type SharedRepository interface {
	Create(headerID int64, models []InsertionModel) error
	MarkHeaderChecked(headerID int64) error
	SetDB(db *postgres.DB)
}

type ForeignKeyValues map[constants.ForeignKeyField]string
type ColumnValues map[string]interface{}
type InsertionModel struct {
	TableName      string   // For MarkHeaderChecked, insert query
	OrderedColumns []string // Defines the fields to insert, and in which order the table expects them
	// ColumnValues needs to be typed interface{}, since `raw_log` is a slice of bytes and not a string
	ColumnValues     ColumnValues     // Associated values for columns, headerID, FKs and event metadata populated automatically
	ForeignKeyValues ForeignKeyValues // FK name and value to get/create ID for
}

var modelToQuery = map[string]string{}

func getMemoizedQuery(model InsertionModel) string {
	// The table name uniquely determines the insertion query, use that for memoization
	query, queryMemoized := modelToQuery[model.TableName]
	if !queryMemoized {
		query = generateInsertionQuery(model)
		modelToQuery[model.TableName] = query
	}
	return query
}

// Creates an insertion query from an insertion model. This is called through getMemoizedQuery, so the query is not
// generated on each call to Create.
// Note: With extraction of event metadata, one would not have to supply header_id, tx_idx, etc in InsertionModel.OrderedColumns?
// Note: I have a feeling we can actually do away with the OrderedColumns field, but the tricky part is that some fields
//       needed aren't present in the map in the beginning
func generateInsertionQuery(model InsertionModel) string {
	var valuePlaceholders []string
	var updateOnConflict []string
	for i := 0; i < len(model.OrderedColumns); i++ {
		valuePlaceholder := fmt.Sprintf("$%d", 1+i)
		valuePlaceholders = append(valuePlaceholders, valuePlaceholder)
		updateOnConflict = append(updateOnConflict,
			fmt.Sprintf("%s = %s", model.OrderedColumns[i], valuePlaceholder))
	}

	baseQuery := `INSERT INTO maker.%v (%v) VALUES(%v)
		ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET %v;`
	return fmt.Sprintf(baseQuery,
		model.TableName,
		strings.Join(model.OrderedColumns, ", "),
		strings.Join(valuePlaceholders, ", "),
		strings.Join(updateOnConflict, ", "))
}

/* Given an instance of InsertionModel, example below, generates an insertion query and fills in
foreign keys automatically after getting from the DB. These "special fields" are populated in the
columnToValue mapping, and are treated like any other in the insertion.

testModel = shared.InsertionModel{
			TableName:      "testEvent",
			OrderedColumns: []string{"header_id", "log_idx", "tx_idx", "raw_log", constants.IlkFK, constants.UrnFK, "variable1"},
			ColumnValues: ColumnValues{
				"log_idx":   "1",
				"tx_idx":    "2",
				"raw_log":   fakeLog,
				"variable1": "value1",
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.IlkFK: test_helpers.FakeIlk.Hex,
				constants.UrnFK: "0x12345",
			},
		}
*/
func Create(headerID int64, models []InsertionModel, db *postgres.DB) error {
	if len(models) == 0 {
		return fmt.Errorf("repository got empty model slice")
	}

	// Quick 'n' dirty solution to these not being declared in config.
	// a) Couldn't we somewhere create the table and add a checked column inside the plugin, instead of a migration?
	checkedHeaderColumn := models[0].TableName

	tx, dbErr := db.Beginx()
	if dbErr != nil {
		return dbErr
	}

	for _, model := range models {
		fkErr := populateForeignKeyIDs(model.ForeignKeyValues, model.ColumnValues, tx)
		if fkErr != nil {
			return fmt.Errorf("error gettings FK ids: %s", fkErr.Error())
		}

		// Save headerId in mapping for insertion query
		model.ColumnValues["header_id"] = strconv.FormatInt(headerID, 10)

		// Maps can't be iterated over in a reliable manner, so we rely on OrderedColumns to define the order to insert
		// tx.Exec is variadically typed in the args, so if we wrap in []interface{} we can apply them all automatically
		var args []interface{}
		for _, col := range model.OrderedColumns {
			args = append(args, model.ColumnValues[col])
		}

		insertionQuery := getMemoizedQuery(model)
		_, execErr := tx.Exec(insertionQuery, args...) //couldn't do this trick with args :: []string

		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := repository.MarkHeaderCheckedInTransaction(headerID, tx, checkedHeaderColumn)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

// Gets or creates the FK for the key/values supplied, and inserts the resulting ID into the columnToValue mapping
func populateForeignKeyIDs(fkToValue ForeignKeyValues, columnToValue ColumnValues, tx *sqlx.Tx) error {
	var dbErr error
	var fkID int
	for fk, value := range fkToValue {
		switch fk {
		case constants.IlkFK:
			fkID, dbErr = GetOrCreateIlkInTransaction(value, tx)
		case constants.UrnFK:
			fkID, dbErr = GetOrCreateUrnInTransaction(value, fkToValue[constants.IlkFK], tx)
		default:
			return fmt.Errorf("repository got unrecognised FK: %s", fk)
		}

		if dbErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("couldn't get or create FK (%s, %s): %s", fk, value, dbErr.Error())
		} else {
			columnName := string(fk)
			columnToValue[columnName] = fkID
		}
	}

	return nil
}

func GetOrCreateIlk(ilk string, db *postgres.DB) (int, error) {
	var ilkID int
	err := db.Get(&ilkID, getIlkIdQuery, ilk)
	if err == sql.ErrNoRows {
		ilkIdentifier := DecodeHexToText(ilk)
		insertErr := db.QueryRow(InsertIlkQuery, ilk, ilkIdentifier).Scan(&ilkID)
		return ilkID, insertErr
	}
	return ilkID, err
}

func GetOrCreateIlkInTransaction(ilk string, tx *sqlx.Tx) (int, error) {
	var ilkID int
	err := tx.Get(&ilkID, getIlkIdQuery, ilk)
	if err == sql.ErrNoRows {
		ilkIdentifier := DecodeHexToText(ilk)
		insertErr := tx.QueryRow(InsertIlkQuery, ilk, ilkIdentifier).Scan(&ilkID)
		return ilkID, insertErr
	}
	return ilkID, err
}

func GetOrCreateUrn(guy string, hexIlk string, db *postgres.DB) (urnID int, err error) {
	ilkID, ilkErr := GetOrCreateIlk(hexIlk, db)
	if ilkErr != nil {
		return 0, fmt.Errorf("error getting ilkID for urn: %s", ilkErr.Error())
	}

	err = db.Get(&urnID, getUrnIdQuery, guy, ilkID)
	if err != nil {
		if err == sql.ErrNoRows {
			insertErr := db.QueryRow(InsertUrnQuery, guy, ilkID).Scan(&urnID)
			return urnID, insertErr
		}
	}

	return urnID, err
}

func GetOrCreateUrnInTransaction(guy string, hexIlk string, tx *sqlx.Tx) (urnID int, err error) {
	ilkID, ilkErr := GetOrCreateIlkInTransaction(hexIlk, tx)
	if ilkErr != nil {
		return 0, fmt.Errorf("error getting ilkID for urn")
	}
	err = tx.Get(&urnID, getUrnIdQuery, guy, ilkID)

	if err != nil {
		if err == sql.ErrNoRows {
			insertErr := tx.QueryRow(InsertUrnQuery, guy, ilkID).Scan(&urnID)
			return urnID, insertErr
		}
	}

	return urnID, err
}
