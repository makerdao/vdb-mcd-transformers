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

package vat_heal

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	repo "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type VatHealRepository struct {
	db *postgres.DB
}

func (repository *VatHealRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository VatHealRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}

	for _, model := range models {
		vatHeal, ok := model.(VatHealModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, VatHealModel{})
		}

		_, execErr := tx.Exec(`INSERT INTO maker.vat_heal (header_id, rad, log_idx, tx_idx, raw_log)
		VALUES($1, $2::NUMERIC, $3, $4, $5)
		ON CONFlICT (header_id, tx_idx, log_idx) DO UPDATE SET rad = $2, raw_log = $5;`,
			headerID, vatHeal.Rad, vatHeal.LogIndex, vatHeal.TransactionIndex, vatHeal.Raw)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}

	checkHeaderErr := repo.MarkHeaderCheckedInTransaction(headerID, tx, constants.VatHealChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repository VatHealRepository) MissingHeaders(startingBlock, endingBlock int64) ([]core.Header, error) {
	return repo.MissingHeaders(startingBlock, endingBlock, repository.db, constants.VatHealChecked)
}

func (repository VatHealRepository) RecheckHeaders(startingBlock, endingBlock int64) ([]core.Header, error) {
	return repo.RecheckHeaders(startingBlock, endingBlock, repository.db, constants.VatHealChecked)
}

func (repository VatHealRepository) MarkHeaderChecked(headerId int64) error {
	return repo.MarkHeaderChecked(headerId, repository.db, constants.VatHealChecked)
}
