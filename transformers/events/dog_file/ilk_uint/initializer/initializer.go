package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/ilk_uint"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.DogFileIlkUintTable, constants.DogFileIlkUintSignature()),
	Transformer: ilk_uint.Transformer{},
}.NewTransformer
