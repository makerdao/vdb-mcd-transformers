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
	"math"

	"github.com/ethereum/go-ethereum/common"
)

func GetLogNoteAddresses(arrayLength uint64, eventLogData []byte) ([]string, error) {
	startingIndex := 2
	maxVisibleAddresses := 4
	endIndex := startingIndex + int(math.Min(float64(maxVisibleAddresses), float64(arrayLength))-1)
	var addresses []string
	for i := startingIndex; i <= endIndex; i++ {
		logDataBytes, logNoteErr := GetLogNoteArgumentAtIndex(i, eventLogData)
		if logNoteErr != nil {
			return nil, logNoteErr
		}
		addressHex := common.BytesToAddress(logDataBytes).Hex()
		addresses = append(addresses, addressHex)
	}
	return addresses, nil
}
