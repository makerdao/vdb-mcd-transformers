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
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	InsertSpotPokeQuery = `INSERT INTO maker.spot_poke (header_id, ilk_id, value, spot, log_id)
		VALUES($1, $2, $3::NUMERIC, $4::NUMERIC, $5)
		ON CONFLICT (header_id, log_id) DO UPDATE SET ilk_id = $2, value = $3, spot = $5;`
)

type SpotPokeRepository struct {
	db *postgres.DB
}

func (repository *SpotPokeRepository) Create(models []interface{}) error {
	tx, dbErr := repository.db.Beginx()
	if dbErr != nil {
		return dbErr
	}

	for _, model := range models {
		spotPokeModel, ok := model.(SpotPokeModel)
		if !ok {
			wrongTypeErr := fmt.Errorf("model of type %T, not %T", model, SpotPokeModel{})
			return rollbackErr(tx, wrongTypeErr)
		}

		ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(spotPokeModel.Ilk, tx)
		if ilkErr != nil {
			return rollbackErr(tx, ilkErr)
		}

		_, insertErr := tx.Exec(
			InsertSpotPokeQuery,
			spotPokeModel.HeaderID, ilkID, spotPokeModel.Value, spotPokeModel.Spot, spotPokeModel.LogID,
		)
		if insertErr != nil {
			return rollbackErr(tx, insertErr)
		}

		_, logErr := tx.Exec(`UPDATE public.header_sync_logs SET transformed = true WHERE id = $1`, spotPokeModel.LogID)
		if logErr != nil {
			return rollbackErr(tx, logErr)
		}
	}

	return tx.Commit()
}

func (repository *SpotPokeRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func rollbackErr(tx *sqlx.Tx, err error) error {
	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		logrus.Error("failed to rollback ", rollbackErr)
	}
	return err
}
