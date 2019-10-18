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

package cdp_manager

import (
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertVatQuery      = `INSERT INTO maker.cdp_manager_vat (block_number, block_hash, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertCdpiQuery     = `INSERT INTO maker.cdp_manager_cdpi (block_number, block_hash, cdpi) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertUrnsQuery     = `INSERT INTO maker.cdp_manager_urns (block_number, block_hash, cdpi, urn) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertListPrevQuery = `INSERT INTO maker.cdp_manager_list_prev (block_number, block_hash, cdpi, prev) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertListNextQuery = `INSERT INTO maker.cdp_manager_list_next (block_number, block_hash, cdpi, next) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertOwnsQuery     = `INSERT INTO maker.cdp_manager_owns (block_number, block_hash, cdpi, owner) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertIlksQuery     = `INSERT INTO maker.cdp_manager_ilks (block_number, block_hash, cdpi, ilk_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFirstQuery    = `INSERT INTO maker.cdp_manager_first (block_number, block_hash, owner, first) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertLastQuery     = `INSERT INTO maker.cdp_manager_last (block_number, block_hash, owner, last) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertCountQuery    = `INSERT INTO maker.cdp_manager_count (block_number, block_hash, owner, count) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
)

type CdpManagerStorageRepository struct {
	db *postgres.DB
}

func (repository *CdpManagerStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository CdpManagerStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Vat:
		return repository.insertVat(blockNumber, blockHash, value.(string))
	case Cdpi:
		return repository.insertCdpi(blockNumber, blockHash, value.(string))
	case Urns:
		return repository.insertUrns(blockNumber, blockHash, metadata, value.(string))
	case ListPrev:
		return repository.insertListPrev(blockNumber, blockHash, metadata, value.(string))
	case ListNext:
		return repository.insertListNext(blockNumber, blockHash, metadata, value.(string))
	case Owns:
		return repository.insertOwns(blockNumber, blockHash, metadata, value.(string))
	case Ilks:
		return repository.insertIlks(blockNumber, blockHash, metadata, value.(string))
	case First:
		return repository.insertFirst(blockNumber, blockHash, metadata, value.(string))
	case Last:
		return repository.insertLast(blockNumber, blockHash, metadata, value.(string))
	case Count:
		return repository.insertCount(blockNumber, blockHash, metadata, value.(string))
	default:
		panic("unrecognized storage metadata name")
	}
}

func (repository CdpManagerStorageRepository) insertVat(blockNumber int, blockHash string, vat string) error {
	_, err := repository.db.Exec(insertVatQuery, blockNumber, blockHash, vat)
	return err
}

func (repository CdpManagerStorageRepository) insertCdpi(blockNumber int, blockHash string, cdpi string) error {
	_, err := repository.db.Exec(InsertCdpiQuery, blockNumber, blockHash, cdpi)
	return err
}

func (repository CdpManagerStorageRepository) insertUrns(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, urns string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertUrnsQuery, blockNumber, blockHash, cdpi, urns)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertListPrev(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, prev string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertListPrevQuery, blockNumber, blockHash, cdpi, prev)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertListNext(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, next string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertListNextQuery, blockNumber, blockHash, cdpi, next)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertOwns(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, owner string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(InsertOwnsQuery, blockNumber, blockHash, cdpi, owner)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertIlks(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, ilks string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	ilkId, ilkErr := shared.GetOrCreateIlk(ilks, repository.db)
	if ilkErr != nil {
		return ilkErr
	}
	_, writeErr := repository.db.Exec(insertIlksQuery, blockNumber, blockHash, cdpi, ilkId)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertFirst(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, first string) error {
	owner, keyErr := getOwner(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertFirstQuery, blockNumber, blockHash, owner, first)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertLast(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, last string) error {
	owner, keyErr := getOwner(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertLastQuery, blockNumber, blockHash, owner, last)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertCount(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, count string) error {
	owner, keyErr := getOwner(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertCountQuery, blockNumber, blockHash, owner, count)
	return writeErr
}

func getCdpi(keys map[utils.Key]string) (string, error) {
	cdpi, ok := keys[constants.Cdpi]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.Cdpi}
	}
	return cdpi, nil
}

func getOwner(keys map[utils.Key]string) (string, error) {
	owner, ok := keys[constants.Owner]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.Owner}
	}
	return owner, nil
}
