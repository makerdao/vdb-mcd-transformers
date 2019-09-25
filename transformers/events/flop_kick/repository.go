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

package flop_kick

import (
	"fmt"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const InsertFlopKickQuery = `INSERT into maker.flop_kick
	(header_id, bid_id, lot, bid, gal, address_id, log_id)
	VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5, $6, $7)
	ON CONFLICT (header_id, log_id)
	DO UPDATE SET bid_id = $2, lot = $3, bid = $4, gal = $5, address_id = $6;`

type FlopKickRepository struct {
	db *postgres.DB
}

func (repo FlopKickRepository) Create(models []interface{}) error {
	tx, dBaseErr := repo.db.Beginx()
	if dBaseErr != nil {
		return dBaseErr
	}
	for _, model := range models {
		flopKickModel, ok := model.(Model)
		if !ok {
			wrongTypeErr := fmt.Errorf("model of type %T, not %T", model, Model{})
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return shared.FormatRollbackError("flop kick", wrongTypeErr.Error())
			}
			return wrongTypeErr
		}

		addressId, addressErr := shared.GetOrCreateAddressInTransaction(flopKickModel.ContractAddress, tx)
		if addressErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return shared.FormatRollbackError("flop address", addressErr.Error())
			}
			return addressErr
		}

		_, execErr := tx.Exec(InsertFlopKickQuery, flopKickModel.HeaderID, flopKickModel.BidId, flopKickModel.Lot, flopKickModel.Bid,
			flopKickModel.Gal, addressId, flopKickModel.LogID)
		if execErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return shared.FormatRollbackError("flop kick", execErr.Error())
			}
			return execErr
		}

		_, logErr := tx.Exec(`UPDATE public.header_sync_logs SET transformed = true WHERE id = $1`, flopKickModel.LogID)
		if logErr != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return shared.FormatRollbackError("flop kick", logErr.Error())
			}
			return logErr
		}
	}
	return tx.Commit()
}

func (repo *FlopKickRepository) SetDB(db *postgres.DB) {
	repo.db = db
}
