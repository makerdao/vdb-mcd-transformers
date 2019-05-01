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
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	math2 "github.com/ethereum/go-ethereum/common/math"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
)

var ErrInvalidIndex = func(index int) error {
	return errors.New(fmt.Sprintf("unsupported log data index: %d", index))
}

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

func GetDSNoteThirdArgument(logData []byte) []byte {
	return getDataWithIndexOffset(1, logData)
}

func GetVatNoteDataBytesAtIndex(index int, logData []byte) ([]byte, error) {
	indexOffset, err := getVatNoteArgIndexOffset(index)
	if err != nil {
		return nil, err
	}
	return getDataWithIndexOffset(indexOffset, logData), nil
}

func getVatNoteArgIndexOffset(index int) (int, error) {
	minArgIndex := 4
	maxArgIndex := 6
	if index < minArgIndex || index > maxArgIndex {
		return 0, ErrInvalidIndex(index)
	}
	offsets := map[int]int{4: 3, 5: 2, 6: 1}
	return offsets[index], nil
}

func getDataWithIndexOffset(offset int, logData []byte) []byte {
	zeroPaddedSignatureOffset := 28
	dataBegin := len(logData) - (offset * constants.DataItemLength) - zeroPaddedSignatureOffset
	dataEnd := len(logData) - ((offset - 1) * constants.DataItemLength) - zeroPaddedSignatureOffset
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

func DecodeIlkName(hexIlk string) (string, error) {
	// Remove possible 0x intro
	nakedHex := GetHexWithoutPrefix(common.FromHex(hexIlk))
	ilkNamePadded, err := hex.DecodeString(nakedHex)
	if err != nil {
		return "", fmt.Errorf("couldn't convert hex ilk representation to string: %v", err)
	}
	ilkName := fmt.Sprintf("%s", bytes.Trim(ilkNamePadded, "\x00"))
	return ilkName, nil
}
