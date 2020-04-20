package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_trade"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.LogTradeTable, constants.LogTradeSignature()),
	Transformer: log_trade.Transformer{},
}.NewTransformer
