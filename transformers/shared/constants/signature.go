package constants

func GetBiteSignature() string { return GetEventTopicZero(biteMethod()) }
func GetCatFileChopLumpSignature() string {
	return GetLogNoteTopicZeroWithZeroPadding(catFileChopLumpMethod())
}
func GetCatFileFlipSignature() string  { return GetLogNoteTopicZeroWithZeroPadding(catFileFlipMethod()) }
func GetCatFileVowSignature() string   { return GetLogNoteTopicZeroWithZeroPadding(catFileVowMethod()) }
func GetDealSignature() string         { return GetLogNoteTopicZeroWithZeroPadding(dealMethod()) }
func GetDentFunctionSignature() string { return GetLogNoteTopicZeroWithZeroPadding(dentMethod()) }
func GetFlapKickSignature() string     { return GetEventTopicZero(flapKickMethod()) }
func GetFlipKickSignature() string     { return GetEventTopicZero(flipKickMethod()) }
func GetFlopKickSignature() string     { return GetEventTopicZero(flopKickMethod()) }
func GetJugDripSignature() string      { return GetLogNoteTopicZeroWithZeroPadding(jugDripMethod()) }
func GetJugFileIlkSignature() string   { return GetLogNoteTopicZeroWithZeroPadding(jugFileIlkMethod()) }
func GetJugFileBaseSignature() string  { return GetLogNoteTopicZeroWithZeroPadding(jugFileBaseMethod()) }
func GetJugFileVowSignature() string   { return GetLogNoteTopicZeroWithZeroPadding(jugFileVowMethod()) }
func GetPipLogValueSignature() string  { return GetEventTopicZero(pipLogValueMethod()) }
func GetTendFunctionSignature() string { return GetLogNoteTopicZeroWithZeroPadding(tendMethod()) }
func GetVatFileDebtCeilingSignature() string {
	return GetLogNoteTopicZeroWithLeadingZeros(vatFileDebtCeilingMethod())
}
func GetVatFileIlkSignature() string { return GetLogNoteTopicZeroWithLeadingZeros(vatFileIlkMethod()) }
func GetVatFluxSignature() string    { return GetLogNoteTopicZeroWithLeadingZeros(vatFluxMethod()) }
func GetVatFoldSignature() string    { return GetLogNoteTopicZeroWithLeadingZeros(vatFoldMethod()) }
func GetVatFrobSignature() string    { return GetLogNoteTopicZeroWithLeadingZeros(vatFrobMethod()) }
func GetVatGrabSignature() string    { return GetLogNoteTopicZeroWithLeadingZeros(vatGrabMethod()) }
func GetVatHealSignature() string    { return GetLogNoteTopicZeroWithLeadingZeros(vatHealMethod()) }
func GetVatInitSignature() string    { return GetLogNoteTopicZeroWithLeadingZeros(vatInitMethod()) }
func GetVatMoveSignature() string    { return GetLogNoteTopicZeroWithLeadingZeros(vatMoveMethod()) }
func GetVatSlipSignature() string    { return GetLogNoteTopicZeroWithLeadingZeros(vatSlipMethod()) }
func GetVowFessSignature() string    { return GetLogNoteTopicZeroWithZeroPadding(vowFessMethod()) }
func GetVowFlogSignature() string    { return GetLogNoteTopicZeroWithZeroPadding(vowFlogMethod()) }
