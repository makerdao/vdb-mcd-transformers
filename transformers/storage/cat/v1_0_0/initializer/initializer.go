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

package initializer

import (
	"github.com/ethereum/go-ethereum/common"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat/v1_0_0"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
)

var catAddress = constants.GetContractAddress("MCD_CAT_1_0_0")
var StorageTransformerInitializer storage.TransformerInitializer = storage.Transformer{
	Address:           common.HexToAddress(catAddress),
	StorageKeysLookup: storage.NewKeysLookup(v1_0_0.NewKeysLoader(&mcdStorage.MakerStorageRepository{}, catAddress)),
	Repository:        &cat.StorageRepository{ContractAddress: catAddress},
}.NewTransformer
