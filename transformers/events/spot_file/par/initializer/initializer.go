package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/spot_file/par"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.SpotFileParTable, constants.SpotFileParSignature()),
	Transformer: par.Transformer{},
}.NewTransformer
