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
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
)

// Storage keys
const (
	A         types.Key = "a"
	Address   types.Key = "address"
	BidId     types.Key = "bid_id"
	Cdpi      types.Key = "cdpi"
	Flip      types.Key = "flip"
	Guy       types.Key = "guy"
	Ilk       types.Key = "ilk"
	MsgSender types.Key = "msg_sender"
	Owner     types.Key = "owner"
	SaleId    types.Key = "sale_id"
	SlotId    types.Key = "slot_id"
	Timestamp types.Key = "timestamp"
	User      types.Key = "usr"
)

// Table column names
const (
	AColumn         event.ColumnName = "a"
	ALengthColumn   event.ColumnName = "a_length"
	AgeColumn       event.ColumnName = "age"
	ArtColumn       event.ColumnName = "art"
	BidColumn       event.ColumnName = "bid"
	BidIDColumn     event.ColumnName = "bid_id"
	CdpColumn       event.ColumnName = "cdp"
	ClipColumn      event.ColumnName = "clip"
	ClipIDColumn    event.ColumnName = "clip_id"
	CoinColumn      event.ColumnName = "coin"
	DartColumn      event.ColumnName = "dart"
	DataColumn      event.ColumnName = "data"
	DinkColumn      event.ColumnName = "dink"
	DstColumn       event.ColumnName = "dst"
	DueColumn       event.ColumnName = "due"
	EndColumn       event.ColumnName = "end"
	EraColumn       event.ColumnName = "era"
	FlipColumn      event.ColumnName = "flip"
	GalColumn       event.ColumnName = "gal"
	GuyColumn       event.ColumnName = "guy"
	IlkColumn       event.ColumnName = "ilk_id"
	InkColumn       event.ColumnName = "ink"
	KprColumn       event.ColumnName = "kpr"
	LotColumn       event.ColumnName = "lot"
	MaxColumn       event.ColumnName = "max"
	MsgSenderColumn event.ColumnName = "msg_sender"
	OweColumn       event.ColumnName = "owe"
	OwnColumn       event.ColumnName = "own"
	PipColumn       event.ColumnName = "pip"
	PriceColumn     event.ColumnName = "price"
	RadColumn       event.ColumnName = "rad"
	RateColumn      event.ColumnName = "rate"
	SaleIDColumn    event.ColumnName = "sale_id"
	SpotColumn      event.ColumnName = "spot"
	SrcColumn       event.ColumnName = "src"
	TabColumn       event.ColumnName = "tab"
	TicColumn       event.ColumnName = "tic"
	TopColumn       event.ColumnName = "top"
	UColumn         event.ColumnName = "u"
	UrnColumn       event.ColumnName = "urn_id"
	UsrColumn       event.ColumnName = "usr"
	VColumn         event.ColumnName = "v"
	ValColumn       event.ColumnName = "val"
	WColumn         event.ColumnName = "w"
	WadColumn       event.ColumnName = "wad"
	WhatColumn      event.ColumnName = "what"
	ValueColumn     event.ColumnName = "value"
)
