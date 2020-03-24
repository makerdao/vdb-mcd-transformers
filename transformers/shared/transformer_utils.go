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
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vulcanizedb/libraries/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	TwoTopicsRequired   = 2
	ThreeTopicsRequired = 3
	FourTopicsRequired  = 4
	LogDataRequired     = true
	LogDataNotRequired  = false
	emptyLogDataSlot    = "0x0000000000000000000000000000000000000000"
	lastDataIndex       = 5
)

var (
	ErrLogMissingTopics = func(expectedNumTopics, actualNumTopics int) error {
		return fmt.Errorf("log missing topics: has %d, want %d", actualNumTopics, expectedNumTopics)
	}
	ErrLogMissingData   = errors.New("log missing data")
	ErrCouldNotCreateFK = func(err error) error {
		return fmt.Errorf("transformer could not create FK: %v", err)
	}
)

func VerifyLog(log types.Log, expectedNumTopics int, isDataRequired bool) error {
	actualNumTopics := len(log.Topics)
	if actualNumTopics < expectedNumTopics {
		return ErrLogMissingTopics(expectedNumTopics, actualNumTopics)
	}
	if isDataRequired && len(log.Data) < constants.DataItemLength {
		return ErrLogMissingData
	}
	return nil
}

func GetLogNoteData(startingIndex int, eventLogData []byte, db *postgres.DB) ([]int64, error) {
	var addressIDs []int64
	for j := startingIndex; j <= lastDataIndex; j++ {
		logDataBytes, _ := GetLogNoteArgumentAtIndex(j, eventLogData)
		logData := common.BytesToAddress(logDataBytes).String()
		addressID, addressErr := GetOrCreateAddress(logData, db)
		if addressErr != nil {
			return nil, addressErr
		}
		if logData == emptyLogDataSlot {
			addressIDs = append(addressIDs, 0)
		} else {
			addressIDs = append(addressIDs, addressID)
		}

	}
	return addressIDs, nil
}
