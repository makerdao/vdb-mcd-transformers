package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/dunk"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.CatFileDunkTable, constants.CatFileDunkSignature()),
	Transformer: dunk.Transformer{},
}.NewTransformer
