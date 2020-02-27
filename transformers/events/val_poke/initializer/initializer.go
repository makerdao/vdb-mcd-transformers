package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/val_poke"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.ValPokeTable, constants.ValPokeSignature()),
	Transformer: val_poke.Transformer{},
}.NewTransformer
