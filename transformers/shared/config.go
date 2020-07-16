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
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

// Creates a transformer config by pulling values from configuration environment
func GetEventTransformerConfig(transformerLabel, signature string) event.TransformerConfig {
	contractNames := constants.GetTransformerContractNames(transformerLabel)
	return event.TransformerConfig{
		TransformerName:     transformerLabel,
		ContractAddresses:   constants.GetContractAddresses(contractNames),
		ContractAbi:         constants.GetFirstABI(contractNames),
		Topic:               signature,
		StartingBlockNumber: constants.GetMinDeploymentBlock(contractNames),
		EndingBlockNumber:   -1, // TODO Generalise endingBlockNumber
	}
}
