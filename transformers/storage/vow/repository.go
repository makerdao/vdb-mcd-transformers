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

package vow

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertVatQuery        = `INSERT INTO maker.vow_vat (diff_id, header_id, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlapperQuery    = `INSERT INTO maker.vow_flapper (diff_id, header_id, flapper) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlopperQuery    = `INSERT INTO maker.vow_flopper (diff_id, header_id, flopper) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertSinMappingQuery = `INSERT INTO maker.vow_sin_mapping (diff_id, header_id, era, tab) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertSinIntegerQuery = `INSERT INTO maker.vow_sin_integer (diff_id, header_id, sin) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertAshQuery        = `INSERT INTO maker.vow_ash (diff_id, header_id, ash) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertWaitQuery       = `INSERT INTO maker.vow_wait (diff_id, header_id, wait) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertDumpQuery       = `INSERT INTO maker.vow_dump (diff_id, header_id, dump) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertSumpQuery       = `INSERT INTO maker.vow_sump (diff_id, header_id, sump) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertBumpQuery       = `INSERT INTO maker.vow_bump (diff_id, header_id, bump) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertHumpQuery       = `INSERT INTO maker.vow_hump (diff_id, header_id, hump) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertLiveQuery       = `INSERT INTO maker.vow_live (diff_id, header_id, live) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository *StorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repository.ContractAddress, value.(string), repository.db)
	case Vat:
		return repository.insertVowVat(diffID, headerID, value.(string))
	case Flapper:
		return repository.insertVowFlapper(diffID, headerID, value.(string))
	case Flopper:
		return repository.insertVowFlopper(diffID, headerID, value.(string))
	case SinMapping:
		return repository.insertSinMapping(diffID, headerID, metadata, value.(string))
	case SinInteger:
		return repository.insertSinInteger(diffID, headerID, value.(string))
	case Ash:
		return repository.insertVowAsh(diffID, headerID, value.(string))
	case Wait:
		return repository.insertVowWait(diffID, headerID, value.(string))
	case Dump:
		return repository.insertVowDump(diffID, headerID, value.(string))
	case Sump:
		return repository.insertVowSump(diffID, headerID, value.(string))
	case Bump:
		return repository.insertVowBump(diffID, headerID, value.(string))
	case Hump:
		return repository.insertVowHump(diffID, headerID, value.(string))
	case Live:
		return repository.insertVowLive(diffID, headerID, value.(string))
	default:
		return fmt.Errorf("unrecognized storage metadata name: %s", metadata.Name)
	}
}

func (repository StorageRepository) insertVowVat(diffID, headerID int64, vat string) error {
	_, err := repository.db.Exec(insertVatQuery, diffID, headerID, vat)
	if err != nil {
		return fmt.Errorf("error inserting vow vat %s from diff ID %d: %w", vat, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertVowFlapper(diffID, headerID int64, flapper string) error {
	_, err := repository.db.Exec(insertFlapperQuery, diffID, headerID, flapper)
	if err != nil {
		return fmt.Errorf("error inserting vow flapper %s from diff ID %d: %w", flapper, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertVowFlopper(diffID, headerID int64, flopper string) error {
	_, err := repository.db.Exec(insertFlopperQuery, diffID, headerID, flopper)
	if err != nil {
		return fmt.Errorf("error inserting vow flopper %s from diff ID %d: %w", flopper, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertSinInteger(diffID, headerID int64, sin string) error {
	_, err := repository.db.Exec(insertSinIntegerQuery, diffID, headerID, sin)
	if err != nil {
		return fmt.Errorf("error inserting vow sin integer %s, from diff ID %d: %w", sin, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertSinMapping(diffID, headerID int64, metadata types.ValueMetadata, sin string) error {
	timestamp, err := getTimestamp(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting timestamp for vow sin mapping: %w", err)
	}
	_, insertErr := repository.db.Exec(insertSinMappingQuery, diffID, headerID, timestamp, sin)
	if insertErr != nil {
		msgToFormat := "error inserting vow sin mapping with timestamp %s and sin %s from diff ID %d"
		msg := fmt.Sprintf(msgToFormat, timestamp, sin, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertVowAsh(diffID, headerID int64, ash string) error {
	_, err := repository.db.Exec(insertAshQuery, diffID, headerID, ash)
	if err != nil {
		return fmt.Errorf("error inserting vow ash %s from diff ID %d: %w", ash, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertVowWait(diffID, headerID int64, wait string) error {
	_, err := repository.db.Exec(insertWaitQuery, diffID, headerID, wait)
	if err != nil {
		return fmt.Errorf("error inserting vow wait %s from diff ID %d: %w", wait, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertVowDump(diffID, headerID int64, dump string) error {
	_, err := repository.db.Exec(insertDumpQuery, diffID, headerID, dump)
	if err != nil {
		return fmt.Errorf("error inserting vow dump %s from diff ID %d: %w", dump, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertVowSump(diffID, headerID int64, sump string) error {
	_, err := repository.db.Exec(insertSumpQuery, diffID, headerID, sump)
	if err != nil {
		return fmt.Errorf("error inserting vow sump %s from diff ID %d: %w", sump, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertVowBump(diffID, headerID int64, bump string) error {
	_, err := repository.db.Exec(insertBumpQuery, diffID, headerID, bump)
	if err != nil {
		return fmt.Errorf("error inserting vow bump %s from diff ID %d: %w", bump, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertVowHump(diffID, headerID int64, hump string) error {
	_, err := repository.db.Exec(insertHumpQuery, diffID, headerID, hump)
	if err != nil {
		return fmt.Errorf("error inserting vow hump %s from diff ID %d: %w", hump, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertVowLive(diffID, headerID int64, live string) error {
	_, err := repository.db.Exec(insertLiveQuery, diffID, headerID, live)
	if err != nil {
		return fmt.Errorf("error inserting vow live %s from diff ID %d: %w", live, diffID, err)
	}
	return nil
}

func getTimestamp(keys map[types.Key]string) (string, error) {
	timestamp, ok := keys[constants.Timestamp]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.Timestamp}
	}
	return timestamp, nil
}
