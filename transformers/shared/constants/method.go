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
	"fmt"
)

// TODO Figure out signatures automatically from config somehow :(
func Cat100ABI() string     { return getContractABI("MCD_CAT_1_0_0") }
func Cat110ABI() string     { return getContractABI("MCD_CAT_1_1_0") }
func CdpManagerABI() string { return getContractABI("CDP_MANAGER") }
func FlapABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"MCD_FLAP_1_0_0",
		"MCD_FLAP_1_0_9",
	})
}
func FlipV100ABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"MCD_FLIP_BAT_A_1_0_0",
		"MCD_FLIP_BAT_A_1_0_9",
		"MCD_FLIP_ETH_A_1_0_0",
		"MCD_FLIP_ETH_A_1_0_9",
		"MCD_FLIP_KNC_A_1_0_8",
		"MCD_FLIP_KNC_A_1_0_9",
		"MCD_FLIP_MANA_A_1_0_9",
		"MCD_FLIP_SAI_1_0_0",
		"MCD_FLIP_TUSD_A_1_0_7",
		"MCD_FLIP_TUSD_A_1_0_9",
		"MCD_FLIP_USDC_A_1_0_4",
		"MCD_FLIP_USDC_A_1_0_9",
		"MCD_FLIP_USDC_B_1_0_7",
		"MCD_FLIP_USDC_B_1_0_9",
		"MCD_FLIP_WBTC_A_1_0_6",
		"MCD_FLIP_WBTC_A_1_0_9",
		"MCD_FLIP_ZRX_A_1_0_8",
		"MCD_FLIP_ZRX_A_1_0_9",
	})
}
func FlipV110ABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"MCD_FLIP_BAL_A_1_1_14",
		"MCD_FLIP_BAT_A_1_1_0",
		"MCD_FLIP_COMP_A_1_1_2",
		"MCD_FLIP_ETH_A_1_1_0",
		"MCD_FLIP_ETH_B_1_1_3",
		"MCD_FLIP_GUSD_A_1_1_5",
		"MCD_FLIP_KNC_A_1_1_0",
		"MCD_FLIP_LINK_A_1_1_2",
		"MCD_FLIP_LRC_A_1_1_2",
		"MCD_FLIP_MANA_A_1_1_0",
		"MCD_FLIP_PAXUSD_A_1_1_1",
		"MCD_FLIP_TUSD_A_1_1_0",
		"MCD_FLIP_USDC_A_1_1_0",
		"MCD_FLIP_USDC_B_1_1_0",
		"MCD_FLIP_USDT_A_1_1_1",
		"MCD_FLIP_WBTC_A_1_1_0",
		"MCD_FLIP_YFI_A_1_1_14",
		"MCD_FLIP_ZRX_A_1_1_0",
	})
}
func FlopABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"MCD_FLOP_1_0_1",
		"MCD_FLOP_1_0_9",
	})
}

func JugABI() string { return getContractABI("MCD_JUG") }

func MedianV100ABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"MEDIAN_BAT_1_0_0",
		"MEDIAN_COMP_1_1_2",
		"MEDIAN_ETH_1_0_0",
		"MEDIAN_KNC_1_0_8",
		"MEDIAN_LINK_1_1_2",
		"MEDIAN_LRC_1_1_2",
		"MEDIAN_MANA_1_0_9",
		"MEDIAN_USDT_1_0_4",
		"MEDIAN_WBTC_1_0_6",
		"MEDIAN_ZRX_1_0_8",
	})
}

func OasisABI() string {
	return GetABIFromContractsWithMatchingABI([]string{"OASIS_MATCHING_MARKET_ONE", "OASIS_MATCHING_MARKET_TWO"})
}

func MedianV114ABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"MEDIAN_BAL_1_1_14",
		"MEDIAN_YFI_1_1_14",
	})
}

func OsmABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"OSM_BAL",
		"OSM_BAT",
		"OSM_COMP",
		"OSM_ETH",
		"OSM_KNC",
		"OSM_LINK",
		"OSM_LRC",
		"OSM_MANA",
		"OSM_USDT",
		"OSM_WBTC",
		"OSM_YFI",
		"OSM_ZRX",
	})
}
func PotABI() string  { return getContractABI("MCD_POT") }
func SpotABI() string { return getContractABI("MCD_SPOT") }
func VatABI() string  { return getContractABI("MCD_VAT") }
func VowABI() string  { return getContractABI("MCD_VOW") }

func auctionFileMethod() string { return getSolidityFunctionSignature(FlipV100ABI(), "file") }
func biteMethod() string        { return getSolidityFunctionSignature(Cat100ABI(), "Bite") }
func catClawMethod() string     { return getSolidityFunctionSignature(Cat110ABI(), "claw") }
func catFileBoxMethod() string {
	return getOverloadedFunctionSignature(Cat110ABI(), "file", []string{"bytes32", "uint256"})
}
func catFileChopLumpMethod() string {
	return getOverloadedFunctionSignature(Cat110ABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func catFileFlipMethod() string {
	return getOverloadedFunctionSignature(Cat100ABI(), "file", []string{"bytes32", "bytes32", "address"})
}
func catFileVowMethod() string {
	return getOverloadedFunctionSignature(Cat100ABI(), "file", []string{"bytes32", "address"})
}
func dealMethod() string     { return getSolidityFunctionSignature(FlipV100ABI(), "deal") }
func dentMethod() string     { return getSolidityFunctionSignature(FlipV100ABI(), "dent") }
func denyMethod() string     { return getSolidityFunctionSignature(Cat100ABI(), "deny") }
func flapKickMethod() string { return getSolidityFunctionSignature(FlapABI(), "Kick") }
func flipFileCatMethod() string {
	return getOverloadedFunctionSignature(FlipV110ABI(), "file", []string{"bytes32", "address"})
}
func flipKickMethod() string { return getSolidityFunctionSignature(FlipV100ABI(), "Kick") }
func flopKickMethod() string { return getSolidityFunctionSignature(FlopABI(), "Kick") }
func jugDripMethod() string  { return getSolidityFunctionSignature(JugABI(), "drip") }
func jugFileBaseMethod() string {
	return getOverloadedFunctionSignature(JugABI(), "file", []string{"bytes32", "uint256"})
}
func jugFileIlkMethod() string {
	return getOverloadedFunctionSignature(JugABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func jugFileVowMethod() string {
	return getOverloadedFunctionSignature(JugABI(), "file", []string{"bytes32", "address"})
}
func jugInitMethod() string      { return getSolidityFunctionSignature(JugABI(), "init") }
func logBumpEvent() string       { return getSolidityFunctionSignature(OasisABI(), "LogBump") }
func logBuyEnabledEvent() string { return getSolidityFunctionSignature(OasisABI(), "LogBuyEnabled") }
func logDeleteEvent() string     { return getSolidityFunctionSignature(OasisABI(), "LogDelete") }
func logInsertEvent() string     { return getSolidityFunctionSignature(OasisABI(), "LogInsert") }
func logItemUpdateEvent() string { return getSolidityFunctionSignature(OasisABI(), "LogItemUpdate") }
func logKillEvent() string       { return getSolidityFunctionSignature(OasisABI(), "LogKill") }
func logMakeEvent() string       { return getSolidityFunctionSignature(OasisABI(), "LogMake") }
func logMatchingEnabledEvent() string {
	return getSolidityFunctionSignature(OasisABI(), "LogMatchingEnabled")
}

func getMedianFunctionSignature(name string) string {
	signature := getSolidityFunctionSignature(MedianV114ABI(), name)
	if oldSignature := getSolidityFunctionSignature(MedianV100ABI(), name); signature != oldSignature {
		panic(fmt.Sprintf("ABI Function signature has changed! was %s but is now %s", oldSignature, signature))
	}
	return signature
}

func getOverloadedMedianFunctionSignature(name string, paramTypes []string) string {
	signature := getOverloadedFunctionSignature(MedianV114ABI(), name, paramTypes)
	if oldSignature := getOverloadedFunctionSignature(MedianV100ABI(), name, paramTypes); signature != oldSignature {
		panic(fmt.Sprintf("ABI Function signature has changed! was %s but is now %s", oldSignature, signature))
	}
	return signature
}

func logMedianPriceEvent() string {
	return getMedianFunctionSignature("LogMedianPrice")
}
func logMinSellEvent() string      { return getSolidityFunctionSignature(OasisABI(), "LogMinSell") }
func logSortedOfferMethod() string { return getSolidityFunctionSignature(OasisABI(), "LogSortedOffer") }
func logTakeEvent() string         { return getSolidityFunctionSignature(OasisABI(), "LogTake") }
func logTradeEvent() string        { return getSolidityFunctionSignature(OasisABI(), "LogTrade") }
func logUnsortedOfferMethod() string {
	return getSolidityFunctionSignature(OasisABI(), "LogUnsortedOffer")
}
func logValueMethod() string { return getSolidityFunctionSignature(OsmABI(), "LogValue") }

func medianDissBatchMethod() string {
	return getOverloadedMedianFunctionSignature("diss", []string{"address[]"})
}
func medianDissSingleMethod() string {
	return getOverloadedMedianFunctionSignature("diss", []string{"address"})
}
func medianDropMethod() string {
	return getMedianFunctionSignature("drop")
}
func medianLiftMethod() string {
	return getMedianFunctionSignature("lift")
}
func medianKissBatchMethod() string {
	return getOverloadedMedianFunctionSignature("kiss", []string{"address[]"})
}
func medianKissSingleMethod() string {
	return getOverloadedMedianFunctionSignature("kiss", []string{"address"})
}
func newCdpMethod() string    { return getSolidityFunctionSignature(CdpManagerABI(), "NewCdp") }
func osmChangeMethod() string { return getSolidityFunctionSignature(OsmABI(), "change") }
func potCageMethod() string   { return getSolidityFunctionSignature(PotABI(), "cage") }
func potDripMethod() string   { return getSolidityFunctionSignature(PotABI(), "drip") }
func potExitMethod() string   { return getSolidityFunctionSignature(PotABI(), "exit") }
func potFileDSRMethod() string {
	return getOverloadedFunctionSignature(PotABI(), "file", []string{"bytes32", "uint256"})
}
func potFileVowMethod() string {
	return getOverloadedFunctionSignature(PotABI(), "file", []string{"bytes32", "address"})
}
func potJoinMethod() string    { return getSolidityFunctionSignature(PotABI(), "join") }
func relyMethod() string       { return getSolidityFunctionSignature(Cat100ABI(), "rely") }
func setMinSellMethod() string { return getSolidityFunctionSignature(OasisABI(), "setMinSell") }
func spotFileMatMethod() string {
	return getOverloadedFunctionSignature(SpotABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func spotFileParMethod() string {
	return getOverloadedFunctionSignature(SpotABI(), "file", []string{"bytes32", "uint256"})
}
func spotFilePipMethod() string {
	return getOverloadedFunctionSignature(SpotABI(), "file", []string{"bytes32", "bytes32", "address"})
}
func spotPokeMethod() string { return getSolidityFunctionSignature(SpotABI(), "Poke") }
func tendMethod() string     { return getSolidityFunctionSignature(FlipV100ABI(), "tend") }
func tickMethod() string     { return getSolidityFunctionSignature(FlipV100ABI(), "tick") }
func vatFileDebtCeilingMethod() string {
	return getOverloadedFunctionSignature(VatABI(), "file", []string{"bytes32", "uint256"})
}
func vatFileIlkMethod() string {
	return getOverloadedFunctionSignature(VatABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func vatFluxMethod() string { return getSolidityFunctionSignature(VatABI(), "flux") }
func vatFoldMethod() string { return getSolidityFunctionSignature(VatABI(), "fold") }
func vatForkMethod() string { return getSolidityFunctionSignature(VatABI(), "fork") }
func vatFrobMethod() string { return getSolidityFunctionSignature(VatABI(), "frob") }
func vatGrabMethod() string { return getSolidityFunctionSignature(VatABI(), "grab") }
func vatHealMethod() string { return getSolidityFunctionSignature(VatABI(), "heal") }
func vatHopeMethod() string { return getSolidityFunctionSignature(VatABI(), "hope") }
func vatInitMethod() string { return getSolidityFunctionSignature(VatABI(), "init") }
func vatMoveMethod() string { return getSolidityFunctionSignature(VatABI(), "move") }
func vatNopeMethod() string { return getSolidityFunctionSignature(VatABI(), "nope") }
func vatSlipMethod() string { return getSolidityFunctionSignature(VatABI(), "slip") }
func vatSuckMethod() string { return getSolidityFunctionSignature(VatABI(), "suck") }
func vowFessMethod() string { return getSolidityFunctionSignature(VowABI(), "fess") }
func vowFileMethod() string {
	return getOverloadedFunctionSignature(VowABI(), "file", []string{"bytes32", "uint256"})
}
func vowFileAuctionAddressMethod() string {
	return getOverloadedFunctionSignature(VowABI(), "file", []string{"bytes32", "address"})
}
func vowFlogMethod() string { return getSolidityFunctionSignature(VowABI(), "flog") }
func vowHealMethod() string { return getSolidityFunctionSignature(VowABI(), "heal") }
func yankMethod() string    { return getSolidityFunctionSignature(FlipV100ABI(), "yank") }
