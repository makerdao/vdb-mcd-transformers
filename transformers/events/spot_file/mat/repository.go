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

package mat

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const InsertSpotFileMatQuery = `INSERT INTO maker.spot_file_mat (header_id, ilk_id, what, data, log_idx, tx_idx, raw_log)
	VALUES($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (header_id, tx_idx, log_idx)
	DO UPDATE SET ilk_id = $2, what = $3, data = $4, raw_log = $7;`

type SpotFileMatRepository struct {
	db *postgres.DB
}

func (repo SpotFileMatRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repo.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		spotFileMatModel, ok := model.(SpotFileMatModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, SpotFileMatModel{})
		}

		ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(spotFileMatModel.Ilk, tx)
		if ilkErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return ilkErr
		}

		_, execErr := tx.Exec(InsertSpotFileMatQuery, headerID, ilkID, spotFileMatModel.What, spotFileMatModel.Data,
			spotFileMatModel.LogIndex, spotFileMatModel.TransactionIndex, spotFileMatModel.Raw)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}
	checkHeaderErr := repository.MarkHeaderCheckedInTransaction(headerID, tx, constants.SpotFileMatChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logrus.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repo SpotFileMatRepository) MarkHeaderChecked(headerID int64) error {
	return repository.MarkHeaderChecked(headerID, repo.db, constants.SpotFileMatChecked)
}

func (repo SpotFileMatRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repository.MissingHeaders(startingBlockNumber, endingBlockNumber, repo.db, constants.SpotFileMatChecked)
}

func (repo SpotFileMatRepository) RecheckHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repository.RecheckHeaders(startingBlockNumber, endingBlockNumber, repo.db, constants.SpotFileMatChecked)
}

func (repo *SpotFileMatRepository) SetDB(db *postgres.DB) {
	repo.db = db
}
