//  VulcanizeDB
//  Copyright © 2019 Vulcanize
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package spot_poke

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	shared_repo "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	InsertSpotPokeQuery = `INSERT INTO maker.spot_poke (header_id, ilk_id, value, spot, tx_idx, log_idx, raw_log)
		VALUES($1, $2, $3::NUMERIC, $4::NUMERIC, $5, $6, $7)
		ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET ilk_id = $2, value = $3, spot = $5, raw_log = $7;`
)

type SpotPokeRepository struct {
	db *postgres.DB
}

func (repository *SpotPokeRepository) Create(headerID int64, models []interface{}) error {
	tx, dbErr := repository.db.Beginx()
	if dbErr != nil {
		return dbErr
	}

	for _, model := range models {
		spotPokeModel, ok := model.(SpotPokeModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, SpotPokeModel{})
		}

		ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(spotPokeModel.Ilk, tx)
		if ilkErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return ilkErr
		}

		_, insertErr := tx.Exec(
			InsertSpotPokeQuery,
			headerID, ilkID, spotPokeModel.Value, spotPokeModel.Spot, spotPokeModel.TransactionIndex, spotPokeModel.LogIndex, spotPokeModel.Raw,
		)
		if insertErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return insertErr
		}
	}

	checkHeaderErr := shared_repo.MarkHeaderChecked(headerID, repository.db, constants.SpotPokeLabel)
	if checkHeaderErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return checkHeaderErr
	}

	return tx.Commit()
}

func (repository SpotPokeRepository) MarkHeaderChecked(headerID int64) error {
	return shared_repo.MarkHeaderChecked(headerID, repository.db, constants.BiteLabel)
}

func (repository *SpotPokeRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
