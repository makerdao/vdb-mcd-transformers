package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_exit"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/libraries/shared/transformer"
)

var EventTransformerInitializer transformer.EventTransformerInitializer = event.Transformer{
	Config:    shared.GetEventTransformerConfig(constants.PotExitTable, constants.PotExitSignature()),
	Converter: pot_exit.Converter{},
}.NewTransformer
