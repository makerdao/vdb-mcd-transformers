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
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	InsertJugIlkRhoQuery  = `INSERT INTO maker.jug_ilk_rho (diff_id, header_id, ilk_id, rho) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertJugIlkDutyQuery = `INSERT INTO maker.jug_ilk_duty (diff_id, header_id, ilk_id, duty) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertJugVatQuery     = `INSERT INTO maker.jug_vat (diff_id, header_id, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertJugVowQuery     = `INSERT INTO maker.jug_vow (diff_id, header_id, vow) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertJugBaseQuery    = `INSERT INTO maker.jug_base (diff_id, header_id, base) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type JugStorageRepository struct {
	db *postgres.DB
}

func (repository JugStorageRepository) Create(diffID, headerID int64, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
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

func (repository *JugStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository JugStorageRepository) insertIlkRho(diffID, headerID int64, metadata utils.StorageValueMetadata, rho string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertFieldWithIlk(diffID, headerID, ilk, IlkRho, InsertJugIlkRhoQuery, rho)
}

func (repository JugStorageRepository) insertIlkDuty(diffID, headerID int64, metadata utils.StorageValueMetadata, duty string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(diffID, headerID, ilk, IlkDuty, InsertJugIlkDutyQuery, duty)
}

func (repository JugStorageRepository) insertJugVat(diffID, headerID int64, vat string) error {
	_, err := repository.db.Exec(insertJugVatQuery, diffID, headerID, vat)
	return err
}

func (repository JugStorageRepository) insertJugVow(diffID, headerID int64, vow string) error {
	_, err := repository.db.Exec(insertJugVowQuery, diffID, headerID, vow)
	return err
}

func (repository JugStorageRepository) insertJugBase(diffID, headerID int64, repo string) error {
	_, err := repository.db.Exec(insertJugBaseQuery, diffID, headerID, repo)
	return err
}

func (repository *JugStorageRepository) insertFieldWithIlk(diffID, headerID int64, ilk, variableName, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("ilk", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(query, diffID, headerID, ilkID, value)

	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError(variableName, writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func getIlk(keys map[utils.Key]string) (string, error) {
	ilk, ok := keys[constants.Ilk]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.Ilk}
	}
	return ilk, nil
}
