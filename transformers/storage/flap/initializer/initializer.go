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

package initializer

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	mcdStorage "github.com/makerdao/vdb-mcd-transformers/transformers/storage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flap"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/storage"
)

var StorageTransformerInitializer storage.TransformerInitializer = storage.Transformer{
	Address: common.HexToAddress(constants.GetContractAddress("MCD_FLAP")),
	StorageKeysLookup: storage.NewKeysLookup(flap.NewKeysLoader(
		&mcdStorage.MakerStorageRepository{},
		constants.GetContractAddress("MCD_FLAP"))),
	Repository: &flap.StorageRepository{ContractAddress: constants.GetContractAddress("MCD_FLAP")},
}.NewTransformer
