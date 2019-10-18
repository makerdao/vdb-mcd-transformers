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
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	InsertJugIlkRhoQuery  = `INSERT INTO maker.jug_ilk_rho (block_number, block_hash, ilk_id, rho) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertJugIlkDutyQuery = `INSERT INTO maker.jug_ilk_duty (block_number, block_hash, ilk_id, duty) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertJugVatQuery     = `INSERT INTO maker.jug_vat (block_number, block_hash, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertJugVowQuery     = `INSERT INTO maker.jug_vow (block_number, block_hash, vow) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertJugBaseQuery    = `INSERT INTO maker.jug_base (block_number, block_hash, base) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type JugStorageRepository struct {
	db *postgres.DB
}

func (repository *JugStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository JugStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case IlkRho:
		return repository.insertIlkRho(blockNumber, blockHash, metadata, value.(string))
	case IlkDuty:
		return repository.insertIlkDuty(blockNumber, blockHash, metadata, value.(string))
	case Vat:
		return repository.insertJugVat(blockNumber, blockHash, value.(string))
	case Vow:
		return repository.insertJugVow(blockNumber, blockHash, value.(string))
	case Base:
		return repository.insertJugBase(blockNumber, blockHash, value.(string))

	default:
		panic(fmt.Sprintf("unrecognized jug contract storage name: %s", metadata.Name))
	}
}

func (repository JugStorageRepository) insertIlkRho(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, rho string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkRho, InsertJugIlkRhoQuery, rho)
}

func (repository JugStorageRepository) insertIlkDuty(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, duty string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkDuty, InsertJugIlkDutyQuery, duty)
}

func (repository JugStorageRepository) insertJugVat(blockNumber int, blockHash string, vat string) error {
	_, err := repository.db.Exec(insertJugVatQuery, blockNumber, blockHash, vat)
	return err
}

func (repository JugStorageRepository) insertJugVow(blockNumber int, blockHash string, vow string) error {
	_, err := repository.db.Exec(insertJugVowQuery, blockNumber, blockHash, vow)
	return err
}

func (repository JugStorageRepository) insertJugBase(blockNumber int, blockHash string, repo string) error {
	_, err := repository.db.Exec(insertJugBaseQuery, blockNumber, blockHash, repo)
	return err
}

func (repository *JugStorageRepository) insertFieldWithIlk(blockNumber int, blockHash, ilk, variableName, query, value string) error {
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
	_, writeErr := tx.Exec(query, blockNumber, blockHash, ilkID, value)

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
