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

	"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

// TODO Figure out signatures automatically from config somehow :(
func Cat100ABI() string     { return constants.GetContractABI("MCD_CAT_1_0_0") }
func Cat110ABI() string     { return constants.GetContractABI("MCD_CAT_1_1_0") }
func CdpManagerABI() string { return constants.GetContractABI("CDP_MANAGER") }
func ClipABI() string       { return constants.GetContractABI("MCD_CLIP_1_x_x") }
func DogABI() string        { return constants.GetContractABI("MCD_DOG_1_x_x") }
func FlapABI() string {
	return constants.GetABIFromContractsWithMatchingABI([]string{
		"MCD_FLAP_1_0_0",
		"MCD_FLAP_1_0_9",
	})
}
func FlipV100ABI() string {
	return constants.GetABIFromContractsWithMatchingABI([]string{
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
	return constants.GetABIFromContractsWithMatchingABI([]string{
		"MCD_FLIP_AAVE_A_1_2_2",
		"MCD_FLIP_BAL_A_1_1_14",
		"MCD_FLIP_BAT_A_1_1_0",
		"MCD_FLIP_COMP_A_1_1_2",
		"MCD_FLIP_ETH_A_1_1_0",
		"MCD_FLIP_ETH_B_1_1_3",
		"MCD_FLIP_ETH_C_1_2_10",
		"MCD_FLIP_GUSD_A_1_1_5",
		"MCD_FLIP_KNC_A_1_1_0",
		"MCD_FLIP_LINK_A_1_1_2",
		"MCD_FLIP_LRC_A_1_1_2",
		"MCD_FLIP_MANA_A_1_1_0",
		"MCD_FLIP_PAXUSD_A_1_1_1",
		"MCD_FLIP_RENBTC_A_1_2_1",
		"MCD_FLIP_TUSD_A_1_1_0",
		"MCD_FLIP_UNI_A_1_2_1",
		"MCD_FLIP_UNIV2AAVEETH_A_1_2_7",
		"MCD_FLIP_UNIV2DAIETH_A_1_2_2",
		"MCD_FLIP_UNIV2DAIUSDC_A_1_2_5",
		"MCD_FLIP_UNIV2DAIUSDT_A_1_2_8",
		"MCD_FLIP_UNIV2ETHUSDT_A_1_2_5",
		"MCD_FLIP_UNIV2UNIETH_A_1_2_6",
		"MCD_FLIP_UNIV2LINKETH_A_1_2_6",
		"MCD_FLIP_UNIV2USDCETH_A_1_2_4",
		"MCD_FLIP_UNIV2WBTCDAI_A_1_2_7",
		"MCD_FLIP_UNIV2WBTCETH_A_1_2_4",
		"MCD_FLIP_USDC_A_1_1_0",
		"MCD_FLIP_USDC_B_1_1_0",
		"MCD_FLIP_USDT_A_1_1_1",
		"MCD_FLIP_WBTC_A_1_1_0",
		"MCD_FLIP_YFI_A_1_1_14",
		"MCD_FLIP_ZRX_A_1_1_0",
	})
}
func FlopABI() string {
	return constants.GetABIFromContractsWithMatchingABI([]string{
		"MCD_FLOP_1_0_1",
		"MCD_FLOP_1_0_9",
	})
}
func JugABI() string { return constants.GetContractABI("MCD_JUG") }
func MedianV100ABI() string {
	return constants.GetABIFromContractsWithMatchingABI([]string{
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

func MedianV114ABI() string {
	return constants.GetABIFromContractsWithMatchingABI([]string{
		"MEDIAN_AAVE_1_2_2",
		"MEDIAN_BAL_1_1_14",
		"MEDIAN_UNI_1_2_1",
		"MEDIAN_YFI_1_1_14",
	})
}
func OsmABI() string {
	return constants.GetABIFromContractsWithMatchingABI([]string{
		"OSM_AAVE",
		"OSM_BAL",
		"OSM_BAT",
		"OSM_COMP",
		"OSM_ETH",
		"OSM_KNC",
		"OSM_LINK",
		"OSM_LRC",
		"OSM_MANA",
		"OSM_UNI",
		"OSM_USDT",
		"OSM_WBTC",
		"OSM_YFI",
		"OSM_ZRX",
	})
}
func PotABI() string  { return constants.GetContractABI("MCD_POT") }
func SpotABI() string { return constants.GetContractABI("MCD_SPOT") }
func VatABI() string  { return constants.GetContractABI("MCD_VAT") }
func VowABI() string  { return constants.GetContractABI("MCD_VOW") }

func auctionFileMethod() string { return constants.GetSolidityFunctionSignature(FlipV100ABI(), "file") }
func biteMethod() string        { return constants.GetSolidityFunctionSignature(Cat100ABI(), "Bite") }
func catClawMethod() string     { return constants.GetSolidityFunctionSignature(Cat110ABI(), "claw") }
func catFileBoxMethod() string {
	return constants.GetOverloadedFunctionSignature(Cat110ABI(), "file", []string{"bytes32", "uint256"})
}
func catFileChopLumpDunkMethod() string {
	return constants.GetOverloadedFunctionSignature(Cat110ABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func catFileFlipMethod() string {
	return constants.GetOverloadedFunctionSignature(Cat110ABI(), "file", []string{"bytes32", "bytes32", "address"})
}
func catFileVowMethod() string {
	return constants.GetOverloadedFunctionSignature(Cat110ABI(), "file", []string{"bytes32", "address"})
}
func clipKickMethod() string { return constants.GetSolidityFunctionSignature(ClipABI(), "kick") }
func dogFileIlkClipMethod() string {
	return constants.GetOverloadedFunctionSignature(DogABI(), "file", []string{"bytes32", "bytes32", "address"})
}
func dogFileIlkChopHoleMethod() string {
	return constants.GetOverloadedFunctionSignature(DogABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func dogFileVowMethod() string {
	return constants.GetOverloadedFunctionSignature(DogABI(), "file", []string{"bytes32", "address"})
}
func dealMethod() string     { return constants.GetSolidityFunctionSignature(FlipV100ABI(), "deal") }
func dentMethod() string     { return constants.GetSolidityFunctionSignature(FlipV100ABI(), "dent") }
func denyMethod() string     { return constants.GetSolidityFunctionSignature(Cat100ABI(), "deny") }
func dogBarkMethod() string  { return constants.GetSolidityFunctionSignature(DogABI(), "Bark") }
func dogDigsMethod() string  { return constants.GetSolidityFunctionSignature(DogABI(), "Digs") }
func dogDenyMethod() string  { return constants.GetSolidityFunctionSignature(DogABI(), "deny") }
func dogRelyMethod() string  { return constants.GetSolidityFunctionSignature(DogABI(), "rely") }
func flapKickMethod() string { return constants.GetSolidityFunctionSignature(FlapABI(), "Kick") }
func flipKickMethod() string { return constants.GetSolidityFunctionSignature(FlipV100ABI(), "Kick") }
func flipFileCatMethod() string {
	return constants.GetOverloadedFunctionSignature(FlipV110ABI(), "file", []string{"bytes32", "address"})
}
func flopKickMethod() string { return constants.GetSolidityFunctionSignature(FlopABI(), "Kick") }
func jugDripMethod() string  { return constants.GetSolidityFunctionSignature(JugABI(), "drip") }
func jugFileBaseMethod() string {
	return constants.GetOverloadedFunctionSignature(JugABI(), "file", []string{"bytes32", "uint256"})
}
func jugFileIlkMethod() string {
	return constants.GetOverloadedFunctionSignature(JugABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func jugFileVowMethod() string {
	return constants.GetOverloadedFunctionSignature(JugABI(), "file", []string{"bytes32", "address"})
}
func jugInitMethod() string { return constants.GetSolidityFunctionSignature(JugABI(), "init") }

func getMedianFunctionSignature(name string) string {
	signature := constants.GetSolidityFunctionSignature(MedianV114ABI(), name)
	if oldSignature := constants.GetSolidityFunctionSignature(MedianV100ABI(), name); signature != oldSignature {
		panic(fmt.Sprintf("ABI Function signature has changed! was %s but is now %s", oldSignature, signature))
	}
	return signature
}

func getOverloadedMedianFunctionSignature(name string, paramTypes []string) string {
	signature := constants.GetOverloadedFunctionSignature(MedianV114ABI(), name, paramTypes)
	if oldSignature := constants.GetOverloadedFunctionSignature(MedianV100ABI(), name, paramTypes); signature != oldSignature {
		panic(fmt.Sprintf("ABI Function signature has changed! was %s but is now %s", oldSignature, signature))
	}
	return signature
}

func logMedianPriceEvent() string {
	return getMedianFunctionSignature("LogMedianPrice")
}

func logValueMethod() string { return constants.GetSolidityFunctionSignature(OsmABI(), "LogValue") }

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
func newCdpMethod() string    { return constants.GetSolidityFunctionSignature(CdpManagerABI(), "NewCdp") }
func osmChangeMethod() string { return constants.GetSolidityFunctionSignature(OsmABI(), "change") }
func potCageMethod() string   { return constants.GetSolidityFunctionSignature(PotABI(), "cage") }
func potDripMethod() string   { return constants.GetSolidityFunctionSignature(PotABI(), "drip") }
func potExitMethod() string   { return constants.GetSolidityFunctionSignature(PotABI(), "exit") }
func potFileDSRMethod() string {
	return constants.GetOverloadedFunctionSignature(PotABI(), "file", []string{"bytes32", "uint256"})
}
func potFileVowMethod() string {
	return constants.GetOverloadedFunctionSignature(PotABI(), "file", []string{"bytes32", "address"})
}
func potJoinMethod() string { return constants.GetSolidityFunctionSignature(PotABI(), "join") }
func relyMethod() string    { return constants.GetSolidityFunctionSignature(Cat110ABI(), "rely") }

func spotFileMatMethod() string {
	return constants.GetOverloadedFunctionSignature(SpotABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func spotFileParMethod() string {
	return constants.GetOverloadedFunctionSignature(SpotABI(), "file", []string{"bytes32", "uint256"})
}
func spotFilePipMethod() string {
	return constants.GetOverloadedFunctionSignature(SpotABI(), "file", []string{"bytes32", "bytes32", "address"})
}
func spotPokeMethod() string { return constants.GetSolidityFunctionSignature(SpotABI(), "Poke") }
func tendMethod() string     { return constants.GetSolidityFunctionSignature(FlipV100ABI(), "tend") }
func tickMethod() string     { return constants.GetSolidityFunctionSignature(FlipV100ABI(), "tick") }
func vatFileDebtCeilingMethod() string {
	return constants.GetOverloadedFunctionSignature(VatABI(), "file", []string{"bytes32", "uint256"})
}
func vatFileIlkMethod() string {
	return constants.GetOverloadedFunctionSignature(VatABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func vatFluxMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "flux") }
func vatFoldMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "fold") }
func vatForkMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "fork") }
func vatFrobMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "frob") }
func vatGrabMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "grab") }
func vatHealMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "heal") }
func vatHopeMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "hope") }
func vatInitMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "init") }
func vatMoveMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "move") }
func vatNopeMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "nope") }
func vatSlipMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "slip") }
func vatSuckMethod() string { return constants.GetSolidityFunctionSignature(VatABI(), "suck") }
func vowFessMethod() string { return constants.GetSolidityFunctionSignature(VowABI(), "fess") }
func vowFileAuctionAttributesMethod() string {
	return constants.GetOverloadedFunctionSignature(VowABI(), "file", []string{"bytes32", "uint256"})
}
func vowFileAuctionAddressMethod() string {
	return constants.GetOverloadedFunctionSignature(VowABI(), "file", []string{"bytes32", "address"})
}
func vowFlogMethod() string { return constants.GetSolidityFunctionSignature(VowABI(), "flog") }
func vowHealMethod() string { return constants.GetSolidityFunctionSignature(VowABI(), "heal") }
func yankMethod() string    { return constants.GetSolidityFunctionSignature(FlipV100ABI(), "yank") }
