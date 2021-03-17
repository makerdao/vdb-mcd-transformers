package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/dog_file/ilk_chop_hole"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.DogFileIlkChopHoleTable, constants.DogFileIlkChopHoleSignature()),
	Transformer: ilk_chop_hole.Transformer{},
}.NewTransformer
