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

package initializers

import (
	mcdStorage "github.com/vulcanize/mcd_transformers/transformers/storage"
	"github.com/vulcanize/mcd_transformers/transformers/storage/flip"
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
)

func GenerateStorageTransformerInitializer(contractAddress string) transformer.StorageTransformerInitializer {
	return storage.Transformer{
		HashedAddress: utils.HexToKeccak256Hash(contractAddress),
		Mappings:      mcdStorage.NewKeysLookup(flip.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, contractAddress)),
		Repository:    &flip.FlipStorageRepository{ContractAddress: contractAddress},
	}.NewTransformer
}
