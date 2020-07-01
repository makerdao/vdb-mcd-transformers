package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_heal"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.VowHealTable, constants.VowHealSignature()),
	Transformer: vow_heal.Transformer{},
}.NewTransformer
