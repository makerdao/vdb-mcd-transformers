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

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

func GetOrCreateIlk(ilk string, db *postgres.DB) (int, error) {
	var ilkID int
	err := db.Get(&ilkID, `SELECT id FROM maker.ilks WHERE ilk = $1`, ilk)
	if err != nil {
		if err == sql.ErrNoRows {
			insertErr := db.QueryRow(`INSERT INTO maker.ilks (ilk) VALUES ($1) RETURNING id`, ilk).Scan(&ilkID)
			return ilkID, insertErr
		}
	}
	return ilkID, err
}

func GetOrCreateIlkInTransaction(ilk string, tx *sqlx.Tx) (int, error) {
	var ilkID int
	err := tx.Get(&ilkID, `SELECT id FROM maker.ilks WHERE ilk = $1`, ilk)
	if err != nil {
		if err == sql.ErrNoRows {
			insertErr := tx.QueryRow(`INSERT INTO maker.ilks (ilk) VALUES ($1) RETURNING id`, ilk).Scan(&ilkID)
			return ilkID, insertErr
		}
	}
	return ilkID, err
}

func GetTicInTx(headerID int64, tx *sqlx.Tx) (int64, error) {
	var blockTimestamp int64
	err := tx.Get(&blockTimestamp, `SELECT block_timestamp FROM public.headers WHERE id = $1;`, headerID)
	if err != nil {
		return 0, err
	}

	tic := blockTimestamp + constants.TTL
	return tic, nil
}
