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

// TODO Figure out signatures automatically from config somehow :(
func CatABI() string        { return getContractABI("MCD_CAT") }
func CdpManagerABI() string { return getContractABI("CDP_MANAGER") }
func FlapABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"MCD_FLAP_1.0.0",
		"MCD_FLAP_1.0.9",
	})
}
func FlipABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"MCD_FLIP_BAT_A_1.0.0",
		"MCD_FLIP_BAT_A_1.0.9",
		"MCD_FLIP_ETH_A_1.0.0",
		"MCD_FLIP_ETH_A_1.0.9",
		"MCD_FLIP_KNC_A_1.0.8",
		"MCD_FLIP_KNC_A_1.0.9",
		"MCD_FLIP_MANA_A_1.0.9",
		"MCD_FLIP_SAI",
		"MCD_FLIP_TUSD_A_1.0.7",
		"MCD_FLIP_TUSD_A_1.0.9",
		"MCD_FLIP_USDC_A_1.0.4",
		"MCD_FLIP_USDC_A_1.0.9",
		"MCD_FLIP_USDC_B_1.0.7",
		"MCD_FLIP_USDC_B_1.0.9",
		"MCD_FLIP_WBTC_A_1.0.6",
		"MCD_FLIP_WBTC_A_1.0.9",
		"MCD_FLIP_ZRX_A_1.0.8",
		"MCD_FLIP_ZRX_A_1.0.9",
	})
}
func FlopABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"MCD_FLOP_1.0.1",
		"MCD_FLOP_1.0.9",
	})
}
func JugABI() string { return getContractABI("MCD_JUG") }
func MedianABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"MEDIAN_BAT",
		"MEDIAN_ETH",
		"MEDIAN_KNC",
		"MEDIAN_WBTC",
		"MEDIAN_ZRX",
	})
}
func OasisABI() string {
	return GetABIFromContractsWithMatchingABI([]string{"OASIS_MATCHING_MARKET_ONE", "OASIS_MATCHING_MARKET_TWO"})
}
func OsmABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"OSM_BAT",
		"OSM_ETH",
		"OSM_KNC",
		"OSM_WBTC",
		"OSM_ZRX",
	})
}
func PotABI() string  { return getContractABI("MCD_POT") }
func SpotABI() string { return getContractABI("MCD_SPOT") }
func VatABI() string  { return getContractABI("MCD_VAT") }
func VowABI() string  { return getContractABI("MCD_VOW") }

func auctionFileMethod() string { return getSolidityFunctionSignature(FlipABI(), "file") }
func biteMethod() string        { return getSolidityFunctionSignature(CatABI(), "Bite") }
func catFileChopLumpMethod() string {
	return getOverloadedFunctionSignature(CatABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func catFileFlipMethod() string {
	return getOverloadedFunctionSignature(CatABI(), "file", []string{"bytes32", "bytes32", "address"})
}
func catFileVowMethod() string {
	return getOverloadedFunctionSignature(CatABI(), "file", []string{"bytes32", "address"})
}
func dealMethod() string     { return getSolidityFunctionSignature(FlipABI(), "deal") }
func dentMethod() string     { return getSolidityFunctionSignature(FlipABI(), "dent") }
func denyMethod() string     { return getSolidityFunctionSignature(CatABI(), "deny") }
func flapKickMethod() string { return getSolidityFunctionSignature(FlapABI(), "Kick") }
func flipKickMethod() string { return getSolidityFunctionSignature(FlipABI(), "Kick") }
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
func logMedianPriceEvent() string  { return getSolidityFunctionSignature(MedianABI(), "LogMedianPrice") }
func logMinSellEvent() string      { return getSolidityFunctionSignature(OasisABI(), "LogMinSell") }
func logSortedOfferMethod() string { return getSolidityFunctionSignature(OasisABI(), "LogSortedOffer") }
func logTakeEvent() string         { return getSolidityFunctionSignature(OasisABI(), "LogTake") }
func logTradeEvent() string        { return getSolidityFunctionSignature(OasisABI(), "LogTrade") }
func logUnsortedOfferMethod() string {
	return getSolidityFunctionSignature(OasisABI(), "LogUnsortedOffer")
}
func logValueMethod() string { return getSolidityFunctionSignature(OsmABI(), "LogValue") }

func medianDissBatchMethod() string {
	return getOverloadedFunctionSignature(MedianABI(), "diss", []string{"address[]"})
}
func medianDissSingleMethod() string {
	return getOverloadedFunctionSignature(MedianABI(), "diss", []string{"address"})
}
func medianDropMethod() string {
	return getSolidityFunctionSignature(MedianABI(), "drop")
}
func medianLiftMethod() string {
	return getSolidityFunctionSignature(MedianABI(), "lift")
}
func medianKissBatchMethod() string {
	return getOverloadedFunctionSignature(MedianABI(), "kiss", []string{"address[]"})
}
func medianKissSingleMethod() string {
	return getOverloadedFunctionSignature(MedianABI(), "kiss", []string{"address"})
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
func relyMethod() string       { return getSolidityFunctionSignature(CatABI(), "rely") }
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
func tendMethod() string     { return getSolidityFunctionSignature(FlipABI(), "tend") }
func tickMethod() string     { return getSolidityFunctionSignature(FlipABI(), "tick") }
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
func vowFileAuctionAttributesMethod() string {
	return getOverloadedFunctionSignature(VowABI(), "file", []string{"bytes32", "uint256"})
}
func vowFileAuctionAddressMethod() string {
	return getOverloadedFunctionSignature(VowABI(), "file", []string{"bytes32", "address"})
}
func vowFlogMethod() string { return getSolidityFunctionSignature(VowABI(), "flog") }
func vowHealMethod() string { return getSolidityFunctionSignature(VowABI(), "heal") }
func yankMethod() string    { return getSolidityFunctionSignature(FlipABI(), "yank") }
