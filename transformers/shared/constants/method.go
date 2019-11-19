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
	"github.com/sirupsen/logrus"
)

// TODO Figure out signatures automatically from config somehow :(
func CatABI() string {
	abi, err := GetContractsABI([]string{"MCD_CAT"})
	if err != nil {
		logrus.Fatalf("error getting cat abi: %s", err.Error())
	}
	return abi
}
func CdpManagerABI() string {
	abi, err := GetContractsABI([]string{"CDP_MANAGER"})
	if err != nil {
		logrus.Fatalf("error getting cdb manager abi: %s", err.Error())
	}
	return abi
}
func FlapABI() string {
	abi, err := GetContractsABI([]string{"MCD_FLAP"})
	if err != nil {
		logrus.Fatalf("error getting flap abi: %s", err.Error())
	}
	return abi
}
func FlipABI() string {
	abi, err := GetContractsABI([]string{
		"MCD_FLIP_ETH_A", "MCD_FLIP_BAT_A", "MCD_FLIP_SAI",
	})
	if err != nil {
		logrus.Fatalf("error getting flip abi: %s", err.Error())
	}
	return abi
}
func FlopABI() string {
	abi, err := GetContractsABI([]string{"MCD_FLOP"})
	if err != nil {
		logrus.Fatalf("error getting flop abi: %s", err.Error())
	}
	return abi
}
func JugABI() string {
	abi, err := GetContractsABI([]string{"MCD_JUG"})
	if err != nil {
		logrus.Fatalf("error getting jug abi: %s", err.Error())
	}
	return abi
}
func OsmABI() string {
	abi, err := GetContractsABI([]string{"OSM_ETH", "OSM_BAT"})
	if err != nil {
		logrus.Fatalf("error getting osm abi: %s", err.Error())
	}
	return abi
}
func PotABI() string {
	abi, err := GetContractsABI([]string{"MCD_POT"})
	if err != nil {
		logrus.Fatalf("error getting pot abi: %s", err.Error())
	}
	return abi
}
func SpotABI() string {
	abi, err := GetContractsABI([]string{"MCD_SPOT"})
	if err != nil {
		logrus.Fatalf("error getting spot abi: %s", err.Error())
	}
	return abi
}
func VatABI() string {
	abi, err := GetContractsABI([]string{"MCD_VAT"})
	if err != nil {
		logrus.Fatalf("error getting vat abi: %s", err.Error())
	}
	return abi
}
func VowABI() string {
	abi, err := GetContractsABI([]string{"MCD_VOW"})
	if err != nil {
		logrus.Fatalf("error getting vow abi: %s", err.Error())
	}
	return abi
}

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
func jugInitMethod() string  { return getSolidityFunctionSignature(JugABI(), "init") }
func logValueMethod() string { return getSolidityFunctionSignature(OsmABI(), "LogValue") }
func newCdpMethod() string   { return getSolidityFunctionSignature(CdpManagerABI(), "NewCdp") }
func potCageMethod() string  { return getSolidityFunctionSignature(PotABI(), "cage") }
func potDripMethod() string  { return getSolidityFunctionSignature(PotABI(), "drip") }
func potExitMethod() string  { return getSolidityFunctionSignature(PotABI(), "exit") }
func potFileDSRMethod() string {
	return getOverloadedFunctionSignature(PotABI(), "file", []string{"bytes32", "uint256"})
}
func potFileVowMethod() string {
	return getOverloadedFunctionSignature(PotABI(), "file", []string{"bytes32", "address"})
}
func potJoinMethod() string { return getSolidityFunctionSignature(PotABI(), "join") }
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
func vatInitMethod() string { return getSolidityFunctionSignature(VatABI(), "init") }
func vatMoveMethod() string { return getSolidityFunctionSignature(VatABI(), "move") }
func vatSlipMethod() string { return getSolidityFunctionSignature(VatABI(), "slip") }
func vatSuckMethod() string { return getSolidityFunctionSignature(VatABI(), "suck") }
func vowFessMethod() string { return getSolidityFunctionSignature(VowABI(), "fess") }
func vowFileMethod() string { return getSolidityFunctionSignature(VowABI(), "file") }
func vowFlogMethod() string { return getSolidityFunctionSignature(VowABI(), "flog") }
func yankMethod() string    { return getSolidityFunctionSignature(FlipABI(), "yank") }
