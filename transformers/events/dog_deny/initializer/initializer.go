package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_deny"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.DogDenyTable, constants.DogDenySignature()),
	Transformer: dog_deny.Transformer{},
}.NewTransformer
