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

package chop_lump

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

var (
	chop = "chop"
	lump = "lump"
)

type CatFileChopLumpConverter struct{}

func (CatFileChopLumpConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var results []interface{}
	for _, ethLog := range ethLogs {
		verifyErr := verifyLog(ethLog)
		if verifyErr != nil {
			return nil, verifyErr
		}
		ilk := shared.GetHexWithoutPrefix(ethLog.Topics[2].Bytes())
		what := shared.DecodeHexToText(ethLog.Topics[3].Hex())
		dataBytes, parseErr := shared.GetLogNoteArgumentAtIndex(2, ethLog.Data)
		if parseErr != nil {
			return nil, parseErr
		}
		data := shared.ConvertUint256HexToBigInt(hexutil.Encode(dataBytes))

		raw, marshalErr := json.Marshal(ethLog)
		if marshalErr != nil {
			return nil, marshalErr
		}
		result := CatFileChopLumpModel{
			Ilk:              ilk,
			What:             what,
			Data:             data.String(),
			TransactionIndex: ethLog.TxIndex,
			LogIndex:         ethLog.Index,
			Raw:              raw,
		}
		results = append(results, result)
	}
	return results, nil
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	if len(log.Data) < constants.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}
