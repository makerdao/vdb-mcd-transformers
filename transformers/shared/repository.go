// VulcanizeDB
// Copyright © 2018 Vulcanize

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

	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	getIlkIdQuery  = `SELECT id FROM maker.ilks WHERE ilk = $1`
	getUrnIdQuery  = `SELECT id FROM maker.urns WHERE identifier = $1 AND ilk_id = $2`
	InsertIlkQuery = `INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2) RETURNING id`
	InsertUrnQuery = `INSERT INTO maker.urns (identifier, ilk_id) VALUES ($1, $2) RETURNING id`
)

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

func GetOrCreateUrn(guy string, ilkID int, db *postgres.DB) (int, error) {
	var urnID int
	err := db.Get(&urnID, getUrnIdQuery, guy, ilkID)
	if err != nil {
		if err == sql.ErrNoRows {
			insertErr := db.QueryRow(InsertUrnQuery, guy, ilkID).Scan(&urnID)
			return urnID, insertErr
		}
	}

	return urnID, err
}

func GetOrCreateUrnInTransaction(guy string, ilkID int, tx *sqlx.Tx) (int, error) {
	var urnID int
	err := tx.Get(&urnID, getUrnIdQuery, guy, ilkID)

	if err != nil {
		if err == sql.ErrNoRows {
			insertErr := tx.QueryRow(InsertUrnQuery, guy, ilkID).Scan(&urnID)
			return urnID, insertErr
		}
	}

	return urnID, err
}
