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

package constants

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/factories/event"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
)

const (
	BidId     utils.Key = "bid_id"
	Ilk       utils.Key = "ilk"
	Guy       utils.Key = "guy"
	Flip      utils.Key = "flip"
	Timestamp utils.Key = "timestamp"
	Cdpi      utils.Key = "cdpi"
	Owner     utils.Key = "owner"
)

// TODO remove after transition to ColumnName
type ForeignKeyField string

const (
	IlkFK     ForeignKeyField = "ilk_id"
	UrnFK     ForeignKeyField = "urn_id"
	AddressFK ForeignKeyField = "address_id"
)

const (
	HeaderFK = "header_id"
	LogFK    = "log_id"
)

const (
	IlkColumn     event.ColumnName = "ilk_id"
	UrnColumn     event.ColumnName = "urn_id"
	AddressColumn event.ColumnName = "address_id"
)
