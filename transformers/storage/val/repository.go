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

package val

import (
	"fmt"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertValValQuery = `INSERT INTO maker.val_val (diff_id, header_id, val) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type ValStorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository *ValStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *ValStorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Val:
		return repository.insertVal(diffID, headerID, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized val contract storage name: %s", metadata.Name))
	}
}

func (repository *ValStorageRepository) insertVal(diffID, headerID int64, val string) error {
	_, writeErr := repository.db.Exec(insertValValQuery, diffID, headerID, val)
	return writeErr
}
