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

// TODO Figure signatures out automatically from config somehow :(

func flipABI() string { return getContractABI("ALL_MCD_FLIPS") }
func catABI() string  { return getContractABI("MCD_CAT") }
func flapABI() string { return getContractABI("MCD_FLAP") }
func jugABI() string  { return getContractABI("MCD_JUG") }
func spotABI() string { return getContractABI("MCD_SPOT") }
func vatABI() string  { return getContractABI("MCD_VAT") }
func vowABI() string  { return getContractABI("MCD_VOW") }

func biteMethod() string { return getSolidityFunctionSignature(catABI(), "Bite") }
func catFileChopLumpMethod() string {
	return getOverloadedFunctionSignature(catABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func catFileFlipMethod() string {
	return getOverloadedFunctionSignature(catABI(), "file", []string{"bytes32", "bytes32", "address"})
}
func catFileVowMethod() string {
	return getOverloadedFunctionSignature(catABI(), "file", []string{"bytes32", "address"})
}
func dealMethod() string     { return getSolidityFunctionSignature(flipABI(), "deal") }
func dentMethod() string     { return getSolidityFunctionSignature(flipABI(), "dent") }
func flapKickMethod() string { return getSolidityFunctionSignature(flapABI(), "Kick") }
func flipKickMethod() string { return getSolidityFunctionSignature(flipABI(), "Kick") }
func flipTickMethod() string { return getSolidityFunctionSignature(flipABI(), "tick") }
func flopKickMethod() string { return getSolidityFunctionSignature(flapABI(), "Kick") }
func jugDripMethod() string  { return getSolidityFunctionSignature(jugABI(), "drip") }
func jugFileBaseMethod() string {
	return getOverloadedFunctionSignature(jugABI(), "file", []string{"bytes32", "uint256"})
}
func jugFileIlkMethod() string {
	return getOverloadedFunctionSignature(jugABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func jugFileVowMethod() string {
	return getOverloadedFunctionSignature(jugABI(), "file", []string{"bytes32", "address"})
}
func jugInitMethod() string { return getSolidityFunctionSignature(jugABI(), "init") }
func spotFileMatMethod() string {
	return getOverloadedFunctionSignature(spotABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func spotFilePipMethod() string {
	return getOverloadedFunctionSignature(spotABI(), "file", []string{"bytes32", "address"})
}
func spotPokeMethod() string { return getSolidityFunctionSignature(spotABI(), "Poke") }
func tendMethod() string     { return getSolidityFunctionSignature(flipABI(), "tend") }
func vatFileDebtCeilingMethod() string {
	return getOverloadedFunctionSignature(vatABI(), "file", []string{"bytes32", "uint256"})
}
func vatFileIlkMethod() string {
	return getOverloadedFunctionSignature(vatABI(), "file", []string{"bytes32", "bytes32", "uint256"})
}
func vatFluxMethod() string { return getSolidityFunctionSignature(vatABI(), "flux") }
func vatFoldMethod() string { return getSolidityFunctionSignature(vatABI(), "fold") }
func vatForkMethod() string { return getSolidityFunctionSignature(vatABI(), "fork") }
func vatFrobMethod() string { return getSolidityFunctionSignature(vatABI(), "frob") }
func vatGrabMethod() string { return getSolidityFunctionSignature(vatABI(), "grab") }
func vatHealMethod() string { return getSolidityFunctionSignature(vatABI(), "heal") }
func vatInitMethod() string { return getSolidityFunctionSignature(vatABI(), "init") }
func vatMoveMethod() string { return getSolidityFunctionSignature(vatABI(), "move") }
func vatSlipMethod() string { return getSolidityFunctionSignature(vatABI(), "slip") }
func vatSuckMethod() string { return getSolidityFunctionSignature(vatABI(), "suck") }
func vowFessMethod() string { return getSolidityFunctionSignature(vowABI(), "fess") }
func vowFileMethod() string { return getSolidityFunctionSignature(vowABI(), "file") }
func vowFlogMethod() string { return getSolidityFunctionSignature(vowABI(), "flog") }
func yankMethod() string    { return getSolidityFunctionSignature(flipABI(), "yank") }
