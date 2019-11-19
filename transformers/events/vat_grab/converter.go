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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

type VatGrabConverter struct{}

const (
	logDataRequired   = true
	numTopicsRequired = 4
)

func (VatGrabConverter) ToModels(_ string, logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	var models []shared.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if err != nil {
			return nil, err
		}
		ilk := log.Log.Topics[1].Hex()
		urn := common.BytesToAddress(log.Log.Topics[2].Bytes()).String()
		v := common.BytesToAddress(log.Log.Topics[3].Bytes()).String()
		wBytes, wErr := shared.GetLogNoteArgumentAtIndex(3, log.Log.Data)
		if wErr != nil {
			return nil, wErr
		}
		w := common.BytesToAddress(wBytes).String()
		dinkBytes, dinkErr := shared.GetLogNoteArgumentAtIndex(4, log.Log.Data)
		if dinkErr != nil {
			return nil, dinkErr
		}
		dink := shared.ConvertInt256HexToBigInt(hexutil.Encode(dinkBytes))
		dartBytes, dartErr := shared.GetLogNoteArgumentAtIndex(5, log.Log.Data)
		if dartErr != nil {
			return nil, dartErr
		}
		dart := shared.ConvertInt256HexToBigInt(hexutil.Encode(dartBytes))

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "vat_grab",
			OrderedColumns: []string{
				constants.HeaderFK, string(constants.UrnFK), "v", "w", "dink", "dart", constants.LogFK,
			},
			ColumnValues: shared.ColumnValues{
				"v":                v,
				"w":                w,
				"dink":             dink.String(),
				"dart":             dart.String(),
				constants.HeaderFK: log.HeaderID,
				constants.LogFK:    log.ID,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.IlkFK: ilk,
				constants.UrnFK: urn,
			},
		}
		models = append(models, model)
	}
	return models, nil
}
