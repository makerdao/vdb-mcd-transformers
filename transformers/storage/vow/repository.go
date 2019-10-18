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
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertVatQuery        = `INSERT INTO maker.vow_vat (block_number, block_hash, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlapperQuery    = `INSERT INTO maker.vow_flapper (block_number, block_hash, flapper) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFlopperQuery    = `INSERT INTO maker.vow_flopper (block_number, block_hash, flopper) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertSinMappingQuery = `INSERT INTO maker.vow_sin_mapping (block_number, block_hash, era, tab) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertSinIntegerQuery = `INSERT INTO maker.vow_sin_integer (block_number, block_hash, sin) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertAshQuery        = `INSERT INTO maker.vow_ash (block_number, block_hash, ash) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertWaitQuery       = `INSERT INTO maker.vow_wait (block_number, block_hash, wait) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertDumpQuery       = `INSERT INTO maker.vow_dump (block_number, block_hash, dump) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertSumpQuery       = `INSERT INTO maker.vow_sump (block_number, block_hash, sump) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertBumpQuery       = `INSERT INTO maker.vow_bump (block_number, block_hash, bump) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertHumpQuery       = `INSERT INTO maker.vow_hump (block_number, block_hash, hump) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type VowStorageRepository struct {
	db *postgres.DB
}

func (repository *VowStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository VowStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Vat:
		return repository.insertVowVat(blockNumber, blockHash, value.(string))
	case Flapper:
		return repository.insertVowFlapper(blockNumber, blockHash, value.(string))
	case Flopper:
		return repository.insertVowFlopper(blockNumber, blockHash, value.(string))
	case SinMapping:
		return repository.insertSinMapping(blockNumber, blockHash, metadata, value.(string))
	case SinInteger:
		return repository.insertSinInteger(blockNumber, blockHash, value.(string))
	case Ash:
		return repository.insertVowAsh(blockNumber, blockHash, value.(string))
	case Wait:
		return repository.insertVowWait(blockNumber, blockHash, value.(string))
	case Dump:
		return repository.insertVowDump(blockNumber, blockHash, value.(string))
	case Sump:
		return repository.insertVowSump(blockNumber, blockHash, value.(string))
	case Bump:
		return repository.insertVowBump(blockNumber, blockHash, value.(string))
	case Hump:
		return repository.insertVowHump(blockNumber, blockHash, value.(string))
	default:
		panic("unrecognized storage metadata name")
	}
}

func (repository VowStorageRepository) insertVowVat(blockNumber int, blockHash string, vat string) error {
	_, err := repository.db.Exec(insertVatQuery, blockNumber, blockHash, vat)

	return err
}

func (repository VowStorageRepository) insertVowFlapper(blockNumber int, blockHash string, flapper string) error {
	_, err := repository.db.Exec(insertFlapperQuery, blockNumber, blockHash, flapper)

	return err
}

func (repository VowStorageRepository) insertVowFlopper(blockNumber int, blockHash string, flopper string) error {
	_, err := repository.db.Exec(insertFlopperQuery, blockNumber, blockHash, flopper)

	return err
}

func (repository VowStorageRepository) insertSinInteger(blockNumber int, blockHash string, sin string) error {
	_, err := repository.db.Exec(insertSinIntegerQuery, blockNumber, blockHash, sin)

	return err
}

func (repository VowStorageRepository) insertSinMapping(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, sin string) error {
	timestamp, err := getTimestamp(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertSinMappingQuery, blockNumber, blockHash, timestamp, sin)

	return writeErr
}

func (repository VowStorageRepository) insertVowAsh(blockNumber int, blockHash string, ash string) error {
	_, err := repository.db.Exec(insertAshQuery, blockNumber, blockHash, ash)

	return err
}

func (repository VowStorageRepository) insertVowWait(blockNumber int, blockHash string, wait string) error {
	_, err := repository.db.Exec(insertWaitQuery, blockNumber, blockHash, wait)

	return err
}

func (repository VowStorageRepository) insertVowDump(blockNumber int, blockHash string, dump string) error {
	_, err := repository.db.Exec(insertDumpQuery, blockNumber, blockHash, dump)

	return err
}

func (repository VowStorageRepository) insertVowSump(blockNumber int, blockHash string, sump string) error {
	_, err := repository.db.Exec(insertSumpQuery, blockNumber, blockHash, sump)

	return err
}

func (repository VowStorageRepository) insertVowBump(blockNumber int, blockHash string, bump string) error {
	_, err := repository.db.Exec(insertBumpQuery, blockNumber, blockHash, bump)

	return err
}

func (repository VowStorageRepository) insertVowHump(blockNumber int, blockHash string, hump string) error {
	_, err := repository.db.Exec(insertHumpQuery, blockNumber, blockHash, hump)

	return err
}

func getTimestamp(keys map[utils.Key]string) (string, error) {
	timestamp, ok := keys[constants.Timestamp]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.Timestamp}
	}
	return timestamp, nil
}
