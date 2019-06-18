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

package pip

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

const InsertSpotFilePipQuery = `INSERT INTO maker.spot_file_pip (header_id, ilk_id, pip, log_idx, tx_idx, raw_log)
	VALUES($1, $2, $3, $4, $5, $6) ON CONFLICT (header_id, tx_idx, log_idx)
	DO UPDATE SET ilk_id = $2, pip = $3, raw_log = $6;`

type SpotFilePipRepository struct {
	db *postgres.DB
}

func (repo SpotFilePipRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repo.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		spotFilePipModel, ok := model.(SpotFilePipModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, SpotFilePipModel{})
		}

		ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(spotFilePipModel.Ilk, tx)
		if ilkErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return ilkErr
		}

		_, execErr := tx.Exec(InsertSpotFilePipQuery, headerID, ilkID, spotFilePipModel.Pip, spotFilePipModel.LogIndex,
			spotFilePipModel.TransactionIndex, spotFilePipModel.Raw)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}
	checkHeaderErr := repository.MarkHeaderCheckedInTransaction(headerID, tx, constants.SpotFilePipChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logrus.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repo SpotFilePipRepository) MarkHeaderChecked(headerID int64) error {
	return repository.MarkHeaderChecked(headerID, repo.db, constants.SpotFilePipChecked)
}

func (repo *SpotFilePipRepository) SetDB(db *postgres.DB) {
	repo.db = db
}
