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
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/makerdao/vulcanizedb/libraries/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

var ErrInvalidIndex = func(index int) error {
	return errors.New(fmt.Sprintf("unsupported log data index: %d", index))
}

func BigIntToString(value *big.Int) string {
	result := value.String()
	if result == "<nil>" {
		return ""
	} else {
		return result
	}
}

func ConvertIntStringToHex(n string) (string, error) {
	b := big.NewInt(0)
	b, ok := b.SetString(n, 10)
	if !ok {
		return "", errors.New("error converting int to hex")
	}
	leftPaddedBytes := common.LeftPadBytes(b.Bytes(), 32)
	hex := common.Bytes2Hex(leftPaddedBytes)
	return hex, nil
}

func ConvertInt256HexToBigInt(hex string) *big.Int {
	n := ConvertUint256HexToBigInt(hex)
	return math.S256(n)
}

func ConvertUint256HexToBigInt(hex string) *big.Int {
	hexBytes := common.FromHex(hex)
	return big.NewInt(0).SetBytes(hexBytes)
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

func DecodeHexToText(payload string) string {
	return string(bytes.Trim(common.FromHex(payload), "\x00"))
}

func FormatRollbackError(field string, err error) error {
	return fmt.Errorf("failed to rollback transaction after failing to insert %s: %w", field, err)
}

func GetFullTableName(schema, table string) string {
	return schema + "." + table
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
			return FormatRollbackError("ilk", ilkErr)
		}
		return fmt.Errorf("error getting or creating ilk: %w", ilkErr)
	}
	_, writeErr := tx.Exec(query, diffID, headerID, ilkID, value)

	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return FormatRollbackError(variableName, writeErr)
		}
		return fmt.Errorf("error inserting field with ilk: %w", writeErr)
	}
	return tx.Commit()
}

func InsertFieldWithIlkAndAddress(diffID, headerID int64, address, ilk, variableName, query, value string, db *postgres.DB) error {
	tx, txErr := db.Beginx()
	if txErr != nil {
		return fmt.Errorf("error beginning transaction: %w", txErr)
	}

	addressID, addressErr := GetOrCreateAddress(address, db)
	if addressErr != nil {
		return fmt.Errorf("Could not retrieve address id for %s, error: %w", address, addressErr)
	}

	ilkID, ilkErr := GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return FormatRollbackError("ilk", ilkErr)
		}
		return fmt.Errorf("error getting or creating ilk: %w", ilkErr)
	}
	_, writeErr := tx.Exec(query, diffID, headerID, addressID, ilkID, value)

	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return FormatRollbackError(variableName, writeErr)
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

	addressId, addressErr := GetOrCreateAddressInTransaction(contractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return FormatRollbackError("address", addressErr)
		}
		return fmt.Errorf("error getting or creating address: %w", addressErr)
	}
	_, insertErr := tx.Exec(query, diffID, headerID, addressId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return FormatRollbackError("field with address", insertErr)
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
	addressId, addressErr := GetOrCreateAddressInTransaction(contractAddress, tx)
	if addressErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return FormatRollbackError("address", addressErr)
		}
		return addressErr
	}
	_, insertErr := tx.Exec(query, diffID, headerID, addressId, bidId, value)
	if insertErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			errorString := fmt.Sprintf("field with address for bid id %s", bidId)
			return FormatRollbackError(errorString, insertErr)
		}
		return insertErr
	}
	return tx.Commit()
}

func GetChecksumAddressString(address string) string {
	return common.HexToAddress(address).Hex()
}
