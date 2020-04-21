package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_min_sell"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.LogMinSellTable, constants.LogMinSellSignature()),
	Transformer: log_min_sell.Transformer{},
}.NewTransformer
