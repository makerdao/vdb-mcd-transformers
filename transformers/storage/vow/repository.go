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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertVatQuery        = `INSERT INTO maker.vow_vat (header_id, vat) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertFlapperQuery    = `INSERT INTO maker.vow_flapper (header_id, flapper) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertFlopperQuery    = `INSERT INTO maker.vow_flopper (header_id, flopper) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertSinMappingQuery = `INSERT INTO maker.vow_sin_mapping (header_id, era, tab) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertSinIntegerQuery = `INSERT INTO maker.vow_sin_integer (header_id, sin) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertAshQuery        = `INSERT INTO maker.vow_ash (header_id, ash) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertWaitQuery       = `INSERT INTO maker.vow_wait (header_id, wait) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertDumpQuery       = `INSERT INTO maker.vow_dump (header_id, dump) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertSumpQuery       = `INSERT INTO maker.vow_sump (header_id, sump) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertBumpQuery       = `INSERT INTO maker.vow_bump (header_id, bump) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertHumpQuery       = `INSERT INTO maker.vow_hump (header_id, hump) VALUES ($1, $2) ON CONFLICT DO NOTHING`
)

type VowStorageRepository struct {
	db *postgres.DB
}

func (repository *VowStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository VowStorageRepository) Create(diffID, headerID int64, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Vat:
		return repository.insertVowVat(headerID, value.(string))
	case Flapper:
		return repository.insertVowFlapper(headerID, value.(string))
	case Flopper:
		return repository.insertVowFlopper(headerID, value.(string))
	case SinMapping:
		return repository.insertSinMapping(headerID, metadata, value.(string))
	case SinInteger:
		return repository.insertSinInteger(headerID, value.(string))
	case Ash:
		return repository.insertVowAsh(headerID, value.(string))
	case Wait:
		return repository.insertVowWait(headerID, value.(string))
	case Dump:
		return repository.insertVowDump(headerID, value.(string))
	case Sump:
		return repository.insertVowSump(headerID, value.(string))
	case Bump:
		return repository.insertVowBump(headerID, value.(string))
	case Hump:
		return repository.insertVowHump(headerID, value.(string))
	default:
		panic("unrecognized storage metadata name")
	}
}

func (repository VowStorageRepository) insertVowVat(headerID int64, vat string) error {
	_, err := repository.db.Exec(insertVatQuery, headerID, vat)

	return err
}

func (repository VowStorageRepository) insertVowFlapper(headerID int64, flapper string) error {
	_, err := repository.db.Exec(insertFlapperQuery, headerID, flapper)

	return err
}

func (repository VowStorageRepository) insertVowFlopper(headerID int64, flopper string) error {
	_, err := repository.db.Exec(insertFlopperQuery, headerID, flopper)

	return err
}

func (repository VowStorageRepository) insertSinInteger(headerID int64, sin string) error {
	_, err := repository.db.Exec(insertSinIntegerQuery, headerID, sin)

	return err
}

func (repository VowStorageRepository) insertSinMapping(headerID int64, metadata utils.StorageValueMetadata, sin string) error {
	timestamp, err := getTimestamp(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertSinMappingQuery, headerID, timestamp, sin)

	return writeErr
}

func (repository VowStorageRepository) insertVowAsh(headerID int64, ash string) error {
	_, err := repository.db.Exec(insertAshQuery, headerID, ash)

	return err
}

func (repository VowStorageRepository) insertVowWait(headerID int64, wait string) error {
	_, err := repository.db.Exec(insertWaitQuery, headerID, wait)

	return err
}

func (repository VowStorageRepository) insertVowDump(headerID int64, dump string) error {
	_, err := repository.db.Exec(insertDumpQuery, headerID, dump)

	return err
}

func (repository VowStorageRepository) insertVowSump(headerID int64, sump string) error {
	_, err := repository.db.Exec(insertSumpQuery, headerID, sump)

	return err
}

func (repository VowStorageRepository) insertVowBump(headerID int64, bump string) error {
	_, err := repository.db.Exec(insertBumpQuery, headerID, bump)

	return err
}

func (repository VowStorageRepository) insertVowHump(headerID int64, hump string) error {
	_, err := repository.db.Exec(insertHumpQuery, headerID, hump)

	return err
}

func getTimestamp(keys map[utils.Key]string) (string, error) {
	timestamp, ok := keys[constants.Timestamp]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.Timestamp}
	}
	return timestamp, nil
}
