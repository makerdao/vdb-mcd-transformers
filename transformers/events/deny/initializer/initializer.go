package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/deny"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

var EventTransformerInitializer transformer.EventTransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.DenyTable, constants.DenySignature()),
	Transformer: deny.Transformer{},
}.NewTransformer
