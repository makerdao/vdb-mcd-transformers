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

package drip

import (
	"fmt"
	shared2 "github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type DripStorageRepository struct {
	db *postgres.DB
}

func (repository *DripStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository DripStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case IlkRho:
		return repository.insertIlkRho(blockNumber, blockHash, metadata, value.(string))
	case IlkTax:
		return repository.insertIlkTax(blockNumber, blockHash, metadata, value.(string))
	case DripVat:
		return repository.insertDripVat(blockNumber, blockHash, value.(string))
	case DripVow:
		return repository.insertDripVow(blockNumber, blockHash, value.(string))
	case DripRepo:
		return repository.insertDripRepo(blockNumber, blockHash, value.(string))

	default:
		panic(fmt.Sprintf("unrecognized drip contract storage name: %s", metadata.Name))
	}
}

func (repository DripStorageRepository) insertIlkRho(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, rho string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	ilkID, ilkErr := shared2.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk: %s", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(`INSERT INTO maker.drip_ilk_rho (block_number, block_hash, ilk, rho) 
										VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilkID, rho)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert drip ilk rho: %s", writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository DripStorageRepository) insertIlkTax(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, tax string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	ilkID, ilkErr := shared2.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk: %s", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(`INSERT INTO maker.drip_ilk_tax (block_number, block_hash, ilk, tax) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilkID, tax)

	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert pit ilk spot: %s", writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository DripStorageRepository) insertDripVat(blockNumber int, blockHash string, vat string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.drip_vat (block_number, block_hash, vat)  
										VALUES ($1, $2, $3)`, blockNumber, blockHash, vat)
	return err
}

func (repository DripStorageRepository) insertDripVow(blockNumber int, blockHash string, vow string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.drip_vow (block_number, block_hash, vow)  
										VALUES ($1, $2, $3)`, blockNumber, blockHash, vow)
	return err
}

func (repository DripStorageRepository) insertDripRepo(blockNumber int, blockHash string, repo string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.drip_repo (block_number, block_hash, repo)  
										VALUES ($1, $2, $3)`, blockNumber, blockHash, repo)
	return err
}

func getIlk(keys map[utils.Key]string) (string, error) {
	ilk, ok := keys[constants.Ilk]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.Ilk}
	}
	return ilk, nil
}
