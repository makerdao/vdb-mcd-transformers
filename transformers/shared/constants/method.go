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
		"MCD_FLIP_BAT_A_1_1_0",
		"MCD_FLIP_COMP_A_1_1_2",
		"MCD_FLIP_ETH_A_1_1_0",
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
func MedianABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"MEDIAN_BAT",
		"MEDIAN_COMP",
		"MEDIAN_ETH",
		"MEDIAN_KNC",
		"MEDIAN_LINK",
		"MEDIAN_LRC",
		"MEDIAN_MANA",
		"MEDIAN_USDT",
		"MEDIAN_WBTC",
		"MEDIAN_ZRX",
	})
}
func OsmABI() string {
	return GetABIFromContractsWithMatchingABI([]string{
		"OSM_BAT",
		"OSM_COMP",
		"OSM_ETH",
		"OSM_KNC",
		"OSM_LINK",
		"OSM_LRC",
		"OSM_MANA",
		"OSM_USDT",
		"OSM_WBTC",
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
	return GetOverloadedFunctionSignature(Cat110ABI(), "file", []string{"bytes32", "uint256"})
}
func catFileChopLumpDunkMethod() string {
	return GetOverloadedFunctionSignature(Cat110ABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func catFileFlipMethod() string {
	return GetOverloadedFunctionSignature(Cat110ABI(), "file", []string{"bytes32", "bytes32", "address"})
}
func catFileVowMethod() string {
	return GetOverloadedFunctionSignature(Cat110ABI(), "file", []string{"bytes32", "address"})
}
func dealMethod() string     { return getSolidityFunctionSignature(FlipV100ABI(), "deal") }
func dentMethod() string     { return getSolidityFunctionSignature(FlipV100ABI(), "dent") }
func denyMethod() string     { return getSolidityFunctionSignature(Cat100ABI(), "deny") }
func flapKickMethod() string { return getSolidityFunctionSignature(FlapABI(), "Kick") }
func flipKickMethod() string { return getSolidityFunctionSignature(FlipV100ABI(), "Kick") }
func flipFileCatMethod() string {
	return GetOverloadedFunctionSignature(FlipV110ABI(), "file", []string{"bytes32", "address"})
}
func flopKickMethod() string { return getSolidityFunctionSignature(FlopABI(), "Kick") }
func jugDripMethod() string  { return getSolidityFunctionSignature(JugABI(), "drip") }
func jugFileBaseMethod() string {
	return GetOverloadedFunctionSignature(JugABI(), "file", []string{"bytes32", "uint256"})
}
func jugFileIlkMethod() string {
	return GetOverloadedFunctionSignature(JugABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func jugFileVowMethod() string {
	return GetOverloadedFunctionSignature(JugABI(), "file", []string{"bytes32", "address"})
}
func jugInitMethod() string       { return getSolidityFunctionSignature(JugABI(), "init") }
func logMedianPriceEvent() string { return getSolidityFunctionSignature(MedianABI(), "LogMedianPrice") }
func logValueMethod() string      { return getSolidityFunctionSignature(OsmABI(), "LogValue") }

func medianDissBatchMethod() string {
	return GetOverloadedFunctionSignature(MedianABI(), "diss", []string{"address[]"})
}
func medianDissSingleMethod() string {
	return GetOverloadedFunctionSignature(MedianABI(), "diss", []string{"address"})
}
func medianDropMethod() string {
	return getSolidityFunctionSignature(MedianABI(), "drop")
}
func medianLiftMethod() string {
	return getSolidityFunctionSignature(MedianABI(), "lift")
}
func medianKissBatchMethod() string {
	return GetOverloadedFunctionSignature(MedianABI(), "kiss", []string{"address[]"})
}
func medianKissSingleMethod() string {
	return GetOverloadedFunctionSignature(MedianABI(), "kiss", []string{"address"})
}
func newCdpMethod() string    { return getSolidityFunctionSignature(CdpManagerABI(), "NewCdp") }
func osmChangeMethod() string { return getSolidityFunctionSignature(OsmABI(), "change") }
func potCageMethod() string   { return getSolidityFunctionSignature(PotABI(), "cage") }
func potDripMethod() string   { return getSolidityFunctionSignature(PotABI(), "drip") }
func potExitMethod() string   { return getSolidityFunctionSignature(PotABI(), "exit") }
func potFileDSRMethod() string {
	return GetOverloadedFunctionSignature(PotABI(), "file", []string{"bytes32", "uint256"})
}
func potFileVowMethod() string {
	return GetOverloadedFunctionSignature(PotABI(), "file", []string{"bytes32", "address"})
}
func potJoinMethod() string { return getSolidityFunctionSignature(PotABI(), "join") }
func relyMethod() string    { return getSolidityFunctionSignature(Cat110ABI(), "rely") }

func spotFileMatMethod() string {
	return GetOverloadedFunctionSignature(SpotABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func spotFileParMethod() string {
	return GetOverloadedFunctionSignature(SpotABI(), "file", []string{"bytes32", "uint256"})
}
func spotFilePipMethod() string {
	return GetOverloadedFunctionSignature(SpotABI(), "file", []string{"bytes32", "bytes32", "address"})
}
func spotPokeMethod() string { return getSolidityFunctionSignature(SpotABI(), "Poke") }
func tendMethod() string     { return getSolidityFunctionSignature(FlipV100ABI(), "tend") }
func tickMethod() string     { return getSolidityFunctionSignature(FlipV100ABI(), "tick") }
func vatFileDebtCeilingMethod() string {
	return GetOverloadedFunctionSignature(VatABI(), "file", []string{"bytes32", "uint256"})
}
func vatFileIlkMethod() string {
	return GetOverloadedFunctionSignature(VatABI(), "file", []string{"bytes32", "bytes32", "uint256"})
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
	return GetOverloadedFunctionSignature(VowABI(), "file", []string{"bytes32", "uint256"})
}
func vowFileAuctionAddressMethod() string {
	return GetOverloadedFunctionSignature(VowABI(), "file", []string{"bytes32", "address"})
}
func vowFlogMethod() string { return getSolidityFunctionSignature(VowABI(), "flog") }
func vowHealMethod() string { return getSolidityFunctionSignature(VowABI(), "heal") }
func yankMethod() string    { return getSolidityFunctionSignature(FlipV100ABI(), "yank") }
