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
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

var InsertFlipKickQuery = `INSERT into maker.flip_kick (header_id, bid_id, lot, bid, tab, usr, gal, address_id, log_id)
				VALUES($1, $2::NUMERIC, $3::NUMERIC, $4::NUMERIC, $5::NUMERIC, $6, $7, $8, $9)
				ON CONFLICT (header_id, log_id) DO UPDATE SET bid_id = $2, lot = $3, bid = $4, tab = $5, usr = $6, gal = $7, address_id = $8;`

type FlipKickRepository struct {
	db *postgres.DB
}

func (repository FlipKickRepository) Create(models []shared.InsertionModel) error {
	return shared.Create(models, repository.db)
}

func (repository *FlipKickRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
