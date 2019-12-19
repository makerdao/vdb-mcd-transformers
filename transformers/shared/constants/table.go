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
	CatLiveTable    = "cat_live"
	CatVatTable     = "cat_vat"
	CatVowTable     = "cat_vow"
	CatIlkFlipTable = "cat_ilk_flip"
	CatIlkChopTable = "cat_ilk_chop"
	CatIlkLumpTable = "cat_ilk_lump"

	CdpManagerVatTable      = "cdp_manager_vat"
	CdpManagerCdpiTable     = "cdp_manager_cdpi"
	CdpManagerUrnsTable     = "cdp_manager_urns"
	CdpManagerListPrevTable = "cdp_manager_list_prev"
	CdpManagerListNextTable = "cdp_manager_list_next"
	CdpManagerOwnsTable     = "cdp_manager_owns"
	CdpManagerFirstTable    = "cdp_manager_first"
	CdpManagerLastTable     = "cdp_manager_last"
	CdpManagerCountTable    = "cdp_manager_count"
	CdpManagerIlksTable     = "cdp_manager_ilks"

	FlapGemTable    = "flap_gem"
	FlapVatTable    = "flap_vat"
	FlapBegTable    = "flap_beg"
	FlapLiveTable   = "flap_live"
	FlapKicksTable  = "flap_kicks"
	FlapBidBidTable = "flap_bid_bid"
	FlapBidEndTable = "flap_bid_end"
	FlapBidGuyTable = "flap_bid_guy"
	FlapBidLotTable = "flap_bid_lot"
	FlapBidTicTable = "flap_bid_tic"
	FlapTtlTable    = "flap_ttl"
	FlapTauTable    = "flap_tau"

	FlipBidGalTable = "flip_bid_gal"
	FlipBidUsrTable = "flip_bid_usr"
	FlipBidTabTable = "flip_bid_tab"
	FlipVatTable    = "flip_vat"
	FlipBegTable    = "flip_beg"
	FlipBidLotTable = "flip_bid_lot"
	FlipKicksTable  = "flip_kicks"
	FlipBidBidTable = "flip_bid_bid"
	FlipBidEndTable = "flip_bid_end"
	FlipBidGuyTable = "flip_bid_guy"
	FlipBidTicTable = "flip_bid_tic"
	FlipIlkTable    = "flip_ilk"
	FlipTtlTable    = "flip_ttl"
	FlipTauTable    = "flip_tau"

	FlopVatTable    = "flop_vat"
	FlopGemTable    = "flop_gem"
	FlopBegTable    = "flop_beg"
	FlopPadTable    = "flop_pad"
	FlopKicksTable  = "flop_kicks"
	FlopLiveTable   = "flop_live"
	FlopBidBidTable = "flop_bid_bid"
	FlopBidEndTable = "flop_bid_end"
	FlopBidGuyTable = "flop_bid_guy"
	FlopBidTicTable = "flop_bid_tic"
	FlopBidLotTable = "flop_bid_lot"
	FlopTauTable    = "flop_tau"
	FlopTtlTable    = "flop_ttl"

	JugIlkRhoTable  = "jug_ilk_rho"
	JugIlkDutyTable = "jug_ilk_duty"
	JugBaseTable    = "jug_base"
	JugVatTable     = "jug_vat"
	JugVowTable     = "jug_vow"

	SpotIlkPipTable = "spot_ilk_pip"
	SpotIlkMatTable = "spot_ilk_mat"
	SpotVatTable    = "spot_vat"
	SpotParTable    = "spot_par"
	SpotLiveTable   = "spot_live"

	VatIlkArtTable  = "vat_ilk_art"
	VatIlkDustTable = "vat_ilk_dust"
	VatIlkRateTable = "vat_ilk_rate"
	VatIlkLineTable = "vat_ilk_line"
	VatIlkSpotTable = "vat_ilk_spot"
	VatDebtTable    = "vat_debt"
	VatSinTable     = "vat_sin"
	VatViceTable    = "vat_vice"
	VatLiveTable    = "vat_live"
	VatLineTable    = "vat_line"
	VatDaiTable     = "vat_dai"
	VatGemTable     = "vat_gem"
	VatUrnArtTable  = "vat_urn_art"
	VatUrnInkTable  = "vat_urn_ink"

	VowBumpTable       = "vow_bump"
	VowHumpTable       = "vow_hump"
	VowDumpTable       = "vow_dump"
	VowSumpTable       = "vow_sump"
	VowWaitTable       = "vow_wait"
	VowVatTable        = "vow_vat"
	VowFlapperTable    = "vow_flapper"
	VowFlopperTable    = "vow_flopper"
	VowSinMappingTable = "vow_sin_mapping"
	VowSinIntegerTable = "vow_sin_integer"
	VowAshTable        = "vow_ash"
)
