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

package vat

import (
	"fmt"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

const (
	insertDaiQuery     = `INSERT INTO maker.vat_dai (block_number, block_hash, guy, dai) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertGemQuery     = `INSERT INTO maker.vat_gem (block_number, block_hash, ilk_id, guy, gem) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertIlkArtQuery  = `INSERT INTO maker.vat_ilk_art (block_number, block_hash, ilk_id, art) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertIlkDustQuery = `INSERT INTO maker.vat_ilk_dust (block_number, block_hash, ilk_id, dust) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertIlkLineQuery = `INSERT INTO maker.vat_ilk_line (block_number, block_hash, ilk_id, line) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertIlkRateQuery = `INSERT INTO maker.vat_ilk_rate (block_number, block_hash, ilk_id, rate) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertIlkSpotQuery = `INSERT INTO maker.vat_ilk_spot (block_number, block_hash, ilk_id, spot) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertSinQuery     = `INSERT INTO maker.vat_sin (block_number, block_hash, guy, sin) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertUrnArtQuery  = `INSERT INTO maker.vat_urn_art (block_number, block_hash, urn_id, art) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertUrnInkQuery  = `INSERT INTO maker.vat_urn_ink (block_number, block_hash, urn_id, ink) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertVatDebtQuery = `INSERT INTO maker.vat_debt (block_number, block_hash, debt) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertVatLineQuery = `INSERT INTO maker.vat_line (block_number, block_hash, line) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertVatLiveQuery = `INSERT INTO maker.vat_live (block_number, block_hash, live) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertVatViceQuery = `INSERT INTO maker.vat_vice (block_number, block_hash, vice) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
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
	case IlkDust:
		return repository.insertIlkDust(blockNumber, blockHash, metadata, value.(string))
	case IlkLine:
		return repository.insertIlkLine(blockNumber, blockHash, metadata, value.(string))
	case IlkRate:
		return repository.insertIlkRate(blockNumber, blockHash, metadata, value.(string))
	case IlkSpot:
		return repository.insertIlkSpot(blockNumber, blockHash, metadata, value.(string))
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
	case VatLine:
		return repository.insertVatLine(blockNumber, blockHash, value.(string))
	case VatLive:
		return repository.insertVatLive(blockNumber, blockHash, value.(string))
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
	_, writeErr := repository.db.Exec(insertDaiQuery, blockNumber, blockHash, guy, dai)
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
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return formatRollbackError("ilk", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(insertGemQuery, blockNumber, blockHash, ilkID, guy, gem)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return formatRollbackError("gem", writeErr.Error())
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
	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkArt, InsertIlkArtQuery, art)
}

func (repository *VatStorageRepository) insertIlkDust(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, dust string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkDust, InsertIlkDustQuery, dust)
}

func (repository *VatStorageRepository) insertIlkLine(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, line string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkLine, InsertIlkLineQuery, line)
}

func (repository *VatStorageRepository) insertIlkRate(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, rate string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkRate, InsertIlkRateQuery, rate)
}

func (repository *VatStorageRepository) insertIlkSpot(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, spot string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(blockNumber, blockHash, ilk, IlkSpot, InsertIlkSpotQuery, spot)
}

func (repository *VatStorageRepository) insertSin(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, sin string) error {
	guy, err := getGuy(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertSinQuery, blockNumber, blockHash, guy, sin)
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
	return repository.insertFieldWithIlkAndUrn(blockNumber, blockHash, ilk, guy, UrnArt, InsertUrnArtQuery, art)
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
	return repository.insertFieldWithIlkAndUrn(blockNumber, blockHash, ilk, guy, UrnInk, InsertUrnInkQuery, ink)
}

func (repository *VatStorageRepository) insertVatDebt(blockNumber int, blockHash, debt string) error {
	_, err := repository.db.Exec(insertVatDebtQuery, blockNumber, blockHash, debt)
	return err
}

func (repository *VatStorageRepository) insertVatLine(blockNumber int, blockHash, line string) error {
	_, err := repository.db.Exec(insertVatLineQuery, blockNumber, blockHash, line)
	return err
}

func (repository *VatStorageRepository) insertVatLive(blockNumber int, blockHash, live string) error {
	_, err := repository.db.Exec(insertVatLiveQuery, blockNumber, blockHash, live)
	return err
}

func (repository *VatStorageRepository) insertVatVice(blockNumber int, blockHash, vice string) error {
	_, err := repository.db.Exec(insertVatViceQuery, blockNumber, blockHash, vice)
	return err
}

func (repository *VatStorageRepository) insertFieldWithIlk(blockNumber int, blockHash, ilk, variableName, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return formatRollbackError("ilk", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(query, blockNumber, blockHash, ilkID, value)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return formatRollbackError(variableName, writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository *VatStorageRepository) insertFieldWithIlkAndUrn(blockNumber int, blockHash, ilk, urn, variableName, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	urnID, urnErr := shared.GetOrCreateUrnInTransaction(urn, ilk, tx)
	if urnErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return formatRollbackError("urn", urnErr.Error())
		}
		return urnErr
	}
	_, writeErr := tx.Exec(query, blockNumber, blockHash, urnID, value)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return formatRollbackError(variableName, writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
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

func formatRollbackError(field, err string) error {
	return fmt.Errorf("failed to rollback transaction after failing to insert %s: %s", field, err)
}
