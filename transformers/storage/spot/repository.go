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

package spot

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/utilities/wards"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	InsertSpotIlkPipQuery = `INSERT INTO maker.spot_ilk_pip (diff_id, header_id, ilk_id, pip) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	InsertSpotIlkMatQuery = `INSERT INTO maker.spot_ilk_mat (diff_id, header_id, ilk_id, mat) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	insertSpotVatQuery    = `INSERT INTO maker.spot_vat (diff_id, header_id, vat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertSpotParQuery    = `INSERT INTO maker.spot_par (diff_id, header_id, par) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertSpotLiveQuery   = `INSERT INTO maker.spot_live (diff_id, header_id, live) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
)

type StorageRepository struct {
	db              *postgres.DB
	ContractAddress string
}

func (repository StorageRepository) Create(diffID, headerID int64, metadata types.ValueMetadata, value interface{}) error {
	switch metadata.Name {
	case wards.Wards:
		return wards.InsertWards(diffID, headerID, metadata, repository.ContractAddress, value.(string), repository.db)
	case IlkPip:
		return repository.insertIlkPip(diffID, headerID, metadata, value.(string))
	case IlkMat:
		return repository.insertIlkMat(diffID, headerID, metadata, value.(string))
	case Vat:
		return repository.insertSpotVat(diffID, headerID, value.(string))
	case Par:
		return repository.insertSpotPar(diffID, headerID, value.(string))
	case Live:
		return repository.insertSpotLive(diffID, headerID, value.(string))

	default:
		return fmt.Errorf("unrecognized spot contract storage name: %s", metadata.Name)
	}
}

func (repository *StorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository StorageRepository) insertIlkPip(diffID, headerID int64, metadata types.ValueMetadata, pip string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk pip: %w", err)
	}

	insertErr := shared.InsertFieldWithIlk(diffID, headerID, ilk, IlkPip, InsertSpotIlkPipQuery, pip, repository.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s pip %s from diff ID %d: %w", ilk, pip, diffID, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertIlkMat(diffID, headerID int64, metadata types.ValueMetadata, mat string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return fmt.Errorf("error getting ilk for ilk mat: %w", err)
	}

	insertErr := shared.InsertFieldWithIlk(diffID, headerID, ilk, IlkMat, InsertSpotIlkMatQuery, mat, repository.db)
	if insertErr != nil {
		return fmt.Errorf("error inserting ilk %s mat %s from diff ID %d: %w", ilk, mat, diffID, insertErr)
	}
	return nil
}

func (repository StorageRepository) insertSpotVat(diffID, headerID int64, vat string) error {
	_, err := repository.db.Exec(insertSpotVatQuery, diffID, headerID, vat)
	if err != nil {
		return fmt.Errorf("error inserting spot vat %s from diff ID %d: %w", vat, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertSpotPar(diffID, headerID int64, par string) error {
	_, err := repository.db.Exec(insertSpotParQuery, diffID, headerID, par)
	if err != nil {
		return fmt.Errorf("error inserting spot par %s from diff ID %d: %w", par, diffID, err)
	}
	return nil
}

func (repository StorageRepository) insertSpotLive(diffID, headerID int64, live string) error {
	_, err := repository.db.Exec(insertSpotLiveQuery, diffID, headerID, live)
	if err != nil {
		return fmt.Errorf("error inserting spot live %s from diff ID %d: %w", live, diffID, err)
	}
	return nil
}

func getIlk(keys map[types.Key]string) (string, error) {
	ilk, ok := keys[constants.Ilk]
	if !ok {
		return "", types.ErrMetadataMalformed{MissingData: constants.Ilk}
	}
	return ilk, nil
}
