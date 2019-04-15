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

package pip_log_value

import (
	shared_t "github.com/vulcanize/vulcanizedb/libraries/shared/transformer"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

func GetPipLogValueConfig() shared_t.EventTransformerConfig {
	return shared_t.EventTransformerConfig{
		TransformerName: constants.PipLogValueLabel,
		ContractAddresses: []string{
			constants.PipEthContractAddress(),
			constants.PipCol1ContractAddress(),
			constants.PipCol2ContractAddress(),
			constants.PipCol3ContractAddress(),
			constants.PipCol4ContractAddress(),
			constants.PipCol5ContractAddress(),
		},
		ContractAbi: constants.PipABI(),
		Topic:       constants.GetPipLogValueSignature(),
		StartingBlockNumber: shared.MinInt64([]int64{
			constants.PipEthDeploymentBlock(),
			constants.PipCol1DeploymentBlock(),
			constants.PipCol2DeploymentBlock(),
			constants.PipCol3DeploymentBlock(),
			constants.PipCol4DeploymentBlock(),
			constants.PipCol5DeploymentBlock(),
		}),
		EndingBlockNumber: -1,
	}
}
