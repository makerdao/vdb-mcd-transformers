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

package constants

const (
	MakerSchema  = "maker"
	APISchema    = "api"
	PublicSchema = "public"
)

// event tables
const (
	BiteTable               = "bite"
	CatFileChopLumpTable    = "cat_file_chop_lump"
	CatFileFlipTable        = "cat_file_flip"
	CatFileVowTable         = "cat_file_vow"
	DealTable               = "deal"
	DentTable               = "dent"
	FlapKickTable           = "flap_kick"
	FlipKickTable           = "flip_kick"
	FlopKickTable           = "flop_kick"
	JugDripTable            = "jug_drip"
	JugFileBaseTable        = "jug_file_base"
	JugFileIlkTable         = "jug_file_ilk"
	JugFileVowTable         = "jug_file_vow"
	JugInitTable            = "jug_init"
	LogValueTable           = "log_value"
	NewCdpTable             = "new_cdp"
	PotCageTable            = "pot_cage"
	PotDripTable            = "pot_drip"
	PotExitTable            = "pot_exit"
	PotFileDSRTable         = "pot_file_dsr"
	PotFileVowTable         = "pot_file_vow"
	PotJoinTable            = "pot_join"
	SpotFileMatTable        = "spot_file_mat"
	SpotFilePipTable        = "spot_file_pip"
	SpotFileParTable        = "spot_file_par"
	SpotPokeTable           = "spot_poke"
	TendTable               = "tend"
	TickTable               = "tick"
	VatFileDebtCeilingTable = "vat_file_debt_ceiling"
	VatFileIlkTable         = "vat_file_ilk"
	VatFluxTable            = "vat_flux"
	VatFoldTable            = "vat_fold"
	VatForkTable            = "vat_fork"
	VatFrobTable            = "vat_frob"
	VatGrabTable            = "vat_grab"
	VatHealTable            = "vat_heal"
	VatInitTable            = "vat_init"
	VatMoveTable            = "vat_move"
	VatSlipTable            = "vat_slip"
	VatSuckTable            = "vat_suck"
	VowFessTable            = "vow_fess"
	VowFileTable            = "vow_file"
	VowFlogTable            = "vow_flog"
	YankTable               = "yank"
)

// storage tables
const (
	CdpManagerVat      = "cdp_manager_vat"
	CdpManagerCdpi     = "cdp_manager_cdpi"
	CdpManagerUrns     = "cdp_manager_urns"
	CdpManagerListPrev = "cdp_manager_list_prev"
	CdpManagerListNext = "cdp_manager_list_next"
	CdpManagerOwns     = "cdp_manager_owns"
	CdpManagerFirst    = "cdp_manager_first"
	CdpManagerLast     = "cdp_manager_last"
	CdpManagerCount    = "cdp_manager_count"
	CdpManagerIlks     = "cdp_manager_ilks"

	CatLive    = "cat_live"
	CatVat     = "cat_vat"
	CatVow     = "cat_vow"
	CatIlkFlip = "cat_ilk_flip"
	CatIlkChop = "cat_ilk_chop"
	CatIlkLump = "cat_ilk_lump"

	FlapGem    = "flap_gem"
	FlapVat    = "flap_vat"
	FlapBeg    = "flap_beg"
	FlapLive   = "flap_live"
	FlapKicks  = "flap_kicks"
	FlapBidBid = "flap_bid_bid"
	FlapBidLot = "flap_bid_lot"

	FlipBidGal = "flip_bid_gal"
	FlipBidUsr = "flip_bid_usr"
	FlipBidTab = "flip_bid_tab"
	FlipVat    = "flip_vat"
	FlipBeg    = "flip_beg"
	FlipBidLot = "flip_bid_lot"
	FlipKicks  = "flip_kicks"
	FlipBidBid = "flip_bid_bid"

	FlopVat    = "flop_vat"
	FlopGem    = "flop_gem"
	FlopBeg    = "flop_beg"
	FlopPad    = "flop_pad"
	FlopKicks  = "flop_kicks"
	FlopLive   = "flop_live"
	FlopBidBid = "flop_bid_bid"
	FlopBidLot = "flop_bid_lot"

	JugIlkRho  = "jug_ilk_rho"
	JugIlkDuty = "jug_ilk_duty"

	SpotIlkPip = "spot_ilk_pip"
	SpotIlkMat = "spot_ilk_mat"

	VatIlkArt  = "vat_ilk_art"
	VatIlkDust = "vat_ilk_dust"
	VatIlkRate = "vat_ilk_rate"
	VatIlkLine = "vat_ilk_line"
	VatIlkSpot = "vat_ilk_spot"
)
