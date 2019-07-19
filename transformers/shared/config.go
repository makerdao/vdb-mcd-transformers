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
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
)

// Creates a transformer config by pulling values from configuration environment
// TODO How can we get signatures in a general way?
func GetEventTransformerConfig(transformerLabel, signature string) transformer.EventTransformerConfig {
	contractNames := constants.GetTransformerContractNames(transformerLabel)
	return transformer.EventTransformerConfig{
		TransformerName:     transformerLabel,
		ContractAddresses:   constants.GetContractAddresses(contractNames),
		ContractAbi:         constants.GetContractsABI(contractNames),
		Topic:               signature,
		StartingBlockNumber: constants.GetMinDeploymentBlock(contractNames),
		EndingBlockNumber:   -1, // TODO Generalise endingBlockNumber
	}
}
