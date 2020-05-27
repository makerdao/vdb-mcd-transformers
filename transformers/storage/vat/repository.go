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
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
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

type StorageRepository struct {
	db *postgres.DB
}

func (repository *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case wards.Wards:
		contractAddress := constants.GetContractAddress("MCD_VAT")
		return wards.InsertWards(diffID, headerID, metadata, contractAddress, value.(string), repository.db)
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

func (repository *StorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository *StorageRepository) insertDai(diffID, headerID int64, metadata types.ValueMetadata, dai string) error {
	guy, err := getGuy(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting gut for vat dai: %w", err)
	}
	_, insertErr := repository.db.Exec(insertDaiQuery, diffID, headerID, guy, dai)
	if insertErr != nil {
		return fmt.Errorf("error inserting vat dai %s from diff ID %d: %w", dai, diffID, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertGem(diffID, headerID int64, metadata types.ValueMetadata, gem string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for vat gem: %w", err)
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return fmt.Errorf("error getting guy for vat gem: %w", guyErr)
	}
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return fmt.Errorf("error beginning transaction for vat gem: %w", txErr)
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("ilk", ilkErr.Error())
		}
		return fmt.Errorf("error getting or creating ilk for vat gem: %w", ilkErr)
	}
	_, insertErr := tx.Exec(insertGemQuery, diffID, headerID, ilkID, guy, gem)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("gem", insertErr.Error())
		}
		return fmt.Errorf("error inserting vat gem %s from diff ID %d: %w", gem, diffID, insertErr)
	}
	return tx.Commit()
}

func (repository *StorageRepository) insertIlkArt(diffID, headerID int64, metadata types.ValueMetadata, art string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk art: %w", err)
	}
	insertErr := repository.insertFieldWithIlk(diffID, headerID, ilk, IlkArt, InsertIlkArtQuery, art)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s art %s from diff ID %d: %w", ilk, art, diffID, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertIlkDust(diffID, headerID int64, metadata types.ValueMetadata, dust string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk dust: %w", err)
	}
	insertErr := repository.insertFieldWithIlk(diffID, headerID, ilk, IlkDust, InsertIlkDustQuery, dust)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s dust %s from diff ID %d: %w", ilk, dust, diffID, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertIlkLine(diffID, headerID int64, metadata types.ValueMetadata, line string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk line: %w", err)
	}
	insertErr := repository.insertFieldWithIlk(diffID, headerID, ilk, IlkLine, InsertIlkLineQuery, line)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s line %s from diff ID %d: %w", ilk, line, diffID, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertIlkRate(diffID, headerID int64, metadata types.ValueMetadata, rate string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk rate: %w", err)
	}
	insertErr := repository.insertFieldWithIlk(diffID, headerID, ilk, IlkRate, InsertIlkRateQuery, rate)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s rate %s from diff ID %d: %w", ilk, rate, diffID, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertIlkSpot(diffID, headerID int64, metadata types.ValueMetadata, spot string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk spot: %w", err)
	}
	insertErr := repository.insertFieldWithIlk(diffID, headerID, ilk, IlkSpot, InsertIlkSpotQuery, spot)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s spot %s from diff ID %d: %w", ilk, spot, diffID, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertSin(diffID, headerID int64, metadata types.ValueMetadata, sin string) error {
	guy, err := getGuy(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting guy for vat sin: %w", err)
	}
	_, insertErr := repository.db.Exec(insertSinQuery, diffID, headerID, guy, sin)
	if insertErr != nil {
		return fmt.Errorf("error inserting vat guy %s sin %s from diff ID %d: %w", guy, sin, diffID, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertUrnArt(diffID, headerID int64, metadata types.ValueMetadata, art string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for urn art: %w", err)
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return fmt.Errorf("error getting guy for urn art: %w", guyErr)
	}
	insertErr := repository.insertFieldWithIlkAndUrn(diffID, headerID, ilk, guy, UrnArt, InsertUrnArtQuery, art)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting urn %s %s art %s from diff ID %d", ilk, guy, art, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertUrnInk(diffID, headerID int64, metadata types.ValueMetadata, ink string) error {
	ilk, ilkErr := getIlk(metadata.Keys)
	if ilkErr != nil {
		return fmt.Errorf("error getting ilk for urn ink: %w", ilkErr)
	}
	guy, guyErr := getGuy(metadata.Keys)
	if guyErr != nil {
		return fmt.Errorf("error getting guy for urn ink: %w", guyErr)
	}
	insertErr := repository.insertFieldWithIlkAndUrn(diffID, headerID, ilk, guy, UrnInk, InsertUrnInkQuery, ink)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting urn %s %s ink %s from diff ID %d", ilk, guy, ink, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository *StorageRepository) insertVatDebt(diffID, headerID int64, debt string) error {
	_, err := repository.db.Exec(insertVatDebtQuery, diffID, headerID, debt)
	if err != nil {
		return fmt.Errorf("error inserting vat debt %s from diff ID %d: %w", debt, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertVatLine(diffID, headerID int64, line string) error {
	_, err := repository.db.Exec(insertVatLineQuery, diffID, headerID, line)
	if err != nil {
		return fmt.Errorf("error inserting vat line %s from diff ID %d: %w", line, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertVatLive(diffID, headerID int64, live string) error {
	_, err := repository.db.Exec(insertVatLiveQuery, diffID, headerID, live)
	if err != nil {
		return fmt.Errorf("error inserting vat live %s from diff ID %d: %w", live, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertVatVice(diffID, headerID int64, vice string) error {
	_, err := repository.db.Exec(insertVatViceQuery, diffID, headerID, vice)
	if err != nil {
		return fmt.Errorf("error inserting vat vice %s from diff ID %d: %w", vice, diffID, err)
	}
	return nil
}

func (repository *StorageRepository) insertFieldWithIlk(diffID, headerID int64, ilk, variableName, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return fmt.Errorf("error beginning transaction: %w", txErr)
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("ilk", ilkErr.Error())
		}
		return fmt.Errorf("error getting or creating ilk: %w", ilkErr)
	}
	_, writeErr := tx.Exec(query, diffID, headerID, ilkID, value)
	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError(variableName, writeErr.Error())
		}
		return fmt.Errorf("error inserting field with ilk: %w", writeErr)
	}
	return tx.Commit()
}

func (repository *StorageRepository) insertFieldWithIlkAndUrn(diffID, headerID int64, ilk, urn, variableName, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return fmt.Errorf("error beginning transaction: %w", txErr)
	}

	urnID, urnErr := shared.GetOrCreateUrnInTransaction(urn, ilk, tx)
	if urnErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("urn", urnErr.Error())
		}
		return fmt.Errorf("error getting or creating urn: %w", urnErr)
	}
	_, insertErr := tx.Exec(query, diffID, headerID, urnID, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError(variableName, insertErr.Error())
		}
		return fmt.Errorf("error inserting field with urn: %w", insertErr)
	}
	return tx.Commit()
}

func getGuy(keys map[types.Key]string) (string, error) {
	guy, ok := keys[constants.Guy]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.Guy}
	}
	return guy, nil
}

func getIlk(keys map[types.Key]string) (string, error) {
	ilk, ok := keys[constants.Ilk]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.Ilk}
	}
	return ilk, nil
}
