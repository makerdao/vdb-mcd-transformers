package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/median_kiss/batch"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.MedianKissBatchTable, constants.MedianKissBatchSignature()),
	Transformer: batch.Transformer{},
}.NewTransformer
