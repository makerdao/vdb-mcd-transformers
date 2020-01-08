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
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
)

// Storage keys
const (
	BidId     storage.Key = "bid_id"
	Cdpi      storage.Key = "cdpi"
	Flip      storage.Key = "flip"
	Guy       storage.Key = "guy"
	Ilk       storage.Key = "ilk"
	MsgSender storage.Key = "msg_sender"
	Owner     storage.Key = "owner"
	Timestamp storage.Key = "timestamp"
	User      storage.Key = "usr"
)

// Table column names
const (
	ArtColumn       event.ColumnName = "art"
	BidColumn       event.ColumnName = "bid"
	BidIDColumn     event.ColumnName = "bid_id"
	CdpColumn       event.ColumnName = "cdp"
	DartColumn      event.ColumnName = "dart"
	DataColumn      event.ColumnName = "data"
	DinkColumn      event.ColumnName = "dink"
	DstColumn       event.ColumnName = "dst"
	EraColumn       event.ColumnName = "era"
	FlipColumn      event.ColumnName = "flip"
	GalColumn       event.ColumnName = "gal"
	IlkColumn       event.ColumnName = "ilk_id"
	InkColumn       event.ColumnName = "ink"
	LotColumn       event.ColumnName = "lot"
	MsgSenderColumn event.ColumnName = "msg_sender"
	OwnColumn       event.ColumnName = "own"
	PipColumn       event.ColumnName = "pip"
	RadColumn       event.ColumnName = "rad"
	RateColumn      event.ColumnName = "rate"
	SpotColumn      event.ColumnName = "spot"
	SrcColumn       event.ColumnName = "src"
	TabColumn       event.ColumnName = "tab"
	UColumn         event.ColumnName = "u"
	UrnColumn       event.ColumnName = "urn_id"
	UsrColumn       event.ColumnName = "usr"
	VColumn         event.ColumnName = "v"
	WColumn         event.ColumnName = "w"
	WadColumn       event.ColumnName = "wad"
	WhatColumn      event.ColumnName = "what"
	ValueColumn     event.ColumnName = "value"
)
