package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/median_diss/batch"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.MedianDissBatchTable, constants.MedianDissBatchSignature()),
	Transformer: batch.Transformer{},
}.NewTransformer
