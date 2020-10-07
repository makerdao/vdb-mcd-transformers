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
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	insertVatQuery      = `INSERT INTO maker.cdp_manager_vat (diff_id, header_id, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertCdpiQuery     = `INSERT INTO maker.cdp_manager_cdpi (diff_id, header_id, cdpi) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertUrnsQuery     = `INSERT INTO maker.cdp_manager_urns (diff_id, header_id, cdpi, urn) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertListPrevQuery = `INSERT INTO maker.cdp_manager_list_prev (diff_id, header_id, cdpi, prev) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertListNextQuery = `INSERT INTO maker.cdp_manager_list_next (diff_id, header_id, cdpi, next) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertOwnsQuery     = `INSERT INTO maker.cdp_manager_owns (diff_id, header_id, cdpi, owner) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertIlksQuery     = `INSERT INTO maker.cdp_manager_ilks (diff_id, header_id, cdpi, ilk_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertFirstQuery    = `INSERT INTO maker.cdp_manager_first (diff_id, header_id, owner, first) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertLastQuery     = `INSERT INTO maker.cdp_manager_last (diff_id, header_id, owner, last) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertCountQuery    = `INSERT INTO maker.cdp_manager_count (diff_id, header_id, owner, count) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	db *postgres.DB
}

func (repository *StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case Vat:
		return repository.insertVat(diffID, headerID, value.(string))
	case Cdpi:
		return repository.insertCdpi(diffID, headerID, value.(string))
	case Urns:
		return repository.insertUrns(diffID, headerID, metadata, value.(string))
	case ListPrev:
		return repository.insertListPrev(diffID, headerID, metadata, value.(string))
	case ListNext:
		return repository.insertListNext(diffID, headerID, metadata, value.(string))
	case Owns:
		return repository.insertOwns(diffID, headerID, metadata, value.(string))
	case Ilks:
		return repository.insertIlks(diffID, headerID, metadata, value.(string))
	case First:
		return repository.insertFirst(diffID, headerID, metadata, value.(string))
	case Last:
		return repository.insertLast(diffID, headerID, metadata, value.(string))
	case Count:
		return repository.insertCount(diffID, headerID, metadata, value.(string))
	default:
		panic("unrecognized storage metadata name")
	}
}

func (repository *StorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository StorageRepository) insertVat(diffID, headerID int64, vat string) error {
	_, err := repository.db.Exec(insertVatQuery, diffID, headerID, vat)
	if err != nil {
		return fmt.Errorf("error inserting cdp manager vat %s from diff ID %d: %w", vat, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertCdpi(diffID, headerID int64, cdpi string) error {
	_, err := repository.db.Exec(InsertCdpiQuery, diffID, headerID, cdpi)
	if err != nil {
		return fmt.Errorf("error inserting cdp manager cdpi %s from diff ID %d: %w", cdpi, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertUrns(diffID, headerID int64, metadata types.ValueMetadata, urn string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return fmt.Errorf("error getting cdpi for urn: %w", keyErr)
	}
	_, insertErr := repository.db.Exec(insertUrnsQuery, diffID, headerID, cdpi, urn)
	if insertErr != nil {
		return fmt.Errorf("error inserting cdpi %s urn %s from diff ID %d: %w", cdpi, urn, diffID, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertListPrev(diffID, headerID int64, metadata types.ValueMetadata, prev string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return fmt.Errorf("error getting cdpi for list prev: %w", keyErr)
	}
	_, insertErr := repository.db.Exec(insertListPrevQuery, diffID, headerID, cdpi, prev)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting cdpi %s list prev %s from diff ID %d", cdpi, prev, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertListNext(diffID, headerID int64, metadata types.ValueMetadata, next string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return fmt.Errorf("error getting cdpi for list next: %w", keyErr)
	}
	_, insertErr := repository.db.Exec(insertListNextQuery, diffID, headerID, cdpi, next)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting cdpi %s list next %s from diff ID %d", cdpi, next, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertOwns(diffID, headerID int64, metadata types.ValueMetadata, owner string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return fmt.Errorf("error getting cdpi for owns: %w", keyErr)
	}
	_, insertErr := repository.db.Exec(InsertOwnsQuery, diffID, headerID, cdpi, owner)
	if insertErr != nil {
		return fmt.Errorf("error inserting cdpi %s owns %s from diff ID %d: %w", cdpi, owner, diffID, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertIlks(diffID, headerID int64, metadata types.ValueMetadata, ilk string) error {
	cdpi, keyErr := getCdpi(metadata.Keys)
	if keyErr != nil {
		return fmt.Errorf("error getting cdpi for ilk: %w", keyErr)
	}
	ilkId, ilkErr := shared.GetOrCreateIlk(ilk, repository.db)
	if ilkErr != nil {
		return fmt.Errorf("error getting or creating ilk: %w", ilkErr)
	}
	_, insertErr := repository.db.Exec(insertIlksQuery, diffID, headerID, cdpi, ilkId)
	if insertErr != nil {
		return fmt.Errorf("error inserting cdpi %s ilk %s from diff ID %d: %w", cdpi, ilk, diffID, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertFirst(diffID, headerID int64, metadata types.ValueMetadata, first string) error {
	owner, keyErr := getOwner(metadata.Keys)
	if keyErr != nil {
		return fmt.Errorf("error getting owner for first: %w", keyErr)
	}
	_, insertErr := repository.db.Exec(insertFirstQuery, diffID, headerID, owner, first)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting owner %s first %s from diff ID %d", owner, first, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertLast(diffID, headerID int64, metadata types.ValueMetadata, last string) error {
	owner, keyErr := getOwner(metadata.Keys)
	if keyErr != nil {
		return fmt.Errorf("error getting owner for last: %w", keyErr)
	}
	_, insertErr := repository.db.Exec(insertLastQuery, diffID, headerID, owner, last)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting owner %s last %s from diff ID %d", owner, last, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertCount(diffID, headerID int64, metadata types.ValueMetadata, count string) error {
	owner, keyErr := getOwner(metadata.Keys)
	if keyErr != nil {
		return fmt.Errorf("error getting owner for count: %w", keyErr)
	}
	_, insertErr := repository.db.Exec(insertCountQuery, diffID, headerID, owner, count)
	if insertErr != nil {
		msg := fmt.Sprintf("error inserting owner %s count %s from diff ID %d", owner, count, diffID)
		return fmt.Errorf("%s: %w", msg, insertErr)
	}
	return nil
}

func getCdpi(keys map[types.Key]string) (string, error) {
	cdpi, ok := keys[constants.Cdpi]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.Cdpi}
	}
	return cdpi, nil
}

func getOwner(keys map[types.Key]string) (string, error) {
	owner, ok := keys[constants.Owner]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.Owner}
	}
	return owner, nil
}
