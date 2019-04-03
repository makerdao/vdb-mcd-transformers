package constants

//TODO: get cat and jug file method signatures directly from the ABI
func biteMethod() string            { return GetSolidityFunctionSignature(CatABI(), "Bite") }
func catFileChopLumpMethod() string { return "file(bytes32,bytes32,uint256)" }
func catFileFlipMethod() string     { return GetSolidityFunctionSignature(CatABI(), "file") }
func catFileVowMethod() string      { return "file(bytes32,address)" }
func dealMethod() string            { return GetSolidityFunctionSignature(FlipperABI(), "deal") }
func dentMethod() string            { return GetSolidityFunctionSignature(FlipperABI(), "dent") }
func flapKickMethod() string        { return GetSolidityFunctionSignature(FlapperABI(), "Kick") }
func flipKickMethod() string        { return GetSolidityFunctionSignature(FlipperABI(), "Kick") }
func flopKickMethod() string        { return GetSolidityFunctionSignature(FlopperABI(), "Kick") }
func jugDripMethod() string         { return GetSolidityFunctionSignature(JugABI(), "drip") }
func jugFileBaseMethod() string     { return "file(bytes32,uint256)" }
func jugFileIlkMethod() string      { return "file(bytes32,bytes32,uint256)" }
func jugFileVowMethod() string      { return "file(bytes32,bytes32)" }
func logMedianPriceMethod() string {
	return GetSolidityFunctionSignature(MedianizerABI(), "LogMedianPrice")
}
func tendMethod() string               { return GetSolidityFunctionSignature(FlipperABI(), "tend") }
func vatFileDebtCeilingMethod() string { return "file(bytes32,uint256)" }
func vatFileIlkMethod() string         { return "file(bytes32,bytes32,uint256)" }
func vatFluxMethod() string            { return GetSolidityFunctionSignature(VatABI(), "flux") }
func vatFoldMethod() string            { return GetSolidityFunctionSignature(VatABI(), "fold") }
func vatFrobMethod() string            { return GetSolidityFunctionSignature(VatABI(), "frob") }
func vatGrabMethod() string            { return GetSolidityFunctionSignature(VatABI(), "grab") }
func vatHealMethod() string            { return GetSolidityFunctionSignature(VatABI(), "heal") }
func vatInitMethod() string            { return GetSolidityFunctionSignature(VatABI(), "init") }
func vatMoveMethod() string            { return GetSolidityFunctionSignature(VatABI(), "move") }
func vatSlipMethod() string            { return GetSolidityFunctionSignature(VatABI(), "slip") }
func vowFessMethod() string            { return GetSolidityFunctionSignature(VowABI(), "fess") }
func vowFlogMethod() string            { return GetSolidityFunctionSignature(VowABI(), "flog") }
