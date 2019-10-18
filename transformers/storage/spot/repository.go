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

package spot

import (
	"fmt"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	InsertSpotIlkPipQuery = `INSERT INTO maker.spot_ilk_pip (block_number, block_hash, ilk_id, pip) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertSpotIlkMatQuery = `INSERT INTO maker.spot_ilk_mat (block_number, block_hash, ilk_id, mat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertSpotVatQuery    = `INSERT INTO maker.spot_vat (block_number, block_hash, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertSpotParQuery    = `INSERT INTO maker.spot_par (block_number, block_hash, par) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type SpotStorageRepository struct {
	db *postgres.DB
}

func (repository *SpotStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository SpotStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case IlkPip:
		return repository.insertIlkPip(blockNumber, blockHash, metadata, value.(string))
	case IlkMat:
		return repository.insertIlkMat(blockNumber, blockHash, metadata, value.(string))
	case Vat:
		return repository.insertSpotVat(blockNumber, blockHash, value.(string))
	case Par:
		return repository.insertSpotPar(blockNumber, blockHash, value.(string))

	default:
		panic(fmt.Sprintf("unrecognized spot contract storage name: %s", metadata.Name))
	}
}

func (repository SpotStorageRepository) insertIlkPip(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, pip string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkPip, InsertSpotIlkPipQuery, pip)
}

func (repository SpotStorageRepository) insertIlkMat(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, mat string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkMat, InsertSpotIlkMatQuery, mat)
}

func (repository SpotStorageRepository) insertSpotVat(blockNumber int, blockHash string, vat string) error {
	_, err := repository.db.Exec(insertSpotVatQuery, blockNumber, blockHash, vat)
	return err
}

func (repository SpotStorageRepository) insertSpotPar(blockNumber int, blockHash string, par string) error {
	_, err := repository.db.Exec(insertSpotParQuery, blockNumber, blockHash, par)
	return err
}

func (repository *SpotStorageRepository) insertFieldWithIlk(blockNumber int, blockHash, ilk, variableName, query, value string) error {
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
