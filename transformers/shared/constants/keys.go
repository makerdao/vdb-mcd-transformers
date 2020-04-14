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
	BidId     types.Key = "bid_id"
	Cdpi      types.Key = "cdpi"
	Flip      types.Key = "flip"
	Guy       types.Key = "guy"
	Ilk       types.Key = "ilk"
	MsgSender types.Key = "msg_sender"
	Owner     types.Key = "owner"
	Timestamp types.Key = "timestamp"
	User      types.Key = "usr"
)

// Table column names
const (
	AColumn         event.ColumnName = "a"
	ALengthColumn   event.ColumnName = "a_length"
	ArtColumn       event.ColumnName = "art"
	BidColumn       event.ColumnName = "bid"
	BidIDColumn     event.ColumnName = "bid_id"
	BuyAmtColumn    event.ColumnName = "buy_amt"
	BuyGemColumn    event.ColumnName = "buy_gem"
	CdpColumn       event.ColumnName = "cdp"
	DartColumn      event.ColumnName = "dart"
	DataColumn      event.ColumnName = "data"
	DinkColumn      event.ColumnName = "dink"
	DstColumn       event.ColumnName = "dst"
	EndColumn       event.ColumnName = "end"
	EraColumn       event.ColumnName = "era"
	FlipColumn      event.ColumnName = "flip"
	GalColumn       event.ColumnName = "gal"
	GuyColumn       event.ColumnName = "guy"
	IlkColumn       event.ColumnName = "ilk_id"
	InkColumn       event.ColumnName = "ink"
	LotColumn       event.ColumnName = "lot"
	MakerColumn     event.ColumnName = "maker"
	MsgSenderColumn event.ColumnName = "msg_sender"
	OfferId         event.ColumnName = "offer_id"
	OwnColumn       event.ColumnName = "own"
	PairColumn      event.ColumnName = "pair"
	PayAmtColumn    event.ColumnName = "pay_amt"
	PayGemColumn    event.ColumnName = "pay_gem"
	PipColumn       event.ColumnName = "pip"
	RadColumn       event.ColumnName = "rad"
	RateColumn      event.ColumnName = "rate"
	SpotColumn      event.ColumnName = "spot"
	SrcColumn       event.ColumnName = "src"
	TabColumn       event.ColumnName = "tab"
	TicColumn       event.ColumnName = "tic"
	TimestampColumn event.ColumnName = "timestamp"
	UColumn         event.ColumnName = "u"
	UrnColumn       event.ColumnName = "urn_id"
	UsrColumn       event.ColumnName = "usr"
	VColumn         event.ColumnName = "v"
	WColumn         event.ColumnName = "w"
	WadColumn       event.ColumnName = "wad"
	WhatColumn      event.ColumnName = "what"
	ValueColumn     event.ColumnName = "value"
)
