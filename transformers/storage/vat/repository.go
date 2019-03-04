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

package vat

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type VatStorageRepository struct {
	db *postgres.DB
}

func (repository *VatStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Dai:
		return repository.insertDai(blockNumber, blockHash, metadata, value.(string))
	case Gem:
		return repository.insertGem(blockNumber, blockHash, metadata, value.(string))
	case IlkArt:
		return repository.insertIlkArt(blockNumber, blockHash, metadata, value.(string))
	case IlkInk:
		return repository.insertIlkInk(blockNumber, blockHash, metadata, value.(string))
	case IlkRate:
		return repository.insertIlkRate(blockNumber, blockHash, metadata, value.(string))
	case IlkTake:
		return repository.insertIlkTake(blockNumber, blockHash, metadata, value.(string))
	case Sin:
		return repository.insertSin(blockNumber, blockHash, metadata, value.(string))
	case UrnArt:
		return repository.insertUrnArt(blockNumber, blockHash, metadata, value.(string))
	case UrnInk:
		return repository.insertUrnInk(blockNumber, blockHash, metadata, value.(string))
	case VatDebt:
		return repository.insertVatDebt(blockNumber, blockHash, value.(string))
	case VatVice:
		return repository.insertVatVice(blockNumber, blockHash, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized vat contract storage name: %s", metadata.Name))
	}
}

func (repository *VatStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *VatStorageRepository) insertDai(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, dai string) error {
	guy, err := getGuy(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(`INSERT INTO maker.vat_dai (block_number, block_hash, guy, dai) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, guy, dai)
	return writeErr
}

func (repository *VatStorageRepository) insertGem(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, gem string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return guyErr
	}
	tx, txErr := repository.db.Begin()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk: %s", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(`INSERT INTO maker.vat_gem (block_number, block_hash, ilk, guy, gem) VALUES ($1, $2, $3, $4, $5)`, blockNumber, blockHash, ilkID, guy, gem)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert gem: %s", writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository *VatStorageRepository) insertIlkArt(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, art string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	tx, txErr := repository.db.Begin()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk: %s", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(`INSERT INTO maker.vat_ilk_art (block_number, block_hash, ilk, art) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilkID, art)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk art: %s", writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository *VatStorageRepository) insertIlkInk(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, ink string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	tx, txErr := repository.db.Begin()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk: %s", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(`INSERT INTO maker.vat_ilk_ink (block_number, block_hash, ilk, ink) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilkID, ink)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk ink: %s", writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository *VatStorageRepository) insertIlkRate(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, rate string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	tx, txErr := repository.db.Begin()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk: %s", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(`INSERT INTO maker.vat_ilk_rate (block_number, block_hash, ilk, rate) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilkID, rate)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk rate: %s", writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository *VatStorageRepository) insertIlkTake(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, take string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	tx, txErr := repository.db.Begin()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk: %s", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(`INSERT INTO maker.vat_ilk_take (block_number, block_hash, ilk, take) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, ilkID, take)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk take: %s", writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository *VatStorageRepository) insertSin(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, sin string) error {
	guy, err := getGuy(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(`INSERT INTO maker.vat_sin (block_number, block_hash, guy, sin) VALUES ($1, $2, $3, $4)`, blockNumber, blockHash, guy, sin)
	return writeErr
}

func (repository *VatStorageRepository) insertUrnArt(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, art string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return guyErr
	}
	tx, txErr := repository.db.Begin()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk: %s", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(`INSERT INTO maker.vat_urn_art (block_number, block_hash, ilk, urn, art) VALUES ($1, $2, $3, $4, $5)`, blockNumber, blockHash, ilkID, guy, art)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert urn art: %s", writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository *VatStorageRepository) insertUrnInk(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, ink string) error {
	ilk, ilkErr := getIlk(metadata.Keys)
	if ilkErr != nil {
		return ilkErr
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return guyErr
	}
	tx, txErr := repository.db.Begin()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert ilk: %s", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(`INSERT INTO maker.vat_urn_ink (block_number, block_hash, ilk, urn, ink) VALUES ($1, $2, $3, $4, $5)`, blockNumber, blockHash, ilkID, guy, ink)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("failed to rollback transaction after failing to insert urn ink: %s", writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository *VatStorageRepository) insertVatDebt(blockNumber int, blockHash, debt string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vat_debt (block_number, block_hash, debt) VALUES ($1, $2, $3)`, blockNumber, blockHash, debt)
	return err
}

func (repository *VatStorageRepository) insertVatVice(blockNumber int, blockHash, vice string) error {
	_, err := repository.db.Exec(`INSERT INTO maker.vat_vice (block_number, block_hash, vice) VALUES ($1, $2, $3)`, blockNumber, blockHash, vice)
	return err
}

func getGuy(keys map[utils.Key]string) (string, error) {
	guy, ok := keys[constants.Guy]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.Guy}
	}
	return guy, nil
}

func getIlk(keys map[utils.Key]string) (string, error) {
	ilk, ok := keys[constants.Ilk]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.Ilk}
	}
	return ilk, nil
}
