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

package vow_file

import (
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type VowFileRepository struct {
	db *postgres.DB
}

func (repo VowFileRepository) Create(headerID int64, models []shared.InsertionModel) error {
	return shared.Create(headerID, models, repo.db)
}

func (repo VowFileRepository) MarkHeaderChecked(headerID int64) error {
	return repository.MarkHeaderChecked(headerID, repo.db, constants.VowFileChecked)
}

func (repo *VowFileRepository) SetDB(db *postgres.DB) {
	repo.db = db
}
