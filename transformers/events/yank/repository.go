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

package yank

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

const InsertYankQuery = `INSERT INTO maker.yank (header_id, bid_id, contract_address, raw_log, log_idx, tx_idx)
		VALUES($1, $2::NUMERIC, $3, $4, $5::NUMERIC, $6::NUMERIC)
		ON CONFLICT (header_id, tx_idx, log_idx)
		DO UPDATE SET bid_id = $2, contract_address = $3, raw_log = $4;`

type YankRepository struct {
	db *postgres.DB
}

func (repo YankRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repo.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		yankModel, ok := model.(YankModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, YankModel{})
		}

		_, execErr := tx.Exec(InsertYankQuery, headerID, yankModel.BidId, yankModel.ContractAddress, yankModel.Raw,
			yankModel.LogIndex, yankModel.TransactionIndex)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}
	checkHeaderErr := repository.MarkHeaderCheckedInTransaction(headerID, tx, constants.YankChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logrus.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repo YankRepository) MarkHeaderChecked(headerID int64) error {
	return repository.MarkHeaderChecked(headerID, repo.db, constants.YankChecked)
}

func (repo *YankRepository) SetDB(db *postgres.DB) {
	repo.db = db
}
