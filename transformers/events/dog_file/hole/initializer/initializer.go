package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/hole"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.DogFileHoleTable, constants.DogFileHoleSignature()),
	Transformer: hole.Transformer{},
}.NewTransformer
