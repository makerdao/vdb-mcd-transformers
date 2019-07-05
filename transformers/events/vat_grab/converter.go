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

package vat_grab

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	constants2 "github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type VatGrabConverter struct{}

func (VatGrabConverter) ToModels(ethLogs []types.Log) ([]shared.InsertionModel, error) {
	var models []shared.InsertionModel
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}
		ilk := ethLog.Topics[1].Hex()
		urn := common.BytesToAddress(ethLog.Topics[2].Bytes()).String()
		v := common.BytesToAddress(ethLog.Topics[3].Bytes()).String()
		wBytes, wErr := shared.GetLogNoteArgumentAtIndex(3, ethLog.Data)
		if wErr != nil {
			return nil, wErr
		}
		w := common.BytesToAddress(wBytes).String()
		dinkBytes, dinkErr := shared.GetLogNoteArgumentAtIndex(4, ethLog.Data)
		if dinkErr != nil {
			return nil, dinkErr
		}
		dink := shared.ConvertInt256HexToBigInt(hexutil.Encode(dinkBytes))
		dartBytes, dartErr := shared.GetLogNoteArgumentAtIndex(5, ethLog.Data)
		if dartErr != nil {
			return nil, dartErr
		}
		dart := shared.ConvertInt256HexToBigInt(hexutil.Encode(dartBytes))

		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}
		model := shared.InsertionModel{
			TableName: "vat_grab",
			OrderedColumns: []string{
				"header_id", string(constants2.UrnFK), "v", "w", "dink", "dart", "log_idx", "tx_idx", "raw_log",
			},
			ColumnValues: shared.ColumnValues{
				"v":       v,
				"w":       w,
				"dink":    dink.String(),
				"dart":    dart.String(),
				"log_idx": ethLog.Index,
				"tx_idx":  ethLog.TxIndex,
				"raw_log": raw,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants2.IlkFK: ilk,
				constants2.UrnFK: urn,
			},
		}
		models = append(models, model)
	}
	return models, nil
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
