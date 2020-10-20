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

package shared

import (
	"errors"
	"fmt"

	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/repository"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var ErrInvalidIndex = func(index int) error {
	return errors.New(fmt.Sprintf("unsupported log data index: %d", index))
}

func GetLogNoteArgumentAtIndex(index int, logData []byte) ([]byte, error) {
	indexOffset, err := getLogNoteArgumentIndexOffset(index)
	if err != nil {
		return nil, err
	}
	return getDataWithIndexOffset(indexOffset, logData), nil
}

func getLogNoteArgumentIndexOffset(index int) (int, error) {
	minArgIndex := 2
	maxArgIndex := 5
	if index < minArgIndex || index > maxArgIndex {
		return 0, ErrInvalidIndex(index)
	}
	offsets := map[int]int{2: 4, 3: 3, 4: 2, 5: 1}
	return offsets[index], nil
}

func getDataWithIndexOffset(offset int, logData []byte) []byte {
	zeroPaddedSignatureOffset := 28
	dataBegin := len(logData) - (offset * constants.DataItemLength) - zeroPaddedSignatureOffset
	dataEnd := len(logData) - ((offset - 1) * constants.DataItemLength) - zeroPaddedSignatureOffset
	return logData[dataBegin:dataEnd]
}

func InsertFieldWithIlk(diffID, headerID int64, ilk, variableName, query, value string, db *postgres.DB) error {
	tx, txErr := db.Beginx()
	if txErr != nil {
		return fmt.Errorf("error beginning transaction: %w", txErr)
	}
	ilkID, ilkErr := GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("ilk", ilkErr)
		}
		return fmt.Errorf("error getting or creating ilk: %w", ilkErr)
	}
	_, writeErr := tx.Exec(query, diffID, headerID, ilkID, value)

	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError(variableName, writeErr)
		}
		return fmt.Errorf("error inserting field with ilk: %w", writeErr)
	}
	return tx.Commit()
}

func InsertFieldWithIlkAndAddress(diffID, headerID, addressID int64, ilk, variableName, query, value string, db *postgres.DB) error {
	tx, txErr := db.Beginx()
	if txErr != nil {
		return fmt.Errorf("error beginning transaction: %w", txErr)
	}

	ilkID, ilkErr := GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("ilk", ilkErr)
		}
		return fmt.Errorf("error getting or creating ilk: %w", ilkErr)
	}
	_, writeErr := tx.Exec(query, diffID, headerID, addressID, ilkID, value)

	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError(variableName, writeErr)
		}
		return fmt.Errorf("error inserting field with ilk: %w", writeErr)
	}
	return tx.Commit()
}

func InsertRecordWithAddress(diffID, headerID int64, query, value, contractAddress string, db *postgres.DB) error {
	tx, txErr := db.Beginx()
	if txErr != nil {
		return txErr
	}

	addressId, addressErr := repository.GetOrCreateAddressInTransaction(tx, contractAddress)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("address", addressErr)
		}
		return fmt.Errorf("error getting or creating address: %w", addressErr)
	}
	_, insertErr := tx.Exec(query, diffID, headerID, addressId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("field with address", insertErr)
		}
		return fmt.Errorf("error inserting record with address: %w", insertErr)
	}

	return tx.Commit()
}

func InsertRecordWithAddressAndBidID(diffID, headerID int64, query, bidId, value, contractAddress string, db *postgres.DB) error {
	tx, txErr := db.Beginx()
	if txErr != nil {
		return txErr
	}
	addressId, addressErr := repository.GetOrCreateAddressInTransaction(tx, contractAddress)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("address", addressErr)
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, diffID, headerID, addressId, bidId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			errorString := fmt.Sprintf("field with address for bid id %s", bidId)
			return shared.FormatRollbackError(errorString, insertErr)
		}
		return insertErr
	}
	return tx.Commit()
}
