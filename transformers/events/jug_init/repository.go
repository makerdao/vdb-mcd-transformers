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

package jug_init

import (
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
)

type JugInitRepository struct {
	db *postgres.DB
}

func (repository JugInitRepository) Create(models []shared.InsertionModel) error {
	return shared.Create(models, repository.db)
}

func (repository *JugInitRepository) SetDB(db *postgres.DB) {
	repository.db = db
}
