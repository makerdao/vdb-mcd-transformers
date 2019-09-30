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

func (repository *SpotPokeRepository) Create(models []shared.InsertionModel) error {
	return shared.Create(models, repository.db)
}

func (repository *SpotPokeRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
