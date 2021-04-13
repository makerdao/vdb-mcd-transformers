package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_digs"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.DogDigsTable, constants.DogDigsSignature()),
	Transformer: dog_digs.Transformer{},
}.NewTransformer
