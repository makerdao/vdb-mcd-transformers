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

func BiteSignature() string               { return getEventTopicZero(biteMethod()) }
func CatFileChopLumpSignature() string    { return getLogNoteTopicZero(catFileChopLumpMethod()) }
func CatFileFlipSignature() string        { return getLogNoteTopicZero(catFileFlipMethod()) }
func CatFileVowSignature() string         { return getLogNoteTopicZero(catFileVowMethod()) }
func DealSignature() string               { return getLogNoteTopicZero(dealMethod()) }
func DentSignature() string               { return getLogNoteTopicZero(dentMethod()) }
func FlapKickSignature() string           { return getEventTopicZero(flapKickMethod()) }
func FlipKickSignature() string           { return getEventTopicZero(flipKickMethod()) }
func FlopKickSignature() string           { return getEventTopicZero(flopKickMethod()) }
func JugDripSignature() string            { return getLogNoteTopicZero(jugDripMethod()) }
func JugFileBaseSignature() string        { return getLogNoteTopicZero(jugFileBaseMethod()) }
func JugFileIlkSignature() string         { return getLogNoteTopicZero(jugFileIlkMethod()) }
func JugFileVowSignature() string         { return getLogNoteTopicZero(jugFileVowMethod()) }
func JugInitSignature() string            { return getLogNoteTopicZero(jugInitMethod()) }
func NewCdpSignature() string             { return getEventTopicZero(newCdpMethod()) }
func SpotFileMatSignature() string        { return getLogNoteTopicZero(spotFileMatMethod()) }
func SpotFilePipSignature() string        { return getLogNoteTopicZero(spotFilePipMethod()) }
func SpotPokeSignature() string           { return getEventTopicZero(spotPokeMethod()) }
func TendSignature() string               { return getLogNoteTopicZero(tendMethod()) }
func TickSignature() string               { return getLogNoteTopicZero(tickMethod()) }
func VatFileDebtCeilingSignature() string { return getLogNoteTopicZero(vatFileDebtCeilingMethod()) }
func VatFileIlkSignature() string         { return getLogNoteTopicZero(vatFileIlkMethod()) }
func VatFluxSignature() string            { return getLogNoteTopicZero(vatFluxMethod()) }
func VatFoldSignature() string            { return getLogNoteTopicZero(vatFoldMethod()) }
func VatForkSignature() string            { return getLogNoteTopicZero(vatForkMethod()) }
func VatFrobSignature() string            { return getLogNoteTopicZero(vatFrobMethod()) }
func VatGrabSignature() string            { return getLogNoteTopicZero(vatGrabMethod()) }
func VatHealSignature() string            { return getLogNoteTopicZero(vatHealMethod()) }
func VatInitSignature() string            { return getLogNoteTopicZero(vatInitMethod()) }
func VatMoveSignature() string            { return getLogNoteTopicZero(vatMoveMethod()) }
func VatSlipSignature() string            { return getLogNoteTopicZero(vatSlipMethod()) }
func VatSuckSignature() string            { return getLogNoteTopicZero(vatSuckMethod()) }
func VowFessSignature() string            { return getLogNoteTopicZero(vowFessMethod()) }
func VowFileSignature() string            { return getLogNoteTopicZero(vowFileMethod()) }
func VowFlogSignature() string            { return getLogNoteTopicZero(vowFlogMethod()) }
func YankSignature() string               { return getLogNoteTopicZero(yankMethod()) }
