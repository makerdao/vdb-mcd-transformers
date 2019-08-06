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

package flip_kick

import (
	"fmt"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

var InsertFlipKickQuery = `INSERT into maker.flip_kick (header_id, bid_id, lot, bid, tab, usr, gal, address_id, log_id)
				VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5::NUMERIC, $6, $7, $8, $9)
				ON CONFLICT (header_id, log_id) DO UPDATE SET bid_id = $2, lot = $3, bid = $4, tab = $5, usr = $6, gal = $7, address_id = $8;`

type FlipKickRepository struct {
	db *postgres.DB
}

func (repository FlipKickRepository) Create(models []interface{}) error {
	tx, dBaseErr := repository.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		flipKickModel, ok := model.(FlipKickModel)
		if !ok {
			wrongTypeErr := fmt.Errorf("model of type %T, not %T", model, FlipKickModel{})
			return shared.FormatRollbackError("flip kick", wrongTypeErr.Error())
		}

		addressId, addressErr := shared.GetOrCreateAddressInTransaction(flipKickModel.ContractAddress, tx)
		if addressErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return shared.FormatRollbackError("flip address", addressErr.Error())
			}
			return addressErr
		}

		_, execErr := tx.Exec(InsertFlipKickQuery, flipKickModel.HeaderID, flipKickModel.BidId, flipKickModel.Lot, flipKickModel.Bid,
			flipKickModel.Tab, flipKickModel.Usr, flipKickModel.Gal, addressId, flipKickModel.LogID)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return shared.FormatRollbackError("flip kick", execErr.Error())
			}
			return execErr
		}

		_, logErr := tx.Exec(`UPDATE public.header_sync_logs SET transformed = true WHERE id = $1`, flipKickModel.LogID)
		if logErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return shared.FormatRollbackError("flip kick", logErr.Error())
			}
			return logErr
		}
	}
	return tx.Commit()
}

func (repository *FlipKickRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
