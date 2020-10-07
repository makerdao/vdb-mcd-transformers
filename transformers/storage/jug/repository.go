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

package jug

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	InsertJugIlkRhoQuery  = `INSERT INTO maker.jug_ilk_rho (diff_id, header_id, ilk_id, rho) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertJugIlkDutyQuery = `INSERT INTO maker.jug_ilk_duty (diff_id, header_id, ilk_id, duty) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertJugVatQuery     = `INSERT INTO maker.jug_vat (diff_id, header_id, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertJugVowQuery     = `INSERT INTO maker.jug_vow (diff_id, header_id, vow) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertJugBaseQuery    = `INSERT INTO maker.jug_base (diff_id, header_id, base) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repository.ContractAddress, value.(string), repository.db)
	case IlkRho:
		return repository.insertIlkRho(diffID, headerID, metadata, value.(string))
	case IlkDuty:
		return repository.insertIlkDuty(diffID, headerID, metadata, value.(string))
	case Vat:
		return repository.insertJugVat(diffID, headerID, value.(string))
	case Vow:
		return repository.insertJugVow(diffID, headerID, value.(string))
	case Base:
		return repository.insertJugBase(diffID, headerID, value.(string))

	default:
		panic(fmt.Sprintf("unrecognized jug contract storage name: %s", metadata.Name))
	}
}

func (repository *StorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository StorageRepository) insertIlkRho(diffID, headerID int64, metadata types.ValueMetadata, rho string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk rho: %w", err)
	}

	insertErr := shared.InsertFieldWithIlk(diffID, headerID, ilk, IlkRho, InsertJugIlkRhoQuery, rho, repository.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s rho %s from diff ID %d: %w", ilk, rho, diffID, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertIlkDuty(diffID, headerID int64, metadata types.ValueMetadata, duty string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk duty: %w", err)
	}

	insertErr := shared.InsertFieldWithIlk(diffID, headerID, ilk, IlkDuty, InsertJugIlkDutyQuery, duty, repository.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s duty %s from diff ID %d: %w", ilk, duty, diffID, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertJugVat(diffID, headerID int64, vat string) error {
	_, err := repository.db.Exec(insertJugVatQuery, diffID, headerID, vat)
	if err != nil {
		return fmt.Errorf("error inserting jug vat %s from diff ID %d: %w", vat, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertJugVow(diffID, headerID int64, vow string) error {
	_, err := repository.db.Exec(insertJugVowQuery, diffID, headerID, vow)
	if err != nil {
		return fmt.Errorf("error inserting jug vow %s from diff ID %d: %w", vow, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertJugBase(diffID, headerID int64, base string) error {
	_, err := repository.db.Exec(insertJugBaseQuery, diffID, headerID, base)
	if err != nil {
		return fmt.Errorf("error inserting jug base %s from diff ID %d: %w", base, diffID, err)
	}
	return nil
}

func getIlk(keys map[types.Key]string) (string, error) {
	ilk, ok := keys[constants.Ilk]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.Ilk}
	}
	return ilk, nil
}
