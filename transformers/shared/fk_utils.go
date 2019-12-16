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

package shared

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	getOrCreateIlkQuery = `WITH insertedIlkId AS (
		INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id
		)
		SELECT id FROM maker.ilks WHERE ilk = $1
		UNION
		SELECT id FROM insertedIlkId`
	getOrCreateUrnQuery = `WITH insertedUrnId AS (
		INSERT INTO maker.urns (identifier, ilk_id) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING id
		)
		SELECT id FROM maker.urns WHERE identifier = $1 AND ilk_id = $2
		UNION
		SELECT id FROM insertedUrnId`
)

func GetOrCreateIlk(ilk string, db *postgres.DB) (int64, error) {
	var ilkID int64
	uniformIlk := common.HexToHash(ilk).Hex()
	ilkIdentifier := DecodeHexToText(uniformIlk)
	err := db.Get(&ilkID, getOrCreateIlkQuery, uniformIlk, ilkIdentifier)
	return ilkID, err
}

func GetOrCreateIlkInTransaction(ilk string, tx *sqlx.Tx) (int64, error) {
	var ilkID int64
	uniformIlk := common.HexToHash(ilk).Hex()
	ilkIdentifier := DecodeHexToText(uniformIlk)
	err := tx.Get(&ilkID, getOrCreateIlkQuery, uniformIlk, ilkIdentifier)
	return ilkID, err
}

func GetOrCreateUrn(guy string, hexIlk string, db *postgres.DB) (urnID int64, err error) {
	ilkID, ilkErr := GetOrCreateIlk(hexIlk, db)
	if ilkErr != nil {
		return 0, fmt.Errorf("error getting ilkID for urn: %s", ilkErr.Error())
	}

	err = db.Get(&urnID, getOrCreateUrnQuery, guy, ilkID)
	return urnID, err
}

func GetOrCreateUrnInTransaction(guy string, hexIlk string, tx *sqlx.Tx) (urnID int64, err error) {
	ilkID, ilkErr := GetOrCreateIlkInTransaction(hexIlk, tx)
	if ilkErr != nil {
		return 0, fmt.Errorf("error getting ilkID for urn: %v", ilkErr.Error())
	}

	err = tx.Get(&urnID, getOrCreateUrnQuery, guy, ilkID)
	return urnID, err
}

func GetOrCreateAddress(address string, db *postgres.DB) (int64, error) {
	return repository.GetOrCreateAddress(db, address)
}

func GetOrCreateAddressInTransaction(address string, tx *sqlx.Tx) (int64, error) {
	addressId, addressErr := repository.GetOrCreateAddressInTransaction(tx, address)
	return addressId, addressErr
}
