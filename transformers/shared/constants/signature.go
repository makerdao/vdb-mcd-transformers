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

func BiteSignature() string { return getEventTopicZero(biteMethod()) }
func CatFileChopLumpSignature() string {
	return getLogNoteTopicZeroWithZeroPadding(catFileChopLumpMethod())
}
func CatFileFlipSignature() string { return getLogNoteTopicZeroWithZeroPadding(catFileFlipMethod()) }
func CatFileVowSignature() string  { return getLogNoteTopicZeroWithZeroPadding(catFileVowMethod()) }
func DealSignature() string        { return getLogNoteTopicZeroWithZeroPadding(dealMethod()) }
func DentSignature() string        { return getLogNoteTopicZeroWithZeroPadding(dentMethod()) }
func FlapKickSignature() string    { return getEventTopicZero(flapKickMethod()) }
func FlipKickSignature() string    { return getEventTopicZero(flipKickMethod()) }
func FlopKickSignature() string    { return getEventTopicZero(flopKickMethod()) }
func JugDripSignature() string     { return getLogNoteTopicZeroWithZeroPadding(jugDripMethod()) }
func JugFileIlkSignature() string  { return getLogNoteTopicZeroWithZeroPadding(jugFileIlkMethod()) }
func JugFileBaseSignature() string { return getLogNoteTopicZeroWithZeroPadding(jugFileBaseMethod()) }
func JugFileVowSignature() string  { return getLogNoteTopicZeroWithZeroPadding(jugFileVowMethod()) }
func PipLogValueSignature() string { return getEventTopicZero(pipLogValueMethod()) }
func TendSignature() string        { return getLogNoteTopicZeroWithZeroPadding(tendMethod()) }
func VatFileDebtCeilingSignature() string {
	return getLogNoteTopicZeroWithLeadingZeros(vatFileDebtCeilingMethod())
}
func VatFileIlkSignature() string { return getLogNoteTopicZeroWithLeadingZeros(vatFileIlkMethod()) }
func VatFluxSignature() string    { return getLogNoteTopicZeroWithLeadingZeros(vatFluxMethod()) }
func VatFoldSignature() string    { return getLogNoteTopicZeroWithLeadingZeros(vatFoldMethod()) }
func VatFrobSignature() string    { return getLogNoteTopicZeroWithLeadingZeros(vatFrobMethod()) }
func VatGrabSignature() string    { return getLogNoteTopicZeroWithLeadingZeros(vatGrabMethod()) }
func VatHealSignature() string    { return getLogNoteTopicZeroWithLeadingZeros(vatHealMethod()) }
func VatInitSignature() string    { return getLogNoteTopicZeroWithLeadingZeros(vatInitMethod()) }
func VatMoveSignature() string    { return getLogNoteTopicZeroWithLeadingZeros(vatMoveMethod()) }
func VatSlipSignature() string    { return getLogNoteTopicZeroWithLeadingZeros(vatSlipMethod()) }
func VowFessSignature() string    { return getLogNoteTopicZeroWithZeroPadding(vowFessMethod()) }
func VowFileSignature() string    { return getLogNoteTopicZeroWithZeroPadding(vowFileMethod()) }
func VowFlogSignature() string    { return getLogNoteTopicZeroWithZeroPadding(vowFlogMethod()) }
