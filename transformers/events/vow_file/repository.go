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

package vow_file

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type VowFileRepository struct {
	db *postgres.DB
}

func (repo VowFileRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repo.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}

	for _, model := range models {
		vowFileModel, ok := model.(VowFileModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, VowFileModel{})
		}

		_, execErr := tx.Exec(
			`INSERT into maker.vow_file (header_id, what, data, log_idx, tx_idx, raw_log)
        	VALUES($1, $2, $3::NUMERIC, $4, $5, $6)`,
			headerID, vowFileModel.What, vowFileModel.Data, vowFileModel.LogIndex, vowFileModel.TransactionIndex, vowFileModel.Raw,
		)

		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := repository.MarkHeaderCheckedInTransaction(headerID, tx, constants.VowFileChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logrus.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}

	return tx.Commit()
}

func (repo VowFileRepository) MarkHeaderChecked(headerID int64) error {
	return repository.MarkHeaderChecked(headerID, repo.db, constants.VowFileChecked)
}

func (repo VowFileRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repository.MissingHeaders(startingBlockNumber, endingBlockNumber, repo.db, constants.VowFileChecked)
}

func (repo VowFileRepository) RecheckHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repository.RecheckHeaders(startingBlockNumber, endingBlockNumber, repo.db, constants.VowFileChecked)
}

func (repo *VowFileRepository) SetDB(db *postgres.DB) {
	repo.db = db
}
