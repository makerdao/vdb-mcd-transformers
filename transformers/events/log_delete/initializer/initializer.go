package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/log_delete"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.LogDeleteTable, constants.LogDeleteSignature()),
	Transformer: log_delete.Transformer{},
}.NewTransformer
