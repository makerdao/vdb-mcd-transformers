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
	AuctionFileTable              = "auction_file"
	BiteTable                     = "bite"
	CatClawTable                  = "cat_claw"
	CatFileBoxTable               = "cat_file_box"
	CatFileChopLumpDunkTable      = "cat_file_chop_lump_dunk"
	CatFileFlipTable              = "cat_file_flip"
	CatFileVowTable               = "cat_file_vow"
	ClipIlkTable                  = "clip_ilk"
	ClipKickTable                 = "clip_kick"
	ClipTakeTable                 = "clip_take"
	ClipRedoTable                 = "clip_redo"
	ClipVatTable                  = "clip_vat"
	ClipYankTable                 = "clip_yank"
	DogBarkTable                  = "dog_bark"
	DogDenyTable                  = "dog_deny"
	DogRelyTable                  = "dog_rely"
	DealTable                     = "deal"
	DentTable                     = "dent"
	DenyTable                     = "deny"
	DogDigsTable                  = "dog_digs"
	DogFileHoleTable              = "dog_file_hole"
	DogFileIlkClipTable           = "dog_file_ilk_clip"
	DogFileIlkChopHoleTable       = "dog_file_ilk_chop_hole"
	DogFileVowTable               = "dog_file_vow"
	FlapKickTable                 = "flap_kick"
	FlipFileCatTable              = "flip_file_cat"
	FlipKickTable                 = "flip_kick"
	FlopKickTable                 = "flop_kick"
	JugDripTable                  = "jug_drip"
	JugFileBaseTable              = "jug_file_base"
	JugFileIlkTable               = "jug_file_ilk"
	JugFileVowTable               = "jug_file_vow"
	JugInitTable                  = "jug_init"
	LogMedianPriceTable           = "log_median_price"
	LogValueTable                 = "log_value"
	MedianDissBatchTable          = "median_diss_batch"
	MedianDissSingleTable         = "median_diss_single"
	MedianDropTable               = "median_drop"
	MedianKissBatchTable          = "median_kiss_batch"
	MedianKissSingleTable         = "median_kiss_single"
	MedianLiftTable               = "median_lift"
	NewCdpTable                   = "new_cdp"
	OsmChangeTable                = "osm_change"
	PotCageTable                  = "pot_cage"
	PotDripTable                  = "pot_drip"
	PotExitTable                  = "pot_exit"
	PotFileDSRTable               = "pot_file_dsr"
	PotFileVowTable               = "pot_file_vow"
	PotJoinTable                  = "pot_join"
	RelyTable                     = "rely"
	SpotFileMatTable              = "spot_file_mat"
	SpotFileParTable              = "spot_file_par"
	SpotFilePipTable              = "spot_file_pip"
	SpotPokeTable                 = "spot_poke"
	TendTable                     = "tend"
	TickTable                     = "tick"
	VatDenyTable                  = "vat_deny"
	VatFileDebtCeilingTable       = "vat_file_debt_ceiling"
	VatFileIlkTable               = "vat_file_ilk"
	VatFluxTable                  = "vat_flux"
	VatFoldTable                  = "vat_fold"
	VatForkTable                  = "vat_fork"
	VatFrobTable                  = "vat_frob"
	VatGrabTable                  = "vat_grab"
	VatHealTable                  = "vat_heal"
	VatHopeTable                  = "vat_hope"
	VatInitTable                  = "vat_init"
	VatMoveTable                  = "vat_move"
	VatNopeTable                  = "vat_nope"
	VatRelyTable                  = "vat_rely"
	VatSlipTable                  = "vat_slip"
	VatSuckTable                  = "vat_suck"
	VowFessTable                  = "vow_fess"
	VowFileAuctionAttributesTable = "vow_file_auction_attributes"
	VowFileAuctionAddressTable    = "vow_file_auction_address"
	VowFlogTable                  = "vow_flog"
	VowHealTable                  = "vow_heal"
	YankTable                     = "yank"
)

// storage tables
const (
	CatBoxTable             = "cat_box"
	CatIlkChopTable         = "cat_ilk_chop"
	CatIlkDunkTable         = "cat_ilk_dunk"
	CatIlkFlipTable         = "cat_ilk_flip"
	CatIlkLumpTable         = "cat_ilk_lump"
	CatLitterTable          = "cat_litter"
	CatLiveTable            = "cat_live"
	CatVatTable             = "cat_vat"
	CatVowTable             = "cat_vow"
	CdpManagerCdpiTable     = "cdp_manager_cdpi"
	CdpManagerCountTable    = "cdp_manager_count"
	CdpManagerFirstTable    = "cdp_manager_first"
	CdpManagerIlksTable     = "cdp_manager_ilks"
	CdpManagerLastTable     = "cdp_manager_last"
	CdpManagerListNextTable = "cdp_manager_list_next"
	CdpManagerListPrevTable = "cdp_manager_list_prev"
	CdpManagerOwnsTable     = "cdp_manager_owns"
	CdpManagerUrnsTable     = "cdp_manager_urns"
	CdpManagerVatTable      = "cdp_manager_vat"
	DogDirtTable            = "dog_dirt"
	DogHoleTable            = "dog_hole"
	DogIlkChopTable         = "dog_ilk_chop"
	DogIlkClipTable         = "dog_ilk_clip"
	DogIlkDirtTable         = "dog_ilk_dirt"
	DogIlkHoleTable         = "dog_ilk_hole"
	DogLiveTable            = "dog_live"
	DogVatTable             = "dog_vat"
	DogVowTable             = "dog_vow"
	FlapBegTable            = "flap_beg"
	FlapBidBidTable         = "flap_bid_bid"
	FlapBidEndTable         = "flap_bid_end"
	FlapBidGuyTable         = "flap_bid_guy"
	FlapBidLotTable         = "flap_bid_lot"
	FlapBidTicTable         = "flap_bid_tic"
	FlapGemTable            = "flap_gem"
	FlapKicksTable          = "flap_kicks"
	FlapLiveTable           = "flap_live"
	FlapTauTable            = "flap_tau"
	FlapTtlTable            = "flap_ttl"
	FlapVatTable            = "flap_vat"
	FlipBegTable            = "flip_beg"
	FlipCatTable            = "flip_cat"
	FlipBidBidTable         = "flip_bid_bid"
	FlipBidEndTable         = "flip_bid_end"
	FlipBidGalTable         = "flip_bid_gal"
	FlipBidGuyTable         = "flip_bid_guy"
	FlipBidLotTable         = "flip_bid_lot"
	FlipBidTabTable         = "flip_bid_tab"
	FlipBidTicTable         = "flip_bid_tic"
	FlipBidUsrTable         = "flip_bid_usr"
	FlipIlkTable            = "flip_ilk"
	FlipKicksTable          = "flip_kicks"
	FlipTauTable            = "flip_tau"
	FlipTtlTable            = "flip_ttl"
	FlipVatTable            = "flip_vat"
	FlopBegTable            = "flop_beg"
	FlopBidBidTable         = "flop_bid_bid"
	FlopBidEndTable         = "flop_bid_end"
	FlopBidGuyTable         = "flop_bid_guy"
	FlopBidLotTable         = "flop_bid_lot"
	FlopBidTicTable         = "flop_bid_tic"
	FlopGemTable            = "flop_gem"
	FlopKicksTable          = "flop_kicks"
	FlopLiveTable           = "flop_live"
	FlopPadTable            = "flop_pad"
	FlopTauTable            = "flop_tau"
	FlopTtlTable            = "flop_ttl"
	FlopVatTable            = "flop_vat"
	FlopVowTable            = "flop_vow"
	JugBaseTable            = "jug_base"
	JugIlkDutyTable         = "jug_ilk_duty"
	JugIlkRhoTable          = "jug_ilk_rho"
	JugVatTable             = "jug_vat"
	JugVowTable             = "jug_vow"
	MedianValTable          = "median_val"
	MedianAgeTable          = "median_age"
	MedianBarTable          = "median_bar"
	MedianBudTable          = "median_bud"
	MedianOrclTable         = "median_orcl"
	MedianSlotTable         = "median_slot"
	PotChiTable             = "pot_chi"
	PotDsrTable             = "pot_dsr"
	PotLiveTable            = "pot_live"
	PotPieTable             = "pot_pie"
	PotRhoTable             = "pot_rho"
	SpotIlkMatTable         = "spot_ilk_mat"
	SpotIlkPipTable         = "spot_ilk_pip"
	SpotLiveTable           = "spot_live"
	SpotParTable            = "spot_par"
	SpotVatTable            = "spot_vat"
	VatDaiTable             = "vat_dai"
	VatDebtTable            = "vat_debt"
	VatGemTable             = "vat_gem"
	VatIlkArtTable          = "vat_ilk_art"
	VatIlkDustTable         = "vat_ilk_dust"
	VatIlkLineTable         = "vat_ilk_line"
	VatIlkRateTable         = "vat_ilk_rate"
	VatIlkSpotTable         = "vat_ilk_spot"
	VatLineTable            = "vat_line"
	VatLiveTable            = "vat_live"
	VatSinTable             = "vat_sin"
	VatUrnArtTable          = "vat_urn_art"
	VatUrnInkTable          = "vat_urn_ink"
	VatViceTable            = "vat_vice"
	VowAshTable             = "vow_ash"
	VowBumpTable            = "vow_bump"
	VowDumpTable            = "vow_dump"
	VowFlapperTable         = "vow_flapper"
	VowFlopperTable         = "vow_flopper"
	VowHumpTable            = "vow_hump"
	VowLiveTable            = "vow_live"
	VowSinIntegerTable      = "vow_sin_integer"
	VowSinMappingTable      = "vow_sin_mapping"
	VowSumpTable            = "vow_sump"
	VowVatTable             = "vow_vat"
	VowWaitTable            = "vow_wait"
	WardsTable              = "wards"
)

// trigger tables
const (
	FlapTable = "flap"
	FlipTable = "flip"
	FlopTable = "flop"
)
