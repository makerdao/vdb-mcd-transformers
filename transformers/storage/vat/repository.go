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
	insertDaiQuery     = `INSERT INTO maker.vat_dai (diff_id, header_id, guy, dai) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertGemQuery     = `INSERT INTO maker.vat_gem (diff_id, header_id, ilk_id, guy, gem) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	InsertIlkArtQuery  = `INSERT INTO maker.vat_ilk_art (diff_id, header_id, ilk_id, art) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertIlkDustQuery = `INSERT INTO maker.vat_ilk_dust (diff_id, header_id, ilk_id, dust) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertIlkLineQuery = `INSERT INTO maker.vat_ilk_line (diff_id, header_id, ilk_id, line) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertIlkRateQuery = `INSERT INTO maker.vat_ilk_rate (diff_id, header_id, ilk_id, rate) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertIlkSpotQuery = `INSERT INTO maker.vat_ilk_spot (diff_id, header_id, ilk_id, spot) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertSinQuery     = `INSERT INTO maker.vat_sin (diff_id, header_id, guy, sin) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertUrnArtQuery  = `INSERT INTO maker.vat_urn_art (diff_id, header_id, urn_id, art) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertUrnInkQuery  = `INSERT INTO maker.vat_urn_ink (diff_id, header_id, urn_id, ink) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertVatDebtQuery = `INSERT INTO maker.vat_debt (diff_id, header_id, debt) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertVatLineQuery = `INSERT INTO maker.vat_line (diff_id, header_id, line) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertVatLiveQuery = `INSERT INTO maker.vat_live (diff_id, header_id, live) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertVatViceQuery = `INSERT INTO maker.vat_vice (diff_id, header_id, vice) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type VatStorageRepository struct {
	db *postgres.DB
}

func (repository *VatStorageRepository) Create(diffID, headerID int64, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Dai:
		return repository.insertDai(diffID, headerID, metadata, value.(string))
	case Gem:
		return repository.insertGem(diffID, headerID, metadata, value.(string))
	case IlkArt:
		return repository.insertIlkArt(diffID, headerID, metadata, value.(string))
	case IlkDust:
		return repository.insertIlkDust(diffID, headerID, metadata, value.(string))
	case IlkLine:
		return repository.insertIlkLine(diffID, headerID, metadata, value.(string))
	case IlkRate:
		return repository.insertIlkRate(diffID, headerID, metadata, value.(string))
	case IlkSpot:
		return repository.insertIlkSpot(diffID, headerID, metadata, value.(string))
	case Sin:
		return repository.insertSin(diffID, headerID, metadata, value.(string))
	case UrnArt:
		return repository.insertUrnArt(diffID, headerID, metadata, value.(string))
	case UrnInk:
		return repository.insertUrnInk(diffID, headerID, metadata, value.(string))
	case Debt:
		return repository.insertVatDebt(diffID, headerID, value.(string))
	case Vice:
		return repository.insertVatVice(diffID, headerID, value.(string))
	case Line:
		return repository.insertVatLine(diffID, headerID, value.(string))
	case Live:
		return repository.insertVatLive(diffID, headerID, value.(string))
	default:
		panic(fmt.Sprintf("unrecognized vat contract storage name: %s", metadata.Name))
	}
}

func (repository *VatStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *VatStorageRepository) insertDai(diffID, headerID int64, metadata utils.StorageValueMetadata, dai string) error {
	guy, err := getGuy(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertDaiQuery, diffID, headerID, guy, dai)
	return writeErr
}

func (repository *VatStorageRepository) insertGem(diffID, headerID int64, metadata utils.StorageValueMetadata, gem string) error {
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
	_, writeErr := tx.Exec(insertGemQuery, diffID, headerID, ilkID, guy, gem)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("gem", writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func (repository *VatStorageRepository) insertIlkArt(diffID, headerID int64, metadata utils.StorageValueMetadata, art string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(diffID, headerID, ilk, IlkArt, InsertIlkArtQuery, art)
}

func (repository *VatStorageRepository) insertIlkDust(diffID, headerID int64, metadata utils.StorageValueMetadata, dust string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(diffID, headerID, ilk, IlkDust, InsertIlkDustQuery, dust)
}

func (repository *VatStorageRepository) insertIlkLine(diffID, headerID int64, metadata utils.StorageValueMetadata, line string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(diffID, headerID, ilk, IlkLine, InsertIlkLineQuery, line)
}

func (repository *VatStorageRepository) insertIlkRate(diffID, headerID int64, metadata utils.StorageValueMetadata, rate string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(diffID, headerID, ilk, IlkRate, InsertIlkRateQuery, rate)
}

func (repository *VatStorageRepository) insertIlkSpot(diffID, headerID int64, metadata utils.StorageValueMetadata, spot string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(diffID, headerID, ilk, IlkSpot, InsertIlkSpotQuery, spot)
}

func (repository *VatStorageRepository) insertSin(diffID, headerID int64, metadata utils.StorageValueMetadata, sin string) error {
	guy, err := getGuy(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertSinQuery, diffID, headerID, guy, sin)
	return writeErr
}

func (repository *VatStorageRepository) insertUrnArt(diffID, headerID int64, metadata utils.StorageValueMetadata, art string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return guyErr
	}
	return repository.insertFieldWithIlkAndUrn(diffID, headerID, ilk, guy, UrnArt, InsertUrnArtQuery, art)
}

func (repository *VatStorageRepository) insertUrnInk(diffID, headerID int64, metadata utils.StorageValueMetadata, ink string) error {
	ilk, ilkErr := getIlk(metadata.Keys)
	if ilkErr != nil {
		return ilkErr
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return guyErr
	}
	return repository.insertFieldWithIlkAndUrn(diffID, headerID, ilk, guy, UrnInk, InsertUrnInkQuery, ink)
}

func (repository *VatStorageRepository) insertVatDebt(diffID, headerID int64, debt string) error {
	_, err := repository.db.Exec(insertVatDebtQuery, diffID, headerID, debt)
	return err
}

func (repository *VatStorageRepository) insertVatLine(diffID, headerID int64, line string) error {
	_, err := repository.db.Exec(insertVatLineQuery, diffID, headerID, line)
	return err
}

func (repository *VatStorageRepository) insertVatLive(diffID, headerID int64, live string) error {
	_, err := repository.db.Exec(insertVatLiveQuery, diffID, headerID, live)
	return err
}

func (repository *VatStorageRepository) insertVatVice(diffID, headerID int64, vice string) error {
	_, err := repository.db.Exec(insertVatViceQuery, diffID, headerID, vice)
	return err
}

func (repository *VatStorageRepository) insertFieldWithIlk(diffID, headerID int64, ilk, variableName, query, value string) error {
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

func (repository *VatStorageRepository) insertFieldWithIlkAndUrn(diffID, headerID int64, ilk, urn, variableName, query, value string) error {
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
	_, writeErr := tx.Exec(query, diffID, headerID, urnID, value)
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
