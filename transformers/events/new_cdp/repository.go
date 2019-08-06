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

package new_cdp

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const InsertNewCdpQuery = `INSERT INTO maker.new_cdp
	(header_id, usr, own, cdp, log_id)
	VALUES($1, $2, $3, $4::NUMERIC, $5)
	ON CONFLICT (header_id, log_id)
	DO UPDATE SET usr = $2, own = $3, cdp = $4;`

type NewCdpRepository struct {
	db *postgres.DB
}

func (repository *NewCdpRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository NewCdpRepository) Create(models []interface{}) error {
	tx, dbErr := repository.db.Beginx()
	if dbErr != nil {
		return dbErr
	}
	for _, model := range models {
		newCdpModel, ok := model.(NewCdpModel)
		if !ok {
			return rollbackErr(tx, fmt.Errorf("model of type %T, not %T", model, NewCdpModel{}))
		}

		_, execErr := tx.Exec(InsertNewCdpQuery, newCdpModel.HeaderID, newCdpModel.Usr, newCdpModel.Own, newCdpModel.Cdp,
			newCdpModel.LogID)
		if execErr != nil {
			return rollbackErr(tx, execErr)
		}
		_, logErr := tx.Exec(`UPDATE public.header_sync_logs SET transformed = true WHERE id = $1`, newCdpModel.LogID)
		if logErr != nil {
			return rollbackErr(tx, logErr)
		}
	}

	return tx.Commit()
}

func rollbackErr(tx *sqlx.Tx, err error) error {
	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		logrus.Error("failed to rollback ", rollbackErr)
	}
	return err
}
