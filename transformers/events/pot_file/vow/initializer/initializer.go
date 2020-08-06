package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_file/vow"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

// EventTransformerInitializer is the ConfiguredTransformer for pot_file_vow events
var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.PotFileVowTable, constants.PotFileVowSignature()),
	Transformer: vow.Transformer{},
}.NewTransformer
