package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/median_drop"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.MedianDropTable, constants.MedianDropSignature()),
	Transformer: median_drop.Transformer{},
}.NewTransformer
