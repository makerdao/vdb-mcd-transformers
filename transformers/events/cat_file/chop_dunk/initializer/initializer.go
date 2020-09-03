package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_file/chop_dunk"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.CatFileChopDunkTable, constants.CatFileChopDunkSignature()),
	Transformer: chop_dunk.Transformer{},
}.NewTransformer
