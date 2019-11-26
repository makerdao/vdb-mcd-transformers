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

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertDaiQuery     = `INSERT INTO maker.vat_dai (header_id, guy, dai) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertGemQuery     = `INSERT INTO maker.vat_gem (header_id, ilk_id, guy, gem) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertIlkArtQuery  = `INSERT INTO maker.vat_ilk_art (header_id, ilk_id, art) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertIlkDustQuery = `INSERT INTO maker.vat_ilk_dust (header_id, ilk_id, dust) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertIlkLineQuery = `INSERT INTO maker.vat_ilk_line (header_id, ilk_id, line) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertIlkRateQuery = `INSERT INTO maker.vat_ilk_rate (header_id, ilk_id, rate) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertIlkSpotQuery = `INSERT INTO maker.vat_ilk_spot (header_id, ilk_id, spot) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertSinQuery     = `INSERT INTO maker.vat_sin (header_id, guy, sin) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertUrnArtQuery  = `INSERT INTO maker.vat_urn_art (header_id, urn_id, art) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertUrnInkQuery  = `INSERT INTO maker.vat_urn_ink (header_id, urn_id, ink) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertVatDebtQuery = `INSERT INTO maker.vat_debt (header_id, debt) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertVatLineQuery = `INSERT INTO maker.vat_line (header_id, line) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertVatLiveQuery = `INSERT INTO maker.vat_live (header_id, live) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertVatViceQuery = `INSERT INTO maker.vat_vice (header_id, vice) VALUES ($1, $2) ON CONFLICT DO NOTHING`
)

type VatStorageRepository struct {
	db *postgres.DB
}

func (repository *VatStorageRepository) Create(headerID int64, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Dai:
		return repository.insertDai(headerID, metadata, value.(string))
	case Gem:
		return repository.insertGem(headerID, metadata, value.(string))
	case IlkArt:
		return repository.insertIlkArt(headerID, metadata, value.(string))
	case IlkDust:
		return repository.insertIlkDust(headerID, metadata, value.(string))
	case IlkLine:
		return repository.insertIlkLine(headerID, metadata, value.(string))
	case IlkRate:
		return repository.insertIlkRate(headerID, metadata, value.(string))
	case IlkSpot:
		return repository.insertIlkSpot(headerID, metadata, value.(string))
	case Sin:
		return repository.insertSin(headerID, metadata, value.(string))
	case UrnArt:
		return repository.insertUrnArt(headerID, metadata, value.(string))
	case UrnInk:
		return repository.insertUrnInk(headerID, metadata, value.(string))
	case Debt:
		return repository.insertVatDebt(headerID, value.(string))
	case Vice:
		return repository.insertVatVice(headerID, value.(string))
	case Line:
		return repository.insertVatLine(headerID, value.(string))
	case Live:
		return repository.insertVatLive(headerID, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized vat contract storage name: %s", metadata.Name))
	}
}

func (repository *VatStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *VatStorageRepository) insertDai(headerID int64, metadata utils.StorageValueMetadata, dai string) error {
	guy, err := getGuy(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertDaiQuery, headerID, guy, dai)
	return writeErr
}

func (repository *VatStorageRepository) insertGem(headerID int64, metadata utils.StorageValueMetadata, gem string) error {
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
			return shared.FormatRollbackError("ilk", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(insertGemQuery, headerID, ilkID, guy, gem)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("gem", writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository *VatStorageRepository) insertIlkArt(headerID int64, metadata utils.StorageValueMetadata, art string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(headerID, ilk, IlkArt, InsertIlkArtQuery, art)
}

func (repository *VatStorageRepository) insertIlkDust(headerID int64, metadata utils.StorageValueMetadata, dust string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(headerID, ilk, IlkDust, InsertIlkDustQuery, dust)
}

func (repository *VatStorageRepository) insertIlkLine(headerID int64, metadata utils.StorageValueMetadata, line string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(headerID, ilk, IlkLine, InsertIlkLineQuery, line)
}

func (repository *VatStorageRepository) insertIlkRate(headerID int64, metadata utils.StorageValueMetadata, rate string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(headerID, ilk, IlkRate, InsertIlkRateQuery, rate)
}

func (repository *VatStorageRepository) insertIlkSpot(headerID int64, metadata utils.StorageValueMetadata, spot string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(headerID, ilk, IlkSpot, InsertIlkSpotQuery, spot)
}

func (repository *VatStorageRepository) insertSin(headerID int64, metadata utils.StorageValueMetadata, sin string) error {
	guy, err := getGuy(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertSinQuery, headerID, guy, sin)
	return writeErr
}

func (repository *VatStorageRepository) insertUrnArt(headerID int64, metadata utils.StorageValueMetadata, art string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return guyErr
	}
	return repository.insertFieldWithIlkAndUrn(headerID, ilk, guy, UrnArt, InsertUrnArtQuery, art)
}

func (repository *VatStorageRepository) insertUrnInk(headerID int64, metadata utils.StorageValueMetadata, ink string) error {
	ilk, ilkErr := getIlk(metadata.Keys)
	if ilkErr != nil {
		return ilkErr
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return guyErr
	}
	return repository.insertFieldWithIlkAndUrn(headerID, ilk, guy, UrnInk, InsertUrnInkQuery, ink)
}

func (repository *VatStorageRepository) insertVatDebt(headerID int64, debt string) error {
	_, err := repository.db.Exec(insertVatDebtQuery, headerID, debt)
	return err
}

func (repository *VatStorageRepository) insertVatLine(headerID int64, line string) error {
	_, err := repository.db.Exec(insertVatLineQuery, headerID, line)
	return err
}

func (repository *VatStorageRepository) insertVatLive(headerID int64, live string) error {
	_, err := repository.db.Exec(insertVatLiveQuery, headerID, live)
	return err
}

func (repository *VatStorageRepository) insertVatVice(headerID int64, vice string) error {
	_, err := repository.db.Exec(insertVatViceQuery, headerID, vice)
	return err
}

func (repository *VatStorageRepository) insertFieldWithIlk(headerID int64, ilk, variableName, query, value string) error {
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
	_, writeErr := tx.Exec(query, headerID, ilkID, value)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError(variableName, writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository *VatStorageRepository) insertFieldWithIlkAndUrn(headerID int64, ilk, urn, variableName, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}

	urnID, urnErr := shared.GetOrCreateUrnInTransaction(urn, ilk, tx)
	if urnErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("urn", urnErr.Error())
		}
		return urnErr
	}
	_, writeErr := tx.Exec(query, headerID, urnID, value)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError(variableName, writeErr.Error())
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
