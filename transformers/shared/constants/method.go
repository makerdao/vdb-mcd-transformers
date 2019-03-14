package constants

//TODO: get cat, pit, and drip file method signatures directly from the ABI
func biteMethod() string               { return GetSolidityFunctionSignature(CatABI(), "Bite") }
func catFileChopLumpMethod() string    { return "file(bytes32,bytes32,uint256)" }
func catFileFlipMethod() string        { return GetSolidityFunctionSignature(CatABI(), "file") }
func catFilePitVowMethod() string      { return "file(bytes32,address)" }
func dealMethod() string               { return GetSolidityFunctionSignature(FlipperABI(), "deal") }
func dentMethod() string               { return GetSolidityFunctionSignature(FlipperABI(), "dent") }
func dripDripMethod() string           { return GetSolidityFunctionSignature(DripABI(), "drip") }
func dripFileIlkMethod() string        { return "file(bytes32,bytes32,uint256)" }
func dripFileRepoMethod() string       { return GetSolidityFunctionSignature(DripABI(), "file") }
func dripFileVowMethod() string        { return "file(bytes32,bytes32)" }
func flapKickMethod() string           { return GetSolidityFunctionSignature(FlapperABI(), "Kick") }
func flipKickMethod() string           { return GetSolidityFunctionSignature(FlipperABI(), "Kick") }
func flopKickMethod() string           { return GetSolidityFunctionSignature(FlopperABI(), "Kick") }
func logValueMethod() string           { return GetSolidityFunctionSignature(MedianizerABI(), "LogValue") }
func pitFileDebtCeilingMethod() string { return "file(bytes32,uint256)" }
func pitFileIlkMethod() string         { return "file(bytes32,bytes32,uint256)" }
func tendMethod() string               { return GetSolidityFunctionSignature(FlipperABI(), "tend") }
func vatFluxMethod() string            { return GetSolidityFunctionSignature(OldVatABI(), "flux") }
func vatFoldMethod() string            { return GetSolidityFunctionSignature(OldVatABI(), "fold") }
func vatFrobMethod() string            { return GetSolidityFunctionSignature(VatABI(), "frob") }
func vatGrabMethod() string            { return GetSolidityFunctionSignature(OldVatABI(), "grab") }
func vatHealMethod() string            { return GetSolidityFunctionSignature(OldVatABI(), "heal") }
func vatInitMethod() string            { return GetSolidityFunctionSignature(OldVatABI(), "init") }
func vatMoveMethod() string            { return GetSolidityFunctionSignature(OldVatABI(), "move") }
func vatSlipMethod() string            { return GetSolidityFunctionSignature(OldVatABI(), "slip") }
func vatTuneMethod() string            { return GetSolidityFunctionSignature(OldVatABI(), "tune") }
func vowFessMethod() string            { return GetSolidityFunctionSignature(VowABI(), "fess") }
func vowFlogMethod() string            { return GetSolidityFunctionSignature(VowABI(), "flog") }
