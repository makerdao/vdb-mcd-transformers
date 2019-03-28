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

package vat_frob

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	repo "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type VatFrobRepository struct {
	db *postgres.DB
}

func (repository VatFrobRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		vatFrobModel, ok := model.(VatFrobModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, VatFrobModel{})
		}

		ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(vatFrobModel.Ilk, tx)
		if ilkErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return ilkErr
		}

		urnID, urnErr := shared.GetOrCreateUrnInTransaction(vatFrobModel.Urn, ilkID, tx)
		if urnErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback", rollbackErr)
			}
			return urnErr
		}

		_, execErr := tx.Exec(
			`INSERT INTO maker.vat_frob (header_id, urn_id, v, w, dink, dart, raw_log, log_idx, tx_idx)
			VALUES($1, $2::NUMERIC, $3, $4, $5::NUMERIC, $6::NUMERIC, $7, $8, $9)
			ON CONFLICT (header_id, tx_idx, log_idx)
			DO UPDATE SET urn_id = $2, v = $3, w = $4, dink = $5, dart = $6, raw_log = $7;`,
			headerID, urnID, vatFrobModel.V, vatFrobModel.W, vatFrobModel.Dink, vatFrobModel.Dart, vatFrobModel.Raw,
			vatFrobModel.LogIndex, vatFrobModel.TransactionIndex)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return execErr
		}
	}
	checkHeaderErr := repo.MarkHeaderCheckedInTransaction(headerID, tx, constants.VatFrobChecked)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}
	return tx.Commit()
}

func (repository VatFrobRepository) MarkHeaderChecked(headerID int64) error {
	return repo.MarkHeaderChecked(headerID, repository.db, constants.VatFrobChecked)
}

func (repository VatFrobRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repo.MissingHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.VatFrobChecked)
}

func (repository VatFrobRepository) RecheckHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repo.RecheckHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.VatFrobChecked)
}

func (repository *VatFrobRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
