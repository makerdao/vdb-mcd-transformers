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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertVatQuery      = `INSERT INTO maker.cdp_manager_vat (header_id, vat) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	InsertCdpiQuery     = `INSERT INTO maker.cdp_manager_cdpi (header_id, cdpi) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertUrnsQuery     = `INSERT INTO maker.cdp_manager_urns (header_id, cdpi, urn) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertListPrevQuery = `INSERT INTO maker.cdp_manager_list_prev (header_id, cdpi, prev) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertListNextQuery = `INSERT INTO maker.cdp_manager_list_next (header_id, cdpi, next) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertOwnsQuery     = `INSERT INTO maker.cdp_manager_owns (header_id, cdpi, owner) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertIlksQuery     = `INSERT INTO maker.cdp_manager_ilks (header_id, cdpi, ilk_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertFirstQuery    = `INSERT INTO maker.cdp_manager_first (header_id, owner, first) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertLastQuery     = `INSERT INTO maker.cdp_manager_last (header_id, owner, last) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertCountQuery    = `INSERT INTO maker.cdp_manager_count (header_id, owner, count) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type CdpManagerStorageRepository struct {
	db *postgres.DB
}

func (repository *CdpManagerStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository CdpManagerStorageRepository) Create(headerID int64, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Vat:
		return repository.insertVat(headerID, value.(string))
	case Cdpi:
		return repository.insertCdpi(headerID, value.(string))
	case Urns:
		return repository.insertUrns(headerID, metadata, value.(string))
	case ListPrev:
		return repository.insertListPrev(headerID, metadata, value.(string))
	case ListNext:
		return repository.insertListNext(headerID, metadata, value.(string))
	case Owns:
		return repository.insertOwns(headerID, metadata, value.(string))
	case Ilks:
		return repository.insertIlks(headerID, metadata, value.(string))
	case First:
		return repository.insertFirst(headerID, metadata, value.(string))
	case Last:
		return repository.insertLast(headerID, metadata, value.(string))
	case Count:
		return repository.insertCount(headerID, metadata, value.(string))
	default:
		panic("unrecognized storage metadata name")
	}
}

func (repository CdpManagerStorageRepository) insertVat(headerID int64, vat string) error {
	_, err := repository.db.Exec(insertVatQuery, headerID, vat)
	return err
}

func (repository CdpManagerStorageRepository) insertCdpi(headerID int64, cdpi string) error {
	_, err := repository.db.Exec(InsertCdpiQuery, headerID, cdpi)
	return err
}

func (repository CdpManagerStorageRepository) insertUrns(headerID int64, metadata utils.StorageValueMetadata, urns string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertUrnsQuery, headerID, cdpi, urns)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertListPrev(headerID int64, metadata utils.StorageValueMetadata, prev string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertListPrevQuery, headerID, cdpi, prev)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertListNext(headerID int64, metadata utils.StorageValueMetadata, next string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertListNextQuery, headerID, cdpi, next)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertOwns(headerID int64, metadata utils.StorageValueMetadata, owner string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(InsertOwnsQuery, headerID, cdpi, owner)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertIlks(headerID int64, metadata utils.StorageValueMetadata, ilks string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	ilkId, ilkErr := shared.GetOrCreateIlk(ilks, repository.db)
	if ilkErr != nil {
		return ilkErr
	}
	_, writeErr := repository.db.Exec(insertIlksQuery, headerID, cdpi, ilkId)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertFirst(headerID int64, metadata utils.StorageValueMetadata, first string) error {
	owner, keyErr := getOwner(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertFirstQuery, headerID, owner, first)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertLast(headerID int64, metadata utils.StorageValueMetadata, last string) error {
	owner, keyErr := getOwner(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertLastQuery, headerID, owner, last)
	return writeErr
}

func (repository CdpManagerStorageRepository) insertCount(headerID int64, metadata utils.StorageValueMetadata, count string) error {
	owner, keyErr := getOwner(metadata.Keys)
	if keyErr != nil {
		return keyErr
	}

	_, writeErr := repository.db.Exec(insertCountQuery, headerID, owner, count)
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
