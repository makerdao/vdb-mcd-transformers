package initializer

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/cat_claw"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-transformer-utilities/pkg/shared"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
)

var EventTransformerInitializer event.TransformerInitializer = event.ConfiguredTransformer{
	Config:      shared.GetEventTransformerConfig(constants.CatClawTable, constants.CatClawSignature()),
	Transformer: cat_claw.Transformer{},
}.NewTransformer
