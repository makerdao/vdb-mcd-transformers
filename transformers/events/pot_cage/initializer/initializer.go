package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/pot_cage"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

// EventTransformerInitializer is the ConfiguredTransformer for pot_cage events
var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.PotCageTable, constants.PotCageSignature()),
	Transformer: pot_cage.Transformer{},
}.NewTransformer
