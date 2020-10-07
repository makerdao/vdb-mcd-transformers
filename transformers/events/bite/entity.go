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

package bite

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type BiteEntity struct {
	Ilk             [32]byte
	Urn             common.Address
	Ink             *big.Int
	Art             *big.Int
	Tab             *big.Int
	Flip            common.Address
	Id              *big.Int
	ContractAddress common.Address
	HeaderID        int64
	LogID           int64
}
