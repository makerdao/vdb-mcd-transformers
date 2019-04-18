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
	"encoding/binary"
	"math"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	math2 "github.com/ethereum/go-ethereum/common/math"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
)

func BigIntToInt64(value *big.Int) int64 {
	if value == nil {
		return int64(0)
	} else {
		return value.Int64()
	}
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
	strAsInt, err := strconv.Atoi(n)
	if err != nil {
		return "", err
	}
	return ConvertIntToHex(strAsInt)
}

func ConvertIntToHex(n int) (string, error) {
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.BigEndian, uint64(n))
	if err != nil {
		return "", err
	}
	leftPaddedBytes := common.LeftPadBytes(b.Bytes(), 32)
	hex := common.Bytes2Hex(leftPaddedBytes)
	return hex, nil
}

func ConvertInt256HexToBigInt(hex string) *big.Int {
	n := ConvertUint256HexToBigInt(hex)
	return math2.S256(n)
}

func ConvertUint256HexToBigInt(hex string) *big.Int {
	hexBytes := common.FromHex(hex)
	return big.NewInt(0).SetBytes(hexBytes)
}

// Extract relevant bytes from log data emitted by DSNote and Vat Note modifiers.
// For DSNote, index is backward from last argument. For example,
//		- if 4 arguments:
//			- topic 0 is function signature
//			- topic 1 is msg.sender
//			- topic 2 is argument 1
//			- topic 3 is argument 2
//			- argument 3 can be accessed via GetLogNoteDataBytesAtIndex(-2, logData)
//			- argument 4 can be accessed via GetLogNoteDataBytesAtIndex(-1, logData)
//		- if 6 arguments:
//			- topics 0-3 are same as above
//			- argument 3 can be accessed via GetLogNoteDataBytesAtIndex(-3, logData)
//			- argument 4 can be accessed via GetLogNoteDataBytesAtIndex(-2, logData)
//			- argument 5 can be accessed via GetLogNoteDataBytesAtIndex(-1, logData)
// For Vat Note, note is padded at fixed length supporting 6 arguments. For example,
//		- if 4 arguments:
//			- topic 0 is function signature
//			- topic 1 is argument 1
//			- topic 2 is argument 2
//			- topic 3 is argument 3
//			- argument 4 can be accessed via GetLogNoteDataBytesAtIndex(-3, logData)
//		- if 6 arguments:
//			- topics 0-3 are same as above
//			- argument 4 can be accessed via GetLogNoteDataBytesAtIndex(-3, logData)
//			- argument 5 can be accessed via GetLogNoteDataBytesAtIndex(-2, logData)
//			- argument 6 can be accessed via GetLogNoteDataBytesAtIndex(-1, logData)
func GetLogNoteDataBytesAtIndex(n int, logData []byte) []byte {
	zeroPaddedSignatureOffset := 28
	indexOffset := int(math.Abs(float64(n)))
	dataBegin := len(logData) - (indexOffset * constants.DataItemLength) - zeroPaddedSignatureOffset
	dataEnd := len(logData) - ((indexOffset - 1) * constants.DataItemLength) - zeroPaddedSignatureOffset
	return logData[dataBegin:dataEnd]
}

func GetHexWithoutPrefix(raw []byte) string {
	return common.Bytes2Hex(raw)
}

func MinInt64(ints []int64) (min int64) {
	if len(ints) == 0 {
		return 0
	}
	min = ints[0]
	for _, i := range ints {
		if i < min {
			min = i
		}
	}
	return
}
