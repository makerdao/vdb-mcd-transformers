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
	"github.com/sirupsen/logrus"

	repo "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

const InsertNewCdpQuery = `INSERT INTO maker.new_cdp
	(header_id, usr, own, cdp, tx_idx, log_idx, raw_log)
	VALUES($1, $2, $3, $4::NUMERIC, $5, $6, $7)
	ON CONFLICT (header_id, tx_idx, log_idx)
	DO UPDATE SET usr = $2, own = $3, cdp = $4;`

type NewCdpRepository struct {
	db *postgres.DB
}

func (repository *NewCdpRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository NewCdpRepository) Create(headerID int64, models []interface{}) error {
	tx, dbErr := repository.db.Beginx()
	if dbErr != nil {
		return dbErr
	}
	for _, model := range models {
		newCdpModel, ok := model.(NewCdpModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, NewCdpModel{})
		}

		_, execErr := tx.Exec(InsertNewCdpQuery, headerID, newCdpModel.Usr, newCdpModel.Own, newCdpModel.Cdp,
			newCdpModel.TransactionIndex, newCdpModel.LogIndex, newCdpModel.Raw)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := repo.MarkHeaderCheckedInTransaction(headerID, tx, constants.NewCdpLabel)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logrus.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}

	return tx.Commit()
}

func (repository *NewCdpRepository) MarkHeaderChecked(headerID int64) error {
	return repo.MarkHeaderChecked(headerID, repository.db, constants.NewCdpLabel)
}
