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

	"github.com/sirupsen/logrus"
	repo "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type BiteRepository struct {
	db *postgres.DB
}

func (repository *BiteRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository BiteRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		biteModel, ok := model.(BiteModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, BiteModel{})
		}

		urnID, urnErr := shared.GetOrCreateUrnInTransaction(biteModel.Urn, biteModel.Ilk, tx)
		if urnErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback", rollbackErr)
			}
			return urnErr
		}

		_, execErr := tx.Exec(
			`INSERT into maker.bite (header_id, urn_id, ink, art, tab, flip, bite_identifier, log_idx, tx_idx, raw_log)
					VALUES($1, $2, $3::NUMERIC, $4::NUMERIC, $5::NUMERIC, $6, $7::NUMERIC, $8, $9, $10)
					ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET urn_id = $2, ink = $3, art = $4, tab = $5, flip = $6, bite_identifier = $7, raw_log = $10`,
			headerID, urnID, biteModel.Ink, biteModel.Art, biteModel.Tab, biteModel.Flip, biteModel.Id, biteModel.LogIndex, biteModel.TransactionIndex, biteModel.Raw,
		)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := repo.MarkHeaderCheckedInTransaction(headerID, tx, constants.BiteLabel)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logrus.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}

	return tx.Commit()
}

func (repository BiteRepository) MarkHeaderChecked(headerID int64) error {
	return repo.MarkHeaderChecked(headerID, repository.db, constants.BiteLabel)
}
