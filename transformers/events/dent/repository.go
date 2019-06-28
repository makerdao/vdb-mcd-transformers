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

package dent

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type DentRepository struct {
	db *postgres.DB
}

func (repo DentRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repo.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}

	for _, model := range models {
		dent, ok := model.(DentModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, DentModel{})
		}

		_, execErr := tx.Exec(
			`INSERT into maker.dent (header_id, bid_id, lot, bid, contract_address, log_idx, tx_idx, raw_log)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET bid_Id = $2, lot = $3, bid = $4, contract_address = $5, raw_log = $8;`,
			headerID, dent.BidId, dent.Lot, dent.Bid, dent.ContractAddress, dent.LogIndex, dent.TransactionIndex, dent.Raw,
		)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				logrus.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	err := repository.MarkHeaderCheckedInTransaction(headerID, tx, constants.DentChecked)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			logrus.Error("failed to rollback ", rollbackErr)
		}
		return err
	}
	return tx.Commit()
}

func (repo DentRepository) MarkHeaderChecked(headerId int64) error {
	return repository.MarkHeaderChecked(headerId, repo.db, constants.DentChecked)
}

func (repo *DentRepository) SetDB(db *postgres.DB) {
	repo.db = db
}
