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

package pip_log_value

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	repo "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

const (
	InsertPipLogValueQuery = `INSERT INTO maker.pip_log_value (block_number, header_id, contract_address, val, log_idx, tx_idx, raw_log)
		VALUES ($1, $2, $3, $4::NUMERIC, $5::NUMERIC, $6, $7)
		ON CONFLICT (header_id, contract_address, tx_idx, log_idx)
		DO UPDATE SET block_number = $1, val = $4, raw_log = $7;`
)

type PipLogValueRepository struct {
	db *postgres.DB
}

func (repository PipLogValueRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		pipLogValue, ok := model.(PipLogValueModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, PipLogValueModel{})
		}

		_, err := tx.Exec(InsertPipLogValueQuery,
			pipLogValue.BlockNumber, headerID, pipLogValue.ContractAddress, pipLogValue.Value,
			pipLogValue.LogIndex, pipLogValue.TransactionIndex, pipLogValue.Raw)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	checkHeaderErr := repo.MarkHeaderCheckedInTransaction(headerID, tx, constants.PipLogValueChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repository PipLogValueRepository) MarkHeaderChecked(headerID int64) error {
	return repo.MarkHeaderChecked(headerID, repository.db, constants.PipLogValueChecked)
}

func (repository PipLogValueRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repo.MissingHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.PipLogValueChecked)
}

func (repository PipLogValueRepository) RecheckHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repo.RecheckHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.PipLogValueChecked)
}

func (repository *PipLogValueRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
