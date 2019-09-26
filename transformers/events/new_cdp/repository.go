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

package new_cdp

import (
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const InsertNewCdpQuery = `INSERT INTO maker.new_cdp
	(header_id, usr, own, cdp, log_id)
	VALUES($1, $2, $3, $4::NUMERIC, $5)
	ON CONFLICT (header_id, log_id)
	DO UPDATE SET usr = $2, own = $3, cdp = $4;`

type NewCdpRepository struct {
	db *postgres.DB
}

func (repository *NewCdpRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository NewCdpRepository) Create(models []shared.InsertionModel) error {
	return shared.Create(models, repository.db)
}
