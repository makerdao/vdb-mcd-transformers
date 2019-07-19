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
func CatABI() string  { return getContractABI("MCD_CAT") }
func FlapABI() string { return getContractABI("MCD_FLAP") }
func FlipABI() string {
	return GetContractsABI([]string{
		"ETH_FLIP_A", "ETH_FLIP_B", "ETH_FLIP_C",
		"COL1_FLIP", "COL2_FLIP", "COL3_FLIP", "COL4_FLIP", "COL5_FLIP",
	})
}
func FlopABI() string { return getContractABI("MCD_FLOP") }
func JugABI() string  { return getContractABI("MCD_JUG") }
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
func flapKickMethod() string { return getSolidityFunctionSignature(FlapABI(), "Kick") }
func flipKickMethod() string { return getSolidityFunctionSignature(FlipABI(), "Kick") }
func flipTickMethod() string { return getSolidityFunctionSignature(FlipABI(), "tick") }
func flopKickMethod() string { return getSolidityFunctionSignature(FlapABI(), "Kick") }
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
func jugInitMethod() string { return getSolidityFunctionSignature(JugABI(), "init") }
func spotFileMatMethod() string {
	return getOverloadedFunctionSignature(SpotABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func spotFilePipMethod() string {
	return getOverloadedFunctionSignature(SpotABI(), "file", []string{"bytes32", "address"})
}
func spotPokeMethod() string { return getSolidityFunctionSignature(SpotABI(), "Poke") }
func tendMethod() string     { return getSolidityFunctionSignature(FlipABI(), "tend") }
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
func vatInitMethod() string { return getSolidityFunctionSignature(VatABI(), "init") }
func vatMoveMethod() string { return getSolidityFunctionSignature(VatABI(), "move") }
func vatSlipMethod() string { return getSolidityFunctionSignature(VatABI(), "slip") }
func vatSuckMethod() string { return getSolidityFunctionSignature(VatABI(), "suck") }
func vowFessMethod() string { return getSolidityFunctionSignature(VowABI(), "fess") }
func vowFileMethod() string { return getSolidityFunctionSignature(VowABI(), "file") }
func vowFlogMethod() string { return getSolidityFunctionSignature(VowABI(), "flog") }
func yankMethod() string    { return getSolidityFunctionSignature(FlipABI(), "yank") }
