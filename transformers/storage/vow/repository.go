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

type VowStorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository *VowStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository VowStorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
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
		panic(fmt.Sprintf("unrecognized storage metadata name: %s", metadata.Name))
	}
}

func (repository VowStorageRepository) insertVowVat(diffID, headerID int64, vat string) error {
	_, err := repository.db.Exec(insertVatQuery, diffID, headerID, vat)

	return err
}

func (repository VowStorageRepository) insertVowFlapper(diffID, headerID int64, flapper string) error {
	_, err := repository.db.Exec(insertFlapperQuery, diffID, headerID, flapper)

	return err
}

func (repository VowStorageRepository) insertVowFlopper(diffID, headerID int64, flopper string) error {
	_, err := repository.db.Exec(insertFlopperQuery, diffID, headerID, flopper)

	return err
}

func (repository VowStorageRepository) insertSinInteger(diffID, headerID int64, sin string) error {
	_, err := repository.db.Exec(insertSinIntegerQuery, diffID, headerID, sin)

	return err
}

func (repository VowStorageRepository) insertSinMapping(diffID, headerID int64, metadata types.ValueMetadata, sin string) error {
	timestamp, err := getTimestamp(metadata.Keys)
	if err != nil {
		return err
	}
	_, writeErr := repository.db.Exec(insertSinMappingQuery, diffID, headerID, timestamp, sin)

	return writeErr
}

func (repository VowStorageRepository) insertVowAsh(diffID, headerID int64, ash string) error {
	_, err := repository.db.Exec(insertAshQuery, diffID, headerID, ash)

	return err
}

func (repository VowStorageRepository) insertVowWait(diffID, headerID int64, wait string) error {
	_, err := repository.db.Exec(insertWaitQuery, diffID, headerID, wait)

	return err
}

func (repository VowStorageRepository) insertVowDump(diffID, headerID int64, dump string) error {
	_, err := repository.db.Exec(insertDumpQuery, diffID, headerID, dump)

	return err
}

func (repository VowStorageRepository) insertVowSump(diffID, headerID int64, sump string) error {
	_, err := repository.db.Exec(insertSumpQuery, diffID, headerID, sump)

	return err
}

func (repository VowStorageRepository) insertVowBump(diffID, headerID int64, bump string) error {
	_, err := repository.db.Exec(insertBumpQuery, diffID, headerID, bump)

	return err
}

func (repository VowStorageRepository) insertVowHump(diffID, headerID int64, hump string) error {
	_, err := repository.db.Exec(insertHumpQuery, diffID, headerID, hump)

	return err
}

func (repository VowStorageRepository) insertVowLive(diffID, headerID int64, live string) error {
	_, err := repository.db.Exec(insertLiveQuery, diffID, headerID, live)

	return err
}

func getTimestamp(keys map[types.Key]string) (string, error) {
	timestamp, ok := keys[constants.Timestamp]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.Timestamp}
	}
	return timestamp, nil
}
