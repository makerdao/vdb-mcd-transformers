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

package bite

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type BiteRepository struct {
	db *postgres.DB
}

func (repository *BiteRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository BiteRepository) Create(models []interface{}) error {
	tx, dBaseErr := repository.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		biteModel, ok := model.(BiteModel)
		if !ok {
			wrongTypeErr := fmt.Errorf("model of type %T, not %T", model, BiteModel{})
			return rollbackErr(tx, wrongTypeErr)
		}

		urnID, urnErr := shared.GetOrCreateUrnInTransaction(biteModel.Urn, biteModel.Ilk, tx)
		if urnErr != nil {
			return rollbackErr(tx, urnErr)
		}

		_, execErr := tx.Exec(
			`INSERT into maker.bite (header_id, urn_id, ink, art, tab, flip, bite_identifier, log_id)
					VALUES($1, $2, $3::NUMERIC, $4::NUMERIC, $5::NUMERIC, $6, $7::NUMERIC, $8)
					ON CONFLICT (header_id, log_id) DO UPDATE SET urn_id = $2, ink = $3, art = $4, tab = $5, flip = $6, bite_identifier = $7`,
			biteModel.HeaderID, urnID, biteModel.Ink, biteModel.Art, biteModel.Tab, biteModel.Flip, biteModel.Id, biteModel.LogID,
		)
		if execErr != nil {
			return rollbackErr(tx, execErr)
		}

		_, logErr := tx.Exec(`UPDATE public.header_sync_logs SET transformed = true WHERE id = $1`, biteModel.LogID)
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
