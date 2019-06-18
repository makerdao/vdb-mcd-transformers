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

package vat_suck

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	repo "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type VatSuckRepository struct {
	db *postgres.DB
}

func (repository *VatSuckRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository VatSuckRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}

	for _, model := range models {
		vatSuck, ok := model.(VatSuckModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, VatSuckModel{})
		}

		_, execErr := tx.Exec(`INSERT INTO maker.vat_suck (header_id, u, v, rad, log_idx, tx_idx, raw_log)
		VALUES($1, $2, $3, $4::NUMERIC, $5, $6, $7)
		ON CONFlICT (header_id, tx_idx, log_idx) DO UPDATE SET u = $2, v = $3, rad = $4, raw_log = $7;`,
			headerID, vatSuck.U, vatSuck.V, vatSuck.Rad, vatSuck.LogIndex, vatSuck.TransactionIndex, vatSuck.Raw)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := repo.MarkHeaderCheckedInTransaction(headerID, tx, constants.VatSuckChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repository VatSuckRepository) MarkHeaderChecked(headerId int64) error {
	return repo.MarkHeaderChecked(headerId, repository.db, constants.VatSuckChecked)
}
