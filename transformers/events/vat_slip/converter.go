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

package vat_slip

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type VatSlipConverter struct{}

const (
	logDataRequired   = false
	numTopicsRequired = 4
)

func (VatSlipConverter) ToModels(_ string, logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	var models []shared.InsertionModel
	for _, log := range logs {
		err := shared.VerifyLog(log.Log, numTopicsRequired, logDataRequired)
		if err != nil {
			return nil, err
		}

		ilk := log.Log.Topics[1].Hex()
		usr := common.BytesToAddress(log.Log.Topics[2].Bytes()).String()
		wad := shared.ConvertInt256HexToBigInt(log.Log.Topics[3].Hex())

		model := shared.InsertionModel{
			SchemaName: "maker",
			TableName:  "vat_slip",
			OrderedColumns: []string{
				constants.HeaderFK, string(constants.IlkFK), "usr", "wad", constants.LogFK,
			},
			ColumnValues: shared.ColumnValues{
				"usr":              usr,
				"wad":              wad.String(),
				constants.HeaderFK: log.HeaderID,
				constants.LogFK:    log.ID,
			},
			ForeignKeyValues: shared.ForeignKeyValues{
				constants.IlkFK: ilk,
			},
		}
		models = append(models, model)
	}
	return models, nil
}
