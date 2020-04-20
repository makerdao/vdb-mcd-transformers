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
func FlapABI() string       { return getContractABI("MCD_FLAP") }
func FlipABI() string {
	return GetContractsABI([]string{
		"MCD_FLIP_ETH_A", "MCD_FLIP_BAT_A", "MCD_FLIP_SAI",
	})
}
func FlopABI() string { return getContractABI("MCD_FLOP") }
func JugABI() string  { return getContractABI("MCD_JUG") }
func MedianABI() string {
	return GetContractsABI([]string{
		"MEDIAN_ETH", "MEDIAN_BAT",
	})
}
func OasisABI() string {
	return GetContractsABI([]string{"OASIS_MATCHING_MARKET_ONE", "OASIS_MATCHING_MARKET_TWO"})
}
func OsmABI() string {
	return GetContractsABI([]string{"OSM_ETH", "OSM_BAT"})
}
func PotABI() string  { return getContractABI("MCD_POT") }
func SpotABI() string { return getContractABI("MCD_SPOT") }
func VatABI() string  { return getContractABI("MCD_VAT") }
func VowABI() string  { return getContractABI("MCD_VOW") }

func biteMethod() string { return getSolidityFunctionSignature(CatABI(), "Bite") }
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
func jugInitMethod() string        { return getSolidityFunctionSignature(JugABI(), "init") }
func logBumpEvent() string         { return getSolidityFunctionSignature(OasisABI(), "LogBump") }
func logItemUpdateEvent() string   { return getSolidityFunctionSignature(OasisABI(), "LogItemUpdate") }
func logKillEvent() string         { return getSolidityFunctionSignature(OasisABI(), "LogKill") }
func logMakeEvent() string         { return getSolidityFunctionSignature(OasisABI(), "LogMake") }
func logSortedOfferMethod() string { return getSolidityFunctionSignature(OasisABI(), "LogSortedOffer") }
func logTakeEvent() string         { return getSolidityFunctionSignature(OasisABI(), "LogTake") }
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
func potJoinMethod() string { return getSolidityFunctionSignature(PotABI(), "join") }
func relyMethod() string    { return getSolidityFunctionSignature(CatABI(), "rely") }
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
func vowFileMethod() string { return getSolidityFunctionSignature(VowABI(), "file") }
func vowFlogMethod() string { return getSolidityFunctionSignature(VowABI(), "flog") }
func yankMethod() string    { return getSolidityFunctionSignature(FlipABI(), "yank") }
